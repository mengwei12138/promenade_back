package models

import (
	"time"
)

type Requirement struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProjectID uint      `json:"project_id" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Proposer  string    `json:"proposer" gorm:"size:64;not null"`
	Status    string    `json:"status" gorm:"type:varchar(16);default:Pending"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Requirement) TableName() string {
	return "requirements"
}
