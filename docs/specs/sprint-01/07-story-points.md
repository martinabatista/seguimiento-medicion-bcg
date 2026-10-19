---
titulo: Estimar historias con Story Points
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#10"
sprint: 01
autor: Florencia Cabeza
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 07 — Estimar historias con Story Points

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#10](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/10). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero asignar Story Points a una historia, para conocer el tamaño relativo del trabajo. ([#10](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/10))

## Objetivo

Dimensionar las historias para poder planificar Sprints y medir velocidad.

Origen: Consigna — Requerimiento 4 (Estimación mediante Story Points).

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

- CA-1: Se puede asignar y modificar la estimación en Story Points de una historia.
- CA-2: Solo se aceptan valores de la escala definida (la escala se define en esta spec, p. ej. Fibonacci).

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 04 (backlog).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
