package models

import (
	"time"
)

type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	ProjectID   uint       `json:"project_id"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Status      string     `json:"status" gorm:"default:Pending"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (t *Task) TableName() string {
	return "tasks"
}
