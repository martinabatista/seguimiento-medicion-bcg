---
titulo: Crear Sprints, definir Sprint Goal y asignar historias
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#8"
sprint: 01
autor: Martina Batista
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 05 — Crear Sprints, definir Sprint Goal y asignar historias

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#8](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/8). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero crear un Sprint con su Sprint Goal y asignarle historias del backlog, para planificar el trabajo de la iteración. ([#8](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/8))

## Objetivo

Soportar la Sprint Planning: definir la iteración, su objetivo y su alcance.

Origen: Consigna — Requerimiento 3 (Gestión de Sprints: crear, Sprint Goal, asignar historias).

## Alcance

**Entra:**
- _TODO: completar al tomar la historia (antes de implementar)._

**NO entra:**
- _TODO: completar al tomar la historia (antes de implementar)._

## Entradas

_TODO: completar al tomar la historia (antes de implementar)._

## Salidas esperadas

_TODO: completar al tomar la historia (antes de implementar)._

## Reglas de negocio

- _TODO: completar al tomar la historia (antes de implementar)._

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- _TODO: completar al tomar la historia (antes de implementar)._

## Casos límite

- _TODO: completar al tomar la historia (antes de implementar)._

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| _TODO: completar al tomar la historia (antes de implementar)._ | |

## Criterios de aceptación

- CA-1: Se puede crear un Sprint con fechas de inicio y fin.
- CA-2: Se puede definir el Sprint Goal.
- CA-3: Se pueden asignar historias del Product Backlog al Sprint.
- CA-4: Una historia no puede estar asignada a dos Sprints abiertos a la vez.
- CA-5: El estado del proyecto (`GET /api/proyectos/{id}/estado`, spec 03) informa el Sprint en curso: el Sprint abierto del proyecto cuyo período incluye la fecha de hoy (inicio y fin incluidos); si hay más de uno, el de fecha de inicio más reciente; si no hay ninguno, indica que no hay Sprint en curso.

> El CA-5 viene de la spec 03 (antes era su CA-2): lo implementa esta historia porque es la que crea los Sprints, así la HU-03 no queda esperando a la HU-05.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 04 (backlog): implementada.
- Spec 03 (estado del proyecto): implementada; esta historia le agrega el Sprint en curso (CA-5).
- La usan: spec 06 (cerrar Sprint), spec 10 (métricas de avance) y spec 11 (defectos: Sprint de detección y resolución).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
