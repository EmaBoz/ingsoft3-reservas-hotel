# 🏨 Sistema de Reservas de Hotel — Arquitectura de Microservicios

Aplicación web para la gestión y reserva de habitaciones de hotel, construida con una
**arquitectura de microservicios** en Go y un frontend en React. Los servicios se
comunican de forma **event-driven** mediante RabbitMQ y usan cachés (Memcached + caché
local) y un motor de búsqueda (Apache Solr) para escalar las consultas.

> App del semestre para la materia **Ingeniería de Software 3** (UCC) — este repositorio
> acompaña toda la cursada: Git colaborativo (TP1), contenedores (TP2), planificación
> (TP3), CI (TP4) y lo que sigue.

---

## 🧰 Stack tecnológico

| Capa | Tecnologías |
|------|-------------|
| **Backend** | Go 1.21 · [Gin](https://gin-gonic.com/) |
| **Frontend** | React 18 · TypeScript · Vite · TailwindCSS · React Router · Axios |
| **Bases de datos** | MySQL 8 (users, rooms) · MongoDB 7 (reservations) |
| **Búsqueda** | Apache Solr 9 |
| **Caché** | Memcached · caché local en memoria |
| **Mensajería** | RabbitMQ (eventos) |
| **Infraestructura** | Docker · Docker Compose |
| **Auth** | JWT (roles: usuario / admin) |

---

## 🏗️ Arquitectura

```
                            ┌──────────────┐
                            │   Frontend   │  React + Vite, servido por nginx (:3000)
                            └──────┬───────┘
                                   │  HTTP (proxy interno de nginx)
        ┌──────────────┬──────────┴──────────┬──────────────────┐
        │              │                      │                  │
 ┌──────▼──────┐ ┌─────▼──────┐      ┌────────▼────────┐ ┌───────▼──────┐
 │  users-api  │ │ rooms-api  │      │ reservations-api│ │  search-api  │
 │   (:8080)   │ │  (:8081)   │      │     (:8082)     │ │   (:8083)    │
 └──────┬──────┘ └─────┬──────┘      └────────┬────────┘ └───────┬──────┘
        │              │                      │                  │
 ┌──────▼──────┐ ┌─────▼──────┐      ┌────────▼────────┐ ┌───────▼──────┐
 │ MySQL       │ │ MySQL      │      │    MongoDB      │ │  Solr        │
 │ (usersdb)   │ │ (roomsdb)  │      │                 │ │              │
 └─────────────┘ └─────┬──────┘      └─────────────────┘ └───────▲──────┘
                       │                                         │
                       │        ┌───────────────┐                │
                       └───────►│   RabbitMQ    │────────────────┘
                    publica     │   (eventos)   │   consume + indexa
                   room.*       └───────────────┘

 Cachés: Memcached (distribuida) + caché local  ·  usada por users, rooms y search
```

**Flujo de eventos:** cuando `rooms-api` crea/actualiza/elimina una habitación, publica un
evento en RabbitMQ. `search-api` lo consume y actualiza su índice en Solr, manteniendo la
búsqueda sincronizada sin acoplar los servicios.

---

## 🚀 Microservicios

| Servicio | Puerto (host) | Puerto interno | Base de datos | Responsabilidad |
|----------|:------:|:------:|---------------|-----------------|
| **users-api** | `8080` | `8080` | MySQL | Registro, login y gestión de usuarios (JWT) |
| **rooms-api** | `8081` | `8080` | MySQL | CRUD de habitaciones; publica eventos |
| **reservations-api** | `8082` | `8080` | MongoDB | Alta y consulta de reservas |
| **search-api** | `8083` | `8083` | Solr | Búsqueda de habitaciones; consume eventos |

### Endpoints principales

**users-api** (`:8080`)
```
POST   /login              Iniciar sesión (devuelve JWT)
POST   /users              Registrar usuario
GET    /users              Listar usuarios
GET    /users/:id          Obtener usuario por ID
```

**rooms-api** (`:8081`)
```
GET    /health
GET    /api/v1/rooms                    Listar habitaciones
GET    /api/v1/rooms/available          Disponibles (vía search-api)
GET    /api/v1/rooms/number/:number     Por número
GET    /api/v1/rooms/:id                Por ID
# Rutas admin (requieren JWT con rol admin)
POST   /api/v1/admin/rooms
PUT    /api/v1/admin/rooms/:id
PATCH  /api/v1/admin/rooms/:id/status
DELETE /api/v1/admin/rooms/:id
```

**reservations-api** (`:8082`)
```
GET    /api/reservations                                  Listar todas
POST   /api/reservations                                  Crear reserva
GET    /api/reservations/:id                              Por ID
DELETE /api/reservations/:id                              Eliminar
GET    /api/reservations/users/:user_id/myreservations   Reservas de un usuario
```

**search-api** (`:8083`)
```
GET    /health
GET    /api/search/rooms   Búsqueda de habitaciones (con filtros)
```

---

## ⚙️ Puesta en marcha (todo en Docker)

Requisito único: **Docker** y **Docker Compose** — no hace falta instalar Go, Node ni las
bases de datos en tu máquina.

```bash
git clone https://github.com/EmaBoz/ingsoft3-reservas-hotel.git
cd ingsoft3-reservas-hotel

cp .env.example .env          # y editá los valores si querés
docker compose up -d --build  # levanta TODO: front + 4 backends + BDs + colas + búsqueda
```

Esto construye y levanta:
- el **frontend** (React + Vite, servido por nginx) en **http://localhost:3000**,
- los **4 microservicios** de backend (`users-api`, `rooms-api`, `reservations-api`,
  `search-api`),
- toda la infraestructura: MySQL ×2, MongoDB, Solr, Memcached, RabbitMQ (panel en
  http://localhost:15672).

Comprobar que levantó todo:

```bash
docker compose ps                 # todos "Up" / "healthy"
curl -s localhost:3000/api/v1/rooms
```

### Variante: bajar las imágenes en vez de construirlas

Las 5 imágenes (los 4 backends + el frontend) están publicadas en GitHub Container
Registry. Para levantar el sistema **sin el código fuente**, descargándolas ya
construidas:

```bash
cp .env.example .env
docker compose -f docker-compose.registry.yml up -d
```

---

## 📁 Estructura del proyecto

```
.
├── docker-compose.yml              # Orquestación completa (construye las imágenes)
├── docker-compose.registry.yml     # Igual, pero descarga las imágenes ya publicadas
├── .env.example                    # Plantilla de variables de entorno
├── frontend/                       # SPA en React + TypeScript + Vite
│   ├── Dockerfile                  # Multi-stage: build (Node) + nginx
│   ├── nginx.conf                  # Proxy hacia cada microservicio
│   └── src/
│       ├── components/             # UI, layout y rutas protegidas
│       ├── pages/                  # Vistas (Home, Login, Rooms, Reservas, Admin…)
│       ├── services/                # Clientes HTTP por dominio
│       ├── context/                 # AuthContext (JWT)
│       └── lib/                     # Configuración de axios
└── services/
    ├── users-api/                  # Go · MySQL · Dockerfile multi-stage
    ├── rooms-api/                  # Go · MySQL · RabbitMQ (publisher) · Dockerfile
    ├── reservations-api/           # Go · MongoDB · RabbitMQ · Dockerfile
    ├── search-api/                 # Go · Solr · RabbitMQ (consumer) · Memcached · Dockerfile
    └── solr/                       # Configset del core de Solr
```

Cada microservicio sigue una arquitectura por capas: `controllers` → `services` →
`repositories` → `domain`, con `config` para infraestructura y `utils` para helpers
(JWT, hashing, errores).

---

## 🔐 Notas de seguridad

Las credenciales de `docker-compose.yml` y `.env.example` son **valores de ejemplo para
desarrollo local**. En un entorno real, definí un `.env` propio con secretos fuertes y
nunca lo subas al repositorio (ya está en `.gitignore`).

---

## 📚 Documentación de la materia

- [`decisiones.md`](decisiones.md) — decisiones técnicas y de diseño, TP a TP.
- [`evidencias.md`](evidencias.md) — evidencias de cada entrega.
