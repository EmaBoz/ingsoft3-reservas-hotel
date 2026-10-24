package services

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"reservations-api/domain"
	"reservations-api/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// repoMock reemplaza al repositorio de Mongo: devuelve lo que el test le
// configura y anota qué le pidieron, para poder verificar la interacción.
type repoMock struct {
	hayOverlap bool
	errOverlap error
	errGetByID error
	reserva    domain.Reservation

	llamadasOverlap int
	llamadasCreate  int
	llamadasDelete  int
	roomConsultada  uint
	motivoDelete    string
}

func (m *repoMock) GetmyReservations(ctx context.Context, userID uint) ([]domain.Reservation, error) {
	return nil, nil
}

func (m *repoMock) GetAll(ctx context.Context) ([]domain.Reservation, error) {
	return nil, nil
}

func (m *repoMock) Create(ctx context.Context, r domain.Reservation) (domain.Reservation, error) {
	m.llamadasCreate++
	r.ID = primitive.NewObjectID()
	return r, nil
}

func (m *repoMock) Delete(ctx context.Context, id string, reason string) error {
	m.llamadasDelete++
	m.motivoDelete = reason
	return nil
}

func (m *repoMock) GetByID(ctx context.Context, id string) (domain.Reservation, error) {
	if m.errGetByID != nil {
		return domain.Reservation{}, m.errGetByID
	}
	return m.reserva, nil
}

func (m *repoMock) HasActiveOverlap(ctx context.Context, roomID uint, startDate, endDate string) (bool, error) {
	m.llamadasOverlap++
	m.roomConsultada = roomID
	return m.hayOverlap, m.errOverlap
}

// publisherStub no publica nada: el evento sale en una goroutine y no es lo
// que estos tests verifican.
type publisherStub struct{}

func (p *publisherStub) PublishReservationCreated(ctx context.Context, e domain.ReservationEvent) error {
	return nil
}

func (p *publisherStub) PublishReservationCanceled(ctx context.Context, e domain.ReservationEvent) error {
	return nil
}

func (p *publisherStub) Close() error { return nil }

// publisherMock anota los eventos que le publican en un canal. El servicio
// publica en una goroutine, así que el test espera el evento en vez de
// asumir que ya se publicó.
type publisherMock struct {
	eventos chan domain.ReservationEvent
}

func nuevoPublisherMock() *publisherMock {
	return &publisherMock{eventos: make(chan domain.ReservationEvent, 1)}
}

func (p *publisherMock) PublishReservationCreated(ctx context.Context, e domain.ReservationEvent) error {
	p.eventos <- e
	return nil
}

func (p *publisherMock) PublishReservationCanceled(ctx context.Context, e domain.ReservationEvent) error {
	p.eventos <- e
	return nil
}

func (p *publisherMock) Close() error { return nil }

func esperarEvento(t *testing.T, p *publisherMock) domain.ReservationEvent {
	t.Helper()
	select {
	case e := <-p.eventos:
		return e
	case <-time.After(2 * time.Second):
		t.Fatal("no se publicó ningún evento")
		return domain.ReservationEvent{}
	}
}

func reservaValida() domain.CreateReservationDTO {
	return domain.CreateReservationDTO{
		UserID:    1,
		RoomID:    101,
		StartDate: "2026-11-10",
		EndDate:   "2026-11-12",
	}
}

func TestCrearReserva_SinFechas_EsRechazada(t *testing.T) {
	casos := []struct {
		nombre  string
		entrada string
		salida  string
	}{
		{"sin fecha de entrada", "", "2026-11-12"},
		{"sin fecha de salida", "2026-11-10", ""},
		{"sin ninguna fecha", "", ""},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := &repoMock{}
			svc := NewReservationService(repo, &publisherStub{})
			dto := reservaValida()
			dto.StartDate = c.entrada
			dto.EndDate = c.salida

			_, err := svc.CreateReservation(context.Background(), dto)

			if err == nil || !strings.Contains(err.Error(), "requeridos") {
				t.Fatalf("esperaba error de fechas requeridas, obtuve: %v", err)
			}
			if repo.llamadasCreate != 0 {
				t.Errorf("no tendría que haber guardado nada, Create se llamó %d veces", repo.llamadasCreate)
			}
		})
	}
}

func TestCrearReserva_SalidaNoPosteriorALaEntrada_EsRechazada(t *testing.T) {
	casos := []struct {
		nombre  string
		entrada string
		salida  string
	}{
		{"salida el mismo día que la entrada", "2026-11-10", "2026-11-10"},
		{"salida antes que la entrada", "2026-11-10", "2026-11-09"},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			repo := &repoMock{}
			svc := NewReservationService(repo, &publisherStub{})
			dto := reservaValida()
			dto.StartDate = c.entrada
			dto.EndDate = c.salida

			_, err := svc.CreateReservation(context.Background(), dto)

			if err == nil || !strings.Contains(err.Error(), "posterior") {
				t.Fatalf("esperaba error de salida posterior a la entrada, obtuve: %v", err)
			}
			if repo.llamadasOverlap != 0 {
				t.Errorf("con fechas inválidas no tendría que consultar disponibilidad")
			}
		})
	}
}

