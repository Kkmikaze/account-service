package entity

import (
	_ "gorm.io/gorm"
	"time"
)

type WithdrawHistory struct {
	CreatedAt time.Time `gorm:"type:timestamptz;default:now()" json:"created_at"`
	ID        string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	AccountID string    `gorm:"type:uuid;not null" json:"citizen_card_id"`
	Amount    float64   `gorm:"type:numeric(20,2);not null;default:0" json:"amount"`

	// Foreign Key Relationships
	Account *Account `gorm:"foreignKey:AccountID;references:ID" json:"account"`
}

func (e *WithdrawHistory) TableName() string {
	return "withdraw_histories"
}
