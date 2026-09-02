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
