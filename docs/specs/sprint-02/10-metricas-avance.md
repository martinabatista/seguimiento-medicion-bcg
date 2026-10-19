---
titulo: Métricas de avance del Sprint y del proyecto
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#14"
sprint: 02
autor: Martina Batista
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 10 — Métricas de avance del Sprint y del proyecto

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#14](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/14). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero ver las métricas de avance, para saber cuánto planificamos, cuánto completamos y a qué velocidad trabajamos. ([#14](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/14))

## Objetivo

Calcular las métricas de avance que pide la consigna, desarrolladas con TDD.

Origen: Consigna — Requerimiento 7 (Métricas).

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

- CA-1: Se calculan los Story Points planificados.
- CA-2: Se calculan los Story Points completados.
- CA-3: Se calcula la velocidad del equipo.
- CA-4: Se calcula el porcentaje de historias completadas.
- CA-5: Los cálculos se desarrollan con TDD y están cubiertos por tests unitarios Go.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 05 y Spec 06 (sprints), Spec 07 (story points).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
