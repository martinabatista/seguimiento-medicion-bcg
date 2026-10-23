---
titulo: Consultar el estado de un proyecto
estado: approved     # draft | in-review | approved | implemented
unidad: backend
issue: "#6"
sprint: 01
autor: Luana Guajardo
aprobada_por: "Martina Batista"
fecha: 2026-10-23
specs_relacionadas: ["01-proyectos", "04-product-backlog", "05-crear-sprints", "07-story-points"]
---

# Spec 03 — Consultar el estado de un proyecto

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero consultar el estado de un proyecto, para saber cómo viene avanzando. ([#6](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/6))

## Objetivo

Dar en una sola consulta una vista rápida y confiable de la situación actual de un proyecto: en qué etapa está según sus fechas, qué Sprint se está trabajando y cuánto avanzó su Product Backlog. Es la base del Dashboard (spec 13).

Origen: Consigna — Requerimiento 1 (Gestión de proyectos: consultar el estado).

## Alcance

**Entra:**
- Consultar el estado de un proyecto por su identificador.
- Calcular el estado del proyecto a partir de sus fechas y de la fecha de hoy.
- Informar el Sprint en curso, si lo hay.
- Calcular un resumen del avance del Product Backlog: cantidades por estado, porcentaje de historias completadas y Story Points.

**NO entra:**
- Cambiar el estado a mano: el estado no se carga, se calcula (RN-2).
- Métricas de velocidad, horas, desviación y defectos (specs 10 y 12).
- Gráficos (Dashboard, spec 13).
- Crear, cerrar o modificar Sprints (specs 05 y 06): esta spec solo lee el Sprint en curso.
- Pantallas del frontend (todavía no existe; se hacen cuando se arme la interfaz).

## Entradas

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |

La **fecha de hoy** no la carga la usuaria: la toma el sistema en el momento de la consulta.

## Salidas esperadas

Una sola respuesta con tres partes:

1. **Datos del proyecto** (los de la spec 01: identificador, nombre, descripción, fecha de inicio y fecha de finalización) y su **estado**: `No iniciado`, `En curso` o `Finalizado`.
2. **Sprint en curso:** su identificador, Sprint Goal, fecha de inicio y fecha de fin (los datos los define la spec 05). Si no hay ninguno, se indica que no hay Sprint en curso.
3. **Resumen del backlog:**
   - Total de historias.
   - Cantidad de historias en cada estado: `Pendiente`, `En progreso` y `Completada`.
   - Porcentaje de historias completadas.
   - Story Points totales (suma de las historias estimadas) y Story Points completados (suma de las historias `Completada` estimadas).
   - Cantidad de historias sin estimar.

**Error:** si el proyecto no existe, se informa "El proyecto no existe".

## Reglas de negocio

- RN-1: Solo se puede consultar el estado de un proyecto existente.
- RN-2: El estado del proyecto **se calcula** comparando la fecha de hoy con sus fechas; no se guarda ni se modifica a mano:
  - `No iniciado`: hoy es anterior a la fecha de inicio.
  - `En curso`: hoy está entre la fecha de inicio y la de finalización, **ambas incluidas**.
  - `Finalizado`: hoy es posterior a la fecha de finalización.
- RN-3: El Sprint en curso es el Sprint **abierto** (no cerrado, spec 06) del proyecto cuyo período incluye la fecha de hoy (inicio y fin incluidos). Si hay más de uno, se informa el de fecha de inicio más reciente. Si no hay ninguno, no hay Sprint en curso.
- RN-4: El resumen cuenta todas las historias del backlog del proyecto (spec 04) según su estado.
- RN-5: El porcentaje de historias completadas es `completadas / total × 100`, redondeado al entero más cercano (la mitad redondea hacia arriba: 12,5 % → 13 %). Con el backlog vacío, el porcentaje es 0.
- RN-6: Los Story Points totales suman los de todas las historias estimadas; los completados, solo los de las historias `Completada` estimadas. Las historias sin estimar no suman y se informan aparte. Si no hay ninguna estimada, ambos totales son 0.
- RN-7: El estado, el Sprint en curso y el resumen se calculan en el momento de la consulta: siempre reflejan los datos actuales.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- El cálculo del estado (RN-2), del Sprint en curso (RN-3) y del resumen (RN-4 a RN-6) vive en el dominio y se desarrolla con TDD. La fecha de hoy se le pasa al cálculo para poder probarlo con fechas fijas.
- La fecha de hoy es la fecha local del servidor, sin hora.
- **Dependencia con la spec 05:** el Sprint en curso (CA-2) recién se puede programar cuando exista el código de Sprints. Se implementa en dos entregas: primero CA-1 y CA-3; después CA-2, junto con o después de la spec 05. El Issue #6 se cierra con la segunda entrega.

## Casos límite

- Hoy es exactamente la fecha de inicio: `En curso`.
- Hoy es exactamente la fecha de finalización: `En curso` (es el último día).
- Hoy es el día siguiente a la fecha de finalización: `Finalizado`.
- Hoy es el día anterior a la fecha de inicio: `No iniciado`.
- Backlog vacío: total 0, todas las cantidades en 0, porcentaje 0 y Story Points 0. No es un error.
- Todas las historias completadas: 100 %.
- 1 de 3 completadas: 33 %. 2 de 3: 67 %. 1 de 8: 13 % (12,5 redondea hacia arriba).
- Ninguna historia estimada: Story Points totales y completados en 0, y todas cuentan como sin estimar.
- Un Sprint abierto cuyo período todavía no empezó o ya terminó: no es el Sprint en curso.
- Un Sprint cerrado cuyo período incluye hoy: no es el Sprint en curso.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Consultar el estado de un proyecto que no existe | No encontrado: "El proyecto no existe" |

## Criterios de aceptación

- CA-1: Se muestran los datos del proyecto y su estado (`No iniciado`, `En curso` o `Finalizado`), calculado a partir de sus fechas y de la fecha de hoy.
- CA-2: Se muestra el Sprint en curso del proyecto, si lo hay; si no, se indica que no hay Sprint en curso.
- CA-3: Se muestra un resumen del avance del Product Backlog: total de historias, cantidad por estado, porcentaje de completadas, Story Points totales y completados, y cantidad sin estimar.
- CA-4: Consultar el estado de un proyecto que no existe informa que el proyecto no existe.

Relación con los criterios iniciales del Issue #6: "datos del proyecto y su estado" → CA-1; "Sprint en curso" → CA-2; "resumen del avance del backlog" → CA-3. El CA-4 se agrega por el caso de error.

## Escenarios BDD

```gherkin
# language: es
Característica: Consultar el estado de un proyecto
  Para saber cómo viene avanzando
  Como integrante del equipo
  Quiero consultar el estado de un proyecto

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Ver los datos del proyecto con su estado (CA-1)
    Dado que hoy es "2026-11-20"
    Cuando consulto el estado del proyecto
    Entonces veo el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Y el estado del proyecto es "En curso"

  Esquema del escenario: El estado depende de la fecha de hoy (CA-1)
    Dado que hoy es "<hoy>"
    Cuando consulto el estado del proyecto
    Entonces el estado del proyecto es "<estado>"

    Ejemplos:
      | hoy        | estado      |
      | 2026-10-15 | No iniciado |
      | 2026-10-31 | No iniciado |
      | 2026-11-01 | En curso    |
      | 2026-11-20 | En curso    |
      | 2026-12-15 | En curso    |
      | 2026-12-16 | Finalizado  |
      | 2027-03-01 | Finalizado  |

  Escenario: Ver el Sprint en curso (CA-2)
    Dado que hoy es "2026-11-10"
    Y que el proyecto tiene el Sprint abierto con Sprint Goal "Catálogo navegable" del "2026-11-03" al "2026-11-16"
    Cuando consulto el estado del proyecto
    Entonces el Sprint en curso tiene Sprint Goal "Catálogo navegable" del "2026-11-03" al "2026-11-16"

  Escenario: Proyecto sin Sprints (CA-2)
    Dado que hoy es "2026-11-10"
    Cuando consulto el estado del proyecto
    Entonces no hay Sprint en curso

  Esquema del escenario: Solo cuenta un Sprint abierto que incluye la fecha de hoy (CA-2)
    Dado que hoy es "<hoy>"
    Y que el proyecto tiene el Sprint <situacion> con Sprint Goal "Catálogo navegable" del "2026-11-03" al "2026-11-16"
    Cuando consulto el estado del proyecto
    Entonces <resultado>

    Ejemplos:
      | hoy        | situacion | resultado                                                                                  |
      | 2026-11-03 | abierto   | el Sprint en curso tiene Sprint Goal "Catálogo navegable" del "2026-11-03" al "2026-11-16" |
      | 2026-11-16 | abierto   | el Sprint en curso tiene Sprint Goal "Catálogo navegable" del "2026-11-03" al "2026-11-16" |
      | 2026-11-02 | abierto   | no hay Sprint en curso                                                                     |
      | 2026-11-17 | abierto   | no hay Sprint en curso                                                                     |
      | 2026-11-10 | cerrado   | no hay Sprint en curso                                                                     |

  Escenario: Resumen del backlog con historias en todos los estados (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" en estado "Completada" con 5 Story Points
    Y que el backlog tiene la historia "Ver catálogo" en estado "Completada" con 3 Story Points
    Y que el backlog tiene la historia "Buscar producto" en estado "En progreso" con 8 Story Points
    Y que el backlog tiene la historia "Pagar pedido" en estado "Pendiente" sin estimar
    Cuando consulto el estado del proyecto
    Entonces el resumen del backlog tiene 4 historias: 1 pendiente, 1 en progreso y 2 completadas
    Y el porcentaje de historias completadas es 50
    Y los Story Points son 16 en total y 8 completados, con 1 historia sin estimar

  Escenario: Resumen de un backlog vacío (CA-3)
    Cuando consulto el estado del proyecto
    Entonces el resumen del backlog tiene 0 historias: 0 pendientes, 0 en progreso y 0 completadas
    Y el porcentaje de historias completadas es 0
    Y los Story Points son 0 en total y 0 completados, con 0 historias sin estimar

  Esquema del escenario: El porcentaje de completadas se redondea al entero más cercano (CA-3)
    Dado que el backlog tiene <total> historias de las cuales <completadas> están completadas
    Cuando consulto el estado del proyecto
    Entonces el porcentaje de historias completadas es <porcentaje>

    Ejemplos:
      | total | completadas | porcentaje |
      | 3     | 1           | 33         |
      | 3     | 2           | 67         |
      | 8     | 1           | 13         |
      | 4     | 4           | 100        |
      | 4     | 0           | 0          |

  Escenario: Historias sin estimar no suman Story Points (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" en estado "Completada" sin estimar
    Y que el backlog tiene la historia "Ver catálogo" en estado "Pendiente" sin estimar
    Cuando consulto el estado del proyecto
    Entonces los Story Points son 0 en total y 0 completados, con 2 historias sin estimar

  Escenario: El resumen refleja los cambios del backlog (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" en estado "Pendiente" sin estimar
    Y que consulto el estado del proyecto
    Cuando la historia "Registrar usuario" pasa a estado "Completada"
    Y consulto el estado del proyecto
    Entonces el porcentaje de historias completadas es 100

  Escenario: Cada proyecto tiene su propio resumen (CA-3)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el backlog del proyecto "App Turnos" tiene la historia "Pedir turno" en estado "Completada" sin estimar
    Cuando consulto el estado del proyecto "Tienda Online"
    Entonces el resumen del backlog tiene 0 historias: 0 pendientes, 0 en progreso y 0 completadas

  Escenario: Consultar el estado de un proyecto que no existe (CA-4)
    Cuando consulto el estado del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"
```

Los pasos que no indican proyecto ("consulto el estado del proyecto", "el backlog tiene") se refieren al proyecto de los Antecedentes, "Tienda Online". En los escenarios que no dicen qué día es hoy, la fecha de hoy no influye en lo que se verifica. Los escenarios del CA-2 se automatizan en la segunda entrega, cuando exista el código de Sprints (spec 05); sus pasos pueden ajustarse a cómo la spec 05 defina crear y cerrar un Sprint.

Los escenarios del CA-3 que cargan historias **con** Story Points ("con 5 Story Points") necesitan poder estimar una historia, que llega con la spec 07. Hasta entonces, la suma de Story Points (RN-6) se verifica con tests unitarios del dominio, y esos escenarios se automatizan cuando la spec 07 esté implementada. Los escenarios del CA-3 que usan solo historias sin estimar se automatizan en la primera entrega.

## Dependencias

- Spec 01 (proyectos): datos y fechas del proyecto.
- Spec 04 (backlog): estado y Story Points de las historias.
- Spec 07 (Story Points): hasta que exista, las historias siempre están sin estimar; los escenarios BDD con Story Points esperan a esa spec (la regla RN-6 igual se prueba con tests unitarios).
- Spec 05 (crear Sprints) y spec 06 (cerrar Sprint): qué es un Sprint abierto, sus fechas y su Sprint Goal. **El CA-2 espera a que la spec 05 esté implementada.**
- La spec 13 (Dashboard) va a usar este estado y este resumen.

## Tareas

Primera entrega (CA-1, CA-3, CA-4):

- [ ] `[P1]` Cálculo del estado del proyecto según la fecha de hoy en `backend/dominio/` (RN-2, TDD).
- [ ] `[P1]` Cálculo del resumen del backlog en `backend/dominio/`: cantidades, porcentaje redondeado y Story Points (RN-4 a RN-6, TDD).
- [ ] `[P1]` Caso de uso "consultar estado" que verifica el proyecto (RN-1) y arma la respuesta con la fecha de hoy que recibe.
- [ ] `[P1]` Ruta de la API para consultar el estado (`/api/proyectos/{id}/estado`), con la fecha del día tomada del sistema.
- [ ] `[P1]` Escenarios de `backend/features/estado-proyecto.feature` (CA-1, CA-3 sin Story Points, CA-4) con sus pasos.
- [ ] `[P2]` Actualizar `backend/CLAUDE.md` (endpoint, specs relacionadas).

Segunda entrega (CA-2, después de la spec 05):

- [ ] `[P1]` Sprint en curso en el dominio (RN-3, TDD) y en la respuesta del estado.
- [ ] `[P1]` Escenarios del CA-2 automatizados.
- [ ] `[P2]` Escenarios del CA-3 con Story Points automatizados (cuando exista la spec 07).

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
