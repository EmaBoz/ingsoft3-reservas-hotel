package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// SolrRoomDocument representa el documento indexado en Solr
// Nota: Solr devuelve arrays para campos multivaluados, usamos el primer elemento
type SolrRoomDocument struct {
	ID          []string  `json:"id"`
	Number      []int     `json:"number"`                // Array en Solr
	Type        []string  `json:"type"`                  // Array en Solr
	Status      []string  `json:"status"`                // Array en Solr
	Price       []float64 `json:"price"`                 // Array en Solr
	Description []string  `json:"description,omitempty"` // Array en Solr
	Capacity    []int     `json:"capacity"`              // Array en Solr
	Floor       []int     `json:"floor"`                 // Array en Solr
	HasWifi     []bool    `json:"has_wifi"`              // Array en Solr
	HasAC       []bool    `json:"has_ac"`                // Array en Solr
	HasTV       []bool    `json:"has_tv"`                // Array en Solr
	HasMinibar  []bool    `json:"has_minibar"`           // Array en Solr
	CreatedAt   []string  `json:"created_at"`            // Array de strings RFC3339
	UpdatedAt   []string  `json:"updated_at"`            // Array de strings RFC3339
}

// UnmarshalJSON implementa custom unmarshaling para manejar tanto strings como arrays de Solr
func (d *SolrRoomDocument) UnmarshalJSON(data []byte) error {
	type Alias SolrRoomDocument
	aux := &struct {
		ID          interface{} `json:"id"`
		Number      interface{} `json:"number"`
		Type        interface{} `json:"type"`
		Status      interface{} `json:"status"`
		Price       interface{} `json:"price"`
		Description interface{} `json:"description"`
		Capacity    interface{} `json:"capacity"`
		Floor       interface{} `json:"floor"`
		HasWifi     interface{} `json:"has_wifi"`
		HasAC       interface{} `json:"has_ac"`
		HasTV       interface{} `json:"has_tv"`
		HasMinibar  interface{} `json:"has_minibar"`
		CreatedAt   interface{} `json:"created_at"`
		UpdatedAt   interface{} `json:"updated_at"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convertir ID si es string a array
	if aux.ID != nil {
		d.ID = toStringArray(aux.ID)
	}

	// Convertir Type si es string a array
	if aux.Type != nil {
		d.Type = toStringArray(aux.Type)
	}

	// Convertir Status si es string a array
	if aux.Status != nil {
		d.Status = toStringArray(aux.Status)
	}

	// Convertir Description si es string a array
	if aux.Description != nil {
		d.Description = toStringArray(aux.Description)
	}

	// Convertir Number si es int a array
	if aux.Number != nil {
		d.Number = toIntArray(aux.Number)
	}

	// Convertir Capacity si es int a array
	if aux.Capacity != nil {
		d.Capacity = toIntArray(aux.Capacity)
	}

	// Convertir Floor si es int a array
	if aux.Floor != nil {
		d.Floor = toIntArray(aux.Floor)
	}

	// Convertir Price si es float a array
	if aux.Price != nil {
		d.Price = toFloatArray(aux.Price)
	}

	// Convertir HasWifi si es bool a array
	if aux.HasWifi != nil {
		d.HasWifi = toBoolArray(aux.HasWifi)
	}

	// Convertir HasAC si es bool a array
	if aux.HasAC != nil {
		d.HasAC = toBoolArray(aux.HasAC)
	}

	// Convertir HasTV si es bool a array
	if aux.HasTV != nil {
		d.HasTV = toBoolArray(aux.HasTV)
	}

	// Convertir HasMinibar si es bool a array
	if aux.HasMinibar != nil {
		d.HasMinibar = toBoolArray(aux.HasMinibar)
	}

	// Convertir CreatedAt si es string a array
	if aux.CreatedAt != nil {
		d.CreatedAt = toStringArray(aux.CreatedAt)
	}

	// Convertir UpdatedAt si es string a array
	if aux.UpdatedAt != nil {
		d.UpdatedAt = toStringArray(aux.UpdatedAt)
	}

	return nil
}

// toStringArray convierte un valor (string o []string) a []string
func toStringArray(v interface{}) []string {
	switch val := v.(type) {
	case string:
		return []string{val}
	case []interface{}:
		result := make([]string, len(val))
		for i, item := range val {
			if str, ok := item.(string); ok {
				result[i] = str
			}
		}
		return result
	case []string:
		return val
	default:
		return []string{}
	}
}

// toIntArray convierte un valor (int o []int) a []int
func toIntArray(v interface{}) []int {
	switch val := v.(type) {
	case float64:
		return []int{int(val)}
	case int:
		return []int{val}
	case []interface{}:
		result := make([]int, len(val))
		for i, item := range val {
			switch itemVal := item.(type) {
			case float64:
				result[i] = int(itemVal)
			case int:
				result[i] = itemVal
			}
		}
		return result
	case []int:
		return val
	default:
		return []int{}
	}
}

// toFloatArray convierte un valor (float o []float) a []float64
func toFloatArray(v interface{}) []float64 {
	switch val := v.(type) {
	case float64:
		return []float64{val}
	case int:
		return []float64{float64(val)}
	case []interface{}:
		result := make([]float64, len(val))
		for i, item := range val {
			switch itemVal := item.(type) {
			case float64:
				result[i] = itemVal
			case int:
				result[i] = float64(itemVal)
			}
		}
		return result
	case []float64:
		return val
	default:
		return []float64{}
	}
}

// toBoolArray convierte un valor (bool o []bool) a []bool
func toBoolArray(v interface{}) []bool {
	switch val := v.(type) {
	case bool:
		return []bool{val}
	case []interface{}:
		result := make([]bool, len(val))
		for i, item := range val {
			if boolVal, ok := item.(bool); ok {
				result[i] = boolVal
			}
		}
		return result
	case []bool:
		return val
	default:
		return []bool{}
	}
}

// Room representa el modelo de Room de rooms-api (para mapeo)
type Room struct {
	ID          uint      `json:"id"`
	Number      string    `json:"number"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Capacity    int       `json:"capacity"`
	Floor       int       `json:"floor"`
	HasWifi     bool      `json:"has_wifi"`
	HasAC       bool      `json:"has_ac"`
	HasTV       bool      `json:"has_tv"`
	HasMinibar  bool      `json:"has_minibar"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RoomEvent representa un evento de RabbitMQ para rooms
type RoomEvent struct {
	EventType string `json:"event_type"` // "created", "updated", "deleted"
	RoomID    uint   `json:"room_id"`
	Timestamp string `json:"timestamp"` // RFC3339 format
}

// SolrRoomWrite representa el documento para escritura en Solr (valores simples)
type SolrRoomWrite struct {
	ID          string  `json:"id"`
	Number      int     `json:"number"`
	Type        string  `json:"type"`
	Status      string  `json:"status"`
	Price       float64 `json:"price"`
	Description string  `json:"description,omitempty"`
	Capacity    int     `json:"capacity"`
	Floor       int     `json:"floor"`
	HasWifi     bool    `json:"has_wifi"`
	HasAC       bool    `json:"has_ac"`
	HasTV       bool    `json:"has_tv"`
	HasMinibar  bool    `json:"has_minibar"`
	CreatedAt   string  `json:"created_at"` // RFC3339
	UpdatedAt   string  `json:"updated_at"` // RFC3339
}

// ToSolrDocument convierte Room a SolrRoomWrite para escritura en Solr
func (r *Room) ToSolrDocument() *SolrRoomWrite {
	// Convertir Number string a int
	numberInt, _ := strconv.Atoi(r.Number)

	return &SolrRoomWrite{
		ID:          fmt.Sprintf("%d", r.ID), // Convertir uint a string
		Number:      numberInt,
		Type:        r.Type,
		Status:      r.Status,
		Price:       r.Price,
		Description: r.Description,
		Capacity:    r.Capacity,
		Floor:       r.Floor,
		HasWifi:     r.HasWifi,
		HasAC:       r.HasAC,
		HasTV:       r.HasTV,
		HasMinibar:  r.HasMinibar,
		CreatedAt:   r.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.Format(time.RFC3339),
	}
}
