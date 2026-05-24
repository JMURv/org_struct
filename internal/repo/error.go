package repo

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrDepartmentCycle = errors.New("department cycle")
	ErrInvalidMode     = errors.New("invalid mode")
)
