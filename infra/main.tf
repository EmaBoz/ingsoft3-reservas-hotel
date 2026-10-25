# Entorno preprod (TP8, plan C): la app completa en contenedores, en el Docker
# de mi máquina, declarada acá. Convive con los proyectos de Compose de QA y
# PROD (sus propios nombres, red, volúmenes y puertos) y no los toca.
#
# Cómo se usa (desde infra/):
#   terraform init -backend-config="path=$HOME/tfstate/reservas-preprod.tfstate"
#   terraform plan    -var-file=local.tfvars
#   terraform apply   -var-file=local.tfvars
#   terraform destroy -var-file=local.tfvars
# local.tfvars (ignorado por git) tiene db_password, jwt_secret e image_tag.

terraform {
  required_version = ">= 1.5"
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 4.0"
    }
  }
  # el estado vive FUERA del repo: el path se pasa en el init con -backend-config
  backend "local" {}
}

provider "docker" {} # habla con el Docker local (Docker Desktop)

locals {
  prefijo  = "reservas-preprod"
  registry = "ghcr.io/emaboz/ingsoft3-reservas-hotel"
}

# ---------- red y volúmenes ----------

resource "docker_network" "app" {
  name = local.prefijo
}

resource "docker_volume" "mysql_users" {
  name = "${local.prefijo}-mysql-users"
}

resource "docker_volume" "mysql_rooms" {
  name = "${local.prefijo}-mysql-rooms"
}

resource "docker_volume" "mongo" {
  name = "${local.prefijo}-mongo"
}

resource "docker_volume" "rabbitmq" {
  name = "${local.prefijo}-rabbitmq"
}

resource "docker_volume" "solr" {
  name = "${local.prefijo}-solr"
}

# ---------- imágenes ----------
# keep_locally: al destruir preprod la imagen sale del estado pero NO se borra
# del Docker local, porque son las mismas que están usando QA y PROD.

resource "docker_image" "mysql" {
  name         = "mysql:8.0"
  keep_locally = true
}

resource "docker_image" "mongo" {
  name         = "mongo:7.0"
  keep_locally = true
}

resource "docker_image" "rabbitmq" {
  name         = "rabbitmq:3-management"
  keep_locally = true
}

resource "docker_image" "solr" {
  name         = "solr:9.5"
  keep_locally = true
}

resource "docker_image" "memcached" {
  name         = "memcached:1.6-alpine"
  keep_locally = true
}

# las de la app: las que publicó el pipeline del TP7, por su sha-<commit>
resource "docker_image" "users_api" {
  name         = "${local.registry}-users-api:${var.image_tag}"
  keep_locally = true
}

resource "docker_image" "rooms_api" {
  name         = "${local.registry}-rooms-api:${var.image_tag}"
  keep_locally = true
}

resource "docker_image" "reservations_api" {
  name         = "${local.registry}-reservations-api:${var.image_tag}"
  keep_locally = true
}

resource "docker_image" "search_api" {
  name         = "${local.registry}-search-api:${var.image_tag}"
  keep_locally = true
}

resource "docker_image" "frontend" {
  name         = "${local.registry}-frontend:${var.image_tag}"
  keep_locally = true
}

# ---------- bases, cola, búsqueda y caché ----------
# wait = true: Terraform espera a que el healthcheck dé sano antes de crear lo
# que depende de estos contenedores (depends_on solo espera que se CREEN).

resource "docker_container" "mysql_users" {
  name  = "${local.prefijo}-mysql-users"
  image = docker_image.mysql.image_id
  env = [
    "MYSQL_ROOT_PASSWORD=${var.db_password}",
    "MYSQL_DATABASE=usersdb",
    "MYSQL_USER=user",
    "MYSQL_PASSWORD=${var.db_password}",
  ]
  networks_advanced {
    name = docker_network.app.name
  }
  volumes {
    volume_name    = docker_volume.mysql_users.name
    container_path = "/var/lib/mysql"
  }
  healthcheck {
    test     = ["CMD-SHELL", "mysqladmin ping -h 127.0.0.1 -uroot -p\"$MYSQL_ROOT_PASSWORD\""]
    interval = "5s"
    timeout  = "3s"
    retries  = 30
  }
  wait         = true
  wait_timeout = 180
}

