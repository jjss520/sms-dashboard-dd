package model

import (
	"time"

	"gorm.io/gorm"
)

type SMS struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	SendTime  string         `json:"sendTime"`     // 短信时间
	Content   string         `json:"content"`      // 短信内容
	Sender    string         `json:"sender"`       // 发件人
	Phone     string         `json:"phone"`        // 手机号/卡槽信息(如: SIM2_13345694791)
	Device    string         `json:"device"`       // 设备型号
}

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex" json:"username"`
	Password  string         `json:"-"` // Store hashed password
}
