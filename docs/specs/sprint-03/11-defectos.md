---
titulo: Gestionar defectos
estado: implemented      # draft | in-review | approved | implemented
unidad: backend
issue: "#17"
sprint: 03
autor: Luana Guajardo
aprobada_por: "Martina Batista"
fecha: 2026-10-27
specs_relacionadas: ["04-product-backlog", "05-crear-sprints", "12-metricas-esfuerzo-calidad"]
---

# Spec 11 — Gestionar defectos

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero registrar los defectos encontrados, para hacer seguimiento de la calidad del producto. ([#17](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/17))

## Objetivo

Registrar los defectos que encuentra el equipo, seguir su estado hasta resolverlos y contar cuántos se detectan y cuántos se resuelven, en el proyecto y en cada Sprint. Son las métricas de calidad que pide la consigna y la base de la parte de defectos del reporte (spec 14).

Origen: Consigna — Requerimiento 6 (Gestión de defectos) y Requerimiento 7 (Métricas: cantidad de defectos detectados y resueltos).

## Alcance

**Entra:**
- Registrar un defecto con descripción, severidad, historia relacionada y Sprint de detección. Arranca `Abierto` y sin Sprint de resolución.
- Listar los defectos del proyecto y ver uno.
- Cambiar el estado de un defecto; al pasarlo a `Resuelto` se registra el Sprint de resolución.
- Calcular la cantidad de defectos detectados y resueltos del proyecto y de cada Sprint (CA-4, que viene de la spec 12).
- Impedir eliminar del backlog una historia que tiene defectos registrados, para no perderlos.
- La pantalla "Defectos" en el frontend.

**NO entra:**
- Modificar la descripción, la severidad, la historia o el Sprint de detección de un defecto ya registrado, y eliminar defectos: si se cargó mal, se resuelve como corresponda. Se puede agregar más adelante si el equipo lo necesita.
- Asignar un defecto a un integrante o registrar horas sobre un defecto (el esfuerzo se registra sobre historias, spec 09).
- Mostrar las métricas de defectos en el Dashboard (spec 13: "se pueden sumar en el Sprint 4") o en el reporte (spec 14).

## Entradas

**Registrar:**

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |
| Descripción | texto | Sí | No puede estar vacía ni tener solo espacios. Máximo 1000 caracteres. |
| Severidad | uno de: `Baja`, `Media`, `Alta`, `Crítica` | Sí | Escrita exactamente así. |
| Historia relacionada | identificador | Sí | Tiene que ser una historia del backlog del proyecto (spec 04), en cualquier estado. |
| Sprint de detección | identificador | Sí | Tiene que ser un Sprint del proyecto (spec 05), abierto o cerrado. |

**Cambiar el estado:** el identificador del defecto, el estado nuevo (`Abierto`, `En corrección` o `Resuelto`) y, solo si es `Resuelto`, el Sprint de resolución.

## Salidas esperadas

- **Registrar:** el defecto queda guardado y se devuelve con su identificador, sus datos, el estado `Abierto` y sin Sprint de resolución.
- **Listar:** los defectos del proyecto en el orden en que se registraron, cada uno con su historia y sus Sprints (identificador y número). Si no hay, la lista está vacía.
- **Ver uno:** el defecto con todos sus datos.
- **Cambiar el estado:** se devuelve el defecto con el estado nuevo y, si se resolvió, con su Sprint de resolución.
- **Métricas de defectos:** del proyecto, la cantidad de defectos detectados y de resueltos; y por cada Sprint del proyecto (ordenados por fecha de inicio, incluidos los que no tienen defectos), los detectados en ese Sprint y los resueltos en ese Sprint.
- **Error:** no se registra ni se cambia nada, y se informa un mensaje que explica qué dato está mal o qué no existe.

## Reglas de negocio

- RN-1: Todo defecto pertenece a un proyecto existente. Su historia relacionada y sus Sprints tienen que ser de ese mismo proyecto; si no, se tratan como inexistentes.
- RN-2: La descripción es obligatoria: no puede estar vacía ni tener solo espacios, y tiene como máximo 1000 caracteres.
- RN-3: La severidad es obligatoria y solo puede ser `Baja`, `Media`, `Alta` o `Crítica`.
- RN-4: La historia relacionada es obligatoria y puede estar en cualquier estado (un defecto puede aparecer en una historia ya completada).
- RN-5: El Sprint de detección es obligatorio y puede ser cualquier Sprint del proyecto, abierto o cerrado (se pueden cargar defectos encontrados en Sprints anteriores).
- RN-6: Un defecto nuevo arranca en estado `Abierto` y sin Sprint de resolución.
- RN-7: El estado solo puede ser `Abierto`, `En corrección` o `Resuelto`, y se puede pasar de cualquiera a cualquiera (también reabrir un defecto resuelto).
- RN-8: Al pasar a `Resuelto` hay que indicar el Sprint de resolución: un Sprint del proyecto, abierto o cerrado, que **no empiece antes** que el Sprint de detección (puede ser el mismo). Un defecto que ya estaba `Resuelto` se puede volver a marcar `Resuelto` con otro Sprint de resolución.
- RN-9: Al pasar a `Abierto` o `En corrección`, el defecto queda sin Sprint de resolución (si se indica uno, no se tiene en cuenta).
- RN-10: Detectados y resueltos (CA-4): un defecto cuenta como **detectado** en su Sprint de detección y en el total del proyecto, cualquiera sea su estado. Cuenta como **resuelto** solo si está `Resuelto`, en su Sprint de resolución y en el total del proyecto. Un defecto reabierto deja de contar como resuelto.
- RN-11: Las métricas de defectos incluyen todos los Sprints del proyecto, ordenados por fecha de inicio, con 0 en los que no tienen defectos.
- RN-12: El identificador lo asigna el sistema. Un defecto de otro proyecto se trata como inexistente.
- RN-13: Si hay más de un dato inválido al registrar, se informa el primero en este orden: descripción, severidad, historia, Sprint de detección. Antes que los datos se verifica que exista el proyecto.
- RN-14: No se puede eliminar del backlog (spec 04) una historia que tiene defectos registrados, igual que con el esfuerzo (spec 09 RN-10): los defectos no se pierden sin querer.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- Las validaciones (RN-2 a RN-9) y el conteo de detectados y resueltos (RN-10, RN-11) viven en el dominio y se desarrollan con TDD.
- "No empieza antes" (RN-8) se compara por la fecha de inicio de los Sprints, no por su número (el número sale del orden de creación, spec 05 RN-6). Como los Sprints no se superponen (spec 05 RN-4) y no se eliminan, la comparación no es ambigua.
- Usa las specs 04 y 05, de Sprints anteriores y ya implementadas: esta historia no espera a ninguna otra.

## Casos límite

- Descripción de exactamente 1000 caracteres: se acepta. Con 1001: se rechaza.
- Severidad escrita en minúscula (`alta`) o sin tilde (`Critica`): se rechaza.
- Historia `Completada`: se puede registrar un defecto sobre ella.
- Sprint de detección cerrado: se acepta.
- Resolver en el mismo Sprint de detección: se acepta.
- Resolver en un Sprint que empieza antes que el de detección: se rechaza.
- Resolver sin indicar el Sprint de resolución: se rechaza.
- Reabrir un defecto resuelto: queda sin Sprint de resolución y deja de contar como resuelto.
- Cambiar a `En corrección` indicando un Sprint de resolución: se ignora el Sprint, queda sin resolución.
- Proyecto sin defectos: métricas en 0, y cada Sprint con 0 detectados y 0 resueltos.
- Proyecto sin Sprints: no se puede registrar ningún defecto (no hay Sprint de detección posible); las métricas dan 0 y la lista de Sprints vacía.
- Un defecto detectado en el Sprint 1 y resuelto en el 2: cuenta como detectado en el 1 y como resuelto en el 2.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Descripción vacía o con solo espacios | Dato inválido: "La descripción es obligatoria" (spec 01) |
| Descripción de más de 1000 caracteres | Dato inválido: "La descripción no puede superar los 1000 caracteres" (spec 04) |
| Falta la severidad | Dato inválido: "La severidad es obligatoria" |
| Severidad que no es una de las cuatro | Dato inválido: "La severidad no es válida" |
| Falta la historia relacionada | Dato inválido: "La historia relacionada es obligatoria" |
| Falta el Sprint de detección | Dato inválido: "El Sprint de detección es obligatorio" |
| Estado que no es uno de los tres | Dato inválido: "El estado no es válido" (spec 04) |
| Resolver sin Sprint de resolución | Dato inválido: "Para resolver el defecto hay que indicar el Sprint de resolución" |
| Sprint de resolución que empieza antes que el de detección | Dato inválido: "El Sprint de resolución no puede ser anterior al de detección" |
| El proyecto no existe | No encontrado: "El proyecto no existe" |
| La historia no existe o es de otro proyecto | No encontrado: "La historia no existe" |
| Un Sprint (de detección o de resolución) no existe o es de otro proyecto | No encontrado: "El Sprint no existe" (spec 05) |
| El defecto no existe o es de otro proyecto | No encontrado: "El defecto no existe" |
| Eliminar del backlog una historia con defectos | Conflicto: "La historia tiene defectos registrados" |

## Criterios de aceptación

- CA-1: Cada defecto tiene descripción, severidad, estado, historia relacionada, Sprint de detección y Sprint de resolución; al registrarlo arranca `Abierto` y sin Sprint de resolución.
- CA-2: Se puede cambiar el estado de un defecto entre `Abierto`, `En corrección` y `Resuelto`, también reabrir uno resuelto.
- CA-3: El Sprint de resolución se registra al resolverlo: es obligatorio al pasar a `Resuelto`, no puede ser anterior al de detección, y se borra al reabrir.
- CA-4: Se calcula la cantidad de defectos detectados y la cantidad de defectos resueltos, del proyecto y de cada Sprint (por Sprint de detección y de resolución).
- CA-5: La descripción, la severidad, la historia relacionada y el Sprint de detección son obligatorios y válidos; se pueden listar los defectos del proyecto y ver uno.
- CA-6: No se puede eliminar del backlog una historia con defectos registrados.
- CA-7: Operar con un proyecto, una historia, un Sprint o un defecto que no existe (o es de otro proyecto) informa cuál no existe, sin registrar ni cambiar nada.

Relación con los criterios del Issue #17: CA-1 a CA-4 son los del Issue. CA-5 detalla las validaciones y la consulta; CA-6 protege los defectos al eliminar del backlog (como la spec 09 con el esfuerzo); CA-7 cubre los casos de error.

## Escenarios BDD

```gherkin
# language: es
Característica: Gestionar defectos
  Para hacer seguimiento de la calidad del producto
  Como integrante del equipo
  Quiero registrar los defectos encontrados y seguir su estado

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14"
    Y que existe un Sprint del "2026-11-16" al "2026-11-28"

  Escenario: Registrar un defecto con todos sus datos (CA-1)
    Cuando registro el defecto "El alta acepta emails sin @" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Entonces el defecto queda registrado con un identificador
    Y el defecto tiene descripción "El alta acepta emails sin @", severidad "Alta", historia "Registrar usuario" y Sprint de detección 1
    Y el defecto está "Abierto" y sin Sprint de resolución

  Esquema del escenario: Registrar defectos con cada severidad (CA-1)
    Cuando registro el defecto "Falla" de severidad "<severidad>" en "Registrar usuario" detectado en el Sprint 1
    Entonces el defecto queda registrado con un identificador

    Ejemplos:
      | severidad |
      | Baja      |
      | Media     |
      | Alta      |
      | Crítica   |

  Escenario: Registrar un defecto en una historia completada y en un Sprint cerrado (CA-1)
    Dado que la historia "Registrar usuario" pasa a estado "Completada"
    Y que cerré el Sprint 1
    Cuando registro el defecto "Falla tardía" de severidad "Media" en "Registrar usuario" detectado en el Sprint 1
    Entonces el defecto queda registrado con un identificador

  Escenario: Pasar un defecto a En corrección (CA-2)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando cambio el estado del defecto "Falla" a "En corrección"
    Entonces el defecto está "En corrección" y sin Sprint de resolución

  Escenario: Reabrir un defecto resuelto (CA-2)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que resolví el defecto "Falla" en el Sprint 2
    Cuando cambio el estado del defecto "Falla" a "Abierto"
    Entonces el defecto está "Abierto" y sin Sprint de resolución

  Esquema del escenario: Validar el estado (CA-2)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando cambio el estado del defecto "Falla" a "<estado>"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | estado        | resultado              |
      | En corrección | modificado             |
      | Abierto       | modificado             |
      | Cerrado       | El estado no es válido |
      | resuelto      | El estado no es válido |

  Escenario: Resolver un defecto registra el Sprint de resolución (CA-3)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando resuelvo el defecto "Falla" en el Sprint 2
    Entonces el defecto está "Resuelto" con Sprint de resolución 2

  Escenario: Resolver en el mismo Sprint de detección (CA-3)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando resuelvo el defecto "Falla" en el Sprint 1
    Entonces el defecto está "Resuelto" con Sprint de resolución 1

  Escenario: No se puede resolver en un Sprint anterior al de detección (CA-3)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 2
    Cuando resuelvo el defecto "Falla" en el Sprint 1
    Entonces se informa el error "El Sprint de resolución no puede ser anterior al de detección"
    Y el defecto está "Abierto" y sin Sprint de resolución

  Escenario: No se puede resolver sin indicar el Sprint de resolución (CA-3)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando cambio el estado del defecto "Falla" a "Resuelto"
    Entonces se informa el error "Para resolver el defecto hay que indicar el Sprint de resolución"

  Escenario: Contar defectos detectados y resueltos (CA-4)
    Dado que registré el defecto "A" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que registré el defecto "B" de severidad "Baja" en "Registrar usuario" detectado en el Sprint 1
    Y que registré el defecto "C" de severidad "Media" en "Registrar usuario" detectado en el Sprint 2
    Y que resolví el defecto "A" en el Sprint 2
    Y que resolví el defecto "C" en el Sprint 2
    Cuando consulto las métricas de defectos
    Entonces el proyecto tiene 3 defectos detectados y 2 resueltos
    Y el Sprint 1 tiene 2 defectos detectados y 0 resueltos
    Y el Sprint 2 tiene 1 defecto detectado y 2 resueltos

  Escenario: Un defecto reabierto deja de contar como resuelto (CA-4)
    Dado que registré el defecto "A" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que resolví el defecto "A" en el Sprint 1
    Y que cambié el estado del defecto "A" a "En corrección"
    Cuando consulto las métricas de defectos
    Entonces el proyecto tiene 1 defecto detectado y 0 resueltos
    Y el Sprint 1 tiene 1 defecto detectado y 0 resueltos

  Escenario: Métricas sin defectos (CA-4)
    Cuando consulto las métricas de defectos
    Entonces el proyecto tiene 0 defectos detectados y 0 resueltos
    Y el Sprint 1 tiene 0 defectos detectados y 0 resueltos
    Y el Sprint 2 tiene 0 defectos detectados y 0 resueltos

  Esquema del escenario: Validar los datos al registrar (CA-5)
    Cuando registro el defecto "<descripcion>" de severidad "<severidad>" en "Registrar usuario" detectado en el Sprint 1
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | descripcion     | severidad | resultado                     |
      | (vacío)         | Alta      | La descripción es obligatoria |
      | (solo espacios) | Alta      | La descripción es obligatoria |
      | Falla           | (vacío)   | La severidad es obligatoria   |
      | Falla           | alta      | La severidad no es válida     |
      | Falla           | Critica   | La severidad no es válida     |

  Esquema del escenario: Validar el largo de la descripción (CA-5)
    Cuando registro un defecto con una descripción de <largo> caracteres de severidad "Baja" en "Registrar usuario" detectado en el Sprint 1
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                           |
      | 1000  | registrado                                          |
      | 1001  | La descripción no puede superar los 1000 caracteres |

  Escenario: Registrar sin historia relacionada (CA-5)
    Cuando registro el defecto "Falla" de severidad "Alta" sin historia detectado en el Sprint 1
    Entonces se informa el error "La historia relacionada es obligatoria"

  Escenario: Registrar sin Sprint de detección (CA-5)
    Cuando registro el defecto "Falla" de severidad "Alta" en "Registrar usuario" sin Sprint de detección
    Entonces se informa el error "El Sprint de detección es obligatorio"
    Y no hay defectos registrados

  Escenario: Listar los defectos en el orden en que se registraron (CA-5)
    Dado que registré el defecto "Primero" de severidad "Baja" en "Registrar usuario" detectado en el Sprint 2
    Y que registré el defecto "Segundo" de severidad "Crítica" en "Registrar usuario" detectado en el Sprint 1
    Cuando listo los defectos
    Entonces veo los defectos "Primero" y "Segundo" en ese orden

  Escenario: No se puede eliminar del backlog una historia con defectos (CA-6)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando elimino la historia "Registrar usuario"
    Entonces se informa el error "La historia tiene defectos registrados"
    Y el backlog tiene solo la historia "Registrar usuario"

  Escenario: Registrar en una historia que no existe (CA-7)
    Cuando registro el defecto "Falla" de severidad "Alta" en la historia con identificador 999 detectado en el Sprint 1
    Entonces se informa el error "La historia no existe"
    Y no hay defectos registrados

  Escenario: Registrar en un Sprint que no existe (CA-7)
    Cuando registro el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint con identificador 999
    Entonces se informa el error "El Sprint no existe"

  Escenario: Resolver en un Sprint de otro proyecto (CA-7)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto "App Turnos" tiene un Sprint del "2026-11-05" al "2026-11-20"
    Cuando resuelvo el defecto "Falla" en el Sprint de "App Turnos"
    Entonces se informa el error "El Sprint no existe"

  Escenario: Ver un defecto que no existe (CA-7)
    Cuando veo el defecto con identificador 999
    Entonces se informa el error "El defecto no existe"
```

Los pasos que no indican proyecto se refieren al de los Antecedentes, "Tienda Online". "El Sprint 1" y "el Sprint 2" son los Sprints de los Antecedentes por su número. En las tablas, `(vacío)` representa un texto vacío y `(solo espacios)` uno con solo espacios, como en las demás specs. Los pasos "que existe un Sprint del … al …", "que cerré el Sprint …", "la historia … pasa a estado …" y "elimino la historia …" son los de las specs 04 a 06.

## Pantalla

`frontend/src/paginas/Defectos.tsx` (sección "Defectos" del menú): arriba, las métricas (detectados y resueltos del proyecto y una tabla por Sprint); abajo, la lista de defectos con su severidad, estado, historia y Sprints, con un desplegable para cambiar el estado (al elegir `Resuelto` se elige también el Sprint de resolución) y un formulario para registrar. No valida nada: muestra los mensajes de error del backend.

## Dependencias

- Spec 04 (backlog) y spec 05 (Sprints): implementadas en el Sprint 1. Esta historia no espera a ninguna otra.
- Spec 04: esta spec agrega una regla al eliminar una historia (RN-14, CA-6), igual que las specs 05 y 09.
- Spec 12: le pasó el cálculo de detectados y resueltos (CA-4).
- La usan: spec 14 (reporte), que muestra los defectos y las cantidades del CA-4; spec 13 (Dashboard), si en el Sprint 4 se suman las métricas de calidad.

## Tareas

- [x] `[P1]` Defecto en `backend/dominio/`: validaciones (RN-2 a RN-5), cambio de estado con Sprint de resolución (RN-6 a RN-9) y conteo de detectados y resueltos (RN-10, RN-11) (TDD).
- [x] `[P1]` Casos de uso registrar, listar, ver, cambiar estado y métricas, con la verificación de proyecto, historia y Sprints (RN-1, RN-12, RN-13), y repositorio en memoria.
- [x] `[P1]` Regla del backlog: no eliminar una historia con defectos (RN-14).
- [x] `[P1]` Tabla `10_defectos.sql`, repositorio PostgreSQL y rutas REST (`/api/proyectos/{id}/defectos` y `/api/proyectos/{id}/metricas/defectos`), con escenarios de `backend/features/defectos.feature`.
- [x] `[P1]` Pantalla "Defectos" en el frontend.
- [x] `[P2]` Actualizar `backend/CLAUDE.md` y `frontend/CLAUDE.md`.

## Definition of Done

- [x] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [x] Escenarios BDD automatizados y en verde.
- [x] CI en verde (gofmt / vet / tests / lint / build).
- [x] Criterios de aceptación verificados.
- [x] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [x] Spec en estado `implemented`.
