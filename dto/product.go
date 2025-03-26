package dto

type ProductSearchRequest struct {
	Keyword  string  `json:"keyword" form:"keyword"`
	MinPrice float64 `json:"min_price" form:"min_price"`
	MaxPrice float64 `json:"max_price" form:"max_price"`
}

type AddToCartRequest struct {
	CustomerID uint   `json:"customer_id" binding:"required"`
	CartName   string `json:"cart_name" binding:"required"`
	ProductID  uint   `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
}
