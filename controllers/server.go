package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ประกาศตัวแปร db เป็น package-level variable
var db *gorm.DB

// SetDB ใช้สำหรับตั้งค่า database connection จากภายนอก
func SetDB(database *gorm.DB) {
	db = database
}

func StartServer() {
	router := gin.Default()

	// Health check endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is now working",
		})
	})

	// Include controllers
	SetupAuthRoutes(router, db)

	router.Run(":8080") // ระบุ port ให้ชัดเจน
}
