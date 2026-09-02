package domain

import "time"

type CreateRoomRequest struct {
	Number      string   `json:"number" binding:"required"`
	Type        RoomType `json:"type" binding:"required,oneof=single double suite deluxe standard"`
	Price       float64  `json:"price" binding:"required,min=0"`
	Description string   `json:"description"`
	Capacity    int      `json:"capacity" binding:"required,min=1"`
	Floor       int      `json:"floor" binding:"required,min=1"`
	HasWifi     bool     `json:"has_wifi"`
	HasAC       bool     `json:"has_ac"`
	HasTV       bool     `json:"has_tv"`
	HasMinibar  bool     `json:"has_minibar"`
}

type UpdateRoomRequest struct {
	Number      *string     `json:"number,omitempty"`
	Type        *RoomType   `json:"type,omitempty"`
	Status      *RoomStatus `json:"status,omitempty"`
	Price       *float64    `json:"price,omitempty"`
	Description *string     `json:"description,omitempty"`
	Capacity    *int        `json:"capacity,omitempty"`
	Floor       *int        `json:"floor,omitempty"`
	HasWifi     *bool       `json:"has_wifi,omitempty"`
	HasAC       *bool       `json:"has_ac,omitempty"`
	HasTV       *bool       `json:"has_tv,omitempty"`
	HasMinibar  *bool       `json:"has_minibar,omitempty"`
}

type RoomResponse struct {
	ID          uint       `json:"id"`
	Number      string     `json:"number"`
	Type        RoomType   `json:"type"`
	Status      RoomStatus `json:"status"`
	Price       float64    `json:"price"`
	Description string     `json:"description"`
	Capacity    int        `json:"capacity"`
	Floor       int        `json:"floor"`
	HasWifi     bool       `json:"has_wifi"`
	HasAC       bool       `json:"has_ac"`
	HasTV       bool       `json:"has_tv"`
	HasMinibar  bool       `json:"has_minibar"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type RoomListResponse struct {
	Rooms []RoomResponse `json:"rooms"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}

type RoomFilter struct {
	Type       *RoomType   `json:"type,omitempty"`
	Status     *RoomStatus `json:"status,omitempty"`
	Floor      *int        `json:"floor,omitempty"`
	MinPrice   *float64    `json:"min_price,omitempty"`
	MaxPrice   *float64    `json:"max_price,omitempty"`
	HasWifi    *bool       `json:"has_wifi,omitempty"`
	HasAC      *bool       `json:"has_ac,omitempty"`
	HasTV      *bool       `json:"has_tv,omitempty"`
	HasMinibar *bool       `json:"has_minibar,omitempty"`
}
