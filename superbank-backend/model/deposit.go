package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Deposit struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CustomerID   uuid.UUID `json:"customer_id" gorm:"type:uuid;not null"`
	Amount       float64   `json:"amount" gorm:"not null"`
	InterestRate float64   `json:"interest_rate" gorm:"not null"`
	TermMonths   int       `json:"term_months" gorm:"not null"`
	StartDate    time.Time `json:"start_date" gorm:"not null"`
	MaturityDate time.Time `json:"maturity_date" gorm:"not null"`
	IsWithdrawn  bool      `json:"is_withdrawn" gorm:"default:false"`
	Note         string    `json:"note"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	Customer Customer `json:"-" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (d *Deposit) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.NewV4()
	return
}
