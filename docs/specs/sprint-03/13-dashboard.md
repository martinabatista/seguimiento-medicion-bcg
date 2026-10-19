---
titulo: Dashboard del proyecto
estado: draft        # draft | in-review | approved | implemented
unidad: global
issue: "#19"
sprint: 03
autor: Florencia Cabeza
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 13 — Dashboard del proyecto

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#19](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/19). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero un Dashboard con el estado del proyecto y gráficos de las métricas, para entender de un vistazo cómo avanzamos. ([#19](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/19))

## Objetivo

Visualizar el estado y las métricas del proyecto (Sprint 3: implementar la visualización).

Origen: Consigna — Requerimiento 8 (Dashboard).

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

- CA-1: Muestra información del estado del proyecto.
- CA-2: Incluye representaciones gráficas de algunas métricas (las que se elijan en esta spec).

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 03 (estado), Spec 10 y Spec 12 (métricas).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
