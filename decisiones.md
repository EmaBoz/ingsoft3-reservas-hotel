# Enlaces de entrega (TP6 y TP7)

> ⚠️ **QA y PROD corren en el FALLBACK LOCAL** (TP6 §3.6 / TP7 §3.6): los despliega mi self-hosted
> runner (`notebook-emanuel`) en mi notebook, así que sus URLs son de `localhost` y están vivas
> cuando el runner y Docker están prendidos (en la defensa, en mi máquina). Por qué el fallback:
> sección TP6 más abajo.

## Enlaces de este TP (TP6)

**Paquetes públicos** (ghcr.io, sin credenciales). Los cinco tienen la etiqueta del merge que llegó
a PROD, `sha-eb2a449ae1914863edf5db4158786968ba455ae7` (release `v6.0.0`):

- https://github.com/users/EmaBoz/packages/container/package/ingsoft3-reservas-hotel-users-api
- https://github.com/users/EmaBoz/packages/container/package/ingsoft3-reservas-hotel-rooms-api
- https://github.com/users/EmaBoz/packages/container/package/ingsoft3-reservas-hotel-reservations-api
- https://github.com/users/EmaBoz/packages/container/package/ingsoft3-reservas-hotel-search-api
- https://github.com/users/EmaBoz/packages/container/package/ingsoft3-reservas-hotel-frontend

```bash
docker pull ghcr.io/emaboz/ingsoft3-reservas-hotel-frontend:sha-eb2a449ae1914863edf5db4158786968ba455ae7
```

**La cadena** (por qué se puede confiar en el registry):

1. PR con los tests en verde y `Entrar al registry` **salteado** (sólo `main` publica):
   https://github.com/EmaBoz/ingsoft3-reservas-hotel/actions/runs/40130368155/job/120980601808
2. Corrida de `main` donde `Construir y publicar la imagen` es el **último** paso, después de los
   tests: https://github.com/EmaBoz/ingsoft3-reservas-hotel/actions/runs/40133943984/job/120991360531

**Entornos** (fallback local):

| | Front | users-api | rooms-api | reservations-api | search-api |
|---|---|---|---|---|---|
| **QA** | http://localhost:3100 | :8181 | :8182 | :8183 | :8184 |
| **PROD** | http://localhost:3200 | :8281 | :8282 | :8283 | :8284 |

---

# Decisiones — TP1 (Git colaborativo)

## 1. Por qué Git no pudo resolver el conflicto solo

Provoqué el conflicto a propósito siguiendo la guía (§4.6): dos ramas, `feature/titulo-a` y
`feature/titulo-b`, nacieron **las dos desde el mismo commit de `main`** y modificaron **la misma
línea** del `README.md` (la primera línea, el título) con contenidos distintos ("versión A" vs.
"versión B").

Git resuelve automáticamente cuando dos ramas tocan **partes distintas** de un archivo (usa un
merge de 3 vías comparando cada rama contra el ancestro común). Pero cuando las dos ramas cambian
**la misma línea**, Git no tiene forma de saber cuál de las dos versiones es "la correcta": las dos
son ediciones válidas del mismo lugar, y no hay una regla mecánica que decida por mí. Por eso me
delega la decisión marcando el archivo con `<<<<<<<`, `=======` y `>>>>>>>`, y me deja a mí elegir
qué queda.

**Qué habría tenido que pasar para que nunca apareciera:** que las ramas no tocaran la misma línea
(por ejemplo, si "versión B" hubiera modificado otra sección del README), o que una de las dos
ramas se hubiera integrado a `main` y la otra hubiera hecho `git pull`/rebase **antes** de tocar esa
misma línea, viendo el cambio de la primera. En un equipo real, esto se mitiga con ramas de vida
corta e integración frecuente (§3.4/3.5 de la guía): cuanto más tiempo vive una rama sin
sincronizarse con `main`, más probable y más grande es el conflicto.

## 2. Problemas encontrados y cómo los solucioné

- **`gh pr merge` bloqueado por el modo automático de la herramienta que usé (Claude Code).** El
  clasificador de seguridad de la herramienta frena por default cualquier acción que fusione código
  a una rama compartida/visible (como un merge de PR), porque son acciones difíciles de revertir.
  Lo resolví agregando una regla de permiso explícita para ese comando en la configuración local del
  proyecto (no del repositorio de la materia), después de decidir conscientemente que quería
  automatizar los merges en lugar de aprobarlos uno por uno desde la web.
- **CRLF/LF en Windows.** Git avisó (`warning: … LF will be replaced by CRLF`) al editar el README
  desde una terminal en Windows. No afectó el resultado (son solo finales de línea), pero lo dejo
  anotado porque es la clase de advertencia que en un equipo con Windows y Linux mezclados conviene
  resolver con un `.gitattributes` (no fue necesario para este TP).
- **Timing de `mergeable` en la API de GitHub.** Al consultar el estado de mergeabilidad del PR de
  la "versión B" justo después de mergear el de la "versión A", la API todavía devolvía
  `mergeable: UNKNOWN` (GitHub recalcula el estado de forma asíncrona). Esperando unos segundos y
  reconsultando, pasó a `CONFLICTING`, que es la evidencia real.

## 3. Declaración de uso de IA

Usé **Claude Code** (agente de IA, modelo Claude Sonnet 5) como herramienta de productividad para
ejecutar este TP, con supervisión y decisión mía en cada paso:

- **Qué hizo la IA:** ejecutó los comandos de `git` y `gh` (crear el repositorio, configurar la
  protección de rama vía la API de GitHub, crear y fusionar los Pull Requests, resolver el conflicto
  de merge, crear el tag y la release) siguiendo al pie de la letra la guía y el enunciado oficial
  del TP1 que yo le pasé (los leyó completos del repositorio público de la cátedra
  `ingsoft3ucc/TPs_2026`). Redactó también este archivo y `evidencias.md`.
- **Qué hice/decidí yo:** aprobé la creación del repositorio y su nombre, aprobé instalar `gh` CLI y
  Node.js, hice el login interactivo de `gh auth login` (necesariamente manual: requiere confirmar
  en el navegador con mi cuenta), decidí explícitamente habilitar el comando `gh pr merge` sin
  confirmación repetida, y elegí qué versión del título ("A" o "B") quedaba al resolver el
  conflicto.
