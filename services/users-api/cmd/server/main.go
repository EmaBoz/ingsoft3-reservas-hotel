package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"users-api/config"
	"users-api/controllers"
	"users-api/repositories"
	"users-api/services"

	"github.com/gin-gonic/gin"
)

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	// DB
	db := config.InitDatabase()

	// Repo DB
	userRepo := repositories.NewUserRepository(db)

	// Cache repo (Memcached)
	mHost := getenv("MEMCACHED_HOST", "memcached")
	mPort := getenv("MEMCACHED_PORT", "11211")
	ttlStr := getenv("MEMCACHED_TTL", "600") // segundos
	ttlSec, _ := strconv.Atoi(ttlStr)
	userCache := repositories.NewUserCacheRepository(mHost, mPort, time.Duration(ttlSec)*time.Second)

	// Service (IMPORTANTE: 2 argumentos)
	userSvc := services.NewUserService(userRepo, userCache)

	// HTTP
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	uc := controllers.NewUserController(userSvc)

	r.GET("/users", uc.GetAllUsers)
	r.POST("/users", uc.CreateUser)
	r.GET("/users/:id", uc.GetUserByID)
	r.POST("/login", uc.Login)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
