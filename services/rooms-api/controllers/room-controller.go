package controllers

import (
	"net/http"
	"rooms-api/domain"
	"rooms-api/services"
	"rooms-api/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	roomService *services.RoomService
}

func NewRoomController(roomService *services.RoomService) *RoomController {
	return &RoomController{
		roomService: roomService,
	}
}

// CreateRoom godoc
// @Summary Create a new room
// @Description Create a new room with the provided details
// @Tags rooms
// @Accept json
// @Produce json
// @Param room body domain.CreateRoomRequest true "Room data"
// @Success 201 {object} domain.RoomResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /rooms [post]
func (c *RoomController) CreateRoom(ctx *gin.Context) {
	var req domain.CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "Datos de solicitud inválidos",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	room, err := c.roomService.CreateRoom(ctx.Request.Context(), req)
	if err != nil {
		statusCode := utils.GetHTTPStatus(err)
		ctx.JSON(statusCode, utils.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	ctx.JSON(http.StatusCreated, room)
}

// GetRoomByID godoc
// @Summary Get room by ID
// @Description Get a room by its ID
// @Tags rooms
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} domain.RoomResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /rooms/{id} [get]
func (c *RoomController) GetRoomByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "ID de habitación inválido",
			Message: "El ID de la habitación es obligatorio",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Convertir id (string) a uint
	idUint64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "ID de habitación inválido",
			Message: "El ID de la habitación debe ser un número",
			Code:    http.StatusBadRequest,
		})
		return
	}
	idUint := uint(idUint64)

	room, err := c.roomService.GetRoomByID(ctx.Request.Context(), idUint)
	if err != nil {
		statusCode := utils.GetHTTPStatus(err)
		ctx.JSON(statusCode, utils.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	ctx.JSON(http.StatusOK, room)
}

// GetRoomByNumber godoc
// @Summary Get room by number
// @Description Get a room by its number
// @Tags rooms
// @Produce json
// @Param number path string true "Room number"
// @Success 200 {object} domain.RoomResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /rooms/number/{number} [get]
func (c *RoomController) GetRoomByNumber(ctx *gin.Context) {
	number := ctx.Param("number")
	if number == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "Número de habitación inválido",
			Message: "El número de habitación es obligatorio",
			Code:    http.StatusBadRequest,
		})
		return
	}

	room, err := c.roomService.GetRoomByNumber(ctx.Request.Context(), number)
	if err != nil {
		statusCode := utils.GetHTTPStatus(err)
		ctx.JSON(statusCode, utils.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	ctx.JSON(http.StatusOK, room)
}

// GetRooms godoc
// @Summary Get all rooms
// @Description Get all rooms with optional filtering and pagination
// @Tags rooms
// @Produce json
// @Param type query string false "Room type (single, double, suite, deluxe, standard)"
// @Param status query string false "Room status (available, occupied, maintenance, reserved)"
// @Param floor query int false "Floor number"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param has_wifi query bool false "Has WiFi"
// @Param has_ac query bool false "Has AC"
// @Param has_tv query bool false "Has TV"
// @Param has_minibar query bool false "Has minibar"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} domain.RoomListResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /rooms [get]
func (c *RoomController) GetRooms(ctx *gin.Context) {
	// Parsear parámetros de consulta
	filter := domain.RoomFilter{}

	if roomType := ctx.Query("type"); roomType != "" {
		filter.Type = (*domain.RoomType)(&roomType)
	}
	if status := ctx.Query("status"); status != "" {
		filter.Status = (*domain.RoomStatus)(&status)
	}
	if floorStr := ctx.Query("floor"); floorStr != "" {
		if floor, err := strconv.Atoi(floorStr); err == nil {
			filter.Floor = &floor
		}
	}
	if minPriceStr := ctx.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = &minPrice
		}
	}
	if maxPriceStr := ctx.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = &maxPrice
		}
	}
	if hasWifiStr := ctx.Query("has_wifi"); hasWifiStr != "" {
		if hasWifi, err := strconv.ParseBool(hasWifiStr); err == nil {
			filter.HasWifi = &hasWifi
		}
	}
	if hasACStr := ctx.Query("has_ac"); hasACStr != "" {
		if hasAC, err := strconv.ParseBool(hasACStr); err == nil {
			filter.HasAC = &hasAC
		}
	}
	if hasTVStr := ctx.Query("has_tv"); hasTVStr != "" {
		if hasTV, err := strconv.ParseBool(hasTVStr); err == nil {
			filter.HasTV = &hasTV
		}
	}
	if hasMinibarStr := ctx.Query("has_minibar"); hasMinibarStr != "" {
		if hasMinibar, err := strconv.ParseBool(hasMinibarStr); err == nil {
			filter.HasMinibar = &hasMinibar
		}
	}

	// Parsear parámetros de paginación
	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	rooms, err := c.roomService.GetRooms(ctx.Request.Context(), filter, page, limit)
	if err != nil {
		statusCode := utils.GetHTTPStatus(err)
		ctx.JSON(statusCode, utils.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	ctx.JSON(http.StatusOK, rooms)
}

// UpdateRoom godoc
// @Summary Update a room
// @Description Update a room by ID with the provided data
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Param room body domain.UpdateRoomRequest true "Room update data"
// @Success 200 {object} domain.RoomResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /rooms/{id} [put]
func (c *RoomController) UpdateRoom(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "ID de habitación inválido",
			Message: "El ID de la habitación es obligatorio",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req domain.UpdateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "Datos de solicitud inválidos",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	idUint64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "ID de habitación inválido",
			Message: "El ID de la habitación debe ser un número",
			Code:    http.StatusBadRequest,
		})
		return
	}
	idUint := uint(idUint64)
	room, err := c.roomService.UpdateRoom(ctx.Request.Context(), idUint, req)
	if err != nil {
		statusCode := utils.GetHTTPStatus(err)
		ctx.JSON(statusCode, utils.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	ctx.JSON(http.StatusOK, room)
}

// DeleteRoom godoc
// @Summary Delete a room
// @Description Delete a room by ID
// @Tags rooms
// @Param id path string true "Room ID"
// @Success 204 "Room deleted successfully"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /rooms/{id} [delete]
func (c *RoomController) DeleteRoom(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "ID de habitación inválido",
			Message: "El ID de la habitación es obligatorio",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Convertir id (string) a uint
	idUint64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "ID de habitación inválido",
			Message: "El ID de la habitación debe ser un número",
			Code:    http.StatusBadRequest,
		})
		return
	}
	idUint := uint(idUint64)

	err = c.roomService.DeleteRoom(ctx.Request.Context(), idUint)
	if err != nil {
		statusCode := utils.GetHTTPStatus(err)
		ctx.JSON(statusCode, utils.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UpdateRoomStatus godoc
// @Summary Update room status
// @Description Update the status of a room
// @Tags rooms
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Param status body object{status=string} true "Room status"
// @Success 200 {object} domain.RoomResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /rooms/{id}/status [patch]
func (c *RoomController) UpdateRoomStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "Invalid room ID",
			Message: "Room ID is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var req struct {
		Status domain.RoomStatus `json:"status" binding:"required,oneof=available occupied maintenance reserved"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "Datos de solicitud inválidos",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	idUint64, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Error:   "ID de habitación inválido",
			Message: "El ID de la habitación debe ser un número",
			Code:    http.StatusBadRequest,
		})
		return
	}
	idUint := uint(idUint64)
	room, err := c.roomService.UpdateRoomStatus(ctx.Request.Context(), idUint, req.Status)
	if err != nil {
		statusCode := utils.GetHTTPStatus(err)
		ctx.JSON(statusCode, utils.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	ctx.JSON(http.StatusOK, room)
}

// GetAvailableRooms godoc
// @Summary Get available rooms
// @Description Get all available rooms with optional filtering
// @Tags rooms
// @Produce json
// @Param type query string false "Room type"
// @Param floor query int false "Floor number"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} domain.RoomListResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /rooms/available [get]
func (c *RoomController) GetAvailableRooms(ctx *gin.Context) {
	// Parsear parámetros (similar a GetRooms) pero sólo para disponibles
	filter := domain.RoomFilter{}

	if roomType := ctx.Query("type"); roomType != "" {
		filter.Type = (*domain.RoomType)(&roomType)
	}
	if floorStr := ctx.Query("floor"); floorStr != "" {
		if floor, err := strconv.Atoi(floorStr); err == nil {
			filter.Floor = &floor
		}
	}
	if minPriceStr := ctx.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = &minPrice
		}
	}
	if maxPriceStr := ctx.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = &maxPrice
		}
	}

	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	rooms, err := c.roomService.GetAvailableRooms(ctx.Request.Context(), filter, page, limit)
	if err != nil {
		statusCode := utils.GetHTTPStatus(err)
		ctx.JSON(statusCode, utils.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}

	ctx.JSON(http.StatusOK, rooms)
}

func (c *RoomController) GetRoomsViaSearch(ctx *gin.Context) {
	// Parsear parámetros de consulta
	filter := domain.RoomFilter{}
	if roomType := ctx.Query("type"); roomType != "" {
		filter.Type = (*domain.RoomType)(&roomType)
	}
	if status := ctx.Query("status"); status != "" {
		filter.Status = (*domain.RoomStatus)(&status)
	}
	if floorStr := ctx.Query("floor"); floorStr != "" {
		if floor, err := strconv.Atoi(floorStr); err == nil {
			filter.Floor = &floor
		}
	}
	if minPriceStr := ctx.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = &minPrice
		}
	}
	if maxPriceStr := ctx.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = &maxPrice
		}
	}
	if hasWifiStr := ctx.Query("has_wifi"); hasWifiStr != "" {
		if hasWifi, err := strconv.ParseBool(hasWifiStr); err == nil {
			filter.HasWifi = &hasWifi
		}
	}
	if hasACStr := ctx.Query("has_ac"); hasACStr != "" {
		if hasAC, err := strconv.ParseBool(hasACStr); err == nil {
			filter.HasAC = &hasAC
		}
	}
	if hasTVStr := ctx.Query("has_tv"); hasTVStr != "" {
		if hasTV, err := strconv.ParseBool(hasTVStr); err == nil {
			filter.HasTV = &hasTV
		}
	}
	if hasMinibarStr := ctx.Query("has_minibar"); hasMinibarStr != "" {
		if hasMinibar, err := strconv.ParseBool(hasMinibarStr); err == nil {
			filter.HasMinibar = &hasMinibar
		}
	}
	// Parsear parámetros de paginación
	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	limit := 10
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	rooms, err := c.roomService.GetRoomsViaSearch(ctx.Request.Context(), filter, page, limit)
	if err != nil {
		statusCode := utils.GetHTTPStatus(err)
		ctx.JSON(statusCode, utils.ErrorResponse{
			Error:   err.Error(),
			Message: err.Error(),
			Code:    statusCode,
		})
		return
	}
	ctx.JSON(http.StatusOK, rooms)
}
