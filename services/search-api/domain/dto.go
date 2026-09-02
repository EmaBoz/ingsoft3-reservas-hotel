package domain

// SearchRoomsRequest representa los parámetros de búsqueda
type SearchRoomsRequest struct {
	Q          string  `form:"q"`           // Texto libre
	Type       string  `form:"type"`        // Filtro por tipo
	Status     string  `form:"status"`      // Filtro por estado
	Floor      *int    `form:"floor"`       // Filtro por piso (pointer para distinguir 0 de no-enviado)
	MinPrice   *float64 `form:"min_price"`  // Precio mínimo
	MaxPrice   *float64 `form:"max_price"`  // Precio máximo
	HasWifi    *bool   `form:"has_wifi"`    // Filtro WiFi
	HasAC      *bool   `form:"has_ac"`      // Filtro AC
	HasTV      *bool   `form:"has_tv"`      // Filtro TV
	HasMinibar *bool   `form:"has_minibar"` // Filtro Minibar
	Sort       string  `form:"sort"`        // Campo de ordenamiento
	Page       int     `form:"page"`        // Página (default 1)
	Limit      int     `form:"limit"`       // Tamaño página (default 10, max 50)
}

// SearchRoomsResponse representa la respuesta de búsqueda paginada
type SearchRoomsResponse struct {
	Page    int                 `json:"page"`
	Limit   int                 `json:"limit"`
	Total   int64               `json:"total"`
	Results []SolrRoomDocument  `json:"results"`
}

// RoomSearchResult representa un resultado individual simplificado
type RoomSearchResult struct {
	ID         string  `json:"id"`
	Number     string  `json:"number"`
	Type       string  `json:"type"`
	Status     string  `json:"status"`
	Price      float64 `json:"price"`
	Capacity   int     `json:"capacity"`
	Floor      int     `json:"floor"`
	HasWifi    bool    `json:"has_wifi"`
	HasAC      bool    `json:"has_ac"`
	HasTV      bool    `json:"has_tv"`
	HasMinibar bool    `json:"has_minibar"`
}
