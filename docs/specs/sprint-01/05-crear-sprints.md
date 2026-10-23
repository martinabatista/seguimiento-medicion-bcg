---
titulo: Crear Sprints, definir Sprint Goal y asignar historias
estado: implemented  # draft | in-review | approved | implemented
unidad: backend
issue: "#8"
sprint: 01
autor: Martina Batista
aprobada_por: "Florencia Cabeza"
fecha: 2026-10-23
specs_relacionadas: ["01-proyectos", "03-estado-proyecto", "04-product-backlog"]
---

# Spec 05 — Crear Sprints, definir Sprint Goal y asignar historias

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero crear un Sprint con su Sprint Goal y asignarle historias del backlog, para planificar el trabajo de la iteración. ([#8](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/8))

## Objetivo

Soportar la Sprint Planning: definir la iteración (sus fechas), su objetivo (Sprint Goal) y su alcance (las historias del backlog que entran). Además, informar en el estado del proyecto cuál es el Sprint en curso.

Origen: Consigna — Requerimiento 3 (Gestión de Sprints: crear, Sprint Goal, asignar historias).

## Alcance

**Entra:**
- Crear un Sprint en un proyecto con fecha de inicio, fecha de fin y Sprint Goal.
- Modificar el Sprint Goal mientras el Sprint esté abierto.
- Listar los Sprints de un proyecto y ver un Sprint con sus historias.
- Asignar historias del backlog a un Sprint abierto y quitarlas.
- Impedir eliminar del backlog una historia asignada a un Sprint abierto (pendiente de la spec 04).
- Informar el Sprint en curso en el estado del proyecto (viene de la spec 03).

**NO entra:**
- Cerrar un Sprint, registrar historias completadas y consultar Sprints anteriores (spec 06).
- Modificar las fechas de un Sprint ya creado: si están mal, se crea otro Sprint.
- Eliminar Sprints.
- Calcular métricas del Sprint (spec 10).
- Pantallas del frontend (todavía no existe; se hacen cuando se arme la interfaz).

## Entradas

**Crear un Sprint:**

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |
| Fecha de inicio | fecha (AAAA-MM-DD) | Sí | Fecha válida, dentro del período del proyecto. |
| Fecha de fin | fecha (AAAA-MM-DD) | Sí | Fecha válida, posterior a la de inicio, dentro del período del proyecto. |
| Sprint Goal | texto | Sí | No puede estar vacío ni tener solo espacios. Máximo 300 caracteres. |

**Modificar el Sprint Goal:** identificador del Sprint y nuevo Sprint Goal (mismas validaciones).

**Asignar o quitar una historia:** identificador del Sprint e identificador de la historia del backlog.

## Salidas esperadas

- **Crear:** el Sprint queda registrado con un identificador asignado por el sistema, un **número** correlativo dentro del proyecto (Sprint 1, Sprint 2…), sus fechas, su Sprint Goal, el estado `Abierto` y sin historias.
- **Modificar el Sprint Goal:** se devuelve el Sprint con el objetivo actualizado.
- **Listar:** se devuelven los Sprints del proyecto ordenados por fecha de inicio. Si no tiene, la lista está vacía.
- **Ver uno:** se devuelve el Sprint con sus datos y la lista de sus historias (identificador, título, prioridad y estado), ordenadas como el backlog (spec 04 RN-8).
- **Asignar / quitar:** se devuelve el Sprint con su lista de historias actualizada.
- **Estado del proyecto** (`spec 03`): incluye el Sprint en curso (identificador, número, Sprint Goal, fechas) o indica que no hay Sprint en curso.
- **Error:** no se crea ni se cambia nada, y se informa un mensaje que explica qué está mal o qué no existe.

## Reglas de negocio

- RN-1: Todo Sprint pertenece a un único proyecto existente.
- RN-2: Las fechas de inicio y de fin son obligatorias y válidas, y la de fin es **posterior** a la de inicio (no igual ni anterior), igual que en la spec 01.
- RN-3: El Sprint tiene que estar **dentro del período del proyecto**: su inicio no puede ser anterior al inicio del proyecto y su fin no puede ser posterior al fin del proyecto.
- RN-4: Dos Sprints del mismo proyecto **no se pueden superponer**: ningún día puede pertenecer a dos Sprints (sean abiertos o cerrados). Un Sprint puede empezar el día siguiente al fin de otro.
- RN-5: El Sprint Goal es obligatorio: no puede estar vacío ni tener solo espacios, y tiene como máximo 300 caracteres.
- RN-6: Un Sprint nuevo arranca en estado `Abierto`, sin historias, y recibe el número siguiente al del último Sprint del proyecto (el primero es el 1).
- RN-7: El Sprint Goal se puede modificar, y se pueden asignar o quitar historias, solo mientras el Sprint está `Abierto`. Un Sprint `Cerrado` (spec 06) no admite cambios.
- RN-8: Solo se pueden asignar historias del **mismo proyecto** que no estén `Completada`.
- RN-9: Una historia no puede estar asignada a dos Sprints abiertos a la vez. Sí puede haber estado en un Sprint ya cerrado y asignarse a uno nuevo.
- RN-10: Solo se puede quitar una historia que está asignada a ese Sprint.
- RN-11: No se puede eliminar del backlog (spec 04) una historia asignada a un Sprint abierto: primero hay que quitarla del Sprint.
- RN-12: El **Sprint en curso** de un proyecto es su Sprint `Abierto` cuyo período incluye la fecha de hoy (inicio y fin incluidos). Por RN-4 puede haber como máximo uno. Si no hay, no hay Sprint en curso.
- RN-13: El identificador y el número los asigna el sistema y no se modifican. Un Sprint de otro proyecto se trata como inexistente.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- Las validaciones de RN-2 a RN-12 viven en el dominio y se desarrollan con TDD.
- Los Sprints cuelgan del proyecto: se accede siempre a través del proyecto al que pertenecen (decisiones técnicas en `backend/CLAUDE.md`).
- La fecha de hoy (RN-12) la toma el sistema, igual que en la spec 03.

## Casos límite

- Sprint que empieza el mismo día que el proyecto o termina el mismo día que el proyecto: se acepta.
- Sprint que empieza un día antes del proyecto o termina un día después: se rechaza.
- Sprint que dura un solo día (fin = inicio + 1): se acepta.
- Sprint que empieza el día siguiente al fin de otro: se acepta. Que empieza el mismo día en que termina otro: se rechaza (se superponen en ese día).
- Sprint Goal de exactamente 300 caracteres: se acepta. Con 301: se rechaza.
- Asignar la misma historia dos veces al mismo Sprint: se rechaza (ya está asignada).
- Asignar una historia `En progreso`: se acepta. Una `Completada`: se rechaza.
- Quitar una historia y volver a asignarla: se acepta.
- Hoy es el último día del Sprint: sigue siendo el Sprint en curso. Hoy cae entre dos Sprints: no hay Sprint en curso.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Falta alguna fecha | Dato inválido: "Las fechas de inicio y finalización son obligatorias" |
| Fecha con formato inválido o inexistente | Dato inválido: "La fecha no es válida" |
| Fin igual o anterior al inicio | Dato inválido: "La fecha de finalización debe ser posterior a la de inicio" |
| Sprint fuera del período del proyecto | Dato inválido: "El Sprint debe estar dentro de las fechas del proyecto" |
| Sprint Goal vacío o con solo espacios | Dato inválido: "El Sprint Goal es obligatorio" |
| Sprint Goal de más de 300 caracteres | Dato inválido: "El Sprint Goal no puede superar los 300 caracteres" |
| El Sprint se superpone con otro del proyecto | Conflicto: "El Sprint se superpone con otro Sprint del proyecto" |
| Cambiar un Sprint cerrado | Conflicto: "El Sprint está cerrado" |
| Asignar una historia completada | Conflicto: "La historia ya está completada" |
| Asignar una historia que ya está en un Sprint abierto | Conflicto: "La historia ya está asignada a un Sprint abierto" |
| Quitar una historia que no está en ese Sprint | Conflicto: "La historia no está asignada a este Sprint" |
| Eliminar del backlog una historia asignada a un Sprint abierto | Conflicto: "La historia está asignada a un Sprint abierto" |
| Proyecto inexistente | No encontrado: "El proyecto no existe" |
| Sprint inexistente (o de otro proyecto) | No encontrado: "El Sprint no existe" |
| Historia inexistente (o de otro proyecto) | No encontrado: "La historia no existe" |

## Criterios de aceptación

- CA-1: Se puede crear un Sprint con fechas de inicio y fin dentro del período del proyecto, sin superponerse con otro Sprint del proyecto; queda registrado con identificador, número correlativo y estado `Abierto`.
- CA-2: El Sprint Goal es obligatorio (máximo 300 caracteres) y se puede modificar mientras el Sprint está abierto.
- CA-3: Se pueden asignar a un Sprint abierto historias del backlog del mismo proyecto que no estén completadas, y quitarlas; se puede listar los Sprints del proyecto y ver un Sprint con sus historias.
- CA-4: Una historia no puede estar asignada a dos Sprints abiertos a la vez.
- CA-5: El estado del proyecto (`GET /api/proyectos/{id}/estado`, spec 03) informa el Sprint en curso: el Sprint abierto cuyo período incluye la fecha de hoy; si no hay, indica que no hay Sprint en curso.
- CA-6: No se puede eliminar del backlog una historia asignada a un Sprint abierto.
- CA-7: Operar sobre un proyecto, un Sprint o una historia que no existe informa cuál no existe, sin cambiar nada.

Relación con el Issue #8: "crear un Sprint con fechas" → CA-1; "definir el Sprint Goal" → CA-2; "asignar historias" → CA-3; "no en dos Sprints abiertos" → CA-4; "Sprint en curso" (de la spec 03) → CA-5. CA-6 resuelve el punto pendiente de la spec 04 y CA-7 los casos de error.

## Escenarios BDD

```gherkin
# language: es
Característica: Crear Sprints, definir Sprint Goal y asignar historias
  Para planificar el trabajo de cada iteración
  Como integrante del equipo
  Quiero crear Sprints con su objetivo y asignarles historias del backlog

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-31"

  Escenario: Crear un Sprint (CA-1)
    Cuando creo un Sprint del "2026-11-01" al "2026-11-14" con objetivo "Registrar clientes"
    Entonces el Sprint queda registrado como el Sprint 1, abierto y sin historias

  Escenario: Los Sprints se numeran en orden (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando creo un Sprint del "2026-11-15" al "2026-11-28" con objetivo "Catálogo"
    Entonces el Sprint queda registrado como el Sprint 2, abierto y sin historias

  Esquema del escenario: Validar las fechas del Sprint (CA-1)
    Dado que existe un Sprint del "2026-11-15" al "2026-11-28"
    Cuando creo un Sprint del "<inicio>" al "<fin>" con objetivo "Objetivo"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | inicio     | fin        | resultado                                                  |
      | 2026-11-01 | 2026-11-14 | registrado                                                 |
      | 2026-11-29 | 2026-12-31 | registrado                                                 |
      | 2026-11-29 | 2026-11-30 | registrado                                                 |
      | 2026-10-31 | 2026-11-10 | El Sprint debe estar dentro de las fechas del proyecto     |
      | 2026-12-20 | 2027-01-02 | El Sprint debe estar dentro de las fechas del proyecto     |
      | 2026-11-10 | 2026-11-15 | El Sprint se superpone con otro Sprint del proyecto        |
      | 2026-11-28 | 2026-12-05 | El Sprint se superpone con otro Sprint del proyecto        |
      | 2026-12-01 | 2026-12-01 | La fecha de finalización debe ser posterior a la de inicio |
      | (vacío)    | 2026-12-01 | Las fechas de inicio y finalización son obligatorias       |

  Esquema del escenario: Validar el Sprint Goal (CA-2)
    Cuando creo un Sprint del "2026-11-01" al "2026-11-14" con un objetivo de <largo> caracteres
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                          |
      | 300   | registrado                                         |
      | 301   | El Sprint Goal no puede superar los 300 caracteres |
      | 0     | El Sprint Goal es obligatorio                      |

  Escenario: Modificar el Sprint Goal (CA-2)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando cambio su objetivo a "Registrar y loguear clientes"
    Entonces el Sprint tiene objetivo "Registrar y loguear clientes"

  Escenario: Asignar historias a un Sprint (CA-3)
    Dado que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando asigno al Sprint las historias "Ver catálogo" y "Registrar usuario"
    Entonces el Sprint tiene las historias "Registrar usuario" y "Ver catálogo" en ese orden

  Escenario: Quitar una historia de un Sprint (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando quito del Sprint la historia "Registrar usuario"
    Entonces el Sprint no tiene historias

  Escenario: No se puede asignar una historia completada (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1 en estado "Completada"
    Y que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando asigno al Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está completada"

  Escenario: No se puede quitar una historia que no está en el Sprint (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando quito del Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia no está asignada a este Sprint"

  Escenario: Listar los Sprints por fecha de inicio (CA-3)
    Dado que existe un Sprint del "2026-11-15" al "2026-11-28"
    Y que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando listo los Sprints del proyecto
    Entonces veo los Sprints que empiezan el "2026-11-01" y el "2026-11-15" en ese orden

  Escenario: Una historia no puede estar en dos Sprints abiertos (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que existe un Sprint del "2026-11-15" al "2026-11-28"
    Cuando asigno al segundo Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está asignada a un Sprint abierto"

  Escenario: Asignar dos veces la misma historia al mismo Sprint (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando asigno al Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está asignada a un Sprint abierto"

  Esquema del escenario: Sprint en curso en el estado del proyecto (CA-5)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Y que existe un Sprint del "2026-11-16" al "2026-11-28"
    Y que hoy es "<hoy>"
    Cuando consulto el estado del proyecto
    Entonces el Sprint en curso es "<sprint>"

    Ejemplos:
      | hoy        | sprint          |
      | 2026-11-01 | Sprint 1        |
      | 2026-11-14 | Sprint 1        |
      | 2026-11-15 | ninguno         |
      | 2026-11-20 | Sprint 2        |
      | 2026-12-10 | ninguno         |

  Escenario: No se puede eliminar del backlog una historia de un Sprint abierto (CA-6)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando elimino la historia "Registrar usuario"
    Entonces se informa el error "La historia está asignada a un Sprint abierto"

  Escenario: Se puede eliminar una historia después de quitarla del Sprint (CA-6)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando quito del Sprint la historia "Registrar usuario"
    Y elimino la historia "Registrar usuario"
    Entonces el backlog está vacío

  Escenario: Crear un Sprint en un proyecto que no existe (CA-7)
    Cuando creo un Sprint en el proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Ver un Sprint que no existe (CA-7)
    Cuando veo el Sprint con identificador 999
    Entonces se informa el error "El Sprint no existe"

  Escenario: Asignar una historia que no existe (CA-7)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando asigno al Sprint la historia con identificador 999
    Entonces se informa el error "La historia no existe"
```

Los pasos que no indican proyecto se refieren al de los Antecedentes, "Tienda Online". "Que existe un Sprint del … al …" crea el Sprint con un objetivo cualquiera. En las tablas, `(vacío)` representa un texto vacío, igual que en las specs 01 y 04. El paso "que hoy es …" es el de la spec 03.

## Dependencias

- Spec 01 (proyectos), spec 04 (backlog) y spec 03 (estado del proyecto): implementadas. Esta historia agrega el Sprint en curso al estado (CA-5) y la regla de eliminación del backlog (CA-6).
- La usan: spec 06 (cerrar Sprint), spec 10 (métricas de avance) y spec 11 (defectos: Sprint de detección y resolución).

## Tareas

- [x] `[P1]` Entidad Sprint en `backend/dominio/` con las reglas RN-2 a RN-13 (TDD).
- [x] `[P1]` Tablas de Sprints y de historias asignadas en `backend/postgres/esquema/` (`04_sprints.sql`: `sprints` y `sprint_historias`).
- [x] `[P1]` Repositorio de Sprints en `backend/postgres/` y en `backend/memoria/`.
- [x] `[P1]` Rutas de la API bajo `/api/proyectos/{id}/sprints` en `backend/api/`.
- [x] `[P1]` Sprint en curso en el estado del proyecto (spec 03) y regla de eliminación en el backlog (spec 04).
- [x] `[P1]` Escenarios de `backend/features/sprints.feature` con sus pasos.
- [x] `[P2]` Actualizar `backend/CLAUDE.md` (endpoints, specs relacionadas).

## Definition of Done

- [x] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [x] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [x] Criterios de aceptación verificados.
- [x] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [x] Spec en estado `implemented`.
