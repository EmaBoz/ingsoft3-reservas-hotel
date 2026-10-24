package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"reservations-api/domain"
	"reservations-api/utils"

	"github.com/gin-gonic/gin"
)

// servicioStub contesta lo que el test le configura, sin tocar Mongo.
type servicioStub struct {
	errCrear      error
	errCancelar   error
	reserva       domain.ReservationResponseDTO
	errBuscar     error
	llamadasCrear int
}

func (s *servicioStub) GetmyReservations(ctx context.Context, userID uint) ([]domain.Reservation, error) {
	return []domain.Reservation{}, nil
}

func (s *servicioStub) GetAllReservations(ctx context.Context) ([]domain.Reservation, error) {
	return []domain.Reservation{}, nil
}

func (s *servicioStub) CreateReservation(ctx context.Context, dto domain.CreateReservationDTO) (domain.Reservation, error) {
	s.llamadasCrear++
	return domain.Reservation{}, s.errCrear
}

func (s *servicioStub) DeleteReservation(ctx context.Context, id string, reason string) error {
	return s.errCancelar
}

func (s *servicioStub) GetReservationByID(ctx context.Context, id string) (domain.ReservationResponseDTO, error) {
	return s.reserva, s.errBuscar
}

func armarRouter(svc *servicioStub) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	c := NewReservationController(svc)
	r.POST("/api/reservations", c.CreateReservation)
	r.GET("/api/reservations/:id", c.GetReservationByID)
	r.DELETE("/api/reservations/:id", c.DeleteReservation)
	r.GET("/api/reservations/users/:user_id/myreservations", c.GetmyReservations)
	return r
}

func hacerPedido(r *gin.Engine, metodo, ruta, cuerpo string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCrear_BodyIncompleto_Devuelve400SinLlamarAlServicio(t *testing.T) {
	svc := &servicioStub{}
	r := armarRouter(svc)

	w := hacerPedido(r, http.MethodPost, "/api/reservations", `{"user_id": 1, "room_id": 101}`)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, obtuve %d", w.Code)
	}
	if svc.llamadasCrear != 0 {
		t.Errorf("con el body incompleto no tendría que llegar al servicio")
	}
}

func TestCrear_HabitacionOcupada_Devuelve409(t *testing.T) {
	svc := &servicioStub{errCrear: utils.ErrReservationConflict}
	r := armarRouter(svc)
	cuerpo := `{"user_id": 1, "room_id": 101, "start_date": "2026-11-10", "end_date": "2026-11-12"}`

	w := hacerPedido(r, http.MethodPost, "/api/reservations", cuerpo)

	if w.Code != http.StatusConflict {
		t.Fatalf("esperaba 409, obtuve %d: %s", w.Code, w.Body.String())
	}
}

func TestCrear_Valida_Devuelve201(t *testing.T) {
	svc := &servicioStub{}
	r := armarRouter(svc)
	cuerpo := `{"user_id": 1, "room_id": 101, "start_date": "2026-11-10", "end_date": "2026-11-12"}`

	w := hacerPedido(r, http.MethodPost, "/api/reservations", cuerpo)

	if w.Code != http.StatusCreated {
		t.Fatalf("esperaba 201, obtuve %d: %s", w.Code, w.Body.String())
	}
}

func TestCancelar_ReservaInexistente_Devuelve404(t *testing.T) {
	svc := &servicioStub{errCancelar: utils.ErrReservationNotFound}
	r := armarRouter(svc)

	w := hacerPedido(r, http.MethodDelete, "/api/reservations/abc123", `{"reason": "cambio de planes"}`)

	if w.Code != http.StatusNotFound {
		t.Fatalf("esperaba 404, obtuve %d", w.Code)
	}
}

func TestBuscar_Existente_Devuelve200ConLaReserva(t *testing.T) {
	svc := &servicioStub{reserva: domain.ReservationResponseDTO{ID: "abc123", RoomID: 101, Status: "active"}}
	r := armarRouter(svc)

	w := hacerPedido(r, http.MethodGet, "/api/reservations/abc123", "")

	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"id":"abc123"`) {
		t.Fatalf("esperaba 200 con la reserva abc123, obtuve %d: %s", w.Code, w.Body.String())
	}
}

func TestMisReservas_UserIDNoNumerico_Devuelve400(t *testing.T) {
	r := armarRouter(&servicioStub{})

	w := hacerPedido(r, http.MethodGet, "/api/reservations/users/pepe/myreservations", "")

	if w.Code != http.StatusBadRequest {
		t.Fatalf("esperaba 400, obtuve %d", w.Code)
	}
}
