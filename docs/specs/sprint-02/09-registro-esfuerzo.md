---
titulo: Registrar esfuerzo
estado: implemented  # draft | in-review | approved | implemented
unidad: backend
issue: "#13"
sprint: 02
autor: Luana Guajardo
aprobada_por: "Martina Batista"
fecha: 2026-10-23
specs_relacionadas: ["01-proyectos", "02-integrantes", "03-estado-proyecto", "04-product-backlog"]
---

# Spec 09 — Registrar esfuerzo

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero registrar las horas que trabajé en una historia, para comparar después el esfuerzo estimado con el real. ([#13](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/13))

## Objetivo

Registrar el esfuerzo real que el equipo dedica a cada historia del backlog: quién trabajó, qué día, en qué actividad y cuántas horas. Es la base de las métricas de horas reales y de desviación entre esfuerzo estimado y real (spec 12) y del esfuerzo registrado del reporte (spec 14).

Origen: Consigna — Requerimiento 5 (Registro de esfuerzo).

## Alcance

**Entra:**
- Registrar esfuerzo sobre una historia del backlog: integrante, fecha, actividad realizada y horas trabajadas.
- Consultar el esfuerzo registrado de una historia, con el total de horas.
- Eliminar un registro de esfuerzo (para corregir errores de carga).
- Impedir eliminar del backlog una historia que tiene esfuerzo registrado, para no perder esas horas.

**NO entra:**
- **Tareas:** la consigna habla de "historia o tarea", pero el sistema no tiene tareas; el esfuerzo se registra sobre la historia. Si el equipo agrega tareas más adelante, esta spec se amplía.
- Modificar un registro: para corregirlo, se elimina y se vuelve a cargar.
- Horas estimadas, horas reales totales y desviación (spec 12).
- Consultar el esfuerzo por integrante o por Sprint (lo pueden agregar las specs 12 y 14).
- Pantallas del frontend (Sprint 2, issue #15).

## Entradas

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |
| Historia | identificador | Sí | Tiene que ser una historia del backlog de ese proyecto (spec 04), en cualquier estado. |
| Integrante | identificador | Sí | Tiene que ser un integrante de ese proyecto (spec 02), con cualquier rol. |
| Fecha | fecha (AAAA-MM-DD) | Sí | Fecha válida, entre la fecha de inicio y la de finalización del proyecto (ambas incluidas), y no posterior a hoy. |
| Actividad realizada | texto | Sí | No puede estar vacía ni tener solo espacios. Máximo 300 caracteres. |
| Horas trabajadas | número decimal | Sí | Mayor a 0 y como máximo 24, con hasta 2 decimales (por ejemplo, 1,5 o 0,25). |

Para eliminar, además se indica el **identificador** del registro.

## Salidas esperadas

- **Registrar:** el registro queda guardado y se devuelve con su identificador asignado por el sistema, la historia, el integrante (identificador y nombre), la fecha, la actividad y las horas.
- **Consultar por historia:** se devuelven los registros de la historia ordenados por fecha (los más antiguos primero; a igual fecha, en orden de carga) y el **total de horas** registradas. Si no hay registros, la lista está vacía y el total es 0.
- **Eliminar:** el registro deja de existir y ya no suma al total.
- **Error:** no se registra ni se elimina nada, y se informa un mensaje que explica qué dato está mal o qué no existe.

## Reglas de negocio

- RN-1: El esfuerzo se registra sobre una historia existente del proyecto (spec 04). La historia puede estar en cualquier estado, también `Completada` (las horas se pueden cargar después de terminarla).
- RN-2: El integrante tiene que pertenecer al proyecto (spec 02). Puede tener cualquier rol.
- RN-3: La fecha es obligatoria, tiene que estar entre la fecha de inicio y la de finalización del proyecto (ambas incluidas) y no puede ser posterior a hoy. La fecha de hoy la toma el sistema, igual que en la spec 03.
- RN-4: La actividad realizada es obligatoria: no puede estar vacía ni tener solo espacios, y tiene como máximo 300 caracteres.
- RN-5: Las horas trabajadas son obligatorias, mayores a 0 y como máximo 24 por registro, con hasta 2 decimales.
- RN-6: El total de horas de una historia es la suma de las horas de todos sus registros, y conserva los 2 decimales (no se redondea).
- RN-7: Los registros de una historia se listan por fecha ascendente; a igual fecha, en el orden en que se cargaron.
- RN-8: El identificador del registro lo asigna el sistema. Un registro de otra historia o de otro proyecto se trata como inexistente.
- RN-9: Si hay más de un dato inválido, se informa el primero en este orden: integrante, fecha, actividad, horas. Antes que los datos se verifica que existan el proyecto y la historia.
- RN-10: No se puede eliminar del backlog (spec 04) una historia que tiene esfuerzo registrado: primero hay que eliminar sus registros. Así las horas reales no se pierden sin querer.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- Las validaciones (RN-3 a RN-5), la suma (RN-6) y el orden (RN-7) viven en el dominio y se desarrollan con TDD. La fecha de hoy se le pasa al dominio (reloj inyectado, como en la spec 03) para poder probarla con fechas fijas.
- Las horas se guardan sin perder precisión: 0,1 + 0,2 tiene que sumar exactamente 0,3.
- La historia y el integrante se toman de las specs 02 y 04, ya implementadas: esta historia no depende de ninguna otra pendiente.

## Casos límite

- Horas exactamente 24: se acepta. 24,01: se rechaza.
- Horas 0,01 (el mínimo con 2 decimales): se acepta. 0: se rechaza. Negativas: se rechazan.
- Horas con 3 decimales (1,125): se rechaza.
- Fecha igual a la de inicio o a la de fin del proyecto: se acepta.
- Fecha igual a hoy: se acepta. Mañana: se rechaza.
- Fecha anterior al inicio del proyecto o posterior a su fin: se rechaza.
- Actividad de exactamente 300 caracteres: se acepta. Con 301: se rechaza.
- Historia completada: se acepta registrar esfuerzo.
- Varios registros del mismo integrante el mismo día en la misma historia: se aceptan (se suman al total).
- Historia sin registros: lista vacía y total 0, no es un error.
- Total con decimales: 1,5 + 2,25 = 3,75.
- Eliminar el último registro de una historia: la historia vuelve a poder eliminarse del backlog.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Falta el integrante | Dato inválido: "El integrante es obligatorio" |
| Falta la fecha | Dato inválido: "La fecha es obligatoria" |
| Fecha con formato inválido o inexistente | Dato inválido: "La fecha no es válida" (spec 01) |
| Fecha fuera de las fechas del proyecto | Dato inválido: "La fecha debe estar dentro de las fechas del proyecto" |
| Fecha posterior a hoy | Dato inválido: "La fecha no puede ser posterior a hoy" |
| Actividad vacía o con solo espacios | Dato inválido: "La actividad es obligatoria" |
| Actividad de más de 300 caracteres | Dato inválido: "La actividad no puede superar los 300 caracteres" |
| Faltan las horas | Dato inválido: "Las horas son obligatorias" |
| Horas que no son un número mayor a 0 y hasta 24, o con más de 2 decimales | Dato inválido: "Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales" |
| El proyecto no existe | No encontrado: "El proyecto no existe" |
| La historia no existe o es de otro proyecto | No encontrado: "La historia no existe" |
| El integrante no existe o es de otro proyecto | No encontrado: "El integrante no existe" (spec 08) |
| Eliminar un registro que no existe o es de otra historia | No encontrado: "El registro de esfuerzo no existe" |
| Eliminar del backlog una historia con esfuerzo registrado | Conflicto: "La historia tiene esfuerzo registrado" |

## Criterios de aceptación

- CA-1: Se puede registrar esfuerzo con integrante, fecha, actividad realizada y horas trabajadas, y queda guardado con un identificador asignado por el sistema.
- CA-2: El esfuerzo se registra sobre una historia del backlog del proyecto (en cualquier estado), por un integrante del proyecto.
- CA-3: Las horas deben ser mayores a cero y como máximo 24, con hasta 2 decimales.
- CA-4: Se puede consultar el esfuerzo registrado por historia, ordenado por fecha, con el total de horas.
- CA-5: La fecha es obligatoria, está dentro de las fechas del proyecto y no es posterior a hoy; la actividad es obligatoria con un máximo de 300 caracteres.
- CA-6: Se puede eliminar un registro de esfuerzo, y deja de sumar al total.
- CA-7: No se puede eliminar del backlog una historia con esfuerzo registrado.
- CA-8: Operar con un proyecto, una historia, un integrante o un registro que no existe (o es de otro proyecto o historia) informa cuál no existe, sin registrar ni eliminar nada.

Relación con los criterios iniciales del Issue #13: "integrante, fecha, actividad y horas" → CA-1 y CA-5; "sobre una historia o tarea" → CA-2 (solo historias, ver Alcance); "horas mayores a cero" → CA-3; "consultar por historia" → CA-4. CA-6 y CA-7 surgen de poder corregir errores sin perder horas; CA-8 cubre los casos de error.

## Escenarios BDD

```gherkin
# language: es
Característica: Registrar esfuerzo
  Para comparar después el esfuerzo estimado con el real
  Como integrante del equipo
  Quiero registrar las horas que trabajé en una historia

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Y que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que hoy es "2026-11-20"

  Escenario: Registrar esfuerzo con todos sus datos (CA-1)
    Cuando registro 2,5 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Diseño del formulario"
    Entonces el registro de esfuerzo queda guardado con un identificador
    Y el registro tiene integrante "Ana Pérez", fecha "2026-11-18", actividad "Diseño del formulario" y 2,5 horas

  Escenario: Registrar esfuerzo en una historia completada (CA-2)
    Dado que la historia "Registrar usuario" pasa a estado "Completada"
    Cuando registro 1 hora de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Pruebas finales"
    Entonces el registro de esfuerzo queda guardado con un identificador

  Escenario: Cualquier rol puede registrar esfuerzo (CA-2)
    Dado que el proyecto tiene al integrante "Bruno Díaz" con email "bruno@mail.com" y rol "Scrum Master"
    Cuando registro 1 hora de "Bruno Díaz" el "2026-11-18" en "Registrar usuario" con la actividad "Refinamiento"
    Entonces el registro de esfuerzo queda guardado con un identificador

  Esquema del escenario: Validar las horas (CA-3)
    Cuando registro <horas> horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Desarrollo"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | horas   | resultado                                                                 |
      | 0,01    | registrado                                                                |
      | 1,5     | registrado                                                                |
      | 24      | registrado                                                                |
      | 24,01   | Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales |
      | 0       | Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales |
      | -2      | Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales |
      | 1,125   | Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales |
      | (vacío) | Las horas son obligatorias                                                |

  Escenario: Consultar el esfuerzo de una historia ordenado por fecha (CA-4)
    Dado que registré 2 horas de "Ana Pérez" el "2026-11-19" en "Registrar usuario" con la actividad "Backend"
    Y que registré 1,5 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Diseño"
    Y que registré 0,25 horas de "Ana Pérez" el "2026-11-19" en "Registrar usuario" con la actividad "Revisión"
    Cuando consulto el esfuerzo de "Registrar usuario"
    Entonces veo las actividades "Diseño", "Backend" y "Revisión" en ese orden
    Y el total de horas es 3,75

  Escenario: Consultar una historia sin esfuerzo registrado (CA-4)
    Cuando consulto el esfuerzo de "Registrar usuario"
    Entonces no hay registros de esfuerzo
    Y el total de horas es 0

  Escenario: El total no pierde precisión con decimales (CA-4)
    Dado que registré 0,1 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "A"
    Y que registré 0,2 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "B"
    Cuando consulto el esfuerzo de "Registrar usuario"
    Entonces el total de horas es 0,3

  Esquema del escenario: Validar la fecha (CA-5)
    Cuando registro 1 hora de "Ana Pérez" el "<fecha>" en "Registrar usuario" con la actividad "Desarrollo"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | fecha      | resultado                                             |
      | 2026-11-01 | registrado                                            |
      | 2026-11-20 | registrado                                            |
      | 2026-11-21 | La fecha no puede ser posterior a hoy                 |
      | 2026-10-31 | La fecha debe estar dentro de las fechas del proyecto |
      | 2026-02-30 | La fecha no es válida                                 |
      | (vacío)    | La fecha es obligatoria                               |

  Escenario: Registrar el último día del proyecto (CA-5)
    Dado que hoy es "2026-12-20"
    Cuando registro 1 hora de "Ana Pérez" el "2026-12-15" en "Registrar usuario" con la actividad "Cierre"
    Entonces el registro de esfuerzo queda guardado con un identificador

  Escenario: Rechazar una fecha posterior al fin del proyecto (CA-5)
    Dado que hoy es "2026-12-20"
    Cuando registro 1 hora de "Ana Pérez" el "2026-12-16" en "Registrar usuario" con la actividad "Cierre"
    Entonces se informa el error "La fecha debe estar dentro de las fechas del proyecto"

  Esquema del escenario: Validar la actividad (CA-5)
    Cuando registro 1 hora de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con una actividad de <largo> caracteres
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                        |
      | 300   | registrado                                       |
      | 301   | La actividad no puede superar los 300 caracteres |
      | 0     | La actividad es obligatoria                      |

  Escenario: Rechazar una actividad con solo espacios (CA-5)
    Cuando registro 1 hora de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "   "
    Entonces se informa el error "La actividad es obligatoria"
    Y no hay esfuerzo registrado en "Registrar usuario"

  Escenario: Si hay varios datos inválidos se informa el primero (CA-5)
    Cuando registro 0 horas de "Ana Pérez" el "2026-11-21" en "Registrar usuario" con la actividad ""
    Entonces se informa el error "La fecha no puede ser posterior a hoy"

  Escenario: Eliminar un registro de esfuerzo (CA-6)
    Dado que registré 2 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Backend"
    Y que registré 1 hora de "Ana Pérez" el "2026-11-19" en "Registrar usuario" con la actividad "Pruebas"
    Cuando elimino el registro de la actividad "Backend"
    Y consulto el esfuerzo de "Registrar usuario"
    Entonces veo solo la actividad "Pruebas"
    Y el total de horas es 1

  Escenario: Eliminar dos veces el mismo registro (CA-6)
    Dado que registré 2 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Backend"
    Cuando elimino el registro de la actividad "Backend"
    Y elimino otra vez el registro de la actividad "Backend"
    Entonces se informa el error "El registro de esfuerzo no existe"

  Escenario: No se puede eliminar del backlog una historia con esfuerzo (CA-7)
    Dado que registré 2 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Backend"
    Cuando elimino la historia "Registrar usuario"
    Entonces se informa el error "La historia tiene esfuerzo registrado"
    Y el backlog tiene solo la historia "Registrar usuario"

  Escenario: Sin registros, la historia se puede eliminar del backlog (CA-7)
    Dado que registré 2 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Backend"
    Y que eliminé el registro de la actividad "Backend"
    Cuando elimino la historia "Registrar usuario"
    Entonces el backlog está vacío

  Escenario: Registrar en una historia que no existe (CA-8)
    Cuando registro 1 hora de "Ana Pérez" el "2026-11-18" en la historia con identificador 999 con la actividad "Desarrollo"
    Entonces se informa el error "La historia no existe"

  Escenario: Registrar con un integrante de otro proyecto (CA-8)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto "App Turnos" tiene al integrante "Carla Gómez" con email "carla@mail.com" y rol "Developer"
    Cuando registro 1 hora de "Carla Gómez" el "2026-11-18" en "Registrar usuario" con la actividad "Desarrollo"
    Entonces se informa el error "El integrante no existe"
    Y no hay esfuerzo registrado en "Registrar usuario"

  Escenario: Registrar sin indicar el integrante (CA-8)
    Cuando registro 1 hora sin integrante el "2026-11-18" en "Registrar usuario" con la actividad "Desarrollo"
    Entonces se informa el error "El integrante es obligatorio"

  Escenario: Consultar el esfuerzo en un proyecto que no existe (CA-8)
    Cuando consulto el esfuerzo de la historia 1 del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"
```

Los pasos que no indican proyecto se refieren al de los Antecedentes, "Tienda Online". Las horas se escriben con coma decimal, como se leen; en la API viajan como número (`2.5`). En las tablas, `(vacío)` representa un texto vacío, igual que en las specs 01 y 04. El paso "que hoy es …" es el de la spec 03 y "elimino la historia …" el de la spec 04.

## Pantalla

`frontend/src/paginas/Esfuerzo.tsx` (sección "Esfuerzo" del menú): se elige una historia del backlog y se ven sus registros en el orden del backend, con el total de horas, y un botón para eliminar cada uno (con confirmación). El formulario carga integrante (desplegable con los del proyecto), fecha, actividad y horas. No valida nada: el campo de horas solo acepta números (`step="any"`, sin límite de decimales ni de rango) y los errores son los mensajes del backend (por ejemplo, "La fecha no puede ser posterior a hoy"). Si falta el integrante o las horas no se mandan, así el backend informa cuál falta.

## Dependencias

- Spec 01 (proyectos): fechas del proyecto. Spec 02 (integrantes) y spec 04 (backlog): el integrante y la historia. Las tres están implementadas: esta historia no espera a ninguna otra.
- Spec 03: reusa el reloj inyectado para la fecha de hoy.
- Spec 04: esta spec agrega una regla al eliminar una historia (RN-10, CA-7), igual que la spec 05 agregó la del Sprint abierto.
- La usan: spec 12 (horas reales y desviación; ahí se definen las horas estimadas) y spec 14 (esfuerzo registrado en el reporte).

## Tareas

- [x] `[P1]` Registro de esfuerzo en `backend/dominio/` con las validaciones RN-3 a RN-5, la suma (RN-6) y el orden (RN-7) (TDD).
- [x] `[P1]` Casos de uso registrar, consultar y eliminar, con la verificación de proyecto, historia e integrante (RN-1, RN-2, RN-8, RN-9), y repositorio en memoria.
- [x] `[P1]` Regla del backlog: no eliminar una historia con esfuerzo (RN-10).
- [x] `[P1]` Tabla `08_esfuerzo.sql`, repositorio PostgreSQL y rutas REST (`/api/proyectos/{id}/backlog/{historiaID}/esfuerzo`), con escenarios de `backend/features/registro-esfuerzo.feature`.
- [x] `[P2]` Actualizar `backend/CLAUDE.md` (endpoints, specs relacionadas, trampas).

## Definition of Done

- [x] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [x] Escenarios BDD automatizados y en verde.
- [x] CI en verde (gofmt / vet / tests / lint / build).
- [x] Criterios de aceptación verificados.
- [x] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [x] Spec en estado `implemented`.
