---
titulo: Administrar el Product Backlog
estado: implemented  # draft | in-review | approved | implemented
unidad: backend
issue: "#7"
sprint: 01
autor: Florencia Cabeza
aprobada_por: "Martina Batista"
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

- CA-1: Se puede crear una historia en el backlog de un proyecto con título, descripción, prioridad y criterios de aceptación; queda registrada con un identificador asignado por el sistema, en estado `Pendiente` y con los Story Points sin estimar.
- CA-2: Se pueden modificar el título, la descripción, la prioridad, los criterios de aceptación y el estado de una historia, con las mismas validaciones que al crear; si alguna falla, la historia no cambia.
- CA-3: Se puede eliminar una historia y deja de aparecer en el backlog.
- CA-4: El backlog de un proyecto se lista ordenado por prioridad (1 primero) y, a igual prioridad, por orden de creación; solo muestra las historias de ese proyecto, y se puede ver una historia por su identificador.
- CA-5: El título y la descripción son obligatorios y respetan su largo máximo (150 y 1000 caracteres); los criterios de aceptación son opcionales con un máximo de 2000.
- CA-6: La prioridad es obligatoria y es un entero mayor o igual a 1; el estado solo admite `Pendiente`, `En progreso` o `Completada`.
- CA-7: Operar sobre un proyecto o una historia que no existe informa cuál de los dos no existe, sin cambiar nada.

Relación con los criterios iniciales del Issue #7: "cada ítem tiene identificador, título, descripción, prioridad, estado, Story Points y criterios" → CA-1, CA-5 y CA-6; "crear, modificar, eliminar y listar" → CA-1 a CA-4; "ver ordenado por prioridad" → CA-4.

## Escenarios BDD

