package models

import (
	"time"
)

type Project struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"unique;not null"`
	DueDate     *time.Time `json:"due_date"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Tasks       []Task     `json:"tasks" gorm:"foreignKey:ProjectID"`
	SortOrder   int        `json:"sort_order"`
}

func (p *Project) TableName() string {
	return "projects"
}
