package config

import (
	"fmt"
	"log"
	"os"
	"time"
	"users-api/domain"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getenvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func InitDatabase() *gorm.DB {
	// Carga .env si existe (útil en local; en Docker no hace falta)
	_ = godotenv.Load()

	user := getenvOrDefault("DB_USER", "user")
	pass := getenvOrDefault("DB_PASSWORD", "userpass")
	host := getenvOrDefault("DB_HOST", "mysql") // nombre del servicio en Compose
	port := getenvOrDefault("DB_PORT", "3306")
	name := getenvOrDefault("DB_NAME", "usersdb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		user, pass, host, port, name,
	)

	var db *gorm.DB
	var err error

	// Backoff exponencial simple (hasta ~1 min)
	for attempt := 1; attempt <= 10; attempt++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			// verifica conexión real
			if sqlDB, e := db.DB(); e == nil {
				if pingErr := sqlDB.Ping(); pingErr == nil {
					log.Printf("✅ Conectado a MySQL en %s:%s (intento %d)", host, port, attempt)
					goto MIGRATE
				}
				err = e
			}
		}
		wait := time.Duration(attempt*2) * time.Second
		log.Printf("⏳ MySQL no listo (intento %d): %v. Reintentando en %s...", attempt, err, wait)
		time.Sleep(wait)
	}

	log.Fatalf("❌ No se pudo conectar a la base tras reintentos: %v", err)

MIGRATE:
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("❌ Error al migrar la base de datos: %v", err)
	}
	log.Println("✅ Conexión y migración de base de datos exitosas")
	return db
}