```gherkin
# language: es
Característica: Administrar el Product Backlog
  Para mantener priorizado el trabajo pendiente
  Como integrante del equipo
  Quiero crear, modificar, eliminar y listar las historias del backlog de un proyecto

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Crear una historia con todos sus datos (CA-1)
    Cuando creo en el backlog la historia "Registrar usuario" con descripción "Alta de clientes", prioridad 1 y criterios "Valida el email"
    Entonces la historia queda registrada con un identificador
    Y la historia tiene título "Registrar usuario", descripción "Alta de clientes", prioridad 1 y criterios "Valida el email"
    Y la historia está en estado "Pendiente" y sin estimar

  Escenario: Crear una historia sin criterios de aceptación (CA-1)
    Cuando creo en el backlog la historia "Ver catálogo" con descripción "Listado de productos", prioridad 2 y sin criterios
    Entonces la historia queda registrada con un identificador
    Y la historia no tiene criterios de aceptación

  Escenario: Modificar la prioridad y el estado de una historia (CA-2)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 3
    Cuando cambio su prioridad a 1 y su estado a "En progreso"
    Entonces la historia tiene prioridad 1 y estado "En progreso"

  Escenario: Una modificación inválida no cambia la historia (CA-2)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 3
    Cuando cambio su prioridad a 0
    Entonces se informa el error "La prioridad debe ser un número entero mayor o igual a 1"
    Y la historia tiene prioridad 3 y estado "Pendiente"

  Escenario: Eliminar una historia del backlog (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Cuando elimino la historia "Registrar usuario"
    Entonces el backlog tiene solo la historia "Ver catálogo"

  Escenario: Eliminar dos veces la misma historia (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando elimino la historia "Registrar usuario"
    Y elimino otra vez la historia "Registrar usuario"
    Entonces se informa el error "La historia no existe"

  Escenario: Listar el backlog ordenado por prioridad (CA-4)
    Dado que el backlog tiene la historia "Pagar pedido" con prioridad 3
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Cuando listo el backlog
    Entonces veo las historias "Registrar usuario", "Ver catálogo" y "Pagar pedido" en ese orden

  Escenario: A igual prioridad se respeta el orden de creación (CA-4)
    Dado que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Y que el backlog tiene la historia "Buscar producto" con prioridad 2
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando listo el backlog
    Entonces veo las historias "Registrar usuario", "Ver catálogo" y "Buscar producto" en ese orden

  Escenario: Listar un backlog vacío (CA-4)
    Cuando listo el backlog
    Entonces el backlog está vacío

  Escenario: Cada proyecto ve solo su backlog (CA-4)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando listo el backlog del proyecto "App Turnos"
    Entonces el backlog está vacío

  Escenario: Ver una historia por su identificador (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando veo esa historia por su identificador
    Entonces la historia tiene título "Registrar usuario" y prioridad 1

  Esquema del escenario: Validar título, descripción y criterios (CA-5)
    Cuando creo en el backlog una historia con un título de <titulo> caracteres, una descripción de <descripcion> caracteres y criterios de <criterios> caracteres
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | titulo | descripcion | criterios | resultado                                                          |
      | 150    | 1000        | 2000      | registrada                                                         |
      | 151    | 10          | 0         | El título no puede superar los 150 caracteres                      |
      | 0      | 10          | 0         | El título es obligatorio                                           |
      | 10     | 1001        | 0         | La descripción no puede superar los 1000 caracteres                |
      | 10     | 0           | 0         | La descripción es obligatoria                                      |
      | 10     | 10          | 2001      | Los criterios de aceptación no pueden superar los 2000 caracteres  |

  Escenario: Rechazar un título con solo espacios (CA-5)
    Cuando creo en el backlog la historia "   " con descripción "Alta de clientes", prioridad 1 y sin criterios
    Entonces se informa el error "El título es obligatorio"
    Y el backlog está vacío

  Esquema del escenario: Validar la prioridad (CA-6)
    Cuando creo en el backlog la historia "Registrar usuario" con descripción "Alta de clientes", prioridad "<prioridad>" y sin criterios
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | prioridad | resultado                                                  |
      | 1         | registrada                                                 |
      | 1000      | registrada                                                 |
      | 0         | La prioridad debe ser un número entero mayor o igual a 1   |
      | -2        | La prioridad debe ser un número entero mayor o igual a 1   |
      | 2.5       | La prioridad debe ser un número entero mayor o igual a 1   |
      | (vacío)   | La prioridad es obligatoria                                |

  Esquema del escenario: Validar el estado al modificar (CA-6)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando cambio su estado a "<estado>"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | estado      | resultado              |
      | Pendiente   | modificada             |
      | En progreso | modificada             |
      | Completada  | modificada             |
      | Terminada   | El estado no es válido |
      | pendiente   | El estado no es válido |

  Escenario: Crear una historia en un proyecto que no existe (CA-7)
    Cuando creo una historia en el backlog del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Ver una historia que no existe (CA-7)
    Cuando veo la historia con identificador 999
    Entonces se informa el error "La historia no existe"

  Escenario: Una historia de otro proyecto se trata como inexistente (CA-7)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando veo esa historia desde el proyecto "App Turnos"
    Entonces se informa el error "La historia no existe"
```

Los pasos que no indican proyecto ("el backlog", "creo en el backlog") se refieren al proyecto de los Antecedentes, "Tienda Online". En las tablas, `(vacío)` representa un texto vacío (Gherkin no puede escribirlo), igual que en la spec 01.

## Trazabilidad

