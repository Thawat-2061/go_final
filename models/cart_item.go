package models

import "time"

type CartItem struct {
	CartItemID uint      `gorm:"primaryKey;autoIncrement" json:"cart_item_id"`
	CartID     uint      `gorm:"not null" json:"cart_id"`
	ProductID  uint      `gorm:"not null" json:"product_id"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	Product    Product   `gorm:"foreignKey:ProductID" json:"product"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CartItem) TableName() string {
	return "cart_item" // ชื่อตารางในฐานข้อมูล
}
