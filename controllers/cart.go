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

	// ตรวจสอบว่าสินค้ามีในระบบหรือไม่
	var product models.Product
	if err := ctrl.db.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ไม่พบสินค้านี้ในระบบ",
		})
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
	err := ctrl.db.Where("customer_id = ? AND cart_name = ?", req.CustomerID, req.CartName).
		First(&cart).Error

	if err == gorm.ErrRecordNotFound {
		// สร้างตะกร้าใหม่ถ้ายังไม่มี
		cart = models.Cart{
			CustomerID: req.CustomerID,
			CartName:   req.CartName,
		}
		if err := ctrl.db.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ไม่สามารถสร้างตะกร้าใหม่ได้",
			})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "เกิดข้อผิดพลาดในการค้นหาตะกร้า",
		})
		return
	}

	// ตรวจสอบว่าสินค้ามีในตะกร้าแล้วหรือไม่
	var existingItem models.CartItem
	err = ctrl.db.Where("cart_id = ? AND product_id = ?", cart.CartID, req.ProductID).
		First(&existingItem).Error

	if err == nil {
		// อัพเดทจำนวนถ้ามีอยู่แล้ว
		existingItem.Quantity += req.Quantity

		// ตรวจสอบสต็อกอีกครั้งหลังรวมจำนวน
		if product.StockQuantity < existingItem.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "สินค้าในสต็อกไม่เพียงพอสำหรับจำนวนที่ต้องการ",
				"available": product.StockQuantity,
				"requested": existingItem.Quantity,
			})
			return
		}

		if err := ctrl.db.Save(&existingItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ไม่สามารถอัพเดทตะกร้าได้",
			})
			return
		}
	} else if err == gorm.ErrRecordNotFound {
		// เพิ่มใหม่ถ้ายังไม่มีในตะกร้า
		newItem := models.CartItem{
			CartID:    cart.CartID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := ctrl.db.Create(&newItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "ไม่สามารถเพิ่มสินค้าในตะกร้าได้",
			})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "เกิดข้อผิดพลาดในการตรวจสอบตะกร้า",
		})
		return
	}

	// อัพเดทจำนวนสินค้าในสต็อก
	if err := ctrl.db.Model(&product).
		Update("stock_quantity", gorm.Expr("stock_quantity - ?", req.Quantity)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "ไม่สามารถอัพเดทสต็อกสินค้าได้",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "เพิ่มสินค้าในตะกร้าสำเร็จ",
		"cart_id":    cart.CartID,
		"product_id": req.ProductID,
		"quantity":   req.Quantity,
	})
}
func SetupCartRoutes(router *gin.Engine, db *gorm.DB) {
	cartCtrl := NewCartController(db)

	cartGroup := router.Group("/product")
	{
		cartGroup.GET("/search", cartCtrl.SearchProducts)
		cartGroup.POST("/cart/add", cartCtrl.AddToCart)
	}
}
