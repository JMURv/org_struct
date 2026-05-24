package ctrl

import "errors"

// ErrNotFound is returned when a resource is not found.
var ErrNotFound = errors.New("not found")

var ErrDepartmentCycle = errors.New("department cycle")
