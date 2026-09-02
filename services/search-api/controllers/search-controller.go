package controllers

import (
	"log"
	"net/http"

	"search-api/domain"
	"search-api/services"

	"github.com/gin-gonic/gin"
)

// SearchController maneja los endpoints HTTP de búsqueda
type SearchController struct {
	searchService *services.SearchService
}

// NewSearchController crea un nuevo controlador de búsqueda
func NewSearchController(searchService *services.SearchService) *SearchController {
	return &SearchController{
		searchService: searchService,
	}
}

// SearchRooms maneja GET /api/search/rooms
func (c *SearchController) SearchRooms(ctx *gin.Context) {
	var req domain.SearchRoomsRequest

	// Bind query parameters
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Log de búsqueda
	log.Printf("Search request: q=%s, type=%s, status=%s, page=%d, limit=%d",
		req.Q, req.Type, req.Status, req.Page, req.Limit)

	// Realizar búsqueda
	response, err := c.searchService.SearchRooms(&req)
	if err != nil {
		log.Printf("Search failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Search failed",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// HealthCheck maneja GET /health
func (c *SearchController) HealthCheck(ctx *gin.Context) {
	if err := c.searchService.HealthCheck(); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "search-api",
	})
}
