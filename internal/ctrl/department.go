package ctrl

import (
	"context"
	"errors"

	"github.com/JMURv/org-struct/internal/dto"
	md "github.com/JMURv/org-struct/internal/models"
	"github.com/JMURv/org-struct/internal/repo"
	"github.com/opentracing/opentracing-go"
)

type departmentCtrl interface {
	CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*md.Department, error)
	GetDepartment(
		ctx context.Context,
		id uint64,
		req dto.GetDepartmentQuery,
	) (*md.Department, error)
	UpdateDepartment(
		ctx context.Context,
		id uint64,
		req dto.UpdateDepartmentRequest,
	) (*md.Department, error)
	DeleteDepartment(ctx context.Context, id uint64, req dto.DeleteDepartmentQuery) error

	CreateEmployee(
		ctx context.Context,
		id uint64,
		req dto.CreateDepartmentEmployeeRequest,
	) (*md.Employee, error)
}

type departmentRepo interface {
	CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*md.Department, error)
	GetDepartment(
		ctx context.Context,
		id uint64,
		req dto.GetDepartmentQuery,
	) (*md.Department, error)
	UpdateDepartment(
		ctx context.Context,
		id uint64,
		req dto.UpdateDepartmentRequest,
	) (*md.Department, error)
	DeleteDepartment(ctx context.Context, id uint64, req dto.DeleteDepartmentQuery) error

	CreateEmployee(
		ctx context.Context,
		id uint64,
		req dto.CreateDepartmentEmployeeRequest,
	) (*md.Employee, error)
}

func (c *Controller) CreateDepartment(
	ctx context.Context,
	req dto.CreateDepartmentRequest,
) (*md.Department, error) {
	const op = "department.CreateDepartment.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := c.repo.CreateDepartment(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Controller) GetDepartment(
	ctx context.Context,
	id uint64,
	req dto.GetDepartmentQuery,
) (*md.Department, error) {
	const op = "department.GetDepartment.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := c.repo.GetDepartment(ctx, id, req)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return res, nil
}

func (c *Controller) UpdateDepartment(
	ctx context.Context,
	id uint64,
	req dto.UpdateDepartmentRequest,
) (*md.Department, error) {
	const op = "department.UpdateDepartment.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := c.repo.UpdateDepartment(ctx, id, req)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrNotFound
		}

		if errors.Is(err, repo.ErrDepartmentCycle) {
			return nil, ErrDepartmentCycle
		}
		return nil, err
	}

	return res, nil
}

func (c *Controller) DeleteDepartment(
	ctx context.Context,
	id uint64,
	req dto.DeleteDepartmentQuery,
) error {
	const op = "department.DeleteDepartment.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	if err := c.repo.DeleteDepartment(ctx, id, req); err != nil {
		return err
	}

	return nil
}

func (c *Controller) CreateEmployee(
	ctx context.Context,
	id uint64,
	req dto.CreateDepartmentEmployeeRequest,
) (*md.Employee, error) {
	const op = "department.CreateEmployee.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()

	res, err := c.repo.CreateEmployee(ctx, id, req)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return res, nil
}
