package ctrl

import (
	"context"

	"github.com/JMURv/golang-clean-template/internal/dto"
	md "github.com/JMURv/golang-clean-template/internal/models"
	"github.com/opentracing/opentracing-go"
)

type departmentCtrl interface {
	CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (md.Department, error)
	GetDepartment(ctx context.Context, id uint64, req dto.GetDepartmentQuery) (dto.GetDepartmentResponse, error)
	UpdateDepartment(ctx context.Context, id uint64, req dto.UpdateDepartmentRequest) (md.Department, error)
	DeleteDepartment(ctx context.Context, id uint64, req dto.DeleteDepartmentQuery) error

	CreateEmployee(ctx context.Context, id uint64, req dto.CreateDepartmentEmployeeRequest) (md.Employee, error)
}

type departmentRepo interface {
	CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (md.Department, error)
	GetDepartment(ctx context.Context, id uint64, req dto.GetDepartmentQuery) (dto.GetDepartmentResponse, error)
	UpdateDepartment(ctx context.Context, id uint64, req dto.UpdateDepartmentRequest) (md.Department, error)
	DeleteDepartment(ctx context.Context, id uint64, req dto.DeleteDepartmentQuery) error

	CreateEmployee(ctx context.Context, id uint64, req dto.CreateDepartmentEmployeeRequest) (md.Employee, error)
}

func (c *Controller) CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (md.Department, error) {
	const op = "department.CreateDepartment.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()
}

func (c *Controller) GetDepartment(ctx context.Context, id uint64, req dto.GetDepartmentQuery) (dto.GetDepartmentResponse, error) {
	const op = "department.GetDepartment.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()
}

func (c *Controller) UpdateDepartment(ctx context.Context, id uint64, req dto.UpdateDepartmentRequest) (md.Department, error) {
	const op = "department.UpdateDepartment.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()
}

func (c *Controller) DeleteDepartment(ctx context.Context, id uint64, req dto.DeleteDepartmentQuery) error {
	const op = "department.DeleteDepartment.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()
}

func (c *Controller) CreateEmployee(ctx context.Context, id uint64, req dto.CreateDepartmentEmployeeRequest) (md.Employee, error) {
	const op = "department.CreateEmployee.ctrl"
	span, ctx := opentracing.StartSpanFromContext(ctx, op)
	defer span.Finish()
}
