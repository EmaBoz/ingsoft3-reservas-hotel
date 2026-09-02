package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationStatus string

const (
	ReservationStatusActive   ReservationStatus = "active"
	ReservationStatusCanceled ReservationStatus = "canceled"
)

type Reservation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    uint               `bson:"user_id" json:"user_id"`
	StartDate string             `bson:"start_date" json:"start_date"`
	EndDate   string             `bson:"end_date" json:"end_date"`
	RoomID    uint               `bson:"room_id" json:"room_id"`
	Status    ReservationStatus  `bson:"status" json:"status"`

	//Soft Delete de cancelación
	CancelReason *string    `bson:"cancel_reason,omitempty" json:"cancel_reason,omitempty"`
	DeletedAt    *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
