package models

type CartItem struct {
	CartItemID uint    `gorm:"primaryKey;autoIncrement" json:"cart_item_id"`
	CartID     uint    `gorm:"not null" json:"cart_id"`
	ProductID  uint    `gorm:"not null" json:"product_id"`
	Quantity   int     `gorm:"not null" json:"quantity"`
	Product    Product `gorm:"foreignKey:ProductID" json:"product"`
	CreatedAt  int64   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  int64   `gorm:"autoUpdateTime" json:"updated_at"`
}
