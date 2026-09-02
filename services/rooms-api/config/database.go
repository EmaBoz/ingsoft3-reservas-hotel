package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getenvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func InitMySQL() *gorm.DB {
	_ = godotenv.Load()

	// Variables de entorno para conectar a mysql-rooms
	host := getenvOrDefault("DB_HOST", "mysql-rooms")
	port := getenvOrDefault("DB_PORT", "3306")
	user := getenvOrDefault("DB_USER", "roomsuser")
	password := getenvOrDefault("DB_PASSWORD", "roomspass")
	dbName := getenvOrDefault("DB_NAME", "roomsdb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	var db *gorm.DB
	var err error

	// Reintentos con backoff exponencial
	for attempt := 1; attempt <= 10; attempt++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		})

		if err == nil {
			sqlDB, dbErr := db.DB()
			if dbErr == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					// Configuración del pool de conexiones
					sqlDB.SetMaxIdleConns(10)
					sqlDB.SetMaxOpenConns(100)
					sqlDB.SetConnMaxLifetime(time.Hour)

					log.Printf("✅ Conectado a MySQL Rooms: %s@%s:%s/%s (intento %d)",
						user, host, port, dbName, attempt)
					return db
				}
			}
		}

		wait := time.Duration(attempt*2) * time.Second
		log.Printf("⏳ MySQL Rooms no listo (intento %d): %v. Reintentando en %s...",
			attempt, err, wait)
		time.Sleep(wait)
	}

	log.Fatalf("❌ No se pudo conectar a MySQL Rooms tras reintentos: %v", err)
	return nil
}
