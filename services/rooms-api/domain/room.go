package domain

import (
	"time"

	"gorm.io/gorm"
)

type RoomType string

const (
	RoomTypeSingle   RoomType = "single"
	RoomTypeDouble   RoomType = "double"
	RoomTypeSuite    RoomType = "suite"
	RoomTypeDeluxe   RoomType = "deluxe"
	RoomTypeStandard RoomType = "standard"
)

type RoomStatus string

const (
	RoomStatusAvailable   RoomStatus = "available"
	RoomStatusOccupied    RoomStatus = "occupied"
	RoomStatusMaintenance RoomStatus = "maintenance"
	RoomStatusReserved    RoomStatus = "reserved"
)

type Room struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Number      string         `gorm:"uniqueIndex;not null;size:20" json:"number"`
	Type        RoomType       `gorm:"type:varchar(20);not null" json:"type"`
	Status      RoomStatus     `gorm:"type:varchar(20);not null;default:'available'" json:"status"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Description string         `gorm:"type:text" json:"description"`
	Capacity    int            `gorm:"not null" json:"capacity"`
	Floor       int            `gorm:"not null" json:"floor"`
	HasWifi     bool           `gorm:"default:false" json:"has_wifi"`
	HasAC       bool           `gorm:"default:false" json:"has_ac"`
	HasTV       bool           `gorm:"default:false" json:"has_tv"`
	HasMinibar  bool           `gorm:"default:false" json:"has_minibar"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
