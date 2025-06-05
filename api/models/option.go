package models

import (
	"time"

	"gorm.io/gorm"
)

// OptionStatus represents the status of an option
type OptionStatus string

const (
	OptionStatusActive    OptionStatus = "active"
	OptionStatusExpired   OptionStatus = "expired"
	OptionStatusAccepted  OptionStatus = "accepted"
	OptionStatusCancelled OptionStatus = "cancelled"
)

// Option model - represents reservation options
type Option struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	CreationDate   time.Time      `json:"creation_date"`
	ExpirationDate time.Time      `json:"expiration_date"`
	Status         OptionStatus   `json:"status" gorm:"default:'active'"`
	Description    string         `json:"description"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign Keys
	EnseignantID uint `json:"enseignant_id"`
	FamilleID    uint `json:"famille_id"`
	OfferID      uint `json:"offer_id,omitempty"`

	// Relationships
	Enseignant Enseignant `json:"enseignant,omitempty" gorm:"foreignKey:EnseignantID"`
	Famille    Famille    `json:"famille,omitempty" gorm:"foreignKey:FamilleID"`
	Offer      Offer      `json:"offer,omitempty" gorm:"foreignKey:OfferID"`
}

// Option methods
func (o *Option) CreateOption() error {
	// Logic to create option
	o.Status = OptionStatusActive
	o.CreationDate = time.Now()
	// Set expiration date to 7 days from creation by default
	o.ExpirationDate = time.Now().AddDate(0, 0, 7)
	return nil
}

func (o *Option) CancelOption() error {
	// Logic to cancel option
	o.Status = OptionStatusCancelled
	return nil
}

func (o *Option) AcceptOption() error {
	// Logic to accept option
	o.Status = OptionStatusAccepted
	return nil
}

func (o *Option) CheckExpiration() error {
	// Logic to check if option has expired
	if time.Now().After(o.ExpirationDate) && o.Status == OptionStatusActive {
		o.Status = OptionStatusExpired
	}
	return nil
}

// Request/Response structures
type OptionCreateRequest struct {
	EnseignantID   uint      `json:"enseignant_id" binding:"required"`
	FamilleID      uint      `json:"famille_id" binding:"required"`
	OfferID        uint      `json:"offer_id,omitempty"`
	ExpirationDate time.Time `json:"expiration_date"`
	Description    string    `json:"description"`
}

type OptionUpdateRequest struct {
	Status         OptionStatus `json:"status,omitempty"`
	ExpirationDate *time.Time   `json:"expiration_date,omitempty"`
	Description    string       `json:"description,omitempty"`
}

type OptionFilterRequest struct {
	Status       OptionStatus `json:"status,omitempty"`
	EnseignantID uint         `json:"enseignant_id,omitempty"`
	FamilleID    uint         `json:"famille_id,omitempty"`
	OfferID      uint         `json:"offer_id,omitempty"`
	DateFrom     *time.Time   `json:"date_from,omitempty"`
	DateTo       *time.Time   `json:"date_to,omitempty"`
}
