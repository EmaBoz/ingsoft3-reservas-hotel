package domain

// Gin parsea strings JSON a time.Time usando time_format
type CreateReservationDTO struct {
	UserID    uint   `json:"user_id"    binding:"required"`
	StartDate string `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date"   binding:"required,datetime=2006-01-02"`
	RoomID    uint   `json:"room_id"    binding:"required"`
}

type ReservationResponseDTO struct {
	ID        string `json:"id"`
	UserID    uint   `json:"user_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomID    uint   `json:"room_id"`
	Status    string `json:"status"`
}

type CancelReservationDTO struct {
	Reason string `json:"reason" binding:"required"`
}