- **Cómo lo verifiqué:** revisé el resultado de cada paso contra lo que pide el enunciado del TP1
  (protección de `main` probada con un push directo rechazado, dos PRs mergeados con uno de ellos
  resolviendo un conflicto real, tag y release visibles en GitHub) y puedo explicar cada comando
  ejecutado: qué hace `git switch -c`, por qué `enforce_admins: true` es necesario para que la
  protección me alcance a mí también, por qué el merge de la rama A no generó conflicto y el de B
  sí, y qué significan los marcadores `<<<<<<<`/`=======`/`>>>>>>>`.
- **Adaptación de la evidencia:** como la ejecución fue por terminal (CLI) y no clickeando en la web
  de GitHub, las capturas de pantalla del enunciado (📸) las reemplacé por la salida real de los
  comandos/API que muestran exactamente el mismo evento (el rechazo del push, el estado
  `CONFLICTING` del PR, el contenido del archivo con los marcadores). Lo declaro acá porque es una
  desviación consciente del formato sugerido por la guía, no un intento de esconder nada — todo el
  proceso es 100% verificable navegando el historial del repositorio.

---

## TP2 — Contenedores

### Qué app elegí y por qué

Elegí mi **Sistema de Reservas de Hotel** (proyecto propio, hecho para la materia Arquitectura de
Software II): backend en Go organizado como **4 microservicios** (`users-api`, `rooms-api`,
`reservations-api`, `search-api`), frontend en React + TypeScript + Vite, y persistencia en MySQL
(×2) + MongoDB, con Solr para búsqueda, Memcached para caché y RabbitMQ para eventos entre
servicios.

Contra los criterios de la guía (TP2 §3.3 / `elegir-app.md`):

- **¿Buildea y corre localmente hoy?** Sí — la probé completa con `docker compose up -d --build`
  antes de dar el TP por terminado (ver `evidencias.md`).
- **¿La entiendo lo suficiente para modificarla?** Sí: la escribí yo. Sé dónde está cada capa
  (`controllers` → `services` → `repositories` → `domain`) en cada microservicio.
- **¿Tiene lógica de negocio para testear (TP5)?** Sí: reglas de autorización por rol (JWT
  admin/usuario), restricciones de estado de habitaciones, validaciones de reservas — hay bastante
  más que un CRUD plano.
