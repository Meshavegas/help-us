package models

import (
	"time"

	"gorm.io/gorm"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// PaymentType represents the type of payment
type PaymentType string

const (
	PaymentTypeCourse  PaymentType = "course"
	PaymentTypeMission PaymentType = "mission"
	PaymentTypeAdvance PaymentType = "advance"
	PaymentTypeRefund  PaymentType = "refund"
)

// Payment model - represents payments in the system
type Payment struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Amount      float64        `json:"amount" gorm:"not null"`
	PaymentDate time.Time      `json:"payment_date"`
	Status      PaymentStatus  `json:"status" gorm:"default:'pending'"`
	Type        PaymentType    `json:"type" gorm:"not null"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Foreign Keys
	UserID   uint `json:"user_id"`
	CourseID uint `json:"course_id,omitempty"`

	// Relationships
	User   User   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Course Course `json:"course,omitempty" gorm:"foreignKey:CourseID"`
}

// Payment methods
func (p *Payment) ProcessPayment() error {
	// Logic to process payment will be implemented in controllers
	p.Status = PaymentStatusCompleted
	p.PaymentDate = time.Now()
	return nil
}

func (p *Payment) GenerateInvoice() (map[string]interface{}, error) {
	// Logic to generate invoice
	return map[string]interface{}{
		"invoice_number": "INV-" + time.Now().Format("20060102") + "-" + string(rune(p.ID)),
		"amount":         p.Amount,
		"date":           p.PaymentDate,
		"description":    p.Description,
	}, nil
}

func (p *Payment) ViewHistory() ([]Payment, error) {
	// Logic to view payment history will be implemented in controllers
	return nil, nil
}

// Request/Response structures
type PaymentCreateRequest struct {
	Amount      float64     `json:"amount" binding:"required,min=0"`
	Type        PaymentType `json:"type" binding:"required"`
	Description string      `json:"description"`
	CourseID    uint        `json:"course_id,omitempty"`
}

type PaymentUpdateRequest struct {
	Status      PaymentStatus `json:"status,omitempty"`
	Description string        `json:"description,omitempty"`
}

type PaymentFilterRequest struct {
	Status    PaymentStatus `json:"status,omitempty"`
	Type      PaymentType   `json:"type,omitempty"`
	DateFrom  *time.Time    `json:"date_from,omitempty"`
	DateTo    *time.Time    `json:"date_to,omitempty"`
	MinAmount *float64      `json:"min_amount,omitempty"`
	MaxAmount *float64      `json:"max_amount,omitempty"`
}

type PaymentStatsResponse struct {
	TotalAmount     float64 `json:"total_amount"`
	CompletedAmount float64 `json:"completed_amount"`
	PendingAmount   float64 `json:"pending_amount"`
	TotalCount      int64   `json:"total_count"`
	CompletedCount  int64   `json:"completed_count"`
	PendingCount    int64   `json:"pending_count"`
}
