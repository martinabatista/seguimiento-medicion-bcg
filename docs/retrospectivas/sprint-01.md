# Sprint 1 — El MVP

**Equipo:** Martina Batista (Agile Enabler), Florencia Cabeza, Luana Guajardo.
**Período de trabajo:** del 20/10/2026 al 23/10/2026 (Milestone [Sprint 1 - El MVP](https://github.com/martinabatista/seguimiento-medicion-bcg/milestone/2)).
**Modalidad:** las ceremonias se hicieron en reuniones presenciales del equipo, trabajando juntas sobre el proyecto.

## Sprint Planning

**Sprint Goal:** tener la base del sistema funcionando: proyectos con sus integrantes, un Product Backlog estimado con Story Points y Sprints que se planifican y se cierran.

**Sprint Backlog:**

| Historia | Dueña | Resultado |
|---|---|---|
| [HU-01](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/4) Crear y modificar proyectos | Martina | Terminada (#33, #34) |
| [HU-02](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/5) Registrar integrantes | Luana | Terminada (#38, #39) |
| [HU-03](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/6) Consultar el estado de un proyecto | Luana | Casi terminada (#40, #41): falta sacar el CA-2, que pasó a la HU-05 |
| [HU-04](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/7) Administrar el Product Backlog | Florencia | Terminada (#35, #37) |
| [HU-05](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/8) Crear Sprints y asignar historias | Martina | Terminada (#43, #44) |
| [HU-06](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/9) Cerrar Sprint y consultar anteriores | Martina | Terminada (#45, #46) |
| [HU-07](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/10) Estimar con Story Points | Florencia | Terminada (#47, #48) |

Cada historia siguió el flujo del equipo: spec aprobada por otra integrante → código con TDD y BDD → PR con CI en verde.

## Dailies

Al trabajar juntas, cada una contaba qué había terminado, qué seguía y qué la bloqueaba. Los bloqueos que aparecieron y cómo se resolvieron:

- **HU-03 esperaba a la HU-05** (el "Sprint en curso"): se movió ese criterio a la HU-05 y se agregó la regla de historias independientes al `CLAUDE.md` (#42).
- **PRs atrasados respecto de `main`** cuando se mergeaba otro antes: se actualizaron antes de mergear.
- **Tablas nuevas** en cada historia: cada una recreó su base local después de los merges.

## Sprint Review

**Incremento entregado:** la API del backend en Go permite crear y modificar proyectos, registrar integrantes, administrar el Product Backlog, estimar historias con Story Points (escala Fibonacci), crear Sprints con su Sprint Goal, asignarles historias, completarlas, cerrar el Sprint y consultar los anteriores, y ver el estado del proyecto con el Sprint en curso. También se adelantó la HU-08 (Planning Poker) del Sprint 2 (#49, #50).

**Calidad al cierre del Sprint:**

| Control | Resultado |
|---|---|
| Escenarios BDD (godog) | 217 escenarios y 1086 pasos en verde |
| Cobertura de tests | dominio 94%, API 91%, PostgreSQL 81% |
| CI en GitHub Actions | en verde en todos los PRs |
| TDD | commits `test` → `feat` (→ `refactor`) en el historial de cada historia |

**Pendiente para el Sprint 2:** cerrar la HU-03 (sacar su CA-2 de la spec 03).

## Retrospectiva

**Qué salió bien:**

- El trabajo en paralelo por áreas: cada una con sus historias de punta a punta, y se terminaron 6 de las 7 historias del Sprint.
- Escribir y aprobar las specs antes de programar (SDD).
- Los tests y la CI: TDD con evidencia en los commits y la CI frenando errores antes del merge.
- Las reglas del equipo claras (`CLAUDE.md`, skills y `/comando`): las tres trabajan igual con la IA.

**Qué costó:**

- Las dependencias entre historias: la HU-03 quedó a medias esperando a la HU-05.
- El tablero desactualizado: varias tarjetas no se movieron mientras se trabajaba.
- Recrear la base de datos cada vez que una historia agrega tablas.

**Qué vamos a mejorar en el Sprint 2:**

| Acción | Responsable |
|---|---|
| Mover la tarjeta de cada historia al trabajarla (En progreso → En revisión) | Cada dueña; la Agile Enabler lo controla |
| Revisar los PRs el mismo día en que se abren | Las tres |
| Avisar en el grupo cuando un PR agrega o cambia tablas | La autora del PR |
| Chequear la independencia de cada spec antes de pasarla a `in-review` | La autora de la spec |