Historia [#7](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/7) → esta spec → criterios → escenarios de `backend/features/product-backlog.feature` → tests → código Go.

| Criterio | Escenarios BDD | Tests Go | Código Go |
|---|---|---|---|
| CA-1 Crear | Crear una historia con todos sus datos · Crear una historia sin criterios | `TestNuevaHistoria`, `TestBacklogCrear`, `TestCrearHistoria` | `dominio/historia.go` (`NuevaHistoria`), `dominio/servicio_backlog.go` (`Crear`), `api/backlog.go` (`crear`) |
| CA-2 Modificar | Modificar la prioridad y el estado · Una modificación inválida no cambia la historia | `TestBacklogModificar`, `TestBacklogModificarInvalidoNoCambia`, `TestModificarYVerHistoria` | `dominio/servicio_backlog.go` (`Modificar`), `api/backlog.go` (`modificar`) |
| CA-3 Eliminar | Eliminar una historia · Eliminar dos veces la misma historia | `TestBacklogEliminar`, `TestEliminarHistoria`, `TestRepositorioHistorias` | `dominio/servicio_backlog.go` (`Eliminar`), `postgres/historias.go`, `api/backlog.go` (`eliminar`) |
| CA-4 Listar ordenado y ver | Listar ordenado por prioridad · A igual prioridad, orden de creación · Backlog vacío · Cada proyecto ve solo su backlog · Ver una historia | `TestOrdenarBacklog`, `TestBacklogListar`, `TestBacklogObtener`, `TestListarBacklog`, `TestListarHistoriasVacio` | `dominio/historia.go` (`OrdenarBacklog`), `dominio/servicio_backlog.go` (`Listar`, `Obtener`) |
| CA-5 Textos | Validar título, descripción y criterios · Rechazar un título con solo espacios | `TestNuevaHistoria`, `TestBacklogCrearInvalidoNoRegistra`, `TestCrearHistoriaInvalida` | `dominio/historia.go` (`NuevaHistoria`), `dominio/errores.go` |
| CA-6 Prioridad y estado | Validar la prioridad · Validar el estado al modificar | `TestParsearPrioridad`, `TestParsearEstado`, `TestNuevaHistoria`, `TestCrearHistoriaInvalida` | `dominio/historia.go` (`ParsearPrioridad`, `ParsearEstado`) |
| CA-7 Inexistentes | Crear en un proyecto que no existe · Ver una historia que no existe · Historia de otro proyecto | `TestBacklogProyectoInexistente`, `TestBacklogObtener`, `TestBacklogInexistente`, `TestRepositorioHistorias` | `dominio/servicio_backlog.go` (`verificarProyecto`), `memoria/historias.go`, `postgres/historias.go`, `api/proyectos.go` (`responderError`) |

## Pantalla

`frontend/src/paginas/Backlog.tsx` (sección "Backlog" del menú): lista ordenada por prioridad, un solo formulario para crear y modificar, y eliminar con confirmación. No valida nada: muestra los mensajes de error del backend. Los Story Points se estiman en la misma pantalla (spec 07).

## Dependencias

- Spec 01 (proyectos): el backlog pertenece a un proyecto existente.
- La usan: spec 05 (asignar historias a un Sprint), spec 06 (completar historias y cerrar el Sprint), spec 07 (Story Points), spec 08 (Planning Poker) y spec 09 (registro de esfuerzo; agrega que no se puede eliminar una historia con esfuerzo registrado, RN-10), spec 11 (defectos; agrega que no se puede eliminar una historia con defectos registrados, RN-14) y spec 12 (agrega las horas estimadas de la historia, que se cargan aparte, como los Story Points). Si esas specs cambian el estado o los Story Points de forma automática, se actualiza esta.

## Tareas

- [x] `[P1]` Entidad Historia en `backend/dominio/` con las validaciones RN-2 a RN-9 (TDD, un ciclo por regla).
- [x] `[P1]` Tabla de historias en `backend/postgres/tablas/02_historias.sql`, con referencia al proyecto.
- [x] `[P1]` Acceso a datos del backlog en `backend/postgres/` (crear, modificar, eliminar, listar, ver una).
- [x] `[P1]` Rutas de la API bajo `/api/proyectos/{id}/backlog` en `backend/api/`.
- [x] `[P1]` Escenarios de `backend/features/product-backlog.feature` con sus pasos en `features/backlog_test.go`.
- [x] `[P2]` Actualizar `backend/CLAUDE.md` (endpoints disponibles, specs relacionadas, trampas).

## Definition of Done

- [x] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [x] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [x] Criterios de aceptación verificados.
- [x] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [x] Spec en estado `implemented`.
