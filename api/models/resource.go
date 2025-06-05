package models

import (
	"time"

	"gorm.io/gorm"
)

// ResourceType represents the type of resource
type ResourceType string

const (
	ResourceTypeDocument ResourceType = "document"
	ResourceTypeVideo    ResourceType = "video"
	ResourceTypeAudio    ResourceType = "audio"
	ResourceTypeImage    ResourceType = "image"
	ResourceTypeLink     ResourceType = "link"
)

// Resource model - represents educational resources
type Resource struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Type        ResourceType   `json:"type" gorm:"not null"`
	URL         string         `json:"url" gorm:"not null"`
	Description string         `json:"description"`
	FileSize    int64          `json:"file_size"` // Size in bytes
	MimeType    string         `json:"mime_type"`
	UploadDate  time.Time      `json:"upload_date"`
	IsPublic    bool           `json:"is_public" gorm:"default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign Keys
	ManagedByID uint `json:"managed_by_id"`

	// Relationships
	ManagedBy Administrator `json:"managed_by,omitempty" gorm:"foreignKey:ManagedByID"`
	Users     []User        `json:"users,omitempty" gorm:"many2many:user_resources;"`
}

// Resource methods
func (r *Resource) UploadResource() error {
	// Logic to upload resource will be implemented in controllers
	r.UploadDate = time.Now()
	return nil
}

func (r *Resource) DownloadResource() ([]byte, error) {
	// Logic to download resource will be implemented in controllers
	return nil, nil
}

func (r *Resource) AccessResource() error {
	// Logic to access resource will be implemented in controllers
	return nil
}

// Request/Response structures
type ResourceCreateRequest struct {
	Title       string       `json:"title" binding:"required"`
	Type        ResourceType `json:"type" binding:"required"`
	URL         string       `json:"url" binding:"required"`
	Description string       `json:"description"`
	FileSize    int64        `json:"file_size,omitempty"`
	MimeType    string       `json:"mime_type,omitempty"`
	IsPublic    bool         `json:"is_public"`
}

type ResourceUpdateRequest struct {
	Title       string       `json:"title,omitempty"`
	Type        ResourceType `json:"type,omitempty"`
	URL         string       `json:"url,omitempty"`
	Description string       `json:"description,omitempty"`
	IsPublic    *bool        `json:"is_public,omitempty"`
}

type ResourceFilterRequest struct {
	Type        ResourceType `json:"type,omitempty"`
	IsPublic    *bool        `json:"is_public,omitempty"`
	ManagedByID uint         `json:"managed_by_id,omitempty"`
	DateFrom    *time.Time   `json:"date_from,omitempty"`
	DateTo      *time.Time   `json:"date_to,omitempty"`
}

type ResourceAccessRequest struct {
	ResourceID uint `json:"resource_id" binding:"required"`
	UserID     uint `json:"user_id" binding:"required"`
}
