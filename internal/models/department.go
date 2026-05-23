package models

import (
	"time"
)

type Department struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:200;not null"`
	ParentID  *int      `json:"parent_id"`
	CreatedAt time.Time `json:"createdAt"`
}

type Employee struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	DepartmentID int       `json:"department_id"`
	FullName     string    `json:"full_name"`
	Position     string    `json:"position"`
	HiredAt      time.Time `json:"hired_at"`
	CreatedAt    time.Time `json:"created_at"`
}
