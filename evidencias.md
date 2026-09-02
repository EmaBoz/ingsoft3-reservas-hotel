# Evidencias — TP1

> Nota sobre el formato: el TP se ejecutó por terminal (`git` + `gh` CLI) en lugar de clickeando la
> web de GitHub, así que en vez de capturas de pantalla (📸) del navegador, cada evidencia es la
> salida real del comando o de la API de GitHub que muestra el mismo evento que pide la guía. Está
> explicado en `decisiones.md` (declaración de uso de IA).

## 1. Push directo a `main` rechazado

Con la protección de rama activa (`required_pull_request_reviews.required_approving_review_count: 0`,
`enforce_admins: true`), intenté un `git push` directo a `main` — incluso siendo la administradora
del repositorio:

```
remote: error: GH006: Protected branch update failed for refs/heads/main.
remote:
remote: - Changes must be made through a pull request.
To https://github.com/EmaBoz/ingsoft3-reservas-hotel.git
 ! [remote rejected] main -> main (protected branch hook declined)
error: failed to push some refs to 'https://github.com/EmaBoz/ingsoft3-reservas-hotel.git'
```

GitHub rechaza el push porque `main` está protegida y la regla alcanza también al dueño del repo
(`enforce_admins`). El commit local descartado después con `git reset --hard HEAD~1`.

## 2. El PR de la rama B no se puede mergear: conflicto

Después de mergear el PR #2 (`feature/titulo-a`, que cambiaba la primera línea del README a
"versión A"), consulté el estado del PR #3 (`feature/titulo-b`, que cambiaba la misma línea a
"versión B") vía la API de GitHub:

```json
{
  "mergeStateStatus": "DIRTY",
  "mergeable": "CONFLICTING",
  "state": "OPEN",
  "title": "Cambia título del README (versión B)"
}
```

`mergeable: CONFLICTING` es el equivalente exacto del aviso "This branch has conflicts that must be
resolved" que se ve en la web: GitHub ya no puede fusionar automáticamente porque las dos ramas
tocan la misma línea del README.

PR: https://github.com/EmaBoz/ingsoft3-reservas-hotel/pull/3

## 3. Los marcadores de conflicto

Al traer el conflicto a mi máquina (`git merge origin/main` sobre `feature/titulo-b`), Git marcó el
archivo así:

```
<<<<<<< HEAD
# Proyecto IngSoft3 - versión B
=======
# Proyecto IngSoft3 - versión A
>>>>>>> origin/main
Ingeniería de Software 3 (UCC) — Sistema de Reservas de Hotel: ...

## Instalación
```

`HEAD` es la rama en la que estaba parada (`feature/titulo-b`, "versión B"); `origin/main` es lo que
ya había entrado a `main` (la "versión A", mergeada un momento antes). Resolví tomando la versión A,
borré las tres líneas de marcadores, y commiteé (`fix: resuelve conflicto de título tomando la
versión A`).

## 4. La release publicada

`v1.0.0` — https://github.com/EmaBoz/ingsoft3-reservas-hotel/releases/tag/v1.0.0

Incluye: protecciones de rama sobre `main` (con la prueba de rechazo de arriba), el flujo de Pull
Requests funcionando (GitHub Flow + squash merge), el conflicto de merge provocado y resuelto, y
este `decisiones.md`/`evidencias.md`.

---

# Evidencias — TP2

## 1. `docker compose up -d --build` desde cero, sistema funcionando end-to-end

Después de `docker compose down -v` (para partir de cero) y `docker compose up -d --build`, los
12 contenedores (5 propios + 7 de infraestructura) quedan arriba y saludables:

```
NAME               STATUS
frontend           Up (0.0.0.0:3000->80/tcp)
memcached          Up
mongodb            Up (healthy)
mysql-rooms        Up (healthy)
mysql-users        Up (healthy)
rabbitmq           Up (healthy)
reservations-api   Up
rooms-api          Up
search-api         Up (healthy)
solr               Up (healthy)
users-api          Up
```

Prueba end-to-end real, entrando por el frontend (`localhost:3000`) → nginx → microservicio →
base de datos:

```
$ curl -s localhost:3000/api/v1/rooms
{"rooms":[],"total":0,"page":1,"limit":10}

$ curl -s -X POST localhost:3000/users -H "Content-Type: application/json" \
    -d '{"username":"testuser","first_name":"Test","last_name":"User","email":"test@test.com","password":"test1234"}'
{"user":{"id":1,"username":"testuser","email":"test@test.com","first_name":"Test","last_name":"User","role":"normal"}}
```

El segundo comando prueba la cadena completa: el frontend en el puerto 3000 recibe el POST, nginx
lo reenvía por la red interna de compose a `users-api:8080` (sin ningún puerto ni host escrito en
el código del frontend), y `users-api` lo persiste en `mysql-users`.

## 2. Prueba de persistencia

```
=== tras 'down' + 'up' (sin -v): el usuario debe seguir ===
{"users":[{"id":1,"username":"testuser","email":"test@test.com","first_name":"Test","last_name":"User","role":"normal"}]}

=== tras 'down -v' + 'up': el usuario debe haber desaparecido ===
{"users":null}
```

`down` (sin `-v`) recrea los contenedores pero conserva los volúmenes con los datos; `down -v`
además borra los volúmenes, y el usuario creado antes desaparece — el comportamiento esperado.

## 3. Tamaño de imagen final vs. imagen del SDK

```
golang:1.21-alpine                                        337MB   ← compila (SDK completo)
node:22-alpine                                             232MB   ← compila (SDK completo)

ingsoft3-reservas-hotel-users-api:latest                  40.3MB   ← runtime (~8.4× más chica)
ingsoft3-reservas-hotel-rooms-api:latest                  42.1MB   ← runtime (~8× más chica)
ingsoft3-reservas-hotel-reservations-api:latest           38.4MB   ← runtime (~8.8× más chica)
ingsoft3-reservas-hotel-search-api:latest                   42MB   ← runtime (~8× más chica)
ingsoft3-reservas-hotel-frontend:latest                   93.8MB   ← runtime (nginx + estáticos, ~2.5× más chica)
```

Los backends en Go se achican mucho más que el frontend porque el binario compilado no necesita
runtime (Go compila a un ejecutable estático); el frontend en cambio empaqueta nginx completo más
los assets — igual reduce a menos de la mitad del tamaño de la imagen de Node que lo compiló.

## 4. Imágenes publicadas en el registry

Publicadas en GitHub Container Registry (`ghcr.io`), tag `v0.1.0`, visibilidad **pública**:

- https://github.com/EmaBoz?tab=packages (ver los 5 paquetes:
  `ingsoft3-reservas-hotel-users-api`, `-rooms-api`, `-reservations-api`, `-search-api`,
  `-frontend`)

Probado que `docker compose -f docker-compose.registry.yml up -d` levanta el sistema completo
descargando las imágenes ya construidas, sin necesitar el código fuente de los backends/frontend.
