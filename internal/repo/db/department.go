package db

import (
	"context"
	"errors"

	"github.com/JMURv/golang-clean-template/internal/config"
	md "github.com/JMURv/golang-clean-template/internal/models"
	"github.com/JMURv/golang-clean-template/internal/repo"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func (r *Repository) GetDepartment(ctx context.Context, id uint64, depth int) (*md.Department, error) {
	const op = "department.GetDepartment.repo"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	var department md.Department

	err := r.conn.WithContext(ctx).
		Preload("Employees").
		First(&department, id).Error

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

	if depth > 0 {
		if err = r.loadChildren(ctx, &department, depth); err != nil {
			return nil, err
		}
	}

	return &department, nil
}

func (r *Repository) loadChildren(ctx context.Context, department *md.Department, depth int) error {
	const op = "department.loadChildren.repo"

	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	if depth <= 0 {
		return nil
	}

	err := r.conn.WithContext(ctx).
		Preload("Employees").
		Where("parent_id = ?", department.ID).
		Find(&department.Children).Error

	if err != nil {
		zap.L().Error(
			"failed to get department children",
			zap.String("op", op),
			zap.Any("department", department.Name),
			zap.Error(err),
		)
		return err
	}

	for i := range department.Children {
		if err = r.loadChildren(ctx, &department.Children[i], depth-1); err != nil {
			return err
		}
	}

	return nil
}
