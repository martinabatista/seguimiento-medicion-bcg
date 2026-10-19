---
titulo: Consultar el estado de un proyecto
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#6"
sprint: 01
autor: Luana Guajardo
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 03 — Consultar el estado de un proyecto

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#6](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/6). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero consultar el estado de un proyecto, para saber cómo viene avanzando. ([#6](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/6))

## Objetivo

Dar una vista rápida y confiable de la situación actual de un proyecto.

Origen: Consigna — Requerimiento 1 (Gestión de proyectos: consultar el estado).

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

- CA-1: Se muestran los datos del proyecto y su estado (los estados posibles se definen en esta spec).
- CA-2: Se muestra el Sprint en curso, si lo hay.
- CA-3: Se muestra un resumen del avance del Product Backlog.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 01 (proyectos), Spec 04 (backlog), Spec 05 (sprints).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
