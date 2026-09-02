package services

import (
	"context"
	"fmt"
	"log"
	"reservations-api/domain"
	"reservations-api/events"
	"reservations-api/repositories"
	"reservations-api/utils"
	"time"
)

type ReservationService interface {
	GetmyReservations(ctx context.Context, userID uint) ([]domain.Reservation, error)
	GetAllReservations(ctx context.Context) ([]domain.Reservation, error)
	CreateReservation(ctx context.Context, dto domain.CreateReservationDTO) (domain.Reservation, error)
	DeleteReservation(ctx context.Context, id string, reason string) error
	GetReservationByID(ctx context.Context, id string) (domain.ReservationResponseDTO, error)
}
type reservationService struct {
	repository repositories.ReservationRepository
	publisher  events.EventPublisher
}

func NewReservationService(repository repositories.ReservationRepository, publisher events.EventPublisher) ReservationService {
	return &reservationService{
		repository: repository,
		publisher:  publisher,
	}
}

func (s *reservationService) GetmyReservations(ctx context.Context, userID uint) ([]domain.Reservation, error) {
	return s.repository.GetmyReservations(ctx, userID)
}
func (s *reservationService) GetAllReservations(ctx context.Context) ([]domain.Reservation, error) {
	return s.repository.GetAll(ctx)
}
func (s *reservationService) CreateReservation(
	ctx context.Context,
	dto domain.CreateReservationDTO,
) (domain.Reservation, error) {

	// Validaciones simples con strings
	if dto.StartDate == "" || dto.EndDate == "" {
		return domain.Reservation{}, fmt.Errorf("start_date y end_date son requeridos")
	}

	if dto.EndDate <= dto.StartDate {
		return domain.Reservation{}, fmt.Errorf("end_date debe ser posterior a start_date")
	}

	hasOverlap, err := s.repository.HasActiveOverlap(ctx, dto.RoomID, dto.StartDate, dto.EndDate)
	if err != nil {
		return domain.Reservation{}, fmt.Errorf("failed to validar disponibilidad: %w", err)
	}
	if hasOverlap {
		return domain.Reservation{}, utils.ErrReservationConflict
	}

	// Mapear DTO → entidad (sin conversión)
	entity := domain.Reservation{
		UserID:    dto.UserID,
		RoomID:    dto.RoomID,
		StartDate: dto.StartDate, // sigue siendo string
		EndDate:   dto.EndDate,
		Status:    domain.ReservationStatusActive,
	}

	// Guardar en base de datos
	saved, err := s.repository.Create(ctx, entity)
	if err != nil {
		return domain.Reservation{}, fmt.Errorf("failed to create reservation: %w", err)
	}

	// Publicar evento en goroutine para no bloquear la respuesta
	go func() {
		event := domain.ReservationEvent{
			EventType:     domain.EventReservationCreated,
			ReservationID: saved.ID.Hex(),
			UserID:        saved.UserID,
			RoomID:        saved.RoomID,
			StartDate:     saved.StartDate,
			EndDate:       saved.EndDate,
			Status:        string(saved.Status),
			Timestamp:     time.Now(),
		}

		// Crear contexto con timeout para la publicación
		pubCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.publisher.PublishReservationCreated(pubCtx, event); err != nil {
			log.Printf("Error publicando evento de reserva creada: %v", err)
		}
	}()

	// Mapear entidad → DTO de respuesta
	resp := domain.Reservation{
		ID:        saved.ID,
		UserID:    saved.UserID,
		RoomID:    saved.RoomID,
		StartDate: saved.StartDate,
		EndDate:   saved.EndDate,
		Status:    saved.Status,
	}

	return resp, nil
}

func (s *reservationService) DeleteReservation(ctx context.Context, id string, reason string) error {
	if id == "" {
		return fmt.Errorf("id invalido")
	}
	if reason == "" {
		return fmt.Errorf("reason es requerido")
	}

	// Obtener la reserva antes de cancelarla para el evento
	reservation, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Cancelar la reserva
	err = s.repository.Delete(ctx, id, reason)
	if err != nil {
		return err
	}

	// Publicar evento en goroutine para no bloquear la respuesta
	go func() {
		event := domain.ReservationEvent{
			EventType:     domain.EventReservationCanceled,
			ReservationID: reservation.ID.Hex(),
			UserID:        reservation.UserID,
			RoomID:        reservation.RoomID,
			StartDate:     reservation.StartDate,
			EndDate:       reservation.EndDate,
			Status:        string(domain.ReservationStatusCanceled),
			CancelReason:  &reason,
			Timestamp:     time.Now(),
		}

		// Crear contexto con timeout para la publicación
		pubCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.publisher.PublishReservationCanceled(pubCtx, event); err != nil {
			log.Printf("Error publicando evento de reserva cancelada: %v", err)
		}
	}()

	return nil
}

func (s *reservationService) GetReservationByID(ctx context.Context, id string) (domain.ReservationResponseDTO, error) {
	res, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.ReservationResponseDTO{}, err
	}
	dto := domain.ReservationResponseDTO{
		ID:        res.ID.Hex(),
		UserID:    res.UserID,
		RoomID:    res.RoomID,
		StartDate: res.StartDate,
		EndDate:   res.EndDate,
		Status:    string(res.Status),
	}
	return dto, nil
}
