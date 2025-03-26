package models

import (
	"time"
)

type Customer struct {
	CustomerID  uint      `gorm:"primaryKey;column:customer_id" json:"customer_id"`
	FirstName   string    `gorm:"column:first_name;not null" json:"first_name"`
	LastName    string    `gorm:"column:last_name;not null" json:"last_name"`
	Email       string    `gorm:"column:email;unique;not null" json:"email"`
	PhoneNumber string    `gorm:"column:phone_number" json:"phone_number,omitempty"`
	Address     string    `gorm:"column:address" json:"address,omitempty"`
	Password    string    `gorm:"column:password;not null" json:"-"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// ตั้งชื่อตารางให้ตรงกับที่มีอยู่
func (Customer) TableName() string {
	return "customer"
}
