package utils

import (
	"errors"
	"net/http"
)

var (
	ErrRoomNotFound      = errors.New("habitación no encontrada")
	ErrRoomAlreadyExists = errors.New("ya existe una habitación con ese número")
	ErrInvalidRoomData   = errors.New("datos de habitación inválidos")
	ErrDatabaseError     = errors.New("error de base de datos")
	ErrInvalidID         = errors.New("ID de habitación inválido")
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func GetHTTPStatus(err error) int {
	switch err {
	case ErrRoomNotFound:
		return http.StatusNotFound
	case ErrRoomAlreadyExists:
		return http.StatusConflict
	case ErrInvalidRoomData, ErrInvalidID:
		return http.StatusBadRequest
	case ErrDatabaseError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
