package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primarykey"`
	Name          string    `json:"name" gorm:"not null"`
	Phone         string    `json:"phone" gorm:"not null"`
	Address       string    `json:"address" gorm:"not null"`
	ParentName    string    `json:"parent_name" gprm:"not null"`
	AccountNumber string    `json:"account_number" gorm:"uniqueIndex;not null"`
	AccountBranch string    `json:"account_branch" gorm:"not null"`
	AccountType   string    `json:"account_type" gorm:"not null"`
	Balance       float64   `json:"balance" gorm:"default: 500000;not null"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	Deposits []Deposit `json:"-" gorm:"foreignKey:CustomerID"`
	Pockets  []Pocket  `json:"-" gorm:"foreignKey:CustomerID"`
}

func (customer *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	customer.ID = uuid.NewV4()
	return
}