resource "docker_container" "mysql_rooms" {
  name  = "${local.prefijo}-mysql-rooms"
  image = docker_image.mysql.image_id
  env = [
    "MYSQL_ROOT_PASSWORD=${var.db_password}",
    "MYSQL_DATABASE=roomsdb",
    "MYSQL_USER=roomsuser",
    "MYSQL_PASSWORD=${var.db_password}",
  ]
  networks_advanced {
    name = docker_network.app.name
  }
  volumes {
    volume_name    = docker_volume.mysql_rooms.name
    container_path = "/var/lib/mysql"
  }
  healthcheck {
    test     = ["CMD-SHELL", "mysqladmin ping -h 127.0.0.1 -uroot -p\"$MYSQL_ROOT_PASSWORD\""]
    interval = "5s"
    timeout  = "3s"
    retries  = 30
  }
  wait         = true
  wait_timeout = 180
}

resource "docker_container" "mongodb" {
  name  = "${local.prefijo}-mongodb"
  image = docker_image.mongo.image_id
  env = [
    "MONGO_INITDB_ROOT_USERNAME=root",
    "MONGO_INITDB_ROOT_PASSWORD=${var.db_password}",
  ]
  networks_advanced {
    name = docker_network.app.name
  }
  volumes {
    volume_name    = docker_volume.mongo.name
    container_path = "/data/db"
  }
  healthcheck {
    test     = ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
    interval = "5s"
    timeout  = "3s"
    retries  = 30
  }
  wait         = true
  wait_timeout = 180
}

resource "docker_container" "rabbitmq" {
  name  = "${local.prefijo}-rabbitmq"
  image = docker_image.rabbitmq.image_id
  env = [
    "RABBITMQ_DEFAULT_USER=reservas",
    "RABBITMQ_DEFAULT_PASS=${var.db_password}",
  ]
  networks_advanced {
    name = docker_network.app.name
  }
  volumes {
    volume_name    = docker_volume.rabbitmq.name
    container_path = "/var/lib/rabbitmq"
  }
  healthcheck {
    test         = ["CMD", "rabbitmq-diagnostics", "-q", "ping"]
    interval     = "10s"
    timeout      = "10s"
    retries      = 20
    start_period = "40s"
  }
  wait         = true
  wait_timeout = 240
}

resource "docker_container" "solr" {
  name    = "${local.prefijo}-solr"
  image   = docker_image.solr.image_id
  command = ["solr-precreate", "rooms-core"]
  env     = ["SOLR_HEAP=512m"]
  networks_advanced {
    name = docker_network.app.name
  }
  volumes {
    volume_name    = docker_volume.solr.name
    container_path = "/var/solr"
  }
  healthcheck {
    test         = ["CMD", "curl", "-f", "http://localhost:8983/solr/rooms-core/admin/ping"]
    interval     = "10s"
    timeout      = "5s"
    retries      = 20
    start_period = "30s"
  }
  wait         = true
  wait_timeout = 240
}

resource "docker_container" "memcached" {
  name    = "${local.prefijo}-memcached"
  image   = docker_image.memcached.image_id
  command = ["memcached", "-m", "64"]
  networks_advanced {
    name = docker_network.app.name
  }
}

# ---------- las apis ----------
# Las direcciones y cadenas de conexión se DERIVAN de los recursos de arriba
# (sus nombres y la contraseña), no se copian a mano.

resource "docker_container" "users_api" {
  name    = "${local.prefijo}-users-api"
  image   = docker_image.users_api.image_id
  restart = "unless-stopped"
  env = [
    "DB_HOST=${docker_container.mysql_users.name}",
    "DB_PORT=3306",
    "DB_USER=user",
    "DB_PASSWORD=${var.db_password}",
    "DB_NAME=usersdb",
    "JWT_SECRET=${var.jwt_secret}",
    "MEMCACHED_HOST=${docker_container.memcached.name}",
    "MEMCACHED_PORT=11211",
    "MEMCACHED_TTL=600",
  ]
  networks_advanced {
    name = docker_network.app.name
  }
  ports {
    internal = 8080
    external = 8381
  }
}

