# Search API - Microservicio de Búsqueda con Solr

Microservicio de búsqueda para el sistema de reservas de hotel, usando Apache Solr como motor de búsqueda, doble caché (CCache local + Memcached distribuido), y sincronización en tiempo real vía RabbitMQ.

## Características

- **Motor de búsqueda**: Apache Solr 9.5 para búsquedas avanzadas y filtrado
- **Caché multinivel**:
  - CCache (local, in-memory, TTL 60s)
  - Memcached (distribuido, TTL 300s)
- **Sincronización en tiempo real**: Consumer de RabbitMQ que escucha eventos de rooms-api
- **Autenticación**: Soporte JWT (opcional) para proteger endpoints
- **API RESTful**: Endpoint de búsqueda con filtros, ordenamiento y paginación

## Arquitectura

```
┌─────────────┐
│  Frontend   │
└─────┬───────┘
      │ GET /api/search/rooms
      │
┌─────▼──────────────────────────────────────────────┐
│              Search API (Go + Gin)                  │
│                                                      │
│  ┌────────────┐  ┌──────────────┐  ┌─────────────┐│
│  │ Controller │→ │   Service    │→ │ Repositories││
│  └────────────┘  └──────────────┘  └─────────────┘│
│                         │                           │
│         ┌───────────────┼────────────────┐         │
│         │               │                │         │
│    ┌────▼────┐    ┌────▼─────┐    ┌────▼────┐    │
│    │ CCache  │    │ Memcached│    │  Solr   │    │
│    │ (Local) │    │(Distrib.)│    │ (Search)│    │
│    └─────────┘    └──────────┘    └─────────┘    │
│                                                     │
│         RabbitMQ Consumer (rooms events)           │
│                      │                              │
└──────────────────────┼──────────────────────────────┘
                       │
              ┌────────▼─────────┐
              │   RabbitMQ       │
              │  Exchange:rooms  │
              │                  │
              │  • room.created  │
              │  • room.updated  │
              │  • room.deleted  │
              └────────▲─────────┘
                       │
              ┌────────┴─────────┐
              │    Rooms API     │
              │  (Publisher)     │
              └──────────────────┘
```

## Estructura del Proyecto

```
search-api/
├── cmd/server/main.go              # Entry point
├── controllers/
│   └── search-controller.go        # Handlers HTTP
├── services/
│   └── search-service.go           # Lógica de negocio
├── repositories/
│   ├── solr-repository.go          # Acceso a Solr
│   ├── cache-local-repo.go         # CCache
│   └── cache-distributed-repo.go   # Memcached
├── consumers/
│   └── rooms-consumer.go           # Consumer RabbitMQ
├── domain/
│   ├── room_search.go              # Modelos Solr
│   └── dto.go                      # DTOs request/response
├── config/
│   ├── solr.go                     # Config Solr
│   ├── rabbitmq.go                 # Config RabbitMQ
│   ├── cache.go                    # Config cachés
│   └── http_clients.go             # Cliente rooms-api
├── utils/
│   ├── jwt.go                      # Validación JWT
│   └── errors.go                   # Errores custom
├── scripts/
│   ├── schema.json                 # Schema de Solr
│   └── README.md                   # Docs de Solr
├── Dockerfile
├── go.mod
└── README.md
```

## Variables de Entorno

```bash
# Solr
SOLR_URL=http://localhost:8983/solr
SOLR_CORE=rooms-core

# Memcached
MEMCACHED_HOST=localhost
MEMCACHED_PORT=11211
LOCAL_CACHE_TTL_SECONDS=60
DISTRIBUTED_CACHE_TTL_SECONDS=300

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Rooms API (para sincronización)
ROOMS_API_BASE_URL=http://localhost:8081

# JWT (mismo secret que users-api)
JWT_SECRET=supersecreto_cámbialo

# Server
PORT=8083
GIN_MODE=release
```

## API Endpoints

### Búsqueda de Habitaciones

**GET** `/api/search/rooms`

Query Parameters:
- `q` (string): Texto libre para buscar en número, tipo, descripción
- `type` (string): Filtrar por tipo (single, double, suite, deluxe, standard)
- `status` (string): Filtrar por estado (available, occupied, maintenance, reserved)
- `floor` (int): Filtrar por piso
- `min_price` (float): Precio mínimo
- `max_price` (float): Precio máximo
- `has_wifi` (bool): Filtrar por WiFi
- `has_ac` (bool): Filtrar por aire acondicionado
- `has_tv` (bool): Filtrar por TV
- `has_minibar` (bool): Filtrar por minibar
- `sort` (string): Campo de ordenamiento (price, -price, floor, -floor, capacity, -capacity)
- `page` (int): Número de página (default: 1)
- `limit` (int): Tamaño de página (default: 10, max: 50)

Ejemplo:
```bash
curl "http://localhost:8083/api/search/rooms?status=available&has_wifi=true&min_price=50&max_price=200&sort=-price&page=1&limit=10"
```

Respuesta:
```json
{
  "page": 1,
  "limit": 10,
  "total": 45,
  "results": [
    {
      "id": "1",
      "number": "101",
      "type": "double",
      "status": "available",
      "price": 150.00,
      "capacity": 2,
      "floor": 1,
      "has_wifi": true,
      "has_ac": true,
      "has_tv": true,
      "has_minibar": false
    }
  ]
}
```

