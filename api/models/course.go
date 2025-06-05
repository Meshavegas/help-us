package models

import (
	"time"

	"gorm.io/gorm"
)

// CourseStatus represents the status of a course
type CourseStatus string

const (
	CourseStatusScheduled  CourseStatus = "scheduled"
	CourseStatusCompleted  CourseStatus = "completed"
	CourseStatusCancelled  CourseStatus = "cancelled"
	CourseStatusInProgress CourseStatus = "in_progress"
)

// Course model - represents a scheduled course/session
type Course struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	ScheduledTime time.Time      `json:"scheduled_time" gorm:"not null"`
	Duration      int            `json:"duration"` // Duration in minutes
	Location      string         `json:"location"`
	Status        CourseStatus   `json:"status" gorm:"default:'scheduled'"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign Keys
	FamilleID    uint `json:"famille_id"`
	EnseignantID uint `json:"enseignant_id"`
	MissionID    uint `json:"mission_id"`
	AddressID    uint `json:"address_id"`

	// Relationships
	Famille    Famille    `json:"famille,omitempty" gorm:"foreignKey:FamilleID"`
	Enseignant Enseignant `json:"enseignant,omitempty" gorm:"foreignKey:EnseignantID"`
	Mission    Mission    `json:"mission,omitempty" gorm:"foreignKey:MissionID"`
	Address    Address    `json:"address,omitempty" gorm:"foreignKey:AddressID"`
	Payments   []Payment  `json:"payments,omitempty" gorm:"foreignKey:CourseID"`
}

// Course methods
func (c *Course) Schedule() error {
	c.Status = CourseStatusScheduled
	return nil
}

func (c *Course) Cancel() error {
	c.Status = CourseStatusCancelled
	return nil
}

func (c *Course) Validate() error {
	c.Status = CourseStatusCompleted
	return nil
}

func (c *Course) Declare() error {
	c.Status = CourseStatusInProgress
	return nil
}

// Request/Response structures
type CourseCreateRequest struct {
	ScheduledTime time.Time `json:"scheduled_time" binding:"required"`
	Duration      int       `json:"duration" binding:"required,min=30,max=480"` // 30 min to 8 hours
	Location      string    `json:"location" binding:"required"`
	EnseignantID  uint      `json:"enseignant_id" binding:"required"`
	AddressID     uint      `json:"address_id" binding:"required"`
}

type CourseUpdateRequest struct {
	ScheduledTime *time.Time   `json:"scheduled_time,omitempty"`
	Duration      *int         `json:"duration,omitempty"`
	Location      string       `json:"location,omitempty"`
	Status        CourseStatus `json:"status,omitempty"`
}

type CourseFilterRequest struct {
	Status       CourseStatus `json:"status,omitempty"`
	EnseignantID uint         `json:"enseignant_id,omitempty"`
	DateFrom     *time.Time   `json:"date_from,omitempty"`
	DateTo       *time.Time   `json:"date_to,omitempty"`
}
