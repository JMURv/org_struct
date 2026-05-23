package dto

import (
	"time"

	"github.com/JMURv/golang-clean-template/internal/models"
)

type CreateDepartmentRequest struct {
	Name     string `json:"name" validate:"required"`
	ParentID *int   `json:"parent_id"`
}

type CreateDepartmentEmployeeRequest struct {
	FullName string     `json:"full_name" validate:"required"`
	Position string     `json:"position" validate:"required"`
	HiredAt  *time.Time `json:"hired_at"`
}

type GetDepartmentQuery struct {
	Depth            int  `json:"depth" validate:"required"` // 1 by default
	IncludeEmployees bool `json:"include_employees"`         // true by default
}

type GetDepartmentResponse struct {
	Department models.Department   `json:"department"`
	Employees  []models.Employee   `json:"employees"`
	Children   []models.Department `json:"children"`
}

type UpdateDepartmentRequest struct {
	Name     *string `json:"name"`
	ParentID *int    `json:"parent_id"`
}

type DeleteDepartmentQuery struct {
	Mode                   string `json:"name" validate:"required"`  // "cascade" OR "reassign"
	ReassignToDepartmentID int    `json:"reassign_to_department_id"` // required if mode == reassign
}