func TestCrearReserva_HabitacionOcupada_DevuelveConflictoYNoGuarda(t *testing.T) {
	// Arrange: el mock dice que la habitación ya está reservada en esas fechas
	repo := &repoMock{hayOverlap: true}
	svc := NewReservationService(repo, &publisherStub{})

	// Act
	_, err := svc.CreateReservation(context.Background(), reservaValida())

	// Assert: se consultó la habitación correcta, devolvió conflicto y no guardó
	if !errors.Is(err, utils.ErrReservationConflict) {
		t.Fatalf("esperaba ErrReservationConflict, obtuve: %v", err)
	}
	if repo.llamadasOverlap != 1 || repo.roomConsultada != 101 {
		t.Errorf("esperaba una consulta de disponibilidad para la habitación 101, hubo %d (room %d)",
			repo.llamadasOverlap, repo.roomConsultada)
	}
	if repo.llamadasCreate != 0 {
		t.Errorf("no tendría que guardar una reserva superpuesta, Create se llamó %d veces", repo.llamadasCreate)
	}
}

func TestCrearReserva_ErrorAlConsultarDisponibilidad_NoGuarda(t *testing.T) {
	repo := &repoMock{errOverlap: errors.New("mongo no responde")}
	svc := NewReservationService(repo, &publisherStub{})

	_, err := svc.CreateReservation(context.Background(), reservaValida())

	if err == nil || !strings.Contains(err.Error(), "disponibilidad") {
		t.Fatalf("esperaba error al validar disponibilidad, obtuve: %v", err)
	}
	if repo.llamadasCreate != 0 {
		t.Errorf("si no se pudo validar disponibilidad no tendría que guardar")
	}
}

func TestCrearReserva_Valida_SeGuardaActiva(t *testing.T) {
	// Arrange
	repo := &repoMock{hayOverlap: false}
	pub := nuevoPublisherMock()
	svc := NewReservationService(repo, pub)

	// Act
	creada, err := svc.CreateReservation(context.Background(), reservaValida())

	// Assert
	if err != nil {
		t.Fatalf("no esperaba error, obtuve: %v", err)
	}
	if e := esperarEvento(t, pub); e.EventType != domain.EventReservationCreated || e.RoomID != 101 {
		t.Errorf("esperaba el evento reservation.created de la habitación 101, llegó %+v", e)
	}
	if repo.llamadasCreate != 1 {
		t.Errorf("esperaba que se guarde una sola vez, Create se llamó %d veces", repo.llamadasCreate)
	}
	if creada.Status != domain.ReservationStatusActive {
		t.Errorf("una reserva nueva tiene que quedar activa, quedó %q", creada.Status)
	}
	if creada.RoomID != 101 || creada.StartDate != "2026-11-10" || creada.EndDate != "2026-11-12" {
		t.Errorf("la reserva guardada no tiene los datos pedidos: %+v", creada)
	}
}

func TestCancelarReserva_SinMotivo_EsRechazada(t *testing.T) {
	repo := &repoMock{}
	svc := NewReservationService(repo, &publisherStub{})

	err := svc.DeleteReservation(context.Background(), "abc123", "")

	if err == nil || !strings.Contains(err.Error(), "reason") {
		t.Fatalf("esperaba error por falta de motivo, obtuve: %v", err)
	}
	if repo.llamadasDelete != 0 {
		t.Errorf("sin motivo no tendría que cancelar nada")
	}
}

func TestCancelarReserva_SinID_EsRechazada(t *testing.T) {
	repo := &repoMock{}
	svc := NewReservationService(repo, &publisherStub{})

	err := svc.DeleteReservation(context.Background(), "", "cambio de planes")

	if err == nil || !strings.Contains(err.Error(), "id") {
		t.Fatalf("esperaba error por id inválido, obtuve: %v", err)
	}
	if repo.llamadasDelete != 0 {
		t.Errorf("sin id no tendría que cancelar nada")
	}
}

func TestCancelarReserva_Inexistente_NoLlamaAlDelete(t *testing.T) {
	repo := &repoMock{errGetByID: utils.ErrReservationNotFound}
	svc := NewReservationService(repo, &publisherStub{})

	err := svc.DeleteReservation(context.Background(), "abc123", "cambio de planes")

	if !errors.Is(err, utils.ErrReservationNotFound) {
		t.Fatalf("esperaba ErrReservationNotFound, obtuve: %v", err)
	}
	if repo.llamadasDelete != 0 {
		t.Errorf("no tendría que intentar cancelar una reserva que no existe")
	}
}

func TestBuscarReserva_DevuelveElIDComoTexto(t *testing.T) {
	id := primitive.NewObjectID()
	repo := &repoMock{reserva: domain.Reservation{ID: id, RoomID: 101, Status: domain.ReservationStatusActive}}
	svc := NewReservationService(repo, &publisherStub{})

	dto, err := svc.GetReservationByID(context.Background(), id.Hex())

	if err != nil {
		t.Fatalf("no esperaba error, obtuve: %v", err)
	}
	if dto.ID != id.Hex() || dto.Status != "active" || dto.RoomID != 101 {
		t.Errorf("el DTO no refleja la reserva guardada: %+v", dto)
	}
}

func TestCancelarReserva_Valida_CancelaConElMotivo(t *testing.T) {
	repo := &repoMock{reserva: domain.Reservation{ID: primitive.NewObjectID(), RoomID: 101}}
	pub := nuevoPublisherMock()
	svc := NewReservationService(repo, pub)

	err := svc.DeleteReservation(context.Background(), "abc123", "cambio de planes")

	if err != nil {
		t.Fatalf("no esperaba error, obtuve: %v", err)
	}
	if e := esperarEvento(t, pub); e.EventType != domain.EventReservationCanceled || e.CancelReason == nil || *e.CancelReason != "cambio de planes" {
		t.Errorf("esperaba el evento reservation.canceled con el motivo, llegó %+v", e)
	}
	if repo.llamadasDelete != 1 || repo.motivoDelete != "cambio de planes" {
		t.Errorf("esperaba un Delete con el motivo, hubo %d con motivo %q", repo.llamadasDelete, repo.motivoDelete)
	}
}
