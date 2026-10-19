---
titulo: Registrar historias completadas, cerrar Sprint y consultar Sprints anteriores
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#9"
sprint: 01
autor: Martina Batista
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 06 — Registrar historias completadas, cerrar Sprint y consultar Sprints anteriores

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#9](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/9). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero registrar las historias completadas y cerrar el Sprint, para dejar constancia de lo logrado y consultarlo después. ([#9](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/9))

## Objetivo

Cerrar el ciclo del Sprint y conservar su historial, base de las métricas de velocidad.

Origen: Consigna — Requerimiento 3 (Gestión de Sprints: historias completadas, cerrar Sprint, consultar anteriores).

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

- CA-1: Se puede marcar una historia del Sprint como completada.
- CA-2: Se puede cerrar un Sprint (qué pasa con las historias no completadas se define en esta spec).
- CA-3: Un Sprint cerrado no admite cambios.
- CA-4: Se pueden consultar los Sprints anteriores con sus historias.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 05 (sprints).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