resource "docker_container" "search_api" {
  name    = "${local.prefijo}-search-api"
  image   = docker_image.search_api.image_id
  restart = "unless-stopped"
  env = [
    "SOLR_URL=http://${docker_container.solr.name}:8983/solr",
    "SOLR_CORE=rooms-core",
    "MEMCACHED_HOST=${docker_container.memcached.name}",
    "MEMCACHED_PORT=11211",
    "LOCAL_CACHE_TTL_SECONDS=60",
    "DISTRIBUTED_CACHE_TTL_SECONDS=300",
    "RABBITMQ_URL=amqp://reservas:${var.db_password}@${docker_container.rabbitmq.name}:5672/",
    # rooms-api depende de search-api, así que acá no se puede referenciar el
    # recurso (sería un ciclo): se arma con el mismo prefijo que le da su nombre
    "ROOMS_API_BASE_URL=http://${local.prefijo}-rooms-api:8080",
    "JWT_SECRET=${var.jwt_secret}",
    "PORT=8083",
    "GIN_MODE=release",
  ]
  networks_advanced {
    name = docker_network.app.name
  }
  ports {
    internal = 8083
    external = 8384
  }
}

resource "docker_container" "rooms_api" {
  name    = "${local.prefijo}-rooms-api"
  image   = docker_image.rooms_api.image_id
  restart = "unless-stopped"
  env = [
    "DB_HOST=${docker_container.mysql_rooms.name}",
    "DB_PORT=3306",
    "DB_USER=roomsuser",
    "DB_PASSWORD=${var.db_password}",
    "DB_NAME=roomsdb",
    "JWT_SECRET=${var.jwt_secret}",
    "RABBITMQ_URL=amqp://reservas:${var.db_password}@${docker_container.rabbitmq.name}:5672/",
    "MEMCACHED_HOST=${docker_container.memcached.name}",
    "MEMCACHED_PORT=11211",
    "MEMCACHED_TTL=300",
    "SEARCH_API_BASE_URL=http://${docker_container.search_api.name}:8083",
  ]
  networks_advanced {
    name = docker_network.app.name
  }
  ports {
    internal = 8080
    external = 8382
  }
}

resource "docker_container" "reservations_api" {
  name    = "${local.prefijo}-reservations-api"
  image   = docker_image.reservations_api.image_id
  restart = "unless-stopped"
  env = [
    "MONGO_URI=mongodb://root:${var.db_password}@${docker_container.mongodb.name}:27017/?authSource=admin",
    "MONGO_DB_NAME=mongodb",
    "RABBITMQ_URL=amqp://reservas:${var.db_password}@${docker_container.rabbitmq.name}:5672/",
  ]
  networks_advanced {
    name = docker_network.app.name
  }
  ports {
    internal = 8080
    external = 8383
  }
}

# ---------- el front ----------
# La misma imagen de QA y PROD: la dirección de cada api se la pasa este
# entorno por variable (la plantilla de nginx del TP6), derivada de las apis.

resource "docker_container" "frontend" {
  name    = "${local.prefijo}-frontend"
  image   = docker_image.frontend.image_id
  restart = "unless-stopped"
  env = [
    "USERS_API_URL=http://${docker_container.users_api.name}:8080",
    "ROOMS_API_URL=http://${docker_container.rooms_api.name}:8080",
    "RESERVATIONS_API_URL=http://${docker_container.reservations_api.name}:8080",
    "SEARCH_API_URL=http://${docker_container.search_api.name}:8083",
    "DNS_RESOLVER=127.0.0.11",
  ]
  networks_advanced {
    name = docker_network.app.name
  }
  ports {
    internal = 80
    external = 3300
  }
}
