---
titulo: Administrar el Product Backlog
estado: draft        # draft | in-review | approved | implemented
unidad: backend
issue: "#7"
sprint: 01
autor: Florencia Cabeza
aprobada_por: ""
fecha: 2026-10-20
specs_relacionadas: ["01-proyectos"]
---

# Spec 04 — Administrar el Product Backlog

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero administrar el Product Backlog de un proyecto, para mantener priorizado el trabajo pendiente. ([#7](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/7))

## Objetivo

Mantener, para cada proyecto, la lista priorizada de historias que alimenta la estimación (Story Points y Planning Poker) y la planificación de los Sprints. Cada elemento del backlog es una **historia**.

Origen: Consigna — Requerimiento 2 (Product Backlog).

## Alcance

**Entra:**
- Crear una historia en el backlog de un proyecto con título, descripción, prioridad y criterios de aceptación.
- Modificar los datos de una historia, incluido su estado.
- Eliminar una historia del backlog.
- Listar el backlog de un proyecto ordenado por prioridad.
- Ver una historia por su identificador.

**NO entra:**
- Asignar o modificar los Story Points: la historia nace "sin estimar" y la estimación se define en la spec 07 (escala) y la spec 08 (Planning Poker).
- Asignar historias a un Sprint y marcarlas como completadas dentro de un Sprint (specs 05 y 06).
- Reordenar el backlog arrastrando historias: el orden sale de la prioridad.
- Pantallas del frontend (todavía no existe; se hacen cuando se arme la interfaz).

## Entradas

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |
| Título | texto | Sí | No puede estar vacío ni tener solo espacios. Máximo 150 caracteres. Puede repetirse. |
| Descripción | texto | Sí | No puede estar vacía ni tener solo espacios. Máximo 1000 caracteres. |
| Prioridad | número entero | Sí | Mayor o igual a 1. **1 es la más importante.** Puede repetirse entre historias. |
| Criterios de aceptación | texto | No | Máximo 2000 caracteres. Puede quedar vacío mientras la historia se refina. |
| Estado | uno de: `Pendiente`, `En progreso`, `Completada` | Solo al modificar | Al crear no se indica: siempre arranca en `Pendiente`. |

Para modificar, ver o eliminar, además se indica el **identificador** de la historia (y el del proyecto al que pertenece).

## Salidas esperadas

- **Crear:** la historia queda registrada en el backlog del proyecto y se devuelve con su identificador asignado por el sistema, sus datos, el estado `Pendiente` y los Story Points **sin estimar**.
- **Modificar:** se devuelve la historia con los datos actualizados.
- **Eliminar:** la historia deja de existir y ya no aparece en el backlog.
- **Listar:** se devuelven las historias del proyecto ordenadas por prioridad (de 1 en adelante); a igual prioridad, en orden de creación. Si el proyecto no tiene historias, la lista está vacía.
- **Ver una:** se devuelve la historia con su identificador, sus datos, su estado y sus Story Points.
- **Error:** no se crea, modifica ni elimina nada, y se informa un mensaje que explica qué dato está mal o qué no existe.

## Reglas de negocio

- RN-1: Toda historia pertenece a un único proyecto existente. No se puede crear ni consultar el backlog de un proyecto que no existe.
- RN-2: El título es obligatorio: no puede estar vacío ni tener solo espacios, y tiene como máximo 150 caracteres. Puede repetirse.
- RN-3: La descripción es obligatoria: no puede estar vacía ni tener solo espacios, y tiene como máximo 1000 caracteres.
- RN-4: La prioridad es obligatoria y es un número entero mayor o igual a 1, donde 1 es la más importante. Varias historias pueden tener la misma prioridad.
- RN-5: Los criterios de aceptación son opcionales y tienen como máximo 2000 caracteres.
- RN-6: Una historia nueva siempre arranca en estado `Pendiente` y con los Story Points sin estimar.
- RN-7: El estado solo puede ser `Pendiente`, `En progreso` o `Completada`. Al modificar se puede pasar a cualquiera de los tres.
- RN-8: El backlog se ordena por prioridad ascendente; a igual prioridad, por orden de creación.
- RN-9: Al modificar se aplican las mismas reglas que al crear. Si alguna falla, la historia queda como estaba.
- RN-10: El identificador lo asigna el sistema al crear la historia y no se puede modificar. Una historia de otro proyecto se trata como inexistente.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- Las validaciones de RN-2 a RN-9 viven en el dominio y se desarrollan con TDD.
- El backlog cuelga del proyecto: se accede siempre a través del proyecto al que pertenece (ver decisiones técnicas en `backend/CLAUDE.md`).
- El largo de los textos se cuenta en caracteres, no en bytes (los acentos y la ñ cuentan como uno).

## Casos límite

- Título de exactamente 150 caracteres: se acepta. Con 151: se rechaza.
- Descripción de exactamente 1000 caracteres: se acepta. Con 1001: se rechaza.
- Criterios de aceptación vacíos: se acepta. Con exactamente 2000 caracteres: se acepta. Con 2001: se rechaza.
- Título o descripción con solo espacios: se rechaza, igual que vacío.
- Título igual al de otra historia del mismo proyecto: se acepta.
- Prioridad 1: se acepta (es la más alta). Prioridad 0 o negativa: se rechaza. Prioridad con decimales (por ejemplo, 2.5): se rechaza.
- Prioridad muy alta (por ejemplo, 1000): se acepta; no hay máximo.
- Dos historias con la misma prioridad: se listan en el orden en que se crearon.
- Estado escrito distinto a los valores permitidos (por ejemplo, `Terminada` o `pendiente` en minúscula): se rechaza.
- Listar el backlog de un proyecto sin historias: devuelve una lista vacía, no un error.
- Eliminar una historia y volver a listar: ya no aparece; eliminarla otra vez informa que no existe.
- Pedir una historia existente pero indicando otro proyecto: se informa que la historia no existe.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Título vacío o con solo espacios | Dato inválido: "El título es obligatorio" |
| Título de más de 150 caracteres | Dato inválido: "El título no puede superar los 150 caracteres" |
| Descripción vacía o con solo espacios | Dato inválido: "La descripción es obligatoria" |
| Descripción de más de 1000 caracteres | Dato inválido: "La descripción no puede superar los 1000 caracteres" |
| Falta la prioridad | Dato inválido: "La prioridad es obligatoria" |
| Prioridad menor a 1 o no entera | Dato inválido: "La prioridad debe ser un número entero mayor o igual a 1" |
| Criterios de aceptación de más de 2000 caracteres | Dato inválido: "Los criterios de aceptación no pueden superar los 2000 caracteres" |
| Estado distinto de `Pendiente`, `En progreso` o `Completada` | Dato inválido: "El estado no es válido" |
| Operar sobre el backlog de un proyecto que no existe | No encontrado: "El proyecto no existe" |
| Ver, modificar o eliminar una historia que no existe (o que es de otro proyecto) | No encontrado: "La historia no existe" |

## Criterios de aceptación

- CA-1: Cada ítem tiene: identificador, título, descripción, prioridad, estado, Story Points y criterios de aceptación.
- CA-2: Se pueden crear, modificar, eliminar y listar ítems.
- CA-3: El backlog se puede ver ordenado por prioridad.

## Escenarios BDD

_TODO: completar al tomar la historia (antes de implementar)._ Incluir casos normales, alternativos, límite y de error; indicar qué CA cubre cada uno.

```gherkin
Escenario: <caso normal> (CA-1)
  Dado ...
  Cuando ...
  Entonces ...
```

## Dependencias

- Spec 01 (proyectos).

## Tareas

- [ ] _TODO: completar al tomar la historia (antes de implementar)._

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
