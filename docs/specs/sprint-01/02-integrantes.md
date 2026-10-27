---
titulo: Registrar integrantes del proyecto
estado: implemented  # draft | in-review | approved | implemented
unidad: backend
issue: "#5"
sprint: 01
autor: Luana Guajardo
aprobada_por: "Martina Batista"
fecha: 2026-10-22
specs_relacionadas: ["01-proyectos"]
---

# Spec 02 — Registrar integrantes del proyecto

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero registrar a los integrantes de un proyecto, para poder asignarles estimaciones, esfuerzo y trabajo. ([#5](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/5))

## Objetivo

Saber quiénes trabajan en cada proyecto y con qué rol. Es la base del Planning Poker (spec 08: cada integrante vota) y del registro de esfuerzo (spec 09: cada registro de horas es de un integrante).

Origen: Consigna — Requerimiento 1 (Gestión de proyectos: registrar integrantes).

## Alcance

**Entra:**
- Agregar un integrante a un proyecto con nombre, email y rol.
- Listar los integrantes de un proyecto.
- Ver un integrante por su identificador.

**NO entra:**
- Modificar o quitar integrantes: cuando haya esfuerzo y estimaciones registrados a su nombre (specs 08 y 09), quitar a alguien obliga a decidir qué pasa con esos datos. Se agrega en una historia aparte si el equipo lo necesita.
- Usuarios, contraseñas o inicio de sesión: el integrante es un dato del proyecto, no una cuenta del sistema.
- Limitar cuántos integrantes tienen cada rol (por ejemplo, un único Product Owner): se puede cargar cualquier cantidad de cada rol.
- Una misma persona con más de un rol en el mismo proyecto: tiene un único rol.
- Pantallas del frontend (todavía no existe; se hacen cuando se arme la interfaz).

## Entradas

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |
| Nombre | texto | Sí | No puede estar vacío ni tener solo espacios. Máximo 100 caracteres. Puede repetirse. |
| Email | texto | Sí | No puede estar vacío ni tener solo espacios. Formato de email válido. Máximo 254 caracteres. No puede repetirse dentro del mismo proyecto. |
| Rol | uno de: `Product Owner`, `Scrum Master`, `Developer` | Sí | Escrito exactamente así. |

Para ver un integrante, además se indica su **identificador** (y el del proyecto al que pertenece).

## Salidas esperadas

- **Agregar:** el integrante queda registrado en el proyecto y se devuelve con su identificador asignado por el sistema, su nombre, su email (normalizado, ver RN-3) y su rol.
- **Listar:** se devuelven los integrantes del proyecto en el orden en que se agregaron. Si el proyecto no tiene integrantes, la lista está vacía.
- **Ver uno:** se devuelve el integrante con su identificador, nombre, email y rol.
- **Error:** no se registra nada, y se informa un mensaje que explica qué dato está mal, que el integrante ya está en el proyecto o qué no existe.

## Reglas de negocio

- RN-1: Todo integrante pertenece a un único proyecto existente. No se puede agregar ni listar integrantes de un proyecto que no existe.
- RN-2: El nombre es obligatorio: no puede estar vacío ni tener solo espacios, y tiene como máximo 100 caracteres. Puede repetirse (dos personas pueden llamarse igual).
- RN-3: El email es obligatorio y se guarda **normalizado**: sin espacios al principio ni al final y en minúsculas.
- RN-4: El email tiene que tener un formato válido: sin espacios, exactamente una `@`, al menos un carácter antes de la `@` y, después, un dominio con al menos un punto que no esté ni al principio ni al final del dominio. Tiene como máximo 254 caracteres.
- RN-5: El rol es obligatorio y solo puede ser `Product Owner`, `Scrum Master` o `Developer`.
- RN-6: Un mismo email no puede registrarse dos veces en el mismo proyecto. La comparación se hace sobre el email normalizado (RN-3), así que `Ana@Mail.com` y ` ana@mail.com ` son el mismo integrante. El mismo email sí puede estar en proyectos distintos.
- RN-7: Los integrantes se listan en el orden en que se agregaron al proyecto.
- RN-8: El identificador lo asigna el sistema al agregar el integrante. Un integrante de otro proyecto se trata como inexistente.
- RN-9: Si hay más de un dato inválido, se informa el primero en este orden: nombre, email, rol. El email repetido (RN-6) se controla después de que todos los datos sean válidos.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- Las validaciones de RN-2 a RN-6 viven en el dominio y se desarrollan con TDD.
- Los largos se cuentan en caracteres, no en bytes (las tildes y la ñ cuentan como uno).
- La regla de email repetido (RN-6) se garantiza también en la base de datos, para que dos altas simultáneas no puedan registrar el mismo email en el mismo proyecto.

## Casos límite

- Nombre de exactamente 100 caracteres: se acepta. Con 101: se rechaza.
- Nombre igual al de otro integrante del mismo proyecto, con otro email: se acepta.
- Nombre o email con solo espacios: se rechaza, igual que vacío.
- Email con mayúsculas o espacios al borde (` Ana@Mail.COM `): se acepta y se guarda como `ana@mail.com`.
- Email que ya está en el proyecto escrito con otras mayúsculas: se rechaza como repetido.
- Mismo email en otro proyecto: se acepta.
- Email de exactamente 254 caracteres: se acepta. Con 255: se rechaza.
- Email sin `@`, con dos `@`, sin texto antes de la `@`, sin punto en el dominio, con el punto al borde del dominio (`ana@.com`, `ana@mail.`) o con espacios en el medio: se rechaza.
- Rol escrito con otras mayúsculas (`developer`) o un rol que no existe (`Tester`): se rechaza.
- Listar un proyecto sin integrantes: devuelve una lista vacía, no un error.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Nombre vacío o con solo espacios | Dato inválido: "El nombre es obligatorio" |
| Nombre de más de 100 caracteres | Dato inválido: "El nombre no puede superar los 100 caracteres" |
| Email vacío o con solo espacios | Dato inválido: "El email es obligatorio" |
| Email con formato inválido | Dato inválido: "El email no es válido" |
| Email de más de 254 caracteres | Dato inválido: "El email no puede superar los 254 caracteres" |
| Rol vacío | Dato inválido: "El rol es obligatorio" |
| Rol que no es uno de los tres | Dato inválido: "El rol no es válido" |
| El email ya está registrado en ese proyecto | Conflicto: "El integrante ya está registrado en el proyecto" |
| Agregar, listar o ver integrantes de un proyecto que no existe | No encontrado: "El proyecto no existe" |
| Ver un integrante que no existe o es de otro proyecto | No encontrado: "El integrante no existe" |

"El nombre es obligatorio" es el mismo mensaje de la spec 01.

## Criterios de aceptación

- CA-1: Se puede agregar un integrante a un proyecto con nombre, email y rol, y queda registrado con un identificador asignado por el sistema y el email normalizado.
- CA-2: Se pueden listar los integrantes de un proyecto en el orden en que se agregaron, y ver uno por su identificador.
- CA-3: No se permite registrar dos veces al mismo integrante (mismo email, sin distinguir mayúsculas ni espacios al borde) en un proyecto; sí en proyectos distintos.
- CA-4: El nombre es obligatorio con un máximo de 100 caracteres; el email es obligatorio, con formato válido y un máximo de 254 caracteres.
- CA-5: El rol es obligatorio y solo admite `Product Owner`, `Scrum Master` o `Developer`.
- CA-6: Operar sobre un proyecto o un integrante que no existe informa cuál de los dos no existe, sin registrar nada.

Relación con los criterios iniciales del Issue #5: "agregar un integrante" → CA-1, CA-4 y CA-5; "listar los integrantes" → CA-2; "no registrar dos veces al mismo integrante" → CA-3.

## Escenarios BDD

```gherkin
# language: es
Característica: Registrar integrantes del proyecto
  Para poder asignarles estimaciones, esfuerzo y trabajo
  Como integrante del equipo
  Quiero registrar a los integrantes de un proyecto

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Agregar un integrante con todos sus datos (CA-1)
    Cuando agrego al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Entonces el integrante queda registrado con un identificador
    Y el integrante tiene nombre "Ana Pérez", email "ana@mail.com" y rol "Developer"

  Esquema del escenario: Agregar integrantes con cada rol (CA-1)
    Cuando agrego al integrante "Ana Pérez" con email "ana@mail.com" y rol "<rol>"
    Entonces el integrante queda registrado con un identificador
    Y el integrante tiene rol "<rol>"

    Ejemplos:
      | rol           |
      | Product Owner |
      | Scrum Master  |
      | Developer     |

  Escenario: El email se guarda normalizado (CA-1)
    Cuando agrego al integrante "Ana Pérez" con email " Ana@Mail.COM " y rol "Developer"
    Entonces el integrante queda registrado con un identificador
    Y el integrante tiene email "ana@mail.com"

  Escenario: Listar integrantes en el orden en que se agregaron (CA-2)
    Dado que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Y que el proyecto tiene al integrante "Bruno Díaz" con email "bruno@mail.com" y rol "Scrum Master"
    Cuando listo los integrantes
    Entonces veo los integrantes "Ana Pérez" y "Bruno Díaz" en ese orden

  Escenario: Listar un proyecto sin integrantes (CA-2)
    Cuando listo los integrantes
    Entonces la lista de integrantes está vacía

  Escenario: Cada proyecto ve solo sus integrantes (CA-2)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando listo los integrantes del proyecto "App Turnos"
    Entonces la lista de integrantes está vacía

  Escenario: Ver un integrante por su identificador (CA-2)
    Dado que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando veo ese integrante por su identificador
    Entonces el integrante tiene nombre "Ana Pérez", email "ana@mail.com" y rol "Developer"

  Esquema del escenario: No registrar dos veces el mismo email en un proyecto (CA-3)
    Dado que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando agrego al integrante "Ana P." con email "<email>" y rol "Scrum Master"
    Entonces se informa el error "El integrante ya está registrado en el proyecto"
    Y el proyecto tiene 1 integrante

    Ejemplos:
      | email        |
      | ana@mail.com |
      | ANA@Mail.com |

  Escenario: El mismo email con un espacio adelante también está repetido (CA-3)
    Dado que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando agrego al integrante "Ana P." con email " ana@mail.com" y rol "Scrum Master"
    Entonces se informa el error "El integrante ya está registrado en el proyecto"
    Y el proyecto tiene 1 integrante

  Escenario: El mismo email puede estar en otro proyecto (CA-3)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando agrego al proyecto "App Turnos" al integrante "Ana Pérez" con email "ana@mail.com" y rol "Product Owner"
    Entonces el integrante queda registrado con un identificador

  Escenario: Dos integrantes con el mismo nombre y distinto email (CA-3)
    Dado que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando agrego al integrante "Ana Pérez" con email "ana.perez@otro.com" y rol "Developer"
    Entonces el proyecto tiene 2 integrantes

  Esquema del escenario: Validar el nombre (CA-4)
    Cuando agrego un integrante con un nombre de <largo> caracteres, email "ana@mail.com" y rol "Developer"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                       |
      | 100   | registrado                                      |
      | 101   | El nombre no puede superar los 100 caracteres   |
      | 0     | El nombre es obligatorio                        |

  Escenario: Rechazar un nombre con solo espacios (CA-4)
    Cuando agrego al integrante "   " con email "ana@mail.com" y rol "Developer"
    Entonces se informa el error "El nombre es obligatorio"
    Y la lista de integrantes está vacía

  Esquema del escenario: Validar el email (CA-4)
    Cuando agrego al integrante "Ana Pérez" con email "<email>" y rol "Developer"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | email              | resultado               |
      | ana@mail.com       | registrado              |
      | ana.p@mail.com.ar  | registrado              |
      | (vacío)            | El email es obligatorio |
      | (solo espacios)    | El email es obligatorio |
      | ana.mail.com       | El email no es válido   |
      | ana@@mail.com      | El email no es válido   |
      | @mail.com          | El email no es válido   |
      | ana@mail           | El email no es válido   |
      | ana@.com           | El email no es válido   |
      | ana@mail.          | El email no es válido   |
      | ana pérez@mail.com | El email no es válido   |

  Esquema del escenario: Validar el largo del email (CA-4)
    Cuando agrego al integrante "Ana Pérez" con un email de <largo> caracteres y rol "Developer"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                      |
      | 254   | registrado                                     |
      | 255   | El email no puede superar los 254 caracteres   |

  Esquema del escenario: Validar el rol (CA-5)
    Cuando agrego al integrante "Ana Pérez" con email "ana@mail.com" y rol "<rol>"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | rol          | resultado             |
      | Scrum Master | registrado            |
      | (vacío)      | El rol es obligatorio |
      | Tester       | El rol no es válido   |
      | developer    | El rol no es válido   |

  Escenario: Si hay varios datos inválidos se informa el primero (CA-4)
    Cuando agrego al integrante "" con email "ana.mail.com" y rol "Tester"
    Entonces se informa el error "El nombre es obligatorio"

  Escenario: Agregar un integrante a un proyecto que no existe (CA-6)
    Cuando agrego un integrante al proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Listar los integrantes de un proyecto que no existe (CA-6)
    Cuando listo los integrantes del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Ver un integrante que no existe (CA-6)
    Cuando veo el integrante con identificador 999
    Entonces se informa el error "El integrante no existe"

  Escenario: Un integrante de otro proyecto se trata como inexistente (CA-6)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando veo ese integrante desde el proyecto "App Turnos"
    Entonces se informa el error "El integrante no existe"
```

Los pasos que no indican proyecto ("el proyecto tiene", "agrego al integrante", "listo los integrantes") se refieren al proyecto de los Antecedentes, "Tienda Online". En las tablas, `(vacío)` representa un texto vacío y `(solo espacios)` un texto con solo espacios (Gherkin no puede escribirlos), igual que en las specs 01 y 04.

Al implementar, el email con un espacio adelante del CA-3 (`" ana@mail.com"`) pasó de la tabla del `Esquema del escenario` a un escenario propio en `backend/features/integrantes.feature` ("El mismo email con un espacio adelante también está repetido"), porque Gherkin recorta los espacios de las celdas: entre comillas en el paso, el espacio se conserva.

## Pantalla

`frontend/src/paginas/Integrantes.tsx` (sección "Integrantes" del menú): la lista de integrantes en el orden en que se agregaron, con su rol y su email (ya normalizado por el backend), y un formulario para agregar uno con el rol en un desplegable. No valida nada: el email es un campo de texto común (no `type="email"`) para que el formato lo informe el backend, y muestra sus mensajes de error (por ejemplo, "El integrante ya está registrado en el proyecto").

## Dependencias

- Spec 01 (proyectos): el integrante pertenece a un proyecto existente.
- Las specs 08 (Planning Poker) y 09 (registro de esfuerzo) dependen de esta: usan el identificador del integrante.

## Tareas

- [x] `[P1]` Entidad Integrante en `backend/dominio/` con las validaciones RN-2 a RN-5 y la normalización del email (TDD, un ciclo por regla).
- [x] `[P1]` Casos de uso agregar, listar y ver, con la verificación del proyecto (RN-1) y del email repetido (RN-6), y repositorio en memoria en `backend/memoria/`.
- [x] `[P1]` Tabla de integrantes en `backend/postgres/tablas/03_integrantes.sql`, con el email único por proyecto (RN-6).
- [x] `[P1]` Acceso a datos de integrantes en `backend/postgres/` con sus tests de integración.
- [x] `[P1]` Rutas de la API para agregar, listar y ver integrantes en `backend/api/` (`/api/proyectos/{id}/integrantes`), con 409 para el email repetido.
- [x] `[P1]` Escenarios de `backend/features/integrantes.feature` con sus pasos.
- [x] `[P2]` Actualizar `backend/CLAUDE.md` (endpoints disponibles, specs relacionadas, trampas: recrear la base por la tabla nueva).

## Definition of Done

- [x] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [x] Escenarios BDD automatizados y en verde.
- [x] CI en verde (gofmt / vet / tests / lint / build).
- [x] Criterios de aceptación verificados.
- [x] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [x] Spec en estado `implemented`.
