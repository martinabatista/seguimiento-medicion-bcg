---
titulo: Reporte de proyecto o Sprint exportable a PDF
estado: draft        # draft | in-review | approved | implemented
unidad: global
issue: "#21"
sprint: 04
autor: Luana Guajardo
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 14 — Reporte de proyecto o Sprint exportable a PDF

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.
> Borrador generado desde la consigna y el Issue [#21](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/21). La dueña lo completa y lo pasa a `in-review`.

## Historia de usuario

Como integrante del equipo, quiero generar un reporte de un proyecto o de un Sprint y exportarlo a PDF, para compartir y archivar los resultados. ([#21](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/21))

## Objetivo

Resumir el proyecto o Sprint en un informe exportable (Sprint 4: exportar informe a PDF).

Origen: Consigna — Requerimiento 9 (Reportes).

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

- CA-1: El reporte incluye historias planificadas y completadas.
- CA-2: Incluye estimaciones.
- CA-3: Incluye esfuerzo registrado.
- CA-4: Incluye métricas.
- CA-5: Incluye defectos.
- CA-6: Se puede exportar a un archivo PDF.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 10 y Spec 12 (métricas), Spec 09 (esfuerzo), Spec 11 (defectos).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
