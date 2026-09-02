# Rooms API

A microservice for managing hotel rooms using MySQL and Go with Gin framework and GORM.

## Features

- **CRUD Operations**: Create, read, update, and delete rooms
- **Room Management**: Track room status, type, pricing, and amenities
- **Filtering & Search**: Filter rooms by type, status, floor, price range, and amenities
- **Pagination**: Support for paginated room listings
- **MySQL Integration**: Uses MySQL with GORM for data persistence and migrations
- **RESTful API**: Clean REST API design with proper HTTP status codes

## Room Model

### Room Types
- `single` - Single occupancy room
- `double` - Double occupancy room  
- `suite` - Suite room
- `deluxe` - Deluxe room
- `standard` - Standard room

### Room Status
- `available` - Room is available for booking
- `occupied` - Room is currently occupied
- `maintenance` - Room is under maintenance
- `reserved` - Room is reserved but not yet occupied

### Room Properties
- **Basic Info**: Number, type, status, price, description, capacity, floor
- **Amenities**: WiFi, AC, TV, minibar availability
- **Timestamps**: Created and updated timestamps

## API Endpoints

### Health Check
- `GET /health` - Service health status

### Room Management
- `POST /api/v1/rooms` - Create a new room
- `GET /api/v1/rooms` - Get all rooms (with filtering and pagination)
- `GET /api/v1/rooms/available` - Get available rooms only
- `GET /api/v1/rooms/:id` - Get room by ID
- `GET /api/v1/rooms/number/:number` - Get room by number
- `PUT /api/v1/rooms/:id` - Update room
- `PATCH /api/v1/rooms/:id/status` - Update room status
- `DELETE /api/v1/rooms/:id` - Delete room

### Query Parameters for GET /api/v1/rooms
- `type` - Filter by room type
- `status` - Filter by room status
- `floor` - Filter by floor number
- `min_price` - Minimum price filter
- `max_price` - Maximum price filter
- `has_wifi` - Filter by WiFi availability
- `has_ac` - Filter by AC availability
- `has_tv` - Filter by TV availability
- `has_minibar` - Filter by minibar availability
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)

## Environment Variables

- `DB_HOST` - MySQL host (default: localhost)
- `DB_PORT` - MySQL port (default: 3306)
- `DB_USER` - MySQL user (default: root)
- `DB_PASSWORD` - MySQL password (default: root)
- `DB_NAME` - Database name (default: roomsdb)
- `PORT` - Server port (default: 8080)

## Example Usage

### Create a Room
```bash
curl -X POST http://localhost:8081/api/v1/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "number": "101",
    "type": "double",
    "price": 150.00,
    "description": "Comfortable double room with city view",
    "capacity": 2,
    "floor": 1,
    "has_wifi": true,
    "has_ac": true,
    "has_tv": true,
    "has_minibar": false
  }'
```

### Get Available Rooms
```bash
curl "http://localhost:8081/api/v1/rooms/available?type=double&floor=1&min_price=100&max_price=200"
```

### Update Room Status
```bash
curl -X PATCH http://localhost:8081/api/v1/rooms/ROOM_ID/status \
  -H "Content-Type: application/json" \
  -d '{"status": "occupied"}'
```

## Running the Service

### Using Docker Compose
```bash
docker-compose up rooms-api
```

### Running Locally
```bash
# Install dependencies
go mod tidy

# Run the service
go run ./cmd/server
```

## Database Schema

The service uses MySQL with GORM for automatic migrations. The `rooms` table structure:

```sql
CREATE TABLE `rooms` (
  `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `number` VARCHAR(20) NOT NULL UNIQUE,
  `type` VARCHAR(20) NOT NULL,
  `status` VARCHAR(20) NOT NULL DEFAULT 'available',
  `price` DECIMAL(10,2) NOT NULL,
  `description` TEXT,
  `capacity` INT NOT NULL,
  `floor` INT NOT NULL,
  `has_wifi` TINYINT(1) DEFAULT 0,
  `has_ac` TINYINT(1) DEFAULT 0,
  `has_tv` TINYINT(1) DEFAULT 0,
  `has_minibar` TINYINT(1) DEFAULT 0,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  INDEX idx_rooms_deleted_at (`deleted_at`)
);
```

**Note**: The table is automatically created and migrated by GORM on application startup.

## Error Handling

The API returns appropriate HTTP status codes and error messages:

- `400 Bad Request` - Invalid request data
- `404 Not Found` - Room not found
- `409 Conflict` - Room number already exists
- `500 Internal Server Error` - Database or server error

## Development

### Project Structure
```
rooms-api/
├── cmd/server/          # Main application entry point
├── config/              # Database configuration
├── controllers/         # HTTP handlers
├── domain/              # Data models and DTOs
├── repositories/        # Data access layer
├── services/            # Business logic layer
├── utils/               # Utility functions
├── dockerfile           # Docker configuration
└── go.mod              # Go module dependencies
```
