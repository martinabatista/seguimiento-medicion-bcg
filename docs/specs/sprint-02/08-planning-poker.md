---
titulo: Planning Poker
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#12"
sprint: 02
autor: Florencia Cabeza
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 08 — Planning Poker

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#12](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/12). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero estimar historias con Planning Poker, para llegar a una estimación consensuada sin que unos influyan a otros. ([#12](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/12))

## Objetivo

Implementar el mecanismo de estimación colaborativa que pide la consigna.

Origen: Consigna — Requerimiento 4 (mecanismo de Planning Poker).

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

- CA-1: Se registran las estimaciones individuales de cada integrante.
- CA-2: Las estimaciones permanecen ocultas hasta finalizar la votación.
- CA-3: Al finalizar, se muestran todas las estimaciones.
- CA-4: Se detectan diferencias entre las estimaciones.
- CA-5: Se pueden realizar nuevas rondas.
- CA-6: Se registra la estimación acordada en la historia.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 02 (integrantes), Spec 07 (story points).
- Los votos y la estimación acordada usan la escala de la spec 07 (1, 2, 3, 5, 8, 13, 21), y la estimación acordada se registra con la misma operación de la spec 07 (CA-6), con sus mismas reglas (por ejemplo, no se estima una historia completada).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
