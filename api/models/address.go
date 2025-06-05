package models

import (
	"time"

	"gorm.io/gorm"
)

// Address model - represents physical addresses
type Address struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Street     string         `json:"street" gorm:"not null"`
	City       string         `json:"city" gorm:"not null"`
	PostalCode string         `json:"postal_code" gorm:"not null"`
	Country    string         `json:"country" gorm:"not null"`
	Latitude   float64        `json:"latitude"`
	Longitude  float64        `json:"longitude"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign Key
	UserID uint `json:"user_id"`

	// Relationships
	User    User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Courses []Course `json:"courses,omitempty" gorm:"foreignKey:AddressID"`
}

// Address methods
func (a *Address) AddAddress() error {
	// Logic to add address will be implemented in controllers
	return nil
}

func (a *Address) UpdateAddress() error {
	// Logic to update address will be implemented in controllers
	return nil
}

func (a *Address) CalculateRoute(destinationAddress *Address) (map[string]interface{}, error) {
	// Logic to calculate route between addresses
	// This would typically integrate with a mapping service like Google Maps
	return map[string]interface{}{
		"distance": 0,
		"duration": 0,
		"route":    []interface{}{},
	}, nil
}

// Request/Response structures
type AddressCreateRequest struct {
	Street     string  `json:"street" binding:"required"`
	City       string  `json:"city" binding:"required"`
	PostalCode string  `json:"postal_code" binding:"required"`
	Country    string  `json:"country" binding:"required"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

type AddressUpdateRequest struct {
	Street     string  `json:"street,omitempty"`
	City       string  `json:"city,omitempty"`
	PostalCode string  `json:"postal_code,omitempty"`
	Country    string  `json:"country,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

type RouteRequest struct {
	OriginAddressID      uint `json:"origin_address_id" binding:"required"`
	DestinationAddressID uint `json:"destination_address_id" binding:"required"`
}

type RouteResponse struct {
	Distance string                   `json:"distance"`
	Duration string                   `json:"duration"`
	Route    []map[string]interface{} `json:"route"`
}
