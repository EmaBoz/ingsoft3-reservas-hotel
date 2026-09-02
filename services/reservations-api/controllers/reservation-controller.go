package controllers

import (
	"errors"
	"net/http"
	"reservations-api/domain"
	"reservations-api/services"
	"reservations-api/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ReservationController struct {
	service services.ReservationService
}

func NewReservationController(service services.ReservationService) *ReservationController {
	return &ReservationController{service: service}
}

func (c *ReservationController) GetmyReservations(ctx *gin.Context) {
	userIDStr := strings.TrimSpace(ctx.Param("user_id"))
	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id invalido"})
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id invalido"})
		return
	}
	reservations, err := c.service.GetmyReservations(ctx, uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

func (c *ReservationController) GetAllReservations(ctx *gin.Context) {
	reservations, err := c.service.GetAllReservations(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

func (c *ReservationController) CreateReservation(ctx *gin.Context) {
	var req domain.CreateReservationDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := c.service.CreateReservation(ctx, req)
	if err != nil {
		if errors.Is(err, utils.ErrReservationConflict) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"reservation": created})
}
func (c *ReservationController) DeleteReservation(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}

	var body domain.CancelReservationDTO
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.DeleteReservation(ctx, id, body.Reason)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "reserva no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *ReservationController) GetReservationByID(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}

	dto, err := c.service.GetReservationByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "reserva no encontrada"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reservation": dto})
}
