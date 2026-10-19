---
titulo: Administrar el Product Backlog
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#7"
sprint: 01
autor: Florencia Cabeza
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 04 — Administrar el Product Backlog

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#7](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/7). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero administrar el Product Backlog de un proyecto, para mantener priorizado el trabajo pendiente. ([#7](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/7))

## Objetivo

Mantener la lista priorizada de historias que alimenta la planificación de los Sprints.

Origen: Consigna — Requerimiento 2 (Product Backlog).

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

- CA-1: Cada ítem tiene: identificador, título, descripción, prioridad, estado, Story Points y criterios de aceptación.
- CA-2: Se pueden crear, modificar, eliminar y listar ítems.
- CA-3: El backlog se puede ver ordenado por prioridad.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 01 (proyectos).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
