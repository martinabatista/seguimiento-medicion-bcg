---
titulo: <Título de la funcionalidad>
estado: draft        # draft | in-review | approved | implemented
unidad: backend      # backend | frontend | global
issue: "#<N>"        # historia de usuario en GitHub
sprint: <NN>
autor: <Nombre Apellido>
aprobada_por: ""     # otra integrante, al pasar a approved
fecha: <YYYY-MM-DD>
specs_relacionadas: []
---

# Spec <NN> — <Título>

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como <rol>, quiero <acción>, para <beneficio>. (Issue #<N>)

## Objetivo

<Para qué sirve. Qué problema resuelve.>

## Alcance

**Entra:**
- 

**NO entra:**
- 

## Entradas

<Datos que recibe la funcionalidad, con tipo y obligatoriedad.>

## Salidas esperadas

<Qué devuelve o qué cambia en el sistema.>

## Reglas de negocio

- RN-1: 

## Restricciones

- 

## Casos límite

- 

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
|  |  |

## Criterios de aceptación

- CA-1: 

## Escenarios BDD

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

<Incluir casos normales, alternativos, límite y de error. Se automatizan en `backend/features/`.>

## Tareas

- [ ] `[P1]` 

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
