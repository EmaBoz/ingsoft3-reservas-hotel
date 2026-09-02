package main

import (
	"log"
	"reservations-api/config"
	"reservations-api/controllers"
	"reservations-api/events"
	"reservations-api/repositories"
	"reservations-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializar base de datos MongoDB
	db := config.InitDatabase()
	if db == nil {
		log.Fatal("Error inicializando la base de datos")
	}

	// Inicializar RabbitMQ
	rabbitConn := config.InitRabbitMQ()
	if rabbitConn == nil {
		log.Fatal("Error inicializando RabbitMQ")
	}
	defer rabbitConn.Close()

	// Inicializar event publisher
	publisher, err := events.NewEventPublisher(rabbitConn)
	if err != nil {
		log.Fatalf("Error inicializando event publisher: %v", err)
	}
	defer publisher.Close()

	// Inicializar capas
	reservationRepo := repositories.NewReservationRepository(db)
	reservationService := services.NewReservationService(reservationRepo, publisher)
	reservationController := controllers.NewReservationController(reservationService)
	// Configurar router
	router := gin.Default()
	// --- ACA agregás el health check ---
	router.POST("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	// Rutas
	api := router.Group("/api")
	{
		api.GET("/reservations/users/:user_id/myreservations", reservationController.GetmyReservations)
		api.GET("/reservations", reservationController.GetAllReservations)
		api.POST("/reservations", reservationController.CreateReservation)
		api.GET("/reservations/:id", reservationController.GetReservationByID)
		api.DELETE("/reservations/:id", reservationController.DeleteReservation)
	}
	// Iniciar servidor
	log.Println("Reservations API running on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
