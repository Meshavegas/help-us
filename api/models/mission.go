package models

import (
	"time"

	"gorm.io/gorm"
)

// MissionStatus represents the status of a mission
type MissionStatus string

const (
	MissionStatusActive    MissionStatus = "active"
	MissionStatusCompleted MissionStatus = "completed"
	MissionStatusStopped   MissionStatus = "stopped"
	MissionStatusPaused    MissionStatus = "paused"
)

// Mission model - represents a teaching mission
type Mission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	StartDate   time.Time      `json:"start_date" gorm:"not null"`
	EndDate     *time.Time     `json:"end_date"`
	Status      MissionStatus  `json:"status" gorm:"default:'active'"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign Keys
	FamilleID    uint `json:"famille_id"`
	EnseignantID uint `json:"enseignant_id"`

	// Relationships
	Famille    Famille    `json:"famille,omitempty" gorm:"foreignKey:FamilleID"`
	Enseignant Enseignant `json:"enseignant,omitempty" gorm:"foreignKey:EnseignantID"`
	Courses    []Course   `json:"courses,omitempty" gorm:"foreignKey:MissionID"`
	Reports    []Report   `json:"reports,omitempty" gorm:"foreignKey:MissionID"`
}

// Mission methods
func (m *Mission) CreateMission() error {
	// Logic to create mission will be implemented in controllers
	m.Status = MissionStatusActive
	return nil
}

func (m *Mission) StopMission() error {
	// Logic to stop mission
	m.Status = MissionStatusStopped
	endTime := time.Now()
	m.EndDate = &endTime
	return nil
}

func (m *Mission) ExtendMission(newEndDate time.Time) error {
	// Logic to extend mission
	m.EndDate = &newEndDate
	return nil
}

// Request/Response structures
type MissionCreateRequest struct {
	StartDate    time.Time  `json:"start_date" binding:"required"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	Description  string     `json:"description"`
	EnseignantID uint       `json:"enseignant_id" binding:"required"`
}

type MissionUpdateRequest struct {
	EndDate     *time.Time    `json:"end_date,omitempty"`
	Status      MissionStatus `json:"status,omitempty"`
	Description string        `json:"description,omitempty"`
}

type MissionFilterRequest struct {
	Status       MissionStatus `json:"status,omitempty"`
	EnseignantID uint          `json:"enseignant_id,omitempty"`
	FamilleID    uint          `json:"famille_id,omitempty"`
	DateFrom     *time.Time    `json:"date_from,omitempty"`
	DateTo       *time.Time    `json:"date_to,omitempty"`
}
