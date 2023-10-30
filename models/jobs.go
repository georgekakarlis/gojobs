package models

import (
	"time"

	"gorm.io/gorm"
)




type Job struct {
	gorm.Model
	ID          uint
	Title       string
	Description string
	PostedBy    *User       // Foreign key relationship to the User who posted the job
	PostedByID  uint        // Foreign key ID
	Location    string
	Salary      *float64    
	CreatedAt   time.Time
	UpdatedAt   time.Time
}