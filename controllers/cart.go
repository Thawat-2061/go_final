package controllers

import (
	"go-final/dto"
	"go-final/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartController struct {
	db *gorm.DB
}

func NewCartController(db *gorm.DB) *CartController {
	return &CartController{db: db}
}

func (ctrl *CartController) SearchProducts(c *gin.Context) {
	// รับค่าจาก query parameters
	keyword := c.Query("keyword")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")

	query := ctrl.db.Model(&models.Product{})

	// กรองด้วยชื่อหรือคำอธิบาย (ถ้ามี keyword)
	if keyword != "" {
		query = query.Where("product_name LIKE ? OR description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// กรองด้วยราคาขั้นต่ำ (ถ้ามี)
	if minPriceStr != "" {
		minPrice, err := strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ราคาขั้นต่ำไม่ถูกต้อง"})
			return
		}
		query = query.Where("price >= ?", minPrice)
	}

	// กรองด้วยราคาขั้นสูง (ถ้ามี)
	if maxPriceStr != "" {
		maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ราคาขั้นสูงไม่ถูกต้อง"})
			return
		}
		query = query.Where("price <= ?", maxPrice)
	}

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "เกิดข้อผิดพลาดในการค้นหาสินค้า",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (ctrl *CartController) AddToCart(c *gin.Context) {
	// ตรวจสอบและแปลงข้อมูลจาก Request
	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "ข้อมูลไม่ถูกต้อง",
			"details": err.Error(),
		})
		return
	}

	// ตรวจสอบว่าสินค้ามีอยู่ในระบบหรือไม่
	var product models.Product
	if err := ctrl.db.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบสินค้านี้ในระบบ"})
		return
	}

	// ตรวจสอบว่าสินค้ามีพอในสต็อกหรือไม่
	if product.StockQuantity < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "สินค้าในสต็อกไม่เพียงพอ",
			"available": product.StockQuantity,
		})
		return
	}

	// ค้นหาหรือสร้างตะกร้าใหม่
	var cart models.Cart
	err := ctrl.db.Where("customer_id = ? AND cart_name = ?", req.CustomerID, req.CartName).First(&cart).Error

	if err == gorm.ErrRecordNotFound {
		cart = models.Cart{
			CustomerID: req.CustomerID,
			CartName:   req.CartName,
		}
		if err := ctrl.db.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างตะกร้าใหม่ได้"})
			return
		}
	} else if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "เกิดข้อผิดพลาดในการค้นหาตะกร้า",
			"details": err.Error(), // เพิ่มรายละเอียดของ error
		})
		return
	}

	// ตรวจสอบว่าสินค้าอยู่ในตะกร้าแล้วหรือไม่
	var cartItem models.CartItem
	err = ctrl.db.Where("cart_id = ? AND product_id = ?", cart.CartID, req.ProductID).First(&cartItem).Error

	if err == nil {
		// อัพเดทจำนวนถ้ามีอยู่แล้ว
		newQuantity := cartItem.Quantity + req.Quantity
		if product.StockQuantity < newQuantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "สินค้าในสต็อกไม่พอ",
				"available": product.StockQuantity,
				"requested": newQuantity,
			})
			return
		}
		if err := ctrl.db.Model(&cartItem).Update("quantity", newQuantity).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัพเดทจำนวนสินค้าได้"})
			return
		}
	} else if err == gorm.ErrRecordNotFound {
		// เพิ่มสินค้าถ้ายังไม่มี
		newItem := models.CartItem{
			CartID:    cart.CartID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := ctrl.db.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถเพิ่มสินค้าในตะกร้าได้"})
			return
		}
	} else {
		if err != nil && err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "เกิดข้อผิดพลาดในการตรวจสอบสินค้าในตะกร้า",
				"details": err.Error(), // เพิ่มรายละเอียดของ error
			})
			return
		}
	}

	// อัพเดทจำนวนสินค้าในสต็อก
	if err := ctrl.db.Model(&product).Update("stock_quantity", gorm.Expr("stock_quantity - ?", req.Quantity)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัพเดทสต็อกสินค้าได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "เพิ่มสินค้าในตะกร้าสำเร็จ",
		"cart_id":    cart.CartID,
		"product_id": req.ProductID,
		"quantity":   req.Quantity,
	})
}
func (ctrl *CartController) GetAllCarts(c *gin.Context) {
	// รับ customer_id จาก query parameter
	customerID := c.Query("customer_id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "กรุณาระบุ customer_id"})
		return
	}

	// ดึงข้อมูลรถเข็นทั้งหมดของลูกค้า
	var carts []models.Cart
	if err := ctrl.db.Where("customer_id = ?", customerID).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลตะกร้าได้"})
		return
	}

	// ตรวจสอบว่ามีรถเข็นหรือไม่
	if len(carts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบทุกรถเข็นของลูกค้า"})
		return
	}

	// ดึงข้อมูลสินค้าในรถเข็นแต่ละอัน
	var response []gin.H

	for _, cart := range carts {
		var items []models.CartItem
		if err := ctrl.db.Where("cart_id = ?", cart.CartID).Preload("Product").Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "เกิดข้อผิดพลาดในการดึงสินค้าจากตะกร้า"})
			return
		}

		// แปลงข้อมูลเป็น JSON
		var cartItems []gin.H
		for _, item := range items {
			cartItems = append(cartItems, gin.H{
				"product_id":   item.ProductID,
				"product_name": item.Product.ProductName,
				"quantity":     item.Quantity,
				"price":        item.Product.Price,
				"total_price":  float64(item.Quantity) * item.Product.Price,
			})
		}

		response = append(response, gin.H{
			"cart_id":   cart.CartID,
			"cart_name": cart.CartName,
			"items":     cartItems,
		})
	}

	// ส่ง JSON response กลับไป
	c.JSON(http.StatusOK, gin.H{"carts": response})
}

func SetupCartRoutes(router *gin.Engine, db *gorm.DB) {
	cartCtrl := NewCartController(db)

	cartGroup := router.Group("/product")
	{
		cartGroup.GET("/search", cartCtrl.SearchProducts)
		cartGroup.POST("/cart/add", cartCtrl.AddToCart)
		cartGroup.GET("/getAll", cartCtrl.GetAllCarts)
	}
}
