package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Pocket struct {
	ID           uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CustomerID   uuid.UUID  `json:"customer_id" gorm:"type:uuid;not null"`
	Name         string     `json:"name" gorm:"not null"`
	Balance      float64    `json:"balance" gorm:"not null;default:0"`
	TargetAmount *float64   `json:"target_amount" gorm:"default:null"`
	TargetDate   *time.Time `json:"target date" gorm:"default:null"`
	IsActive     bool       `json:"is_active" gorm:"not null;default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Customer Customer `json:"-" gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (d *Pocket) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.NewV4()
	return
}
