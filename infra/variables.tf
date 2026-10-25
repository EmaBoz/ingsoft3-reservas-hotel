variable "db_password" {
  description = "Contraseña de las bases (MySQL x2, Mongo) y de RabbitMQ de preprod"
  type        = string
  sensitive   = true # no se imprime en el plan ni en los logs
}

variable "jwt_secret" {
  description = "Clave con la que las apis firman y validan los JWT"
  type        = string
  sensitive   = true
}

variable "image_tag" {
  description = "Etiqueta de las imágenes del TP7 a correr: sha-<commit> de un merge publicado"
  type        = string # sin default: el pipeline no publica otra etiqueta (no hay latest)
}
