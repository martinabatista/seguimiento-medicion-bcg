---
titulo: Métricas de esfuerzo y calidad
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#18"
sprint: 03
autor: Martina Batista
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 12 — Métricas de esfuerzo y calidad

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#18](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/18). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero ver las métricas de esfuerzo y de defectos, para evaluar la precisión de nuestras estimaciones y la calidad del producto. ([#18](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/18))

## Objetivo

Calcular las métricas de esfuerzo y calidad que pide la consigna, desarrolladas con TDD. Definir dónde se cargan las horas estimadas (historia o tarea).

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

- CA-1: Se calculan las horas estimadas.
- CA-2: Se calculan las horas reales.
- CA-3: Se calcula la desviación entre esfuerzo estimado y real.
- CA-4: Los cálculos se desarrollan con TDD y están cubiertos por tests unitarios Go.

> La cantidad de defectos detectados y resueltos (antes CA-4 y CA-5) pasó a la spec 11, que es la dueña de los defectos: así esta historia no depende de otra del mismo Sprint.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 09 (esfuerzo): del Sprint 2. Las horas estimadas se definen en esta spec.
- La usa: spec 14 (reporte).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
