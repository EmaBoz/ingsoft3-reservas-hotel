package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getenvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func InitDatabase() *mongo.Database {
	// Carga .env si existe (útil en local; en Docker no hace falta)
	_ = godotenv.Load()

	mongoURI := getenvOrDefault("MONGO_URI", "mongodb://mongodb:27017")
	dbName := getenvOrDefault("MONGO_DB_NAME", "mongodb")

	var (
		client *mongo.Client
		err    error
	)

	// Backoff exponencial simple (hasta ~1 min)
	for attempt := 1; attempt <= 10; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		clientOptions := options.Client().
			ApplyURI(mongoURI).
			SetServerSelectionTimeout(5 * time.Second).
			SetConnectTimeout(10 * time.Second)

		client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			// Verificar conexión con ping
			if pingErr := client.Ping(ctx, nil); pingErr != nil {
				err = pingErr
			} else {
				cancel()
				log.Printf("Conectado a MongoDB en %s (intento %d)", mongoURI, attempt)
				db := client.Database(dbName)

				// Crear índices
				if err := createIndexes(db); err != nil {
					log.Printf("Advertencia: Error creando índices: %v", err)
				}

				return db
			}
		}

		cancel()
		wait := time.Duration(attempt*2) * time.Second
		log.Printf("MongoDB no listo (intento %d): %v. Reintentando en %s...", attempt, err, wait)
		time.Sleep(wait)
	}

	log.Fatalf("No se pudo conectar a MongoDB tras reintentos: %v", err)
	return nil
}

// createIndexes crea los índices necesarios para optimizar las queries
func createIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := db.Collection("reservations")

	// Índice compuesto para búsquedas por user_id y status
	userStatusIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "status", Value: 1},
		},
	}

	// Índice compuesto para búsquedas por room_id y fechas (evitar solapamientos)
	roomDatesIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "room_id", Value: 1},
			{Key: "start_date", Value: 1},
			{Key: "end_date", Value: 1},
			{Key: "status", Value: 1},
		},
	}

	// Índice para búsquedas por status
	statusIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "status", Value: 1},
		},
	}

	// Índice para soft deletes
	deletedAtIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "deleted_at", Value: 1},
		},
	}

	// Índice para búsquedas por created_at (ordenamiento temporal)
	createdAtIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "created_at", Value: -1},
		},
	}

	indexes := []mongo.IndexModel{
		userStatusIndex,
		roomDatesIndex,
		statusIndex,
		deletedAtIndex,
		createdAtIndex,
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return err
	}

	log.Println("Índices de MongoDB creados exitosamente")
	return nil
}
