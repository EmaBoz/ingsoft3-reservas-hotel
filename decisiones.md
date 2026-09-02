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
