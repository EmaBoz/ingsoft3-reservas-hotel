output "front_url" {
  value = "http://localhost:${docker_container.frontend.ports[0].external}"
}

output "rooms_api_health" {
  description = "Endpoint de vida"
  value       = "http://localhost:${docker_container.rooms_api.ports[0].external}/health"
}

output "rooms_por_el_front" {
  description = "Endpoint que toca la base (MySQL de rooms), pasando por el front"
  value       = "http://localhost:${docker_container.frontend.ports[0].external}/api/v1/rooms"
}
