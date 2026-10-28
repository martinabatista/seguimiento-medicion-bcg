# Sprint 3 — Funcionalidad y Calidad

**Equipo:** Martina Batista (Agile Enabler), Florencia Cabeza, Luana Guajardo.
**Período de trabajo:** del 24/10/2026 al 27/10/2026, adelantado y en paralelo con el Sprint 2 (Milestone [Sprint 3 - Funcionalidad y Calidad](https://github.com/martinabatista/seguimiento-medicion-bcg/milestone/4)).
**Modalidad:** las ceremonias se hicieron en reuniones presenciales del equipo, trabajando juntas sobre el proyecto.

## Sprint Planning

**Sprint Goal:** añadir las funcionalidades clave que faltaban y robustecer el sistema: registrar los defectos, medir el esfuerzo estimado contra el real e implementar la visualización con el Dashboard.

**Sprint Backlog:**

| Historia | Dueña | Resultado |
|---|---|---|
| [HU-11](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/17) Gestionar defectos | Luana | Terminada (#71, #72) |
| [HU-12](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/18) Métricas de esfuerzo y calidad | Martina | Terminada (#68, #69) |
| [HU-13](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/19) Dashboard del proyecto | Florencia | Terminada en dos etapas: backend (#57, #58) y pantalla con los gráficos (#64) |

Cada historia siguió el flujo del equipo: spec aprobada por otra integrante → código con TDD y BDD → PR con CI en verde, con su pantalla incluida.

## Dailies

Al trabajar juntas, cada una contaba qué había terminado, qué seguía y qué la bloqueaba. Los bloqueos que aparecieron y cómo se resolvieron:

- **Reglas que la consigna no define:** de dónde salen las horas estimadas, cómo se calcula la desviación, qué estados tiene un defecto o qué muestra el Dashboard. Se decidieron en cada spec antes de programar y quedaron escritas como reglas de negocio.
- **El Dashboard necesitaba la interfaz:** se dividió en dos etapas, el backend primero y la pantalla cuando estuvo la interfaz web (#15), para no dejar la historia esperando.
- **Un error encontrado por los tests:** al modificar una historia desde el backlog se perdían sus horas estimadas. Se corrigió antes de mergear (#69).
- **Tablero desactualizado:** tarjetas de Issues ya cerrados (#17, #21) seguían en "Sprint Backlog" y se movieron a "Hecho".

## Sprint Review

**Incremento entregado:**

- **Defectos:** se registran con su severidad, historia relacionada y Sprint de detección; cambian de estado hasta resolverse, con su Sprint de resolución, y se cuentan los detectados y los resueltos por Sprint y del proyecto.
- **Métricas de esfuerzo:** cada historia tiene sus horas estimadas, y se comparan con las horas reales registradas. La desviación sale en horas y en %, por historia, por Sprint y del proyecto.
- **Dashboard:** el estado del proyecto de un vistazo, con 6 indicadores y 3 gráficos (Story Points por Sprint, velocidad y backlog por estado), hechos sin librerías de gráficos.

Todos los cálculos (conteos, desviación, datos de los gráficos) están en Go y las pantallas solo los muestran.

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

- Se adelantó el Sprint completo: las 3 historias se terminaron antes de lo planificado, en paralelo con el Sprint 2.
- Las historias fueron independientes: ninguna quedó esperando a otra, porque los criterios que dependían de otra historia ya se habían movido (por ejemplo, los conteos de defectos pasaron de la HU-12 a la HU-11).
- Los tests encontraron errores antes de llegar a `main`, como el de las horas estimadas que se perdían al modificar una historia.

**Qué costó:**

- Definir las reglas nuevas que la consigna deja abiertas, antes de poder programar.
- El tablero desactualizado: tarjetas cerradas que no se movían a "Hecho".

**Qué vamos a mejorar en el Sprint 4:**

| Acción | Responsable |
|---|---|
| Ensayar la presentación y la demo, incluido el recorrido de trazabilidad de una funcionalidad | Las tres |
| Al aprobar un PR, confirmar en GitHub que la aprobación quedó registrada | La que revisa |
