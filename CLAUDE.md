# Software Metrics & Estimation — seguimiento-medicion-bcg

> Raíz **delgado** del monorepo. Claude Code carga este archivo + el `CLAUDE.md` del directorio
> donde trabaja. Acá va lo transversal (qué es, stack, reglas duras, índice); el detalle de cada
> unidad vive en su propio `CLAUDE.md`. El **estado de los Sprints vive en GitHub** (Issues +
> Milestones + Projects), no acá.

## Qué es

Aplicación web para administrar proyectos de software y medir su **estimación, planificación,
seguimiento y calidad**. Trabajo Práctico Integrador de Ingeniería y Calidad de Software 2026 —
UTN Facultad Regional San Rafael.

Flujo end-to-end: crear proyecto y registrar integrantes → cargar el Product Backlog → estimar
historias con Planning Poker (Story Points) → planificar el Sprint (Sprint Goal + historias) →
registrar esfuerzo (horas) y defectos durante el Sprint → cerrar el Sprint → consultar métricas,
Dashboard y reporte del proyecto/Sprint (exportable a PDF).

Metodología: **Scrum** + **SDD** (specs antes de implementar) + **BDD** (Given–When–Then) +
**TDD** (RED → GREEN → REFACTOR). Sprints en **GitHub Milestones**, historias como **Issues**,
tablero en **GitHub Projects**.

### Equipo y roles Scrum

| Rol | Quién |
|---|---|
| Product Architect | Profesores de la cátedra |
| Agile Enabler | Martina Batista |
| Product Builders | Martina Batista, Florencia Cabeza, Luana Guajardo |

## Stack

Go (versión fijada en `backend/go.mod`; API REST, reglas de negocio) · PostgreSQL · React + Vite +
TypeScript (frontend) · godog (BDD con Gherkin) · `go test` (TDD) · GitHub Actions (CI).

**Regla de la cátedra:** el núcleo de la solución y **todas** las reglas de negocio y cálculos
(métricas, estimación, validaciones) viven en Go. El frontend solo presenta y captura datos.

## Reglas duras (no negociables)

- **Git:** nunca commit directo a `main`. Toda rama sale de `main` actualizado, por tipo (`feat/`,
  `fix/`, `docs/`, `test/`, `refactor/`, `chore/`) → PR a `main` con CI verde y revisión de otra
  integrante. **Conventional Commits con scope = unidad** (`backend`, `frontend`, `docs`, `ci`).
- **Dos ramas por historia:** 1) `docs/spec-NN-slug` con la spec → PR que otra integrante aprueba;
  2) recién mergeada la spec, `feat/slug` con BDD + TDD + pantallas → PR con `Closes #N`.
  Tareas sin spec (fix, chore, docs sueltas) usan una sola rama.
- **SDD:** funcionalidad no trivial → spec **aprobada antes de codear** (`docs/specs/`). Las specs
  describen qué y para qué, **sin código**. Se copian de `docs/specs/_template.md` y se guardan en
  `docs/specs/sprint-NN/NN-slug.md` (carpeta = Sprint del Milestone; `NN` = número de la HU).
  Estados: `draft` (se escribe) → `in-review` (PR abierto) → `approved` (la aprobó **otra
  integrante**, recién ahí se programa) → `implemented` (su PR de código se mergeó a `main`).
- **BDD:** cada funcionalidad seleccionada tiene escenarios Given–When–Then (normales,
  alternativos, límite y error), automatizados con godog siempre que sea posible.
- **TDD con evidencia en el historial:** reglas de negocio, cálculos y validaciones se desarrollan
  RED → GREEN → REFACTOR, y **cada fase es un commit separado**:
  `test(backend): ...` (rojo) → `feat(backend): ...` (verde) → `refactor(backend): ...`.
  No mezclar test e implementación en el mismo commit.
- **Trazabilidad:** Historia (Issue) → Spec SDD → Criterios de aceptación → Escenarios BDD →
  Tests → Código Go. La spec referencia el Issue, el PR cierra el Issue (`Closes #N`).
- **IA:** todo lo generado con IA se entiende, revisa y valida antes de mergear.
- **Cero hardcoding, no parches, resolver de raíz.** Errores manejados explícitamente.
- **CI verde antes del PR. Sin `--no-verify`** (lo fuerza un hook).

## Mapa de la documentación

- [`docs/specs/_template.md`](docs/specs/_template.md) — plantilla de spec (campos que pide la cátedra).
- [`.github/pull_request_template.md`](.github/pull_request_template.md) — checklist de cada PR.
- `.documentacion-clase/` — consigna de la cátedra (local, **no se versiona**).

## Layout

```
.
├── CLAUDE.md
├── backend/          # Go: dominio, reglas de negocio, métricas, API REST
│   └── CLAUDE.md
├── docs/specs/       # specs SDD, una carpeta por Sprint (sprint-NN/)
├── .github/          # CI + PR template
└── .claude/          # hooks de Claude Code y el comando /comando (ponerse al día)
```

El frontend (React + Vite + TS) todavía no está creado: se arma cuando el equipo lo defina.

Regla de capas: el **dominio Go** (entidades, reglas, cálculos) no depende de HTTP ni de la base de
datos; la API (`backend/api`) y el acceso a datos (`backend/postgres`) lo envuelven. El frontend consume solo la API REST.

## Índice de módulos

<!-- Uno por unidad con CLAUDE.md propio. Crece al mergear specs. -->

- [`backend/`](backend/CLAUDE.md) — núcleo Go + API REST.
