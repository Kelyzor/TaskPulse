package models

import "time"

type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UserID      uint       `json:"user_id" gorm:"index"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status" gorm:"default:'pending'"`
	Priority    string     `json:"priority" gorm:"default:'medium'"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
