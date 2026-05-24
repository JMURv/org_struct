package dto

import (
	"time"

	"github.com/JMURv/org-struct/internal/models"
)

type CreateDepartmentRequest struct {
	Name     string  `json:"name"      validate:"required"`
	ParentID *uint64 `json:"parent_id"`
}

type CreateDepartmentEmployeeRequest struct {
	FullName string     `json:"full_name" validate:"required"`
	Position string     `json:"position"  validate:"required"`
	HiredAt  *time.Time `json:"hired_at"`
}

type GetDepartmentQuery struct {
	Depth            int   `json:"depth"             validate:"required"`
	IncludeEmployees *bool `json:"include_employees"`
}

type GetDepartmentResponse struct {
	Department models.Department   `json:"department"`
	Employees  []models.Employee   `json:"employees"`
	Children   []models.Department `json:"children"`
}

type UpdateDepartmentRequest struct {
	Name     *string `json:"name"`
	ParentID *uint64 `json:"parent_id"`
}

type DeleteDepartmentQuery struct {
	Mode                   string `json:"name"                      validate:"required,oneof=cascade reassign"`
	ReassignToDepartmentID int    `json:"reassign_to_department_id"`
}
