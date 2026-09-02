package domain

import "time"

type EventType string

const (
	EventReservationCreated  EventType = "reservation.created"
	EventReservationCanceled EventType = "reservation.canceled"
)

type ReservationEvent struct {
	EventType     EventType `json:"event_type"`
	ReservationID string    `json:"reservation_id"`
	UserID        uint      `json:"user_id"`
	RoomID        uint      `json:"room_id"`
	StartDate     string    `json:"start_date"`
	EndDate       string    `json:"end_date"`
	Status        string    `json:"status"`
	CancelReason  *string   `json:"cancel_reason,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}
