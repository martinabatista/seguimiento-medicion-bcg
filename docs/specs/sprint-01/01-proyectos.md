---
titulo: Crear y modificar proyectos
estado: in-review    # draft | in-review | approved | implemented
unidad: backend
issue: "#4"
sprint: 01
autor: Martina Batista
aprobada_por: ""
fecha: 2026-10-20
specs_relacionadas: []
---

# Spec 01 — Crear y modificar proyectos

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero crear y modificar un proyecto con sus fechas de inicio y finalización, para administrar su seguimiento. ([#4](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/4))

## Objetivo

Permitir dar de alta y mantener los proyectos sobre los que se realiza toda la estimación, el seguimiento y la medición. Es la base del resto del sistema: integrantes, backlog, Sprints, esfuerzo y métricas pertenecen a un proyecto.

Origen: Consigna — Requerimiento 1 (Gestión de proyectos).

## Alcance

**Entra:**
- Crear un proyecto con nombre, descripción, fecha de inicio y fecha de finalización.
- Modificar los datos de un proyecto existente.
- Listar todos los proyectos.
- Ver un proyecto por su identificador.

**NO entra:**
- Eliminar proyectos (la consigna no lo pide).
- El estado del proyecto y su avance (spec 03, HU-03).
- Integrantes, backlog y Sprints del proyecto (specs 02, 04 y 05).
- Pantallas del frontend (todavía no existe; se hacen cuando se arme la interfaz).

## Entradas

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Nombre | texto | Sí | No puede estar vacío ni tener solo espacios. Sin largo máximo. Puede repetirse entre proyectos. |
| Descripción | texto | Sí | No puede estar vacía ni tener solo espacios. Máximo 500 caracteres. |
| Fecha de inicio | fecha (AAAA-MM-DD) | Sí | Tiene que ser una fecha válida del calendario. |
| Fecha de finalización | fecha (AAAA-MM-DD) | Sí | Tiene que ser una fecha válida y **posterior** a la fecha de inicio. |

Para modificar, además se indica el **identificador** del proyecto. Para ver uno, solo el identificador.

## Salidas esperadas

- **Crear:** el proyecto queda registrado y se devuelve con su identificador asignado por el sistema y sus cuatro datos.
- **Modificar:** se devuelve el proyecto con los datos actualizados.
- **Listar:** se devuelven todos los proyectos en orden de creación. Si no hay ninguno, la lista está vacía.
- **Ver uno:** se devuelve el proyecto con su identificador y sus cuatro datos.
- **Error:** no se crea ni se modifica nada, y se informa un mensaje que explica qué dato está mal o que el proyecto no existe.

## Reglas de negocio

- RN-1: El nombre es obligatorio: no puede estar vacío ni tener solo espacios. No tiene largo máximo y puede repetirse entre proyectos.
- RN-2: La descripción es obligatoria: no puede estar vacía ni tener solo espacios, y tiene como máximo 500 caracteres.
- RN-3: Las fechas de inicio y de finalización son obligatorias y tienen que ser fechas válidas.
- RN-4: La fecha de finalización tiene que ser posterior a la de inicio: no puede ser igual ni anterior.
- RN-5: Al modificar se aplican las mismas reglas que al crear. Si alguna falla, el proyecto queda como estaba.
- RN-6: El identificador lo asigna el sistema al crear el proyecto y no se puede modificar.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- Las validaciones de RN-1 a RN-5 viven en el dominio y se desarrollan con TDD.
- Las fechas se reciben y se devuelven en formato AAAA-MM-DD.

## Casos límite

- Nombre muy largo (por ejemplo, 300 caracteres): se acepta, porque no hay largo máximo.
- Nombre igual al de otro proyecto: se acepta.
- Nombre o descripción con solo espacios: se rechaza, igual que vacío.
- Descripción de exactamente 500 caracteres: se acepta. Con 501: se rechaza.
- Fecha de finalización un día después de la de inicio: se acepta (la duración mínima es un día).
- Fecha de finalización igual a la de inicio: se rechaza.
- Fecha inexistente (por ejemplo, 2026-02-30): se rechaza.
- Listar sin proyectos cargados: devuelve una lista vacía, no un error.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Nombre vacío o con solo espacios | Dato inválido: "El nombre es obligatorio" |
| Descripción vacía o con solo espacios | Dato inválido: "La descripción es obligatoria" |
| Descripción de más de 500 caracteres | Dato inválido: "La descripción no puede superar los 500 caracteres" |
| Falta la fecha de inicio o la de finalización | Dato inválido: "Las fechas de inicio y finalización son obligatorias" |
| Fecha con formato inválido o inexistente | Dato inválido: "La fecha no es válida" |
| Fecha de finalización igual o anterior a la de inicio | Dato inválido: "La fecha de finalización debe ser posterior a la de inicio" |
| Modificar o ver un proyecto que no existe | No encontrado: "El proyecto no existe" |

## Criterios de aceptación

- CA-1: Se puede crear un proyecto con nombre, descripción, fecha de inicio y fecha de finalización, y queda registrado con un identificador asignado por el sistema.
- CA-2: Se pueden modificar los datos de un proyecto existente, con las mismas validaciones que al crear; si alguna falla, el proyecto no cambia.
- CA-3: Las fechas de inicio y finalización son obligatorias, y la de finalización tiene que ser posterior a la de inicio (no igual ni anterior).
- CA-4: El nombre es obligatorio; no tiene largo máximo y puede repetirse.
- CA-5: La descripción es obligatoria y tiene como máximo 500 caracteres.
- CA-6: Se pueden listar todos los proyectos en orden de creación y ver uno por su identificador.
- CA-7: Ver o modificar un proyecto que no existe informa que el proyecto no existe.

## Escenarios BDD

```gherkin
# language: es
Característica: Crear y modificar proyectos
  Para administrar el seguimiento de un proyecto de software
  Como integrante del equipo
  Quiero crear, modificar y consultar proyectos

  Escenario: Crear un proyecto con todos sus datos (CA-1)
    Cuando creo un proyecto con nombre "Tienda Online", descripción "Venta de ropa", inicio "2026-11-01" y fin "2026-12-15"
    Entonces el proyecto queda registrado con un identificador
    Y el proyecto tiene nombre "Tienda Online", descripción "Venta de ropa", inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Crear un proyecto que dura un solo día (CA-1)
    Cuando creo un proyecto con nombre "Hackatón", descripción "Evento de un día", inicio "2026-11-01" y fin "2026-11-02"
    Entonces el proyecto queda registrado con un identificador

  Escenario: Modificar el nombre de un proyecto (CA-2)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Cuando cambio su nombre a "Tienda Online v2"
    Entonces el proyecto tiene nombre "Tienda Online v2"

  Escenario: Modificar las fechas de un proyecto (CA-2)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Cuando cambio sus fechas a inicio "2026-11-10" y fin "2027-01-31"
    Entonces el proyecto tiene inicio "2026-11-10" y fin "2027-01-31"

  Escenario: Una modificación inválida no cambia el proyecto (CA-2)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Cuando cambio sus fechas a inicio "2026-11-10" y fin "2026-11-05"
    Entonces se informa el error "La fecha de finalización debe ser posterior a la de inicio"
    Y el proyecto tiene inicio "2026-11-01" y fin "2026-12-15"

  Esquema del escenario: Rechazar fechas inválidas (CA-3)
    Cuando creo un proyecto con nombre "Tienda Online", descripción "Venta de ropa", inicio "<inicio>" y fin "<fin>"
    Entonces se informa el error "<mensaje>"
    Y no se registra ningún proyecto

    Ejemplos:
      | inicio     | fin        | mensaje                                                     |
      | 2026-11-10 | 2026-11-10 | La fecha de finalización debe ser posterior a la de inicio  |
      | 2026-11-10 | 2026-11-01 | La fecha de finalización debe ser posterior a la de inicio  |
      |            | 2026-11-01 | Las fechas de inicio y finalización son obligatorias        |
      | 2026-11-10 |            | Las fechas de inicio y finalización son obligatorias        |
      | 2026-02-30 | 2026-03-10 | La fecha no es válida                                       |

  Esquema del escenario: Rechazar un proyecto sin nombre (CA-4)
    Cuando creo un proyecto con nombre "<nombre>", descripción "Venta de ropa", inicio "2026-11-01" y fin "2026-12-15"
    Entonces se informa el error "El nombre es obligatorio"
    Y no se registra ningún proyecto

    Ejemplos:
      | nombre |
      |        |
      |        |

  Escenario: Aceptar un nombre muy largo (CA-4)
    Cuando creo un proyecto con un nombre de 300 caracteres, descripción "Venta de ropa", inicio "2026-11-01" y fin "2026-12-15"
    Entonces el proyecto queda registrado con un identificador

  Escenario: Aceptar un nombre repetido (CA-4)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Cuando creo un proyecto con nombre "Tienda Online", descripción "Otra tienda", inicio "2027-01-01" y fin "2027-02-01"
    Entonces hay 2 proyectos registrados

  Esquema del escenario: Validar el largo de la descripción (CA-5)
    Cuando creo un proyecto con nombre "Tienda Online", una descripción de <largo> caracteres, inicio "2026-11-01" y fin "2026-12-15"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                            |
      | 500   | registrado                                           |
      | 501   | La descripción no puede superar los 500 caracteres   |
      | 0     | La descripción es obligatoria                        |

  Escenario: Listar proyectos en orden de creación (CA-6)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Y que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Cuando listo los proyectos
    Entonces veo los proyectos "Tienda Online" y "App Turnos" en ese orden

  Escenario: Listar sin proyectos cargados (CA-6)
    Cuando listo los proyectos
    Entonces la lista de proyectos está vacía

  Escenario: Ver un proyecto por su identificador (CA-6)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Cuando veo ese proyecto por su identificador
    Entonces el proyecto tiene nombre "Tienda Online", descripción "Venta de ropa", inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Ver un proyecto que no existe (CA-7)
    Cuando veo el proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Modificar un proyecto que no existe (CA-7)
    Cuando cambio el nombre del proyecto con identificador 999 a "Fantasma"
    Entonces se informa el error "El proyecto no existe"
```

En el `Esquema del escenario` del CA-4, el primer ejemplo es un nombre vacío y el segundo un nombre con solo espacios (`"   "`); al pasarlo al `.feature` se escriben explícitamente.

## Dependencias

Ninguna (base del resto del sistema). Las specs 02, 03, 04 y 05 dependen de esta.

## Tareas

- [ ] `[P1]` Entidad Proyecto en `backend/dominio/` con las validaciones RN-1 a RN-5 (TDD, un ciclo por regla).
- [ ] `[P1]` Tabla de proyectos en `backend/postgres/esquema/01_proyectos.sql`.
- [ ] `[P1]` Agregar el driver `pgx` y el acceso a datos de proyectos en `backend/postgres/`.
- [ ] `[P1]` Rutas de la API para crear, modificar, listar y ver proyectos en `backend/api/`.
- [ ] `[P1]` Escenarios de `backend/features/proyectos.feature` con sus pasos en `features_test.go`.
- [ ] `[P2]` Actualizar `backend/CLAUDE.md` (endpoints disponibles, specs relacionadas).

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
