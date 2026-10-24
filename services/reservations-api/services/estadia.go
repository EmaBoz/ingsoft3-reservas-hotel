package services

import (
	"fmt"
	"time"
)

// MaximoDeNoches es la estadía más larga que se puede reservar de una sola vez.
const MaximoDeNoches = 30

// ValidarEstadia controla que las fechas tengan formato yyyy-MM-dd y que la
// estadía sea de entre 1 y MaximoDeNoches noches. Devuelve la cantidad de noches.
func ValidarEstadia(entrada, salida string) (int, error) {
	desde, err := time.Parse("2006-01-02", entrada)
	if err != nil {
		return 0, fmt.Errorf("fecha de entrada inválida: %w", err)
	}

	hasta, err := time.Parse("2006-01-02", salida)
	if err != nil {
		return 0, fmt.Errorf("fecha de salida inválida: %w", err)
	}

	noches := int(hasta.Sub(desde).Hours() / 24)
	if noches < 1 {
		return 0, fmt.Errorf("la estadía tiene que ser de al menos una noche")
	}
	if noches > MaximoDeNoches {
		return 0, fmt.Errorf("la estadía no puede superar las %d noches", MaximoDeNoches)
	}

	return noches, nil
}
