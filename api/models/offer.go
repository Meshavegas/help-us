package models

import (
	"time"

	"gorm.io/gorm"
)

// OfferStatus represents the status of an offer
type OfferStatus string

const (
	OfferStatusOpen   OfferStatus = "open"
	OfferStatusClosed OfferStatus = "closed"
	OfferStatusFilled OfferStatus = "filled"
	OfferStatusDraft  OfferStatus = "draft"
)

// Offer model - represents job offers for teachers
type Offer struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Title           string         `json:"title" gorm:"not null"`
	Description     string         `json:"description" gorm:"type:text"`
	HourlyRate      float64        `json:"hourly_rate"`
	PublicationDate time.Time      `json:"publication_date"`
	Status          OfferStatus    `json:"status" gorm:"default:'draft'"`
	Requirements    string         `json:"requirements" gorm:"type:text"`
	Subject         string         `json:"subject"`
	Level           string         `json:"level"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign Keys
	CreatedByID uint `json:"created_by_id"`

	// Relationships
	CreatedBy   Administrator `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`
	Enseignants []Enseignant  `json:"enseignants,omitempty" gorm:"many2many:enseignant_offers;"`
	Options     []Option      `json:"options,omitempty" gorm:"foreignKey:OfferID"`
}

// Offer methods
func (o *Offer) CreateOffer() error {
	// Logic to create offer
	o.Status = OfferStatusOpen
	o.PublicationDate = time.Now()
	return nil
}

func (o *Offer) ApplyForOffer() error {
	// Logic for teachers to apply for offer
	return nil
}

func (o *Offer) CloseOffer() error {
	// Logic to close offer
	o.Status = OfferStatusClosed
	return nil
}

// Request/Response structures
type OfferCreateRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	HourlyRate   float64 `json:"hourly_rate" binding:"required,min=0"`
	Requirements string  `json:"requirements"`
	Subject      string  `json:"subject" binding:"required"`
	Level        string  `json:"level" binding:"required"`
}

type OfferUpdateRequest struct {
	Title        string      `json:"title,omitempty"`
	Description  string      `json:"description,omitempty"`
	HourlyRate   float64     `json:"hourly_rate,omitempty"`
	Status       OfferStatus `json:"status,omitempty"`
	Requirements string      `json:"requirements,omitempty"`
	Subject      string      `json:"subject,omitempty"`
	Level        string      `json:"level,omitempty"`
}

type OfferFilterRequest struct {
	Status   OfferStatus `json:"status,omitempty"`
	Subject  string      `json:"subject,omitempty"`
	Level    string      `json:"level,omitempty"`
	MinRate  *float64    `json:"min_rate,omitempty"`
	MaxRate  *float64    `json:"max_rate,omitempty"`
	DateFrom *time.Time  `json:"date_from,omitempty"`
	DateTo   *time.Time  `json:"date_to,omitempty"`
}

type OfferApplicationRequest struct {
	OfferID      uint   `json:"offer_id" binding:"required"`
	CoverLetter  string `json:"cover_letter"`
	Availability string `json:"availability"`
}
