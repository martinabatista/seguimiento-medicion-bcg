---
titulo: Gestionar defectos
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#17"
sprint: 03
autor: Luana Guajardo
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 11 — Gestionar defectos

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#17](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/17). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero registrar los defectos encontrados, para hacer seguimiento de la calidad del producto. ([#17](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/17))

## Objetivo

Registrar y seguir los defectos, base de las métricas de calidad.

Origen: Consigna — Requerimiento 6 (Gestión de defectos).

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

- CA-1: Cada defecto tiene: descripción, severidad, estado, historia relacionada, Sprint de detección y Sprint de resolución.
- CA-2: Se puede cambiar el estado de un defecto.
- CA-3: El Sprint de resolución se registra al resolverlo.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 04 (backlog), Spec 05 (sprints).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