### Health Check

**GET** `/health`

Respuesta:
```json
{
  "status": "healthy",
  "service": "search-api"
}
```

## Instalación y Ejecución

### Con Docker Compose (Recomendado)

```bash
# Desde el directorio raíz del proyecto
docker-compose up -d solr memcached rabbitmq rooms-api search-api
```

### Local (Desarrollo)

```bash
# 1. Asegurar que Solr, Memcached y RabbitMQ estén corriendo
docker-compose up -d solr memcached rabbitmq

# 2. Crear el core de Solr (si no existe)
curl "http://localhost:8983/solr/admin/cores?action=CREATE&name=rooms-core&configSet=_default"

# 3. Aplicar schema
curl -X POST -H 'Content-type:application/json' \
  http://localhost:8983/solr/rooms-core/schema \
  --data-binary @scripts/schema.json

# 4. Instalar dependencias
cd search-api
go mod download

# 5. Ejecutar
go run cmd/server/main.go
```

## Flujo de Sincronización

1. **Evento en rooms-api**: Se crea/actualiza/elimina una habitación
2. **Publicación a RabbitMQ**: rooms-api publica evento al exchange "rooms"
3. **Consumer recibe evento**: search-api consume el mensaje
4. **Obtención de datos**: Consumer llama a rooms-api GET /api/v1/rooms/:id
5. **Indexación en Solr**: Se indexa/actualiza/elimina el documento en Solr
6. **Invalidación de caché**: Se limpia el caché local

## Estrategia de Caché

### Flujo de Búsqueda con Caché

```
Request → CCache (L1)
          ├─ HIT → Return cached
          └─ MISS → Memcached (L2)
                    ├─ HIT → Save to CCache → Return
                    └─ MISS → Solr
                              └─ Save to Memcached → Save to CCache → Return
```

### Cache Keys

Las claves de caché se generan con un hash SHA-256 de todos los parámetros de búsqueda:
```
search:{hash_de_parametros}
```

### TTLs

- **CCache (Local)**: 60 segundos
- **Memcached (Distribuido)**: 300 segundos (5 minutos)

## Testing

### Crear una habitación en rooms-api

```bash
curl -X POST http://localhost:8081/api/v1/rooms \
  -H "Content-Type: application/json" \
  -d '{
    "number": "101",
    "type": "double",
    "price": 120.50,
    "description": "Habitación doble con vista al mar",
    "capacity": 2,
    "floor": 1,
    "has_wifi": true,
    "has_ac": true,
    "has_tv": true,
    "has_minibar": false
  }'
```

### Verificar que se indexó en Solr

```bash
curl "http://localhost:8983/solr/rooms-core/select?q=*:*&wt=json"
```

### Buscar con search-api

```bash
# Buscar habitaciones disponibles con WiFi
curl "http://localhost:8083/api/search/rooms?status=available&has_wifi=true"

# Buscar por rango de precio
curl "http://localhost:8083/api/search/rooms?min_price=100&max_price=200&sort=price"

# Búsqueda de texto libre
curl "http://localhost:8083/api/search/rooms?q=mar"
```

## Logs

Los logs muestran:
- Cache HIT/MISS (local y distribuido)
- Eventos RabbitMQ recibidos y procesados
- Indexaciones/eliminaciones en Solr
- Errores de sincronización

Ejemplo:
```
2024/11/14 10:15:23 Successfully connected to RabbitMQ
2024/11/14 10:15:23 Queue search-api-rooms-queue bound to exchange rooms with routing key room.created
2024/11/14 10:15:23 RabbitMQ consumer started, waiting for messages...
2024/11/14 10:15:23 Starting HTTP server on port 8083
2024/11/14 10:16:45 Received message: {"event_type":"created","room_id":1,"timestamp":"2024-11-14T10:16:45Z"}
2024/11/14 10:16:45 Room 1 indexed successfully in Solr
2024/11/14 10:16:45 Event processed successfully: created for room 1
2024/11/14 10:17:10 Cache MISS: search:a3f2... - Querying Solr
2024/11/14 10:17:11 Cache HIT (local): search:a3f2...
```

## Troubleshooting

### El consumer no recibe eventos

1. Verificar que rooms-api esté publicando eventos
2. Verificar RabbitMQ Management: http://localhost:15672
3. Verificar que el exchange "rooms" exista
4. Verificar que la queue esté bindeada

### Solr no responde

```bash
# Verificar health
curl http://localhost:8983/solr/rooms-core/admin/ping

# Ver logs de Solr
docker logs solr
```

### Memcached no funciona

```bash
# Verificar que esté corriendo
docker ps | grep memcached

# Test manual
telnet localhost 11211
> stats
```

## Próximas Mejoras

- [ ] Soporte para búsqueda por disponibilidad de fechas (integración con reservations-api)
- [ ] Aggregations/facets en Solr (estadísticas por tipo, piso, etc.)
- [ ] Rate limiting en endpoints
- [ ] Metrics y observabilidad (Prometheus/Grafana)
- [ ] Auto-reindexación completa periódica
- [ ] Búsqueda fuzzy y autocompletado
