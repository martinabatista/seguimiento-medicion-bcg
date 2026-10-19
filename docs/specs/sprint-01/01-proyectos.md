---
titulo: Crear y modificar proyectos
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#4"
sprint: 01
autor: Martina Batista
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 01 — Crear y modificar proyectos

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#4](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/4). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero crear y modificar un proyecto con sus fechas de inicio y finalización, para administrar su seguimiento. ([#4](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/4))

## Objetivo

Permitir dar de alta y mantener los proyectos sobre los que se realiza toda la estimación, el seguimiento y la medición.

Origen: Consigna — Requerimiento 1 (Gestión de proyectos).

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

- CA-1: Se puede crear un proyecto con nombre, descripción, fecha de inicio y fecha de finalización.
- CA-2: Se pueden modificar los datos de un proyecto existente.
- CA-3: La fecha de finalización no puede ser anterior a la de inicio.
- CA-4: El nombre es obligatorio.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

Ninguna (base del resto del sistema).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
