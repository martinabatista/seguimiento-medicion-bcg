# Sprint 4 — Cierre y Entrega Final

**Equipo:** Martina Batista (Agile Enabler), Florencia Cabeza, Luana Guajardo.
**Período de trabajo:** del 27/10/2026 al 03/11/2026, adelantado respecto de las fechas del Milestone [Sprint 4 - Cierre y Entrega Final](https://github.com/martinabatista/seguimiento-medicion-bcg/milestone/5).
**Modalidad:** las ceremonias se hicieron en reuniones presenciales del equipo, trabajando juntas sobre el proyecto.

## Sprint Planning

**Sprint Goal:** pulir el producto, documentarlo y preparar la entrega final, con el reporte del proyecto exportable a PDF.

**Sprint Backlog:**

| Tarea | Dueña | Resultado |
|---|---|---|
| [HU-14](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/21) Reporte de proyecto o Sprint exportable a PDF | Luana | Terminada (#73, #74) |
| [#22](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/22) Informe de métricas y cobertura de pruebas | Martina | Terminado (#78) |
| [#23](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/23) Documentación técnica y manual breve de usuario | Florencia | Terminado (#81) |
| [#79](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/79) Mejora de UX/UI de la interfaz web | Martina | Terminada (#80, #82, #83) |
| [#24](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/24) Presentación y demo final | Las tres | Preparada; se presenta en la Sprint Review final |

La HU-14 siguió el flujo del equipo: spec aprobada por otra integrante → código con TDD y BDD → PR con CI en verde, con su pantalla incluida. Las demás tareas no tienen spec porque no agregan funcionalidad: son documentación y diseño.

## Dailies

Al trabajar juntas, cada una contaba qué había terminado, qué seguía y qué la bloqueaba. Los bloqueos que aparecieron y cómo se resolvieron:

- **La app no levantaba:** al probarla, Docker Desktop estaba apagado y el backend no arrancaba sin la base. Se resolvió abriendo Docker antes de `docker compose up -d`, como indica el README.
- **El rediseño no debía cambiar el funcionamiento:** la mejora de UX/UI se hizo solo en los estilos (`estilos.css` e `index.html`), sin tocar ninguna pantalla, y se revisó recorriendo todas las secciones de la app.

## Sprint Review

La Sprint Review final es la presentación formal del producto completo al Cliente (los profesores), con la demo de la app y el recorrido de trazabilidad de al menos una funcionalidad. **Todavía no se hizo**: lo que surja de esa presentación se agrega a esta acta.

**Incremento entregado:**

- **Reporte en PDF:** reporte del proyecto o de un Sprint con historias planificadas y completadas, estimaciones, esfuerzo, métricas y defectos, que el backend arma y la pantalla descarga.
- **Nuevo diseño de la interfaz:** fondo oscuro, menú lateral y un solo color para los datos y gráficos. La app se abre en `http://software-metrics.localhost:5173`.
- **Entregables escritos:** [informe de métricas y cobertura](../informe-metricas.md), [documentación técnica](../documentacion-tecnica.md) y [manual de usuario](../manual-usuario.md).

**Calidad al cierre del Sprint (03/11/2026):**

| Control | Resultado |
|---|---|
| Escenarios BDD (godog) | 359 escenarios y 2309 pasos en verde |
| Cobertura de tests | 88,7 % del backend; dominio 93,6 % (mínimo del equipo: 80 % y 90 %) |
| Frontend | lint, tipos, tests (Vitest) y build en verde |
| CI en GitHub Actions | en verde en todos los PRs |

## Retrospectiva

**Qué salió bien:**

- Se terminó todo antes de la fecha del Milestone: la última historia, los entregables escritos y la mejora de diseño.
- La interfaz quedó más clara y prolija con el rediseño, sin cambiar su funcionamiento.
- La documentación quedó completa: informe de métricas, documentación técnica y manual de usuario.
- Buen trabajo en equipo: cada una cumplió su parte y los PRs se revisaron rápido.

**Qué costó:**

- Levantar la app para probarla: hay que tener Docker Desktop abierto y la base, el backend y el frontend andando a la vez.

**Qué nos llevamos para próximos proyectos:**

| Aprendizaje | Cómo aplicarlo |
|---|---|
| Mantener el tablero al día | Mover cada tarjeta a la par del trabajo (al empezar, al abrir el PR y al mergear), no al final del Sprint |
