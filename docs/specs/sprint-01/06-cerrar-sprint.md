---
titulo: Registrar historias completadas, cerrar Sprint y consultar Sprints anteriores
estado: implemented  # draft | in-review | approved | implemented
unidad: backend
issue: "#9"
sprint: 01
autor: Martina Batista
aprobada_por: "Florencia Cabeza"
fecha: 2026-10-23
specs_relacionadas: ["04-product-backlog", "05-crear-sprints"]
---

# Spec 06 — Registrar historias completadas, cerrar Sprint y consultar Sprints anteriores

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero registrar las historias completadas y cerrar el Sprint, para dejar constancia de lo logrado y consultarlo después. ([#9](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/9))

## Objetivo

Cerrar el ciclo del Sprint y conservar su historial: qué historias se planificaron y cuáles se completaron dentro de ese Sprint. Es la base de las métricas de velocidad (spec 10).

Origen: Consigna — Requerimiento 3 (Gestión de Sprints: historias completadas, cerrar Sprint, consultar anteriores).

## Alcance

**Entra:**
- Marcar como completada una historia de un Sprint abierto, y deshacerlo mientras el Sprint siga abierto.
- Cerrar un Sprint, en cualquier momento.
- Consultar los Sprints anteriores (cerrados) con sus historias planificadas y cuáles se completaron.
- Impedir que el backlog cambie a mano el estado `Completada` de una historia que está en un Sprint abierto.

**NO entra:**
- Reabrir un Sprint cerrado.
- Calcular métricas (Story Points planificados y completados, velocidad, porcentaje): spec 10.
- Pasar automáticamente las historias no completadas al Sprint siguiente: quedan libres en el backlog y se asignan a mano (spec 05).
- Pantallas del frontend (todavía no existe; se hacen cuando se arme la interfaz).

## Entradas

| Operación | Datos |
|---|---|
| Completar o deshacer | Identificador del proyecto, del Sprint y de la historia. |
| Cerrar | Identificador del proyecto y del Sprint. |
| Consultar anteriores | Identificador del proyecto. |

## Salidas esperadas

- **Completar:** la historia pasa al estado `Completada` y el Sprint registra que se completó **en ese Sprint**. Se devuelve el Sprint con sus historias, cada una indicando si se completó en el Sprint.
- **Deshacer:** la historia vuelve al estado `En progreso` y el Sprint deja de registrarla como completada. Se devuelve el Sprint actualizado.
- **Cerrar:** el Sprint pasa a `Cerrado`. Conserva todas las historias que tenía asignadas, marcando cuáles se completaron en él. Se devuelve el Sprint cerrado.
- **Consultar anteriores:** los Sprints cerrados del proyecto, del más reciente (fecha de inicio más tardía) al más antiguo, cada uno con sus historias y cuáles se completaron. Si no hay, la lista está vacía.
- **Error:** no se cambia nada, y se informa un mensaje que explica qué está mal o qué no existe.

## Reglas de negocio

- RN-1: Solo se puede completar una historia **asignada** al Sprint, con el Sprint **abierto**.
- RN-2: Completar una historia la pasa al estado `Completada` (spec 04) y el Sprint registra que se completó en él. Una historia ya completada en el Sprint no se puede completar otra vez.
- RN-3: Mientras el Sprint esté abierto, se puede deshacer una historia completada en él: vuelve a `En progreso` y el Sprint deja de registrarla como completada. Solo se deshace una historia que está completada en ese Sprint.
- RN-4: Mientras una historia está en un Sprint abierto, el backlog (spec 04) **no** puede pasarla a `Completada` ni sacarla de `Completada`: eso se hace desde el Sprint (RN-2, RN-3). Los cambios entre `Pendiente` y `En progreso` sí se pueden hacer desde el backlog.
- RN-5: Un Sprint se puede cerrar en cualquier momento, aunque no haya terminado su período ni se hayan completado sus historias. Un Sprint ya cerrado no se puede volver a cerrar.
- RN-6: Al cerrar, el Sprint conserva **todas** sus historias asignadas (lo planificado) y cuáles se completaron en él. Las historias no completadas mantienen su estado (`Pendiente` o `En progreso`) y quedan libres para asignarse a otro Sprint (spec 05 RN-9).
- RN-7: Un Sprint cerrado no admite cambios: no se puede cambiar su Sprint Goal, asignar, quitar, completar ni deshacer historias, ni cerrarlo otra vez (completa la spec 05 RN-7).
- RN-8: El registro de un Sprint cerrado no cambia aunque después se modifiquen sus historias en el backlog: lo completado en ese Sprint queda como completado en ese Sprint.
- RN-9: Los Sprints anteriores son los Sprints `Cerrado` del proyecto, ordenados del más reciente al más antiguo por fecha de inicio.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- Las reglas RN-1 a RN-9 viven en el dominio y se desarrollan con TDD.
- Se accede a los Sprints a través del proyecto al que pertenecen (decisiones técnicas en `backend/CLAUDE.md`).

## Casos límite

- Cerrar un Sprint sin historias: se acepta.
- Cerrar un Sprint con todas sus historias completadas, o con ninguna: se acepta.
- Cerrar un Sprint antes de su fecha de fin: se acepta.
- Completar, deshacer y volver a completar la misma historia con el Sprint abierto: se acepta.
- Una historia no completada en un Sprint cerrado se asigna a un Sprint nuevo y se completa ahí: en el primero figura como no completada y en el segundo como completada.
- Desde el backlog, pasar a `En progreso` una historia de un Sprint abierto: se acepta. Pasarla a `Completada`: se rechaza.
- Desde el backlog, pasar a `Completada` una historia que no está en ningún Sprint abierto: se acepta (spec 04).
- Consultar anteriores sin Sprints cerrados: lista vacía.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Completar o deshacer una historia que no está en el Sprint | Conflicto: "La historia no está asignada a este Sprint" |
| Completar una historia ya completada en el Sprint | Conflicto: "La historia ya está completada" |
| Deshacer una historia que no está completada en el Sprint | Conflicto: "La historia no está completada" |
| Completar, deshacer o cerrar en un Sprint cerrado | Conflicto: "El Sprint está cerrado" |
| Cambiar desde el backlog el estado `Completada` de una historia en un Sprint abierto | Conflicto: "La historia está en un Sprint abierto: se completa desde el Sprint" |
| Proyecto inexistente | No encontrado: "El proyecto no existe" |
| Sprint inexistente (o de otro proyecto) | No encontrado: "El Sprint no existe" |
| Historia inexistente (o de otro proyecto) | No encontrado: "La historia no existe" |

## Criterios de aceptación

- CA-1: Se puede marcar como completada una historia asignada a un Sprint abierto; pasa a `Completada` y el Sprint registra que se completó en él. Mientras el Sprint esté abierto, se puede deshacer.
- CA-2: Se puede cerrar un Sprint en cualquier momento; conserva sus historias planificadas y cuáles se completaron, y las no completadas quedan libres para otro Sprint.
- CA-3: Un Sprint cerrado no admite cambios (Sprint Goal, asignar, quitar, completar, deshacer ni volver a cerrar).
- CA-4: Se pueden consultar los Sprints anteriores (cerrados) del proyecto, del más reciente al más antiguo, con sus historias y cuáles se completaron en cada uno.
- CA-5: Mientras una historia está en un Sprint abierto, el backlog no puede cambiar su estado `Completada`: se completa y se deshace desde el Sprint.
- CA-6: Operar sobre un proyecto, un Sprint o una historia que no existe informa cuál no existe, sin cambiar nada.

Relación con el Issue #9: "marcar una historia como completada" → CA-1; "cerrar un Sprint" (y qué pasa con las no completadas) → CA-2; "un Sprint cerrado no admite cambios" → CA-3; "consultar Sprints anteriores" → CA-4. CA-5 resuelve el punto pendiente de la spec 04 ("Completada" en dos lugares) y CA-6 los casos de error.

## Escenarios BDD

```gherkin
# language: es
Característica: Registrar historias completadas, cerrar Sprint y consultar Sprints anteriores
  Para dejar constancia de lo logrado en cada iteración
  Como integrante del equipo
  Quiero completar historias, cerrar Sprints y consultar los anteriores

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-31"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2

  Escenario: Completar una historia del Sprint (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Cuando completo en el Sprint la historia "Registrar usuario"
    Entonces en el Sprint la historia "Registrar usuario" está completada y "Ver catálogo" no
    Y la historia "Registrar usuario" del backlog está en estado "Completada"

  Escenario: Deshacer una historia completada (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando deshago en el Sprint la historia "Registrar usuario"
    Entonces en el Sprint ninguna historia está completada
    Y la historia "Registrar usuario" del backlog está en estado "En progreso"

  Escenario: Completar una historia que no está en el Sprint (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Ver catálogo"
    Cuando completo en el Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia no está asignada a este Sprint"

  Escenario: Completar dos veces la misma historia (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando completo en el Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está completada"

  Escenario: Deshacer una historia que no está completada (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando deshago en el Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia no está completada"

  Escenario: Cerrar un Sprint con historias sin terminar (CA-2)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando cierro el Sprint
    Entonces el Sprint está cerrado con las historias "Registrar usuario" y "Ver catálogo"
    Y en el Sprint la historia "Registrar usuario" está completada y "Ver catálogo" no

  Escenario: Las historias no completadas quedan libres para otro Sprint (CA-2)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Ver catálogo"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-15" al "2026-11-28"
    Cuando asigno al Sprint la historia "Ver catálogo"
    Entonces el Sprint tiene la historia "Ver catálogo"

  Escenario: Cerrar un Sprint vacío antes de su fecha de fin (CA-2)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Y que hoy es "2026-11-05"
    Cuando cierro el Sprint
    Entonces el Sprint está cerrado sin historias

  Esquema del escenario: Un Sprint cerrado no admite cambios (CA-3)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que cerré el Sprint
    Cuando <accion>
    Entonces se informa el error "El Sprint está cerrado"

    Ejemplos:
      | accion                                             |
      | cambio su objetivo a "Otro objetivo"               |
      | asigno al Sprint la historia "Ver catálogo"        |
      | quito del Sprint la historia "Registrar usuario"   |
      | completo en el Sprint la historia "Registrar usuario" |
      | cierro el Sprint                                   |

  Escenario: Consultar los Sprints anteriores (CA-4)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-15" al "2026-11-28" con la historia "Ver catálogo"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-29" al "2026-12-12"
    Cuando consulto los Sprints anteriores
    Entonces veo los Sprints anteriores que empiezan el "2026-11-15" y el "2026-11-01" en ese orden
    Y en el Sprint anterior que empieza el "2026-11-01" la historia "Registrar usuario" está completada

  Escenario: Consultar Sprints anteriores sin Sprints cerrados (CA-4)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando consulto los Sprints anteriores
    Entonces no hay Sprints anteriores

  Escenario: El registro del Sprint cerrado no cambia (CA-4)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Cuando la historia "Registrar usuario" pasa a estado "Pendiente"
    Entonces en el Sprint la historia "Registrar usuario" está completada

  Esquema del escenario: El backlog no cambia "Completada" de una historia en un Sprint abierto (CA-5)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando la historia "Registrar usuario" pasa a estado "<estado>"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | estado      | resultado                                                        |
      | En progreso | modificada                                                       |
      | Completada  | La historia está en un Sprint abierto: se completa desde el Sprint |

  Escenario: Cerrar un Sprint que no existe (CA-6)
    Cuando cierro el Sprint con identificador 999
    Entonces se informa el error "El Sprint no existe"

  Escenario: Completar una historia que no existe (CA-6)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando completo en el Sprint la historia con identificador 999
    Entonces se informa el error "La historia no existe"
```

Los pasos que no indican proyecto se refieren al de los Antecedentes, "Tienda Online". "El Sprint" es el último Sprint creado en el escenario. Los pasos "que existe un Sprint del … al …", "asigno al Sprint…", "quito del Sprint…", "cambio su objetivo a…" y "que hoy es…" son los de las specs 03 y 05; "la historia … pasa a estado …" es el de la spec 03.

## Pantalla

`frontend/src/paginas/Sprints.tsx` (sección "Sprints" del menú, la misma de la spec 05). En un Sprint abierto, cada historia tiene "Completar" o "Deshacer", y abajo está "Cerrar Sprint", que pide confirmación. Un Sprint cerrado se ve sin botones: lo planificado y cuáles se completaron en él (CA-4, Sprints anteriores). Todas las reglas y los mensajes de error son los del backend.

## Dependencias

- Spec 04 (backlog) y spec 05 (Sprints): esta historia es de la misma dueña que la 05 y se programa después de que la 05 esté mergeada (regla de independencia del `CLAUDE.md`). Agrega a la spec 04 la regla CA-5 y a la spec 05 el estado `Cerrado`.
- La usan: spec 10 (métricas de avance: Story Points planificados y completados por Sprint, velocidad) y spec 11 (defectos: Sprint de resolución).

## Tareas

- [x] `[P1]` Completar, deshacer y cerrar en el dominio (`Sprint`, `ServicioSprints`) con las reglas RN-1 a RN-9 (TDD).
- [x] `[P1]` Regla del backlog (RN-4) en `ServicioBacklog.Modificar`.
- [x] `[P1]` Registrar en la base qué historias se completaron en cada Sprint (`05_sprint_historias_completada.sql`).
- [x] `[P1]` Rutas de la API para completar, deshacer, cerrar y consultar anteriores.
- [x] `[P1]` Escenarios de `backend/features/cerrar-sprint.feature` con sus pasos.
- [x] `[P2]` Actualizar `backend/CLAUDE.md` (endpoints, specs relacionadas).

## Definition of Done

- [x] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [x] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [x] Criterios de aceptación verificados.
- [x] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [x] Spec en estado `implemented`.
