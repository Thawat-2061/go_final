package models

import "time"

type Cart struct {
	CartID     uint       `gorm:"primaryKey;autoIncrement" json:"cart_id"`
	CustomerID uint       `gorm:"not null" json:"customer_id"`
	CartName   string     `gorm:"size:255" json:"cart_name"`
	Items      []CartItem `gorm:"foreignKey:CartID" json:"items"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
