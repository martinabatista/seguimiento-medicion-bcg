---
titulo: Registrar esfuerzo
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#13"
sprint: 02
autor: Luana Guajardo
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 09 — Registrar esfuerzo

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#13](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/13). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero registrar las horas que trabajé en una historia o tarea, para comparar después el esfuerzo estimado con el real. ([#13](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/13))

## Objetivo

Registrar el esfuerzo real, base de las métricas de horas y desviación.

Origen: Consigna — Requerimiento 5 (Registro de esfuerzo).

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

- CA-1: Cada registro tiene: integrante, fecha, actividad realizada y horas trabajadas.
- CA-2: Se registra sobre una historia o tarea.
- CA-3: Las horas deben ser mayores a cero.
- CA-4: Se puede consultar el esfuerzo registrado por historia.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 02 (integrantes), Spec 04 (backlog).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
