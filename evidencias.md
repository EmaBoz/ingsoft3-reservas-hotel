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