- **Tamaño:** acá me aparto conscientemente de la recomendación de la guía ("2-3 pantallas,
  alcanza con CRUD, sin dependencias exóticas"). Elegí **mantener la arquitectura de
  microservicios completa** (4 backends + Solr + Memcached + RabbitMQ + 2 motores de base de
  datos) en lugar de recortarla a un solo servicio, porque:
  1. Es la app que ya tengo construida y entiendo a fondo — la conozco lo suficiente como
     para modificarla y defenderla en la mesa, que es el criterio #5 de `elegir-app.md`.
  2. Cada pieza (colas, caché, búsqueda) es justamente el tipo de componente que esta materia
     enseña a operar en producción (TP6 en adelante), así que tenerlas desde ahora es más
     representativo del "sistema de entrega profesional" que busca la cursada.
  3. Asumo el costo que la guía advierte: más Dockerfiles, más piezas para explicar en la
     defensa, compilaciones más lentas. Lo documento acá para que quede claro que es una decisión
     tomada con el trade-off a la vista, no una casualidad.

### Decisiones de contenerización

- **Imágenes base:** `golang:1.21-alpine` para compilar cada microservicio (coincide con la
  versión de Go que ya usaban los `go.mod`) y `alpine:latest` para la etapa de runtime — mínima,
  sin SDK. Para el frontend: `node:22-alpine` para build y `nginx:alpine` para servir los
  estáticos.
- **Multi-stage en los 4 backends:** ya venían escritos así en el proyecto original (etapa
  `builder` con el SDK de Go, etapa final solo con el binario compilado + `ca-certificates`).
  Verifiqué que la estructura sigue el mismo principio que enseña la guía: SDK que compila,
  runtime que solo ejecuta.
- **Frontend con nginx + proxy interno:** el frontend es una SPA que llama a rutas relativas
  (`/login`, `/users`, `/api/v1/...`, `/api/reservations`, `/search-api/...`) sin host ni puerto
  hardcodeado — igual que recomienda la guía (TP2 §2.6, opción a). En desarrollo, quien traduce
  esas rutas es el proxy de Vite (`vite.config.ts`); en el contenedor, es `nginx.conf`, que
  reenvía cada prefijo al microservicio correspondiente **por nombre de servicio** (`users-api`,
  `rooms-api`, `reservations-api`, `search-api`) usando variables + `resolver 127.0.0.11` (no
  hardcodeado), tal como explica la guía, para que el contenedor del frontend no dependa de que
  los backends ya existan al arrancar.
- **A diferencia del ejemplo de la guía (un solo backend), acá hay CUATRO**, así que el
  `nginx.conf` tiene un `location` + `proxy_pass` por servicio en vez de uno solo. Es la misma
  idea, aplicada las veces que hacen falta.
- **Qué persiste y qué no:** los datos de MySQL (×2), MongoDB, RabbitMQ y Solr viven en
  **volúmenes nombrados** (`mysql_users_data`, `mysql_rooms_data`, `mongo_data`, `rabbitmq_data`,
  `solr_data`) — sobreviven a `docker compose down`. Memcached es intencionalmente efímero (es
  una caché: perder su contenido no pierde datos, solo obliga a recalcular). Lo verifiqué con la
  prueba de persistencia real (`evidencias.md`).
- **Secretos por variable de entorno:** todas las contraseñas (MySQL ×2, Mongo, RabbitMQ, JWT)
  viven en `.env` (ignorado por git) con `.env.example` commiteado como plantilla — el compose ya
  traía esto resuelto del proyecto original.
- **`depends_on` + `healthcheck`:** todos los servicios con base de datos/cola tienen
  `condition: service_healthy` en lugar de solo esperar a que el contenedor arranque —
  particularmente importante acá porque `rooms-api` depende de que `search-api` (que a su vez
  depende de Solr y RabbitMQ) esté *listo*, no solo *arrancado*.

### Problemas encontrados y cómo los resolví

1. **`npm ci` fallaba solo dentro de Docker, nunca en mi máquina.** El Dockerfile del frontend
   copiaba `package*.json` pero no `.npmrc` (que tiene `legacy-peer-deps=true`) antes de correr
   `npm ci`. Sin ese archivo, `npm` dentro del contenedor usaba resolución de dependencias
   "estricta" en vez de "legacy", y encontraba el lockfile (generado en modo legacy) inconsistente
   — el error mencionaba versiones de `picomatch` que ni siquiera estaban en el lockfile. Lo
   arreglé agregando `.npmrc` al `COPY` que precede a `npm ci`.
2. **`search-api` nunca pasaba a "healthy" y bloqueaba a `rooms-api`.** El Dockerfile de
   `search-api` instalaba `ca-certificates` pero no `curl`, y su healthcheck usa
   `curl -f http://localhost:8083/health`. Sin `curl`, el healthcheck fallaba siempre
   (`/bin/sh: curl: not found`), y como `rooms-api` depende de `search-api` con
   `condition: service_healthy`, todo el `docker compose up` se colgaba esperando algo que nunca
   iba a pasar. Lo arreglé agregando `curl` al `apk add` de la etapa runtime.
3. **El mapeo de puerto de `search-api` estaba mal en el compose original** (`"8083:8080"`, pero
   el proceso escucha en el puerto `8083` adentro del contenedor, no en `8080`). Esto no rompía la
   comunicación *entre* contenedores (que usa el puerto real, no el mapeo de host), pero sí
   hubiera roto el acceso desde afuera (`curl localhost:8083/...` desde la máquina). Lo corregí a
   `"8083:8083"`.
4. **El lockfile del frontend (`package-lock.json`) estaba desincronizado con `package.json`**
   (típico cuando se edita `package.json` a mano sin correr `npm install` después). Lo regeneré
   corriendo `npm install` una vez con Node 22 instalado localmente.
5. **Detalle que documento pero no corregí:** el frontend llama a
   `/search-api/api/search/rooms` (con el prefijo `/search-api` incluido en la URL, ver
   `src/services/SearchServices.ts`), pero la ruta real que expone `search-api` es
   `/api/search/rooms` (sin ese prefijo). En desarrollo, el proxy de Vite tampoco recorta el
   prefijo, así que esto ya fallaba antes de dockerizar. Mi `nginx.conf` reproduce el mismo
   comportamiento (no recorta el prefijo) para no cambiar el comportamiento de la app por mi
   cuenta — es un bug preexistente de la aplicación, no de la contenerización, y queda fuera del
   alcance de este TP (que es sobre Docker, no sobre corregir lógica de negocio del frontend).

### Declaración de uso de IA

Igual que en el TP1: usé **Claude Code** para ejecutar los comandos (crear archivos, correr
`docker build`/`docker compose up`, diagnosticar y corregir los errores de arriba, publicar las
imágenes) siguiendo el enunciado oficial del TP2. Decisiones que tomé yo: mantener la arquitectura
de microservicios completa en vez de recortarla, y qué versión de la app usar como base. Verifiqué
cada corrección leyendo el log del error (`curl: not found`, el mensaje de `npm ci`) y
confirmando que el fix lo resolvía de verdad, no solo que "dejaba de tirar error" — en los tres
casos reproduje el problema, entendí la causa raíz, y **después** de corregir volví a levantar el
sistema completo desde cero para confirmar que las tres cosas seguían funcionando juntas.

---

## TP3 — Planificación y trazabilidad

### Duración del sprint

Elegí **sprints de 1 semana**. Justificación: la cursada entrega un TP por semana (§2 del
reglamento: "la materia está diseñada para que cada TP tome una semana si venís al día"), así que
alinear el sprint con esa cadencia hace que "cerrar el sprint" y "cerrar el TP de la semana"
coincidan — el sprint deja de ser una unidad artificial superpuesta al calendario real de trabajo.

### Límite de trabajo en progreso

Configuré el límite en **2** en la columna *In Progress* del board. La regla de arranque de la
guía es "cantidad de personas + 1"; trabajando sola, eso da 1 + 1 = 2. El "+1" es la válvula para
cuando algo queda esperando una revisión o una respuesta y necesito avanzar en otra cosa sin que
el límite me trabe por completo, pero sigue siendo lo bastante bajo como para forzarme a terminar
antes de empezar algo nuevo. Señal de que está mal calibrado: si nunca lo alcanzo, está demasiado
alto; si lo alcanzo constantemente y me bloquea, está demasiado bajo para mi forma real de
trabajar.

### Diagnóstico de la historia mal escrita

Historia de ejemplo: *"Como desarrollador quiero crear la tabla usuarios para guardar los
datos."*

**Por qué está mal escrita:** es una **tarea técnica disfrazada de historia de usuario**. El
formato "Como... quiero... para..." está completo en la superficie, pero el contenido no cumple
lo que ese formato exige: el "quiero" describe una acción de implementación ("crear una tabla"),
no una capacidad observable por alguien; y el "para" no es un beneficio de negocio, es la
justificación técnica de la acción anterior ("guardar los datos" es *cómo* funciona la tabla, no
*qué* gana un usuario). Nadie "quiere" una tabla — la tabla es un medio, no un fin. Tampoco es
**Testeable** en el sentido de la guía (§2.3, INVEST): no hay un criterio de aceptación
verificable *desde afuera del sistema*, sólo "la tabla existe".

**Cómo la reescribiría:** subo un nivel, a la capacidad real que esa tabla hace posible, por
ejemplo *"Como usuario quiero registrarme con un usuario y una contraseña para poder iniciar
sesión más adelante"*, con criterios de aceptación como "el registro rechaza un email ya
existente" o "la contraseña se guarda hasheada, nunca en texto plano". "Crear la tabla usuarios"
pasa a ser una **tarea** técnica *dentro* de esa historia, no la historia en sí — que es
exactamente la distinción de niveles de la jerarquía (épica → historia → tarea, §2.2 de la guía).

### Problemas encontrados y cómo los resolví

- **La API pública de GitHub Projects no permite configurar la duración/fechas del campo
  Iteration ni el límite de trabajo en progreso del tablero Board.** Automaticé todo lo demás por
  API/CLI (proyecto, visibilidad, labels, issues, jerarquía con sub-issues, agregar los 5 items al
  proyecto, la vista Board), pero confirmé por introspección del schema de GraphQL que no existe
  ninguna mutación pública para esas dos configuraciones puntuales — son features exclusivas de la
  interfaz web. Las completé manualmente desde ahí, siguiendo exactamente los pasos de la guía.
- **`gh issue edit --add-sub-issue` requiere gh ≥ 2.94** (la guía lo advierte). La versión
  instalada (2.98.0) ya lo soporta, así que armé la jerarquía completa (épica → historia → 2
  tareas) sin pasar por la web.
- **El bug lo tomé de mi propia app**, no del ejemplo del video: durante las pruebas
  end-to-end del TP2 encontré que la búsqueda de habitaciones devuelve 404 (el frontend llama a
  `/search-api/api/search/rooms` pero `search-api` expone la ruta real en `/api/search/rooms`, sin
  ese prefijo). Es un bug real, sobre una versión ya entregada del sistema (TP2), así que encaja
  exactamente con la definición de la guía de "cuándo algo es un bug y no trabajo pendiente de una
  historia en curso" (§3.2).

### Declaración de uso de IA

Igual que en los TPs anteriores: usé Claude Code para ejecutar los comandos de `gh` y las
consultas/mutaciones de GraphQL (crear labels, issues, jerarquía, proyecto, vista Board, campo
Sprint) siguiendo el enunciado del TP3. Decisiones que tomé yo: la duración del sprint, el número
del límite de trabajo en progreso (y su justificación), qué bug reportar, y el diagnóstico de la
historia mal escrita. Yo misma configuré manualmente —vía la web, siguiendo instrucciones exactas—
las dos cosas que la API no expone (fechas del sprint y límite de la columna), y verifiqué por
API que quedaron aplicadas correctamente antes de seguir.

---

## TP4 — CI: Pipelines as Code

### Estructura elegida del pipeline

**Cinco jobs en paralelo, uno por Dockerfile** (`build-users-api`, `build-rooms-api`,
`build-reservations-api`, `build-search-api`, `build-frontend`), en vez de los dos
(`build-backend`/`build-frontend`) del ejemplo de la guía. La razón es la misma que ya quedó
documentada en el TP2: mi app es una arquitectura de microservicios con **cuatro** backends
independientes, cada uno con su propio Dockerfile — así que "un job por Dockerfile" son cinco
jobs, no dos. Cada job sigue exactamente el mismo patrón (checkout → buildx → build con cache
scoped al servicio), lo único que cambia es el `context` y el `scope` del cache. Corren en
paralelo porque no dependen entre sí: son builds independientes, cada uno en su propio runner.

El pipeline **no compila con `go build` ni `npm run build` directamente**: usa el `Dockerfile` de
cada servicio (`docker/build-push-action`), la misma definición de build del TP2. Esto evita tener
dos definiciones de "cómo se compila" que puedan divergir — lo que el pipeline verifica es
exactamente lo que después se despliega, no una compilación paralela e independiente.

### Cache

Cachea las **capas de la imagen Docker** de cada uno de los 5 servicios, con
`type=gha` (el cache de GitHub Actions) y un **`scope` distinto por servicio**
(`users-api`, `rooms-api`, `reservations-api`, `search-api`, `frontend`) para que no se pisen
entre sí. En la segunda corrida sobre el mismo PR, el log mostró las capas reutilizadas:

```
#11 CACHED
#12 CACHED
#13 CACHED
#14 CACHED
#15 CACHED
#16 CACHED
#17 CACHED
#18 CACHED
```

Si el cache desaparece (la plataforma lo puede desalojar en cualquier momento), el pipeline sigue
funcionando exactamente igual — sólo que cada capa se reconstruye desde cero, más lento. No hay
ninguna dependencia oculta en que el cache exista: lo comprobé corriendo el pipeline por primera
vez (sin ningún cache previo) y funcionó igual, sólo que sin ningún `CACHED` en el log.

### El pipeline como gate

`main` exige, además del Pull Request obligatorio del TP1: los **5 checks en verde**
(`required_status_checks`) y **`strict: true`** (la rama tiene que estar actualizada con `main`
antes de mergear).

**Demostración completa (PR #14):**
1. Rompí a propósito el build de `users-api` (import a un paquete inexistente) y abrí el PR.
2. El pipeline corrió: 4 de los 5 checks pasaron, `build-users-api` falló, y GitHub marcó el PR
   como bloqueado:
   ```json
   {"checks":[{"conclusion":"FAILURE","name":"build-users-api"},
              {"conclusion":"SUCCESS","name":"build-rooms-api"},
              {"conclusion":"SUCCESS","name":"build-reservations-api"},
              {"conclusion":"SUCCESS","name":"build-search-api"},
              {"conclusion":"SUCCESS","name":"build-frontend"}],
    "mergeStateStatus":"BLOCKED","mergeable":"MERGEABLE"}
   ```
3. Saqué el import roto (`fix: saca el import que no existe`) y agregué además un cambio real
   (un comentario de documentación), para que el PR tuviera contenido más allá de la rotura y el
   arreglo.
4. El pipeline volvió a correr, los 5 checks pasaron, `mergeStateStatus` pasó a `CLEAN`.
5. Antes de mergear, revisé el diff completo del PR (`gh pr diff`) — el gate verifica que
   compile, no que el cambio tenga sentido; eso lo reviso yo.
6. Mergeé el PR #14 (queda en el historial con sus corridas en rojo y en verde).

**Efecto de `strict: true` (PR #15, abierto en paralelo):** con un segundo PR abierto al mismo
tiempo, después de mergear el #14 el #15 pasó a:
```json
{"mergeStateStatus":"BEHIND","mergeable":"MERGEABLE"}
```
— es decir, aunque sus checks seguían en verde, GitHub no lo deja mergear porque quedó
desactualizado respecto de `main`. Actualicé la rama (`Update branch` / `PUT .../update-branch`),
el pipeline corrió de nuevo sobre la mezcla, y recién ahí lo mergeé.

### Problemas encontrados y cómo los resolví

- **El bug más importante de este TP no fue del pipeline: fue de un `.gitignore` de mi propia
  app, y lo descubrió el CI.** `services/search-api/.gitignore` tenía una línea `server` (sin
  barra inicial) para ignorar el binario compilado, pero esa regla también matchea **cualquier
  directorio** llamado `server` en cualquier profundidad — incluido `cmd/server/`, que es donde
  vive el código fuente real de `search-api` (`cmd/server/main.go`). Como resultado, ese archivo
  **nunca había llegado a git**: en mi máquina el build local (TP2) funcionaba porque `docker
  build` usa el directorio de trabajo tal cual está en el disco, con o sin git de por medio; pero
  el runner de GitHub Actions clona el repositorio desde cero (`actions/checkout`) y sólo trae lo
  que está *versionado*. El primer intento de correr el pipeline sobre `search-api` falló con
  `stat /app/cmd/server: directory not found` — un error que no tenía nada que ver con Docker ni
  con el workflow, y que sólo salió a la luz porque el CI, a diferencia de mi máquina, parte de
  una copia limpia. Es la prueba más concreta de por qué "anda en mi máquina" no alcanza. Lo
  arreglé anclando la regla a la raíz del servicio (`/server`, `/main`) y agregando el archivo
  recuperado al repositorio.
- **El `PUT` de `branches/main/protection` reescribe la protección entera.** Cuando agregué
  `required_status_checks`, tuve que volver a declarar también `required_pull_request_reviews` y
  `enforce_admins` (configurados en el TP1) en el mismo cuerpo del request — si los omitía, se
  perdían.

### Declaración de uso de IA

Igual que en los TPs anteriores: usé Claude Code para escribir el workflow, ejecutar los comandos
de `gh`/`docker`/`git`, diagnosticar el problema del `.gitignore` de `search-api`, y armar la
demostración completa del gate (romper, verificar bloqueo, arreglar, verificar desbloqueo,
mergear). Decisiones que tomé yo: la estructura de 5 jobs en paralelo (en vez de 2) y su
justificación, y verificar personalmente —leyendo cada log y cada respuesta de la API de GitHub,
no solo confiando en que "ya no tira error"— que el `.gitignore` roto era la causa raíz real y no
un síntoma de otra cosa, antes de dar el fix por bueno.

### Nota: corrección post-tag

Después de taguear `v4.0.0`, una relectura completa del repositorio encontró dos detalles menores
(una frase que había quedado en inglés en este archivo, y una referencia a una carpeta que nunca
existió en el README). Los corregí en un PR aparte (`chore: correcciones de la auditoría final`) y,
siguiendo la convención del TP1 (§3.7), **moví el tag `v4.0.0`** al commit corregido en vez de
dejarlo apuntando a la versión con el error — son correcciones de documentación, no cambios de
comportamiento del sistema.

---

## TP5 — Testing: unit tests, coverage y el umbral que frena un merge

### Qué testeé y por qué eso

- **Backend: `reservations-api`.** Es donde un bug duele de verdad: si se rompe la regla de
  superposición, dos personas reservan la misma habitación las mismas noches, y eso no se arregla
  con un reintento. Las reglas que cubren los tests (`services/reservation-services_test.go`):
  1. las fechas de entrada y salida son obligatorias;
  2. la salida tiene que ser posterior a la entrada (el borde es el mismo día);
  3. no se puede reservar una habitación que ya está reservada en esas fechas (devuelve conflicto y
     **no guarda**);
  4. para cancelar hacen falta id y motivo, y sólo se cancela una reserva que existe.

  Además, los handlers (`controllers/reservation-controller_test.go`): que cada error del servicio
  salga con su código HTTP (400 si el body está incompleto, 409 si la habitación está ocupada, 404
  si la reserva no existe, 201 si se crea). En total son 16 funciones de test; dos son
  parametrizadas (tabla de casos + `t.Run`, el equivalente en Go del `[Theory]`).
- **Los otros tres microservicios** (`users-api`, `rooms-api`, `search-api`) por ahora sólo tienen
  el gate del TP4 (que compilen). Es lo siguiente que habría que sumar, empezando por `users-api`
  (login y hash de contraseñas).
- **Frontend.** No había lógica testeable: el cálculo de noches, el total y la validación de fechas
  estaban adentro del componente `ReserveRoom.tsx`. Los saqué a `src/lib/reservas.ts` como
  funciones puras, y el componente ahora sólo las llama. Los tests están en
  `src/lib/reservas.test.ts`: `it.each` para el total, casos de error de `validarFechas` (con el
  mensaje que ve el usuario) y el mock del cliente HTTP.

Para saber si los tests verifican algo, invertí cada regla a mano y miré que algo se pusiera rojo:
`<=` por `<` en la fecha de salida (falla el caso de borde "salida el mismo día"), sacar el chequeo
de superposición, `||` por `&&` en las fechas obligatorias, sacar el motivo obligatorio, y en el
front `<= 0` por `< 0` en `validarFechas`. Todos quedaron en rojo.

### El mock, y el refactor que hizo falta

- **Backend:** no hizo falta refactorizar. `reservationService` ya recibía el repositorio y el
  publisher por constructor, como interfaces (`ReservationRepository`, `EventPublisher`). El doble
  es un struct escrito a mano (`repoMock`) que implementa la interfaz, devuelve lo que el test le
  configura y **anota las llamadas**. `TestCrearReserva_HabitacionOcupada_DevuelveConflictoYNoGuarda`
  es un mock y no un stub porque el assert mira la interacción: que se consultó la disponibilidad de
  la habitación 101 una sola vez y que `Create` **no** se llamó. El publisher es un stub que no hace
  nada (el evento sale en una goroutine y no es lo que verifico).
- **Frontend:** sí refactoricé. La carga de "Mis reservas" llamaba a axios adentro del componente,
  así que no había por dónde meter un doble. Ahora `misReservas(userId, traer)` recibe el cliente
  por parámetro: la pantalla le pasa el real (`traerJson`, en `reservationService.ts`) y el test le
  pasa un `vi.fn()`. Un test lo usa como stub (contesta una lista) y otro como mock
  (`toHaveBeenCalledWith` con la ruta del usuario).

### Herramientas (mi stack es Go + React/TS, no el de la cátedra)

| Lo que pide la tabla | Qué usé |
|---|---|
| Test parametrizado | Go: tabla de casos + `t.Run` · TS: `it.each` |
| El doble | Go: struct a mano que implementa la interfaz · TS: `vi.fn()` |
| Medir la cobertura | Go: `go test -coverpkg=./... -coverprofile` · TS: `@vitest/coverage-v8` |
| Umbral que rompe el build | Go: `coverage.sh` compara el total de `go tool cover -func` con el umbral y sale con 1 · TS: `coverage.thresholds` de vitest |
| Qué entra en la cuenta | Go: `grep -v` sobre el perfil (por exclusión) · TS: `include` / `exclude` en `vite.config.ts` |

Los dos corren adentro de la etapa `test` de su Dockerfile (`FROM builder AS test` /
`FROM build AS test`): el pipeline construye esa etapa y la corre con `docker run`, así no hay una
segunda receta de build.

### El umbral

- **Backend: 60 % de sentencias.** Mi suite hoy mide **66,7 %**. Puse 60 y no 66 por algo que vi
  midiendo: `CreateReservation` y `DeleteReservation` publican el evento en una goroutine, y si el
  test termina antes de que corra, esas líneas no cuentan. En una corrida me dio 59 % y en las
  siguientes 62,9 % con el mismo código; con el umbral pegado al número tendría un gate flaky. Para
  subirlo tendría que testear los métodos que sólo pasan el pedido al repositorio (`GetAll`,
  `GetmyReservations`) y las ramas de error 500 de los handlers. **Go no mide ramas**: `-cover` es
  cobertura de sentencias, así que del backend no tengo número de rama.
- **Frontend: 80 % de líneas y 80 % de ramas.** Hoy mide **100 % de líneas y 100 % de ramas**
  sobre `src/lib`. No puse 100 porque me trabaría por cualquier rama defensiva; con 80, el código
  nuevo tiene que entrar testeado (lo comprobé en el PR #21). Puse las dos métricas porque con
  vitest 3 (el mío es 3.2.7) una función que nadie llama no suma ramas, sólo líneas.

Resumen de coverage en la corrida, con los dos reportes descargables (`coverage-reservations-api` y
`coverage-frontend`): https://github.com/EmaBoz/ingsoft3-reservas-hotel/actions/runs/40129216437

### Qué dejé afuera de la cuenta

- **Backend** (`coverage.sh`): `cmd/server` y `config/` (el arranque: arma las conexiones y las
  rutas, sin reglas), `domain/` (structs de datos y DTOs, sin comportamiento), y `repositories/` y
  `events/`. Estos dos últimos sí tienen código, pero son el borde con Mongo y RabbitMQ: un unit
  test con doble no puede verificar si la consulta de superposición está bien escrita, porque
  reemplaza justamente a la base. Eso lo tiene que ver una prueba contra la base de verdad (TP7).
  Lo dejo dicho como hueco conocido, no como "no hace falta testearlo". Se sacan con `grep -v` sobre
  el perfil, o sea por exclusión: un paquete nuevo entra solo a la cuenta.
- **Frontend** (`vite.config.ts`): entra `src/lib/**` menos `src/lib/axios.ts`, que es la
  configuración del cliente HTTP (base URL y header del token), no una regla mía. Los componentes
  React no entran: los tests de componentes son opcionales en la guía, y la pantalla completa la
  verifica la e2e del TP7.

### Por qué un coverage alto no es calidad (con mi ejemplo)

El frontend da 100 %. Si borro todos los `expect` de `reservas.test.ts`, sigue dando 100 %: las
funciones se ejecutan igual, sólo que nadie verifica qué devuelven. En el backend pasa lo mismo:
un test que llame a `CreateReservation` con una reserva válida y no mire el resultado deja la
función casi al 100 % sin comprobar que la reserva queda `active`. Si alguien cambia el `<=` de
`dto.EndDate <= dto.StartDate` por `<`, lo atrapa el caso de borde "salida el mismo día que la
entrada", no el número de cobertura.

### Los dos Pull Requests del freno

1. **El que cuenta la historia (mergeado): https://github.com/EmaBoz/ingsoft3-reservas-hotel/pull/21**
   Agregué `etiquetaDeReserva` (Próxima / En curso / Finalizada / Cancelada en "Mis reservas") sin
   tests. Los 12 tests pasaban y compilaba, pero el check `build-frontend` quedó rojo **por
   líneas**: `ERROR: Coverage for lines (55.26%) does not meet global threshold (80%)`. Las ramas
   seguían en 100 % porque vitest 3 no cuenta las ramas de una función que nunca se ejecuta.
   Corrida roja: https://github.com/EmaBoz/ingsoft3-reservas-hotel/actions/runs/40129453727
   Lo arreglé con un test por cada camino que declara la función (cancelada, sin fechas, próxima,
   en curso, finalizada), con `hoy` fijo por parámetro para no depender del reloj. Volvió a 100 %,
   el check pasó a verde y lo mergeé.
2. **El que queda abierto y en rojo: https://github.com/EmaBoz/ingsoft3-reservas-hotel/pull/22**
   Agrega `ValidarEstadia` en `reservations-api` sin tests. Compila y todos los tests pasan, pero
   `build-reservations-api` falla con `ERROR: la cobertura de sentencias (59.8%) no llega al umbral
   (60%)` y el merge queda bloqueado.
   Corrida roja: https://github.com/EmaBoz/ingsoft3-reservas-hotel/actions/runs/40129742111

La diferencia con el freno del TP4: allá frenaba algo que no compilaba. Acá todo compila y anda, y
lo que frena es un criterio de calidad que elegí yo. Lo que este freno deja pasar igual: un test
que ejecuta el código sin un assert que sirva, o una regla mal entendida desde el requisito (el test
congela lo que yo entendí, no lo que se pidió).

### El camino sin cubrir

La primera vez que medí el frontend, el reporte de `reservas.ts` daba 100 % de líneas y 90,9 % de
ramas, con la línea 9 marcada:

```ts
if (!entrada || !salida) return 0;   // en calcularNoches
```

- **Qué entrada la recorre:** una fecha vacía, por ejemplo `calcularTotal('', '2026-11-12', 100)`.
  Pasa de verdad: el input de fecha deja borrar el valor, y sin ese `if` `new Date('')` es una fecha
  inválida y la pantalla mostraría `$NaN` como total.
- **Qué decidí:** agregar el test (`da 0 y no NaN si el usuario borró la fecha de entrada`), porque
  protege algo que ve el usuario. Con eso las ramas pasaron a 100 %.

### Problemas encontrados

- **El número del backend variaba entre corridas** (59 % y 62,9 % con el mismo código) por las
  goroutines que publican eventos. Lo vi corriendo los tests varias veces seguidas con `-count=1`,
  y por eso el umbral quedó con margen (60 contra 66,7).
- **`go test` imprime otro porcentaje** ("29.5% of statements in ./..."): es la cobertura sobre todo
  el módulo, antes de sacar lo excluido. El número que vale es el que imprime `coverage.sh` después
  del filtro, que es el mismo que va al resumen de la corrida.
- **La clave `test:` de `vite.config.ts` no compilaba** en `tsc -b` hasta que agregué
  `/// <reference types="vitest/config" />` arriba del archivo.
- **Un commit local viejo ("prueba") que nunca había subido** apareció al hacer el primer `pull`
  después de mergear y armó un merge raro en mi `main` local. No llegó a GitHub: lo saqué de `main`
  y lo dejé guardado en una rama local de respaldo.

### Declaración de uso de IA

Usé Claude Code (Claude Opus 5.5) para escribir los tests, el `coverage.sh`, los pasos del
`ci.yml` y este texto, a partir del enunciado del TP5. Decisiones mías: testear `reservations-api`
y no los cuatro servicios, el umbral de cada lado y qué dejar afuera de la cuenta. Cómo lo
verifiqué: corrí las etapas de test de los dos Dockerfiles en mi máquina antes de subir; invertí
cada regla a mano y comprobé que algún test se pusiera rojo (lista arriba); y en el pipeline
confirmé que el rojo del PR #21 era por cobertura y no por un test fallando (el log dice
`12 passed` arriba del `ERROR`). Lo que la suite no cubre: los repositorios y el publisher (Mongo y
RabbitMQ), las ramas de error 500 de los handlers y los otros tres microservicios.

---

## TP6 — CD: environments, aprobaciones y deployment patterns

> Los enlaces de este TP (paquetes, la cadena y las URLs de QA y PROD) están arriba de todo.

### Por qué el fallback local, y cómo cumple el contrato

Mi app no es un back y un front: son 4 microservicios Go, MySQL ×2, MongoDB, RabbitMQ, Solr y
Memcached. Render free no aloja nada de eso salvo los contenedores web, y `reservations-api` no
arranca sin Mongo y RabbitMQ. Para ir a la nube tendría que haber migrado dos servicios de MySQL a
Postgres, Mongo a otro proveedor, RabbitMQ a otro, y dejar Solr afuera, con unos diez servicios
compartiendo las 750 horas. Por eso usé el plan C de la guía (§3.6): mi notebook como self-hosted
runner y un proyecto de Compose por entorno. Los cinco puntos del contrato:

1. **Dos entornos separados**: `compose.qa.yml` (proyecto `qa`) y `compose.prod.yml` (proyecto
   `prod`), corriendo a la vez con puertos distintos (QA 3100 y 8181-8184, PROD 3200 y 8281-8284).
2. **Una base por entorno**: cada proyecto tiene sus propios volúmenes (`_qa` / `_prod`) y su propia
   red. Lo comprobé creando un usuario `soyprod` sólo en PROD (`POST localhost:3200/users`): aparece
   en `localhost:3200/users` y en `localhost:3100/users` no.
3. **Front y back como contenedores**, desde mis Dockerfiles.
4. **El deploy lo dispara el pipeline**: `deploy-qa` y `deploy-prod` corren `docker compose up` en el
   runner. Yo no levanto nada a mano. Antes de mergear probé QA una vez en mi máquina y lo bajé con
   `down -v` justamente para que lo levante la corrida.
5. **PROD detrás de la aprobación** del environment `production`.

Lo que pierdo respecto de la nube: no hay URL pública (sólo `localhost`), el entorno depende de que
mi notebook esté prendida, con Docker y el runner corriendo, y no hay cold start real que manejar.

### Por qué la imagen se publica sólo en verde

Son tres eslabones encadenados: nada entra a `main` sin los checks en verde (protección de rama),
sólo `main` publica (`push` del paso atado a `push` + `refs/heads/main`), y en cada job publicar es
el **último** paso, después de los tests: si `docker run` de los tests sale con error, el job muere
ahí y nunca llega a publicar. No hay un `if` que diga "si pasaron los tests": la condición es el
orden. Si se publicara igual, el registry dejaría de significar "esto pasó la verificación" y una
etiqueta `sha-…` podría ser una imagen con tests en rojo. Lo que la cadena **no** garantiza: alguien
con permiso de escritura puede hacer un `docker push` a mano con una etiqueta `sha-…`; lo único
que nadie puede mover es el digest (`@sha256:…`).

### Continuous Delivery, no Deployment

Implementé **Continuous Delivery**: cada merge llega solo a QA, y a PROD sólo con una aprobación
humana. Para Continuous Deployment me faltaría sacar el gate, y no lo querría hoy: el smoke sólo
dice que la app responde, todavía no tengo pruebas que digan que funciona (llegan en el TP7) ni
observabilidad para darme cuenta si PROD se rompe (TP9).

### La cadena y los secrets

```
build-* (5, sin needs) → deploy-qa (needs: los 5 builds · if: main · environment: qa)
                       → deploy-prod (needs: deploy-qa · environment: production · concurrency)
```

- `needs` arma la cadena con compuertas: sin los 5 builds en verde no hay deploy a QA, y sin QA vivo
  no hay PROD. `deploy-prod` no repite el `if` de la rama porque lo hereda: depende de `deploy-qa`,
  que en un PR sale salteado.
- **Secrets**: `DB_PASSWORD` y `JWT_SECRET` son **environment secrets**, distintos en `qa` y en
  `production`. Los de PROD sólo los puede leer un job que pasó la aprobación. En el repo no hay
  ningún secret: para publicar en ghcr se usa el `GITHUB_TOKEN` de cada corrida, con
  `packages: write` declarado sólo en los jobs de build.
- **El commit desplegado es el verificado**: el `checkout` del job de deploy baja el `github.sha` de
  esa corrida (el equivalente del `&ref=` del hook de Render), no la punta de `main`.

### Qué miro antes de aprobar

1. Que `deploy-qa` esté en verde y abrir QA (`localhost:3100`) para ver el cambio de ese PR.
2. Que la corrida que apruebo sea la del **último** merge: si hay dos esperando, apruebo la más
   nueva y rechazo la vieja (aprobar la vieja haría retroceder PROD).
3. Qué cambió en el PR: si toca un compose, secrets o puertos, lo miro con más tiempo.

Lo que mi aprobador **no** puede ver: si la app funciona más allá de lo que yo pruebe a mano en un
minuto, y qué le pasa a un usuario en PROD después del deploy. Para eso faltan las e2e (TP7) y
monitoreo (TP9).

- Rechazo con motivo: https://github.com/EmaBoz/ingsoft3-reservas-hotel/actions/runs/40132680314
  ("Primer deploy a PROD: antes de aprobar verifico que compose.prod.yml no comparta puertos ni
  volúmenes con QA.").
- Aprobación y PROD con el smoke en verde:
  https://github.com/EmaBoz/ingsoft3-reservas-hotel/actions/runs/40133943984 (el cambio visible es
  el título del inicio, "Bienvenido al Hotel Grand Prix").

### La "letra chica" de mi infraestructura

No uso free tier de nube, pero el fallback tiene su propio contrato:

- **El runner tiene que estar corriendo** (`C:\actions-runner\iniciar-runner.cmd`). Si está
  apagado, los jobs de deploy quedan en cola esperando un runner; no fallan.
- **Docker Desktop prendido.** Los dos entornos juntos son 22 contenedores.
- **Minutos de Actions**: los builds corren en runners de GitHub (gratis en repo público); los
  deploys corren en mi máquina y no consumen minutos.
- **Seguridad del runner**: GitHub desaconseja runners propios en repos públicos porque un PR de un
  fork podría ejecutar código en mi máquina. Puse *Require approval for all external contributors*
  y el runner no está como servicio: lo prendo cuando lo uso.
- El smoke reintenta igual (30 × 20 s): `docker compose up` termina antes de que las apps estén
  listas (MySQL, RabbitMQ y Solr tardan en arrancar).

### Qué garantía pierdo porque el entorno reconstruye

En este TP el deploy hace `docker compose up --build`: el runner **vuelve a construir** las imágenes
desde el commit, en vez de usar las que publicó el pipeline. Despliego el commit verificado, pero no
la imagen verificada: una imagen base que cambió (`alpine:latest`, `nginx:alpine`) o una dependencia
que se resolvió distinto pueden hacer que lo que corre en QA y en PROD no sea bit a bit lo que pasó
los tests, y PROD se construye otra vez, aparte de QA. Es lo que arregla el TP7.

### Qué prueba el smoke y qué no

Prueba que `rooms-api` está viva (`/health`), que a través del front (o sea, nginx con la plantilla)
las rutas de rooms (MySQL), reservas (Mongo) y usuarios (MySQL) contestan 200, y que el front sirve
la página. **No** prueba qué versión está corriendo: si el `up` no hubiera reemplazado algún
contenedor, el smoke daría verde contra la versión vieja. Tampoco prueba que la app funcione (que se
pueda reservar): sólo que responde.

### Qué quedó por variable y qué adentro de la imagen del front

Adentro de la imagen: los estáticos de Vite y la plantilla `default.conf.template`. Por variable, al
arrancar: `USERS_API_URL`, `ROOMS_API_URL`, `RESERVATIONS_API_URL`, `SEARCH_API_URL` y
`DNS_RESOLVER`. En el fallback los valores son iguales en QA y en PROD porque cada proyecto de
Compose tiene su propia red (el `rooms-api` de QA no es el de PROD), pero la plantilla es lo que
permite que **la misma imagen** sirva en cualquier entorno, que es lo que necesita el TP7.

### Deployment pattern para una producción real

Usaría **blue-green**. Mis servicios no guardan estado (el estado está en las bases), así que puedo
levantar la versión nueva completa al lado de la vieja, probarla, y cambiar el tráfico de una vez
desde el proxy. El rollback es volver a apuntar a la vieja, que sigue levantada: segundos, sin
reconstruir nada. El costo es tener dos juegos de contenedores durante el cambio, y la condición es
que las bases sean compatibles con las dos versiones a la vez (cambios de esquema en dos pasos).
Canary me parece mejor para una app con mucho tráfico, pero hoy no tengo con qué decidir si el 5 %
va bien: no tengo métricas de errores ni de latencia por versión (eso falta hasta el TP9).

### Plan de rollback (cronometrado)

Como el deploy siempre baja un commit explícito, volver atrás es **re-correr el job `deploy-prod`
de la corrida que desplegó la versión buena** (la de `v6.0.0`): baja ese commit y levanta PROD con
él. Pasa otra vez por la aprobación, porque es un deploy a producción.

Lo hice de verdad: con PROD en `818f88c` (el PR #28, "¿Por qué elegir el Hotel Grand Prix?"), volví
a `v6.0.0` (`eb2a449`):

- Disparo (re-run del job): 21:35:25 UTC
- Aprobación y arranque del job: 21:38:16 UTC
- Smoke de PROD en verde: 21:38:33 UTC → **3 min 8 s en total**, de los cuales el deploy en sí son
  **17 s** (el resto es la espera de la aprobación).
- Job: https://github.com/EmaBoz/ingsoft3-reservas-hotel/actions/runs/40133943984/job/120994216688
  (intento 2 de esa corrida). En *Deployments* → `production` se ve el deploy de `eb2a449` a esa
  hora.

Fue rápido porque el runner ya tenía en caché las capas de ese commit; sin caché, reconstruir los
cinco servicios lleva un par de minutos. Lo que el rollback de código **no** deshace son los datos:
los volúmenes siguen como estaban (el usuario `soyprod` sigue ahí), y si una versión nueva agregó
una columna con `AutoMigrate`, esa columna queda.

### Problemas encontrados

- **`denied: permission_denied: write_package` en la primera corrida de `main`.** Los paquetes los
  había subido a mano en el TP2, así que no estaban vinculados al repo y el `GITHUB_TOKEN` del
  workflow no podía escribir. Lo resolví dándole al repo rol *Write* en *Manage Actions access* de
  cada paquete (la primera vez quedó en *Read*, que es el default) y re-corriendo los jobs.
- **El runner usaba el `bash` de WSL** (`C:\Windows\system32\bash.exe`) en los pasos con
  `shell: bash`, y fallaba con `execvpe(/bin/bash) failed`. Lo arranco con un `iniciar-runner.cmd`
  que pone `C:\Program Files\Git\bin` primero en el PATH.
- **Choque de puertos con el compose de desarrollo**, que tiene `restart: always` y queda levantado
  en 8080-8083. Por eso QA y PROD usan otros puertos.
- **Contraseña de prueba en los volúmenes**: probé QA en mi máquina con una contraseña cualquiera, y
  MySQL guarda la contraseña al inicializar el volumen. Antes del primer deploy lo bajé con
  `down -v` para que la corrida arranque las bases con la del secret.

### Declaración de uso de IA

Usé Claude Code (Claude Opus 5.5) para escribir los compose de QA y PROD, la plantilla de nginx, los
jobs de deploy y este texto, siguiendo la guía del TP6. También instaló y registró el runner y cargó
los secrets con valores al azar, sin mostrarlos. Decisiones mías: el fallback local en vez de
Render, el motivo del rechazo y las dos aprobaciones (las hice yo desde la web, después de mirar
QA). Cómo lo verifiqué: comprobé el `docker pull` de los cinco paquetes sin sesión, miré en el log
que los deploys corrieron en `notebook-emanuel` con el bash de Git, abrí QA y PROD en el navegador,
y comprobé que un dato creado en PROD no aparece en QA. Entre el merge y PROD pasa esto: los 5
builds (tests y publicación) → `deploy-qa` en mi runner con los secrets de `qa` y su smoke → pausa
del environment `production` hasta que apruebo → `deploy-prod` con los secrets de `production` y su
smoke.
