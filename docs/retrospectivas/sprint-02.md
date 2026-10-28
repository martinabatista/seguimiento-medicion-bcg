# Sprint 2 — La Interfaz

**Equipo:** Martina Batista (Agile Enabler), Florencia Cabeza, Luana Guajardo.
**Período de trabajo:** del 23/10/2026 al 27/10/2026 (Milestone [Sprint 2 - La Interfaz](https://github.com/martinabatista/seguimiento-medicion-bcg/milestone/3)).
**Modalidad:** las ceremonias se hicieron en reuniones presenciales del equipo, trabajando juntas sobre el proyecto.

## Sprint Planning

**Sprint Goal:** dotar al software de una interfaz usable: una página web con una pantalla por funcionalidad, sobre el backend en Go del Sprint 1, y completar el Planning Poker, el registro de esfuerzo y las métricas de avance.

**Sprint Backlog:**

| Historia | Dueña | Resultado |
|---|---|---|
| [HU-08](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/12) Planning Poker | Florencia | Terminada (#49, #50; votaciones abiertas y pantalla en #62, #63) |
| [HU-09](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/13) Registrar esfuerzo | Luana | Terminada (#55, #56) |
| [HU-10](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/14) Métricas de avance | Martina | Terminada (#52, #53) |
| [#15](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/15) Interfaz web | Las tres | Terminada: base y pantallas de Proyectos, Backlog, Planning Poker y Dashboard de Florencia (#60, #61, #63, #64); Sprints, Métricas y modificar proyecto de Martina (#65, #67, #75); Integrantes, Estado y Esfuerzo de Luana (#70) |
| Cierre de la [HU-03](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/6) (pendiente del Sprint 1) | Luana | Terminada (#54) |

Cada historia siguió el flujo del equipo: spec aprobada por otra integrante → código con TDD y BDD → PR con CI en verde. Las pantallas se sumaron a cada spec en su sección "Pantalla".

**Decisión de la interfaz:** la consigna permite web, escritorio o consola. El equipo dudó entre web y consola, lo consultó con el profesor y eligió **web con React + Vite + TypeScript**, lo más simple posible: sin librerías de rutas, de estilos ni de gráficos, y una pantalla por archivo.

## Dailies

Al trabajar juntas, cada una contaba qué había terminado, qué seguía y qué la bloqueaba. Los bloqueos que aparecieron y cómo se resolvieron:

- **Qué tipo de interfaz hacer:** se frenó el arranque de la #15 hasta decidirlo con el profesor.
- **Archivos compartidos** (`App.tsx`, `api.ts`, `Proyectos.tsx`): la base de la interfaz (#60) dejó una guía en `frontend/CLAUDE.md` para sumar pantallas tocando una sola línea del menú, y quien tocó un archivo ajeno lo avisó en el PR.
- **Una parte sin terminar del reparto:** "modificar proyecto" (HU-01) quedó sin pantalla hasta que Luana lo señaló en su PR (#70); se hizo en #75.
- **Columnas y tablas nuevas** (horas estimadas, esfuerzo, defectos): cada una actualizó su base local después de los merges.

## Sprint Review

**Incremento entregado:** una aplicación web usable de punta a punta. Desde el navegador se puede crear y modificar proyectos, registrar integrantes, administrar el backlog y estimarlo con Story Points o con Planning Poker, planificar y cerrar Sprints, registrar horas de esfuerzo y ver el estado del proyecto, las métricas de avance y el Dashboard. Todas las reglas y los cálculos siguen en Go: la interfaz solo muestra y captura datos, y los mensajes de error son los del backend.

**Se adelantaron historias de los Sprints 3 y 4:** HU-11 Defectos (#71, #72), HU-12 Métricas de esfuerzo (#68, #69), HU-13 Dashboard (#57, #58, #64) y HU-14 Reporte PDF (#73, #74). Con eso las 14 historias del Product Backlog quedaron implementadas.

**Calidad al cierre del Sprint (27/10/2026):**

| Control | Resultado |
|---|---|
| Escenarios BDD (godog) | 359 escenarios y 2309 pasos en verde |
| Cobertura de tests | dominio 94 %, API 94 %, PostgreSQL 83 %, PDF 97 % |
| Frontend | lint, tipos, tests (Vitest) y build en verde |
| CI en GitHub Actions | en verde en todos los PRs |
| TDD | commits `test` → `feat` (→ `refactor`) en el historial de cada historia |

## Retrospectiva

**Qué salió bien:**

- Se terminó todo lo planificado: las 3 historias y la interfaz completa, y además se adelantaron las historias de los Sprints siguientes.
- El reparto por áreas y los PRs: cada una hizo sus pantallas y los PRs se revisaron y aprobaron rápido.
- TDD y BDD salieron sin problemas: el flujo spec → tests → código ya es la forma natural de trabajar.

**Qué costó:**

- Decidir el tipo de interfaz: hubo que frenar y consultarlo antes de arrancar.
- Coordinar los archivos que comparten las pantallas (`App.tsx`, `api.ts` y la pantalla de Proyectos).

**Qué vamos a mejorar en el Sprint 3:**

| Acción | Responsable |
|---|---|
| Avisar en el grupo cuando un PR agrega una tabla o una columna, con el comando para actualizar la base | La autora del PR |
| Mover la tarjeta a "Hecho" apenas se mergea el PR que cierra el Issue | La autora del PR; la Agile Enabler lo controla |
| Repasar el reparto antes de cerrar una tarea compartida, para que no quede nada sin dueña | La Agile Enabler |
