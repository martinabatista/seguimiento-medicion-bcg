---
titulo: Estimar historias con Story Points
estado: in-review    # draft | in-review | approved | implemented
unidad: backend
issue: "#10"
sprint: 01
autor: Florencia Cabeza
aprobada_por: ""
fecha: 2026-10-23
specs_relacionadas: ["04-product-backlog", "06-cerrar-sprint", "08-planning-poker"]
---

# Spec 07 — Estimar historias con Story Points

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero asignar Story Points a una historia, para conocer el tamaño relativo del trabajo. ([#10](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/10))

## Objetivo

Dimensionar las historias del backlog con un número que indica su **tamaño relativo** (cuánto trabajo, complejidad e incertidumbre tiene comparada con otras), no horas. Es la base para planificar Sprints y medir la velocidad del equipo (spec 10), y el valor que registra el Planning Poker (spec 08).

Origen: Consigna — Requerimiento 4 (Estimación mediante Story Points).

## Alcance

**Entra:**
- Asignar Story Points a una historia del backlog que está sin estimar.
- Cambiar la estimación de una historia por otro valor.
- Quitar la estimación: la historia vuelve a quedar sin estimar.
- Consultar la escala de valores permitidos.

**NO entra:**
- El Planning Poker (votación del equipo, rondas, diferencias): spec 08. Esa spec registra la estimación acordada usando esta misma operación y esta misma escala.
- Calcular Story Points planificados, completados o la velocidad: spec 10.
- Estimar en horas: el esfuerzo en horas se registra en la spec 09.
- Pantallas del frontend (todavía no existe; se hacen cuando se arme la interfaz).

## Entradas

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |
| Historia | identificador | Sí | Tiene que ser una historia del backlog de ese proyecto (spec 04). |
| Story Points | número entero | Sí, al asignar o cambiar | Tiene que ser un valor de la escala: **1, 2, 3, 5, 8, 13 o 21**. |

Para **quitar** la estimación solo se indican el proyecto y la historia. Para **consultar la escala** no hace falta ningún dato.

## Salidas esperadas

- **Asignar o cambiar:** se devuelve la historia con sus Story Points actualizados y el resto de sus datos sin cambios.
- **Quitar:** se devuelve la historia sin estimar.
- **Consultar la escala:** los valores permitidos de menor a mayor: 1, 2, 3, 5, 8, 13, 21.
- **Ver o listar el backlog** (spec 04): cada historia muestra sus Story Points, o que está sin estimar.
- **Error:** la estimación no cambia y se informa un mensaje que explica qué está mal.

## Reglas de negocio

- RN-1: La escala de Story Points es la de **Fibonacci** hasta 21: **1, 2, 3, 5, 8, 13, 21**. Cada valor es la suma de los dos anteriores: cuanto más grande es una historia, más incertidumbre tiene, y los valores se separan más. No se discute si algo es 14 o 15: es un 13.
- RN-2: Solo se aceptan valores de la escala. No se aceptan 0, negativos, decimales, valores intermedios (por ejemplo, 4) ni mayores a 21. Una historia que parece de más de 21 es señal de que conviene dividirla.
- RN-3: Una historia del backlog arranca sin estimar (spec 04 RN-6). Se le pueden asignar Story Points, cambiarlos por otro valor de la escala o quitarlos (vuelve a sin estimar).
- RN-4: No se puede asignar, cambiar ni quitar la estimación de una historia `Completada`: ya se terminó, y cambiar su tamaño alteraría lo que se midió en el Sprint (velocidad, spec 10). Si se deshace la completada (spec 06 RN-3), vuelve a `En progreso` y se puede volver a estimar.
- RN-5: Estimar no cambia ningún otro dato de la historia (título, descripción, prioridad, criterios ni estado). Y modificar la historia desde el backlog no cambia sus Story Points (spec 04).
- RN-6: La historia tiene que pertenecer al proyecto indicado. Una historia de otro proyecto se trata como inexistente (como en spec 04 RN-10).

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- La escala (RN-1) y sus validaciones (RN-2, RN-4) viven en el dominio y se desarrollan con TDD. El frontend y el Planning Poker (spec 08) obtienen los valores consultando la escala: no la repiten.
- Los Story Points se guardan en la historia del backlog (spec 04), que ya tiene el dato "sin estimar".

## Casos límite

- 1 (el menor de la escala): se acepta. 21 (el mayor): se acepta.
- 0: se rechaza. 22 o 34: se rechaza (fuera de la escala).
- 4, 6 o 10 (números entre valores de la escala): se rechaza.
- Negativo (-3) o con decimales (2.5): se rechaza.
- Falta el valor al asignar: se rechaza.
- Asignar el mismo valor que ya tiene: se acepta y queda igual.
- Quitar la estimación de una historia que ya está sin estimar: se acepta y queda sin estimar (no es un error).
- Estimar una historia `En progreso` o asignada a un Sprint abierto: se acepta.
- Estimar una historia `Completada`: se rechaza, y sus Story Points quedan como estaban.
- Historia completada en un Sprint y después deshecha (vuelve a `En progreso`): se puede estimar.
- Estimar una historia existente indicando otro proyecto: se informa que la historia no existe.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Falta el valor al asignar o cambiar | Dato inválido: "Los Story Points son obligatorios" |
| Valor fuera de la escala, con decimales, cero o negativo | Dato inválido: "Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21" |
| Asignar, cambiar o quitar la estimación de una historia completada | Conflicto: "La historia ya está completada" |
| El proyecto no existe | No encontrado: "El proyecto no existe" |
| La historia no existe (o es de otro proyecto) | No encontrado: "La historia no existe" |

El mensaje "La historia ya está completada" es el mismo que usa la spec 05.

## Criterios de aceptación

- CA-1: Se pueden asignar Story Points de la escala a una historia sin estimar; la historia queda con ese valor y el resto de sus datos no cambia.
- CA-2: Se puede cambiar la estimación por otro valor de la escala y se puede quitar, con lo que la historia vuelve a quedar sin estimar.
- CA-3: Solo se aceptan los valores de la escala de Fibonacci 1, 2, 3, 5, 8, 13 y 21; cualquier otro valor se rechaza y la estimación no cambia.
- CA-4: No se puede asignar, cambiar ni quitar la estimación de una historia completada.
- CA-5: Se puede consultar la escala de valores permitidos, de menor a mayor.
- CA-6: Estimar una historia de un proyecto o una historia que no existe informa cuál de los dos no existe, sin cambiar nada.

Relación con los criterios iniciales del Issue #10: "asignar y modificar la estimación" → CA-1 y CA-2; "solo valores de la escala definida" → CA-3 y CA-5.

## Escenarios BDD

```gherkin
# language: es
Característica: Estimar historias con Story Points
  Para conocer el tamaño relativo del trabajo
  Como integrante del equipo
  Quiero asignar Story Points a las historias del backlog

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Asignar Story Points a una historia sin estimar (CA-1)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando estimo la historia "Registrar usuario" con 5 Story Points
    Entonces la historia "Registrar usuario" tiene 5 Story Points
    Y la historia tiene prioridad 1 y estado "Pendiente"

  Escenario: Estimar una historia en progreso (CA-1)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1 en estado "En progreso"
    Cuando estimo la historia "Registrar usuario" con 8 Story Points
    Entonces la historia "Registrar usuario" tiene 8 Story Points

  Escenario: Cambiar la estimación de una historia (CA-2)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que estimé la historia "Registrar usuario" con 3 Story Points
    Cuando estimo la historia "Registrar usuario" con 13 Story Points
    Entonces la historia "Registrar usuario" tiene 13 Story Points

  Escenario: Quitar la estimación de una historia (CA-2)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que estimé la historia "Registrar usuario" con 3 Story Points
    Cuando quito la estimación de la historia "Registrar usuario"
    Entonces la historia "Registrar usuario" está sin estimar

  Escenario: Quitar la estimación de una historia que ya está sin estimar (CA-2)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando quito la estimación de la historia "Registrar usuario"
    Entonces la historia "Registrar usuario" está sin estimar

  Esquema del escenario: Solo se aceptan valores de la escala (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando estimo la historia "Registrar usuario" con "<valor>" Story Points
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | valor   | resultado                                                                   |
      | 1       | estimada                                                                    |
      | 21      | estimada                                                                    |
      | 0       | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | 4       | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | 22      | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | -3      | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | 2.5     | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | (vacío) | Los Story Points son obligatorios                                           |

  Escenario: Un valor fuera de la escala no cambia la estimación (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que estimé la historia "Registrar usuario" con 3 Story Points
    Cuando estimo la historia "Registrar usuario" con 4 Story Points
    Entonces se informa el error "Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21"
    Y la historia "Registrar usuario" tiene 3 Story Points

  Escenario: No se puede estimar una historia completada (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1 en estado "Completada"
    Cuando estimo la historia "Registrar usuario" con 5 Story Points
    Entonces se informa el error "La historia ya está completada"
    Y la historia "Registrar usuario" está sin estimar

  Escenario: No se puede quitar la estimación de una historia completada (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que estimé la historia "Registrar usuario" con 3 Story Points
    Y que existe un Sprint del "2026-11-02" al "2026-11-15" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando quito la estimación de la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está completada"
    Y la historia "Registrar usuario" tiene 3 Story Points

  Escenario: Una historia completada y deshecha se puede volver a estimar (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-02" al "2026-11-15" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando deshago en el Sprint la historia "Registrar usuario"
    Y estimo la historia "Registrar usuario" con 2 Story Points
    Entonces la historia "Registrar usuario" tiene 2 Story Points

  Escenario: Consultar la escala de Story Points (CA-5)
    Cuando consulto la escala de Story Points
    Entonces la escala es 1, 2, 3, 5, 8, 13 y 21

  Escenario: Estimar en un proyecto que no existe (CA-6)
    Cuando estimo la historia con identificador 1 del proyecto con identificador 999 con 3 Story Points
    Entonces se informa el error "El proyecto no existe"

  Escenario: Estimar una historia que no existe (CA-6)
    Cuando estimo la historia con identificador 999 con 3 Story Points
    Entonces se informa el error "La historia no existe"

  Escenario: Una historia de otro proyecto se trata como inexistente (CA-6)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando estimo la historia "Registrar usuario" desde el proyecto "App Turnos" con 3 Story Points
    Entonces se informa el error "La historia no existe"
    Y la historia "Registrar usuario" está sin estimar
```

Los pasos que no indican proyecto se refieren al proyecto de los Antecedentes, "Tienda Online". En las tablas, `(vacío)` representa un valor que falta (Gherkin no puede escribirlo), igual que en las specs 01 y 04. Los pasos del backlog y del Sprint (`que el backlog tiene la historia…`, `que existe un Sprint…`, `que completé en el Sprint…`, `deshago en el Sprint…`) ya existen en las specs 04, 05 y 06 y se reutilizan.

## Dependencias

- Spec 04 (backlog): implementada. La historia ya tiene el dato de Story Points, que arranca sin estimar.
- Spec 06 (completar y deshacer en el Sprint): implementada. Se usa en los escenarios del CA-4.
- La usan: spec 08 (Planning Poker: los votos y la estimación acordada usan esta escala y esta operación; misma dueña, Sprint 2), spec 10 (métricas de avance: suma los Story Points) y spec 13 (Dashboard).
- Chequeo de independencia (INVEST): todos los criterios usan historias ya implementadas (04, 05 y 06); ninguno depende de una historia pendiente de otra integrante.

## Tareas

- [ ] `[P1]` Escala de Story Points y sus reglas (RN-1, RN-2, RN-4) en `backend/dominio/` (TDD).
- [ ] `[P1]` Casos de uso asignar/cambiar, quitar la estimación y consultar la escala (RN-3, RN-5, RN-6).
- [ ] `[P1]` Rutas de la API para estimar, quitar la estimación y consultar la escala en `backend/api/`.
- [ ] `[P1]` Escenarios de `backend/features/story-points.feature` con sus pasos.
- [ ] `[P2]` Actualizar `backend/CLAUDE.md` (endpoints disponibles, specs relacionadas).

No hace falta tabla nueva: la tabla de historias (spec 04) ya guarda los Story Points.

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
