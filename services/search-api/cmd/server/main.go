package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"search-api/config"
	"search-api/consumers"
	"search-api/controllers"
	"search-api/repositories"
	"search-api/services"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting search-api...")

	// Cargar configuraciones
	solrConfig := config.NewSolrConfig()
	rabbitMQConfig := config.NewRabbitMQConfig()
	cacheConfig := config.NewCacheConfig()

	log.Printf("Solr URL: %s", solrConfig.GetCoreURL())
	log.Printf("RabbitMQ URL: %s", rabbitMQConfig.URL)
	log.Printf("Memcached: %s", cacheConfig.GetMemcachedAddress())

	// Inicializar clientes
	roomsAPIClient := config.NewRoomsAPIClient()

	// Inicializar cachés
	localCache := cacheConfig.NewLocalCache()
	memcachedClient := cacheConfig.NewMemcachedClient()

	// Inicializar repositorios
	solrRepo := repositories.NewSolrRepository(solrConfig)
	localCacheRepo := repositories.NewLocalCacheRepository(localCache)
	distributedCacheRepo := repositories.NewDistributedCacheRepository(memcachedClient)

	// Inicializar servicios
	searchService := services.NewSearchService(
		solrRepo,
		localCacheRepo,
		distributedCacheRepo,
		cacheConfig,
	)

	// Inicializar controladores
	searchController := controllers.NewSearchController(searchService)

	// Inicializar consumer de RabbitMQ
	roomsConsumer := consumers.NewRoomsConsumer(
		rabbitMQConfig,
		roomsAPIClient,
		searchService,
	)

	// Iniciar consumer en goroutine
	go func() {
		if err := roomsConsumer.Start(); err != nil {
			log.Printf("Failed to start RabbitMQ consumer: %v", err)
		}
	}()

	// Configurar Gin
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware de CORS
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", searchController.HealthCheck)

	// Rutas API
	api := router.Group("/api")
	{
		search := api.Group("/search")
		{
			// Endpoint de búsqueda (opcional: agregar auth con utils.OptionalAuthMiddleware())
			search.GET("/rooms", searchController.SearchRooms)

			// Ejemplo con autenticación obligatoria (descomentá si querés protegerlo):
			// search.GET("/rooms", utils.AuthMiddleware(), searchController.SearchRooms)
		}
	}

	// Puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Starting HTTP server on port %s", port)

	// Manejar señales de cierre graceful
	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Esperar señal de cierre
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Cerrar consumer de RabbitMQ
	if err := roomsConsumer.Stop(); err != nil {
		log.Printf("Error stopping consumer: %v", err)
	}

	log.Println("Server stopped")
}

// corsMiddleware configura CORS
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
