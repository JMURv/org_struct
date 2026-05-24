package models

import (
	"time"
)

type Department struct {
	ID        uint64    `gorm:"primaryKey"        json:"id"`
	Name      string    `gorm:"size:200;not null" json:"name"`
	ParentID  *uint64   `                         json:"parent_id"`
	CreatedAt time.Time `                         json:"created_at"`

	Parent    *Department  `gorm:"foreignKey:ParentID" json:"-"`
	Children  []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Employees []Employee   `                           json:"employees,omitempty"`
}

func (Department) TableName() string {
	return "department"
}

type Employee struct {
	ID           uint64     `gorm:"primaryKey"        json:"id"`
	DepartmentID uint64     `                         json:"department_id"`
	FullName     string     `gorm:"size:200;not null" json:"full_name"`
	Position     string     `gorm:"size:200;not null" json:"position"`
	HiredAt      *time.Time `                         json:"hired_at"`
	CreatedAt    time.Time  `                         json:"created_at"`
}

func (Employee) TableName() string {
	return "employee"
}
