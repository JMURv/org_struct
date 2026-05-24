package validation

import (
	"github.com/JMURv/org-struct/internal/dto"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var V = validator.New()

func init() {
	V = validator.New()

	V.RegisterStructValidation(deleteDepartmentValidation, dto.DeleteDepartmentQuery{})
}

func deleteDepartmentValidation(sl validator.StructLevel) {
	query, ok := sl.Current().Interface().(dto.DeleteDepartmentQuery)
	if !ok {
		zap.L().Error("deleteDepartmentValidation: cannot convert to deleteDepartmentQuery")
		sl.ReportError(
			query.ReassignToDepartmentID,
			"reassign_to_department_id",
			"ReassignToDepartmentID",
			"required_when_reassign",
			"",
		)
		return
	}

	if query.Mode == "reassign" && query.ReassignToDepartmentID <= 0 {
		sl.ReportError(
			query.ReassignToDepartmentID,
			"reassign_to_department_id",
			"ReassignToDepartmentID",
			"required_when_reassign",
			"",
		)
	}
}
