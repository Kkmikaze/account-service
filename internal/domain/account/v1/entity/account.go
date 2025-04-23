package entity

import (
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"time"
)

type Account struct {
	CreatedAt time.Time      `gorm:"type:timestamptz;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamptz;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamptz;default:null" json:"deleted_at"`

	ID            string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CitizenID     string `gorm:"type:varchar(16);not null" json:"citizen_id"`
	Name          string `gorm:"type:varchar(150);not null" json:"name"`
	Phone         string `gorm:"type:varchar(16);not null" json:"phone"`
	AccountNumber string `gorm:"type:varchar(14);unique;not null" json:"account_number"`

	// Many-to-many relationship with roles
	SavingHistory   []*SavingHistory   `gorm:"foreignKey:AccountID;references:ID" json:"saving_history"`
	WithdrawHistory []*WithdrawHistory `gorm:"foreignKey:AccountID;references:ID" json:"withdraw_history"`

	Balance float64 `gorm:"-" json:"balance"`
}

func (e *Account) TableName() string {
	return "accounts"
}
