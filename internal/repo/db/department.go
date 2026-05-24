package db

import (
	"context"
	"errors"
	"strings"

	"github.com/JMURv/golang-clean-template/internal/config"
	"github.com/JMURv/golang-clean-template/internal/dto"
	md "github.com/JMURv/golang-clean-template/internal/models"
	"github.com/JMURv/golang-clean-template/internal/repo"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (r *Repository) CreateDepartment(
	ctx context.Context,
	req dto.CreateDepartmentRequest,
) (*md.Department, error) {
	res := &md.Department{
		Name:     strings.TrimSpace(req.Name),
		ParentID: req.ParentID,
	}

	if err := r.conn.WithContext(ctx).Create(res).Error; err != nil {
		zap.L().Error("Unable to create department", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (r *Repository) CreateEmployee(
	ctx context.Context,
	id uint64,
	req dto.CreateDepartmentEmployeeRequest,
) (*md.Employee, error) {
	var exists bool

	if err := r.conn.WithContext(ctx).
		Model(&md.Department{}).
		Select("count(*) > 0").
		Where("id = ?", id).
		Find(&exists).Error; err != nil {
		zap.L().Error("error getting department count", zap.Error(err))
		return nil, err
	}

	if !exists {
		zap.L().Info("department does not exist", zap.Uint64("department_id", id))
		return nil, repo.ErrNotFound
	}

	res := &md.Employee{
		DepartmentID: id,
		FullName:     strings.TrimSpace(req.FullName),
		Position:     strings.TrimSpace(req.Position),
		HiredAt:      req.HiredAt,
	}

	if err := r.conn.WithContext(ctx).
		Create(res).Error; err != nil {
		zap.L().Error("error creating employee", zap.Error(err))
		return nil, err
	}

	return res, nil
}

func (r *Repository) GetDepartment(
	ctx context.Context,
	id uint64,
	req dto.GetDepartmentQuery,
) (*md.Department, error) {
	const op = "department.GetDepartment.repo"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var department md.Department

	query := r.conn.WithContext(ctx)
	if req.IncludeEmployees != nil && *req.IncludeEmployees {
		query = query.Preload("Employees")
	}

	err := query.First(&department, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Info(
				"department not found",
				zap.String("op", op),
				zap.Uint64("ID", id),
			)
			return nil, repo.ErrNotFound
		}

		zap.L().Error(
			"failed to get department",
			zap.String("op", op),
			zap.Uint64("ID", id),
			zap.Error(err),
		)
		span.SetTag(config.ErrorSpanTag, true)
		return nil, err
	}

	if req.Depth > 0 {
		if err = r.loadChildren(ctx, &department, req.Depth, *req.IncludeEmployees); err != nil {
			return nil, err
		}
	}

	return &department, nil
}

func (r *Repository) loadChildren(
	ctx context.Context,
	department *md.Department,
	depth int,
	includeEmployees bool,
) error {
	const op = "department.loadChildren.repo"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	if depth <= 0 {
		return nil
	}

	query := r.conn.WithContext(ctx)
	if includeEmployees {
		query = query.Preload("Employees")
	}

	err := query.
		Where("parent_id = ?", department.ID).
		Order("created_at ASC").
		Find(&department.Children).
		Error
	if err != nil {
		zap.L().Error(
			"failed to get department children",
			zap.String("op", op),
			zap.Any("department", department.Name),
			zap.Error(err),
		)

		span.SetTag(config.ErrorSpanTag, true)
		return err
	}

	for i := range department.Children {
		if err = r.loadChildren(
			ctx,
			&department.Children[i],
			depth-1,
			includeEmployees,
		); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) UpdateDepartment(
	ctx context.Context,
	id uint64,
	req dto.UpdateDepartmentRequest,
) (*md.Department, error) {
	const op = "department.UpdateDepartment.repo"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var department md.Department
	if err := r.conn.WithContext(ctx).
		First(&department, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Info("department not found", zap.Uint64("id", id))
			return nil, repo.ErrNotFound
		}

		zap.L().Error("error getting department", zap.Uint64("id", id), zap.Error(err))
		span.SetTag(config.ErrorSpanTag, true)
		return nil, err
	}

	if req.Name != nil {
		department.Name = strings.TrimSpace(*req.Name)
	}

	if req.ParentID != nil {
		if *req.ParentID == id {
			return nil, repo.ErrDepartmentCycle
		}

		if err := r.checkCycle(ctx, id, *req.ParentID); err != nil {
			return nil, err
		}

		department.ParentID = req.ParentID
	}

	if err := r.conn.WithContext(ctx).
		Save(&department).Error; err != nil {
		return nil, err
	}

	return &department, nil
}

func (r *Repository) checkCycle(ctx context.Context, departmentID uint64, parentID uint64) error {
	current := parentID

	for {
		var department md.Department

		err := r.conn.WithContext(ctx).First(&department, current).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		if department.ID == departmentID {
			return repo.ErrDepartmentCycle
		}

		if department.ParentID == nil {
			return nil
		}

		current = *department.ParentID
	}
}

func (r *Repository) DeleteDepartment(
	ctx context.Context,
	id uint64,
	req dto.DeleteDepartmentQuery,
) error {
	const op = "department.DeleteDepartment.repo"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	switch req.Mode {
	case "cascade":
		return r.conn.WithContext(ctx).Delete(&md.Department{}, id).Error
	case "reassign":
		return r.conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&md.Employee{}).
				Where("department_id = ?", id).
				Update("department_id", req.ReassignToDepartmentID).Error; err != nil {
				return err
			}

			if err := tx.Model(&md.Department{}).
				Where("parent_id = ?", id).
				Update("parent_id", req.ReassignToDepartmentID).Error; err != nil {
				return err
			}

			return tx.Delete(&md.Department{}, id).Error
		})
	default:
		return repo.ErrInvalidMode
	}
}
