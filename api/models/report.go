package models

import (
	"time"

	"gorm.io/gorm"
)

// ReportStatus represents the status of a report
type ReportStatus string

const (
	ReportStatusSubmitted ReportStatus = "submitted"
	ReportStatusValidated ReportStatus = "validated"
	ReportStatusRejected  ReportStatus = "rejected"
	ReportStatusPending   ReportStatus = "pending"
)

// Report model - represents reports submitted by teachers
type Report struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	SubmissionDate time.Time      `json:"submission_date"`
	Content        string         `json:"content" gorm:"type:text"`
	Status         ReportStatus   `json:"status" gorm:"default:'pending'"`
	ValidationDate *time.Time     `json:"validation_date"`
	Comments       string         `json:"comments"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign Keys
	EnseignantID  uint `json:"enseignant_id"`
	MissionID     uint `json:"mission_id"`
	ValidatedByID uint `json:"validated_by_id,omitempty"`

	// Relationships
	Enseignant  Enseignant    `json:"enseignant,omitempty" gorm:"foreignKey:EnseignantID"`
	Mission     Mission       `json:"mission,omitempty" gorm:"foreignKey:MissionID"`
	ValidatedBy Administrator `json:"validated_by,omitempty" gorm:"foreignKey:ValidatedByID"`
}

// Report methods
func (r *Report) SubmitReport() error {
	// Logic to submit report
	r.Status = ReportStatusSubmitted
	r.SubmissionDate = time.Now()
	return nil
}

func (r *Report) ValidateReport() error {
	// Logic to validate report
	r.Status = ReportStatusValidated
	validationTime := time.Now()
	r.ValidationDate = &validationTime
	return nil
}

func (r *Report) ViewReport() (Report, error) {
	// Logic to view report will be implemented in controllers
	return *r, nil
}

// Request/Response structures
type ReportCreateRequest struct {
	Content   string `json:"content" binding:"required"`
	MissionID uint   `json:"mission_id" binding:"required"`
}

type ReportUpdateRequest struct {
	Content  string       `json:"content,omitempty"`
	Status   ReportStatus `json:"status,omitempty"`
	Comments string       `json:"comments,omitempty"`
}

type ReportFilterRequest struct {
	Status       ReportStatus `json:"status,omitempty"`
	EnseignantID uint         `json:"enseignant_id,omitempty"`
	MissionID    uint         `json:"mission_id,omitempty"`
	DateFrom     *time.Time   `json:"date_from,omitempty"`
	DateTo       *time.Time   `json:"date_to,omitempty"`
}
