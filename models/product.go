package models

import "time"

type Product struct {
	ProductID     uint       `gorm:"primaryKey;autoIncrement" json:"product_id"`
	ProductName   string     `gorm:"not null" json:"product_name"`
	Description   string     `gorm:"type:text" json:"description"`
	Price         float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	StockQuantity int        `gorm:"not null" json:"stock_quantity"`
	CartItems     []CartItem `gorm:"foreignKey:ProductID" json:"-"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Product) TableName() string {
	return "product" // matches your CREATE TABLE statement
}
