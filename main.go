package main

import (
	"fmt"
	"go-final/controllers"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	dsn := viper.GetString("mysql.dsn")
	fmt.Println("Database DSN:", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // หรือใช้ logger.Error เพื่อแสดงเฉพาะ error จริงๆ
	})
	controllers.SetDB(db)

	// แก้เฉพาะบรรทัดนี้ จากเดิมที่พยายามใช้ if ตรวจสอบ
	controllers.StartServer() // เรียกใช้เฉยๆ ไม่ต้องตรวจสอบค่า return
}
