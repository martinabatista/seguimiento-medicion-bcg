---
titulo: Métricas de avance del Sprint y del proyecto
estado: in-review    # draft | in-review | approved | implemented
unidad: backend
issue: "#14"
sprint: 02
autor: Martina Batista
aprobada_por: ""
fecha: 2026-10-23
specs_relacionadas: ["05-crear-sprints", "06-cerrar-sprint", "07-story-points"]
---

# Spec 10 — Métricas de avance del Sprint y del proyecto

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero ver las métricas de avance, para saber cuánto planificamos, cuánto completamos y a qué velocidad trabajamos. ([#14](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/14))

## Objetivo

Calcular las métricas de avance que pide la consigna (Story Points planificados y completados, velocidad del equipo y porcentaje de historias completadas), por Sprint y para el proyecto, para que el equipo vea si cumple lo que planifica y pueda planificar mejor los próximos Sprints. Son la base del Dashboard (spec 13) y del reporte (spec 14).

Origen: Consigna — Requerimiento 7 (Métricas).

## Alcance

**Entra:**
- Métricas de cada Sprint del proyecto: historias y Story Points planificados y completados, porcentaje de historias completadas e historias sin estimar.
- Métricas del proyecto: velocidad del equipo y totales de los Sprints cerrados.
- Congelar los Story Points de cada historia al cerrar un Sprint, para que sus métricas no cambien después.

**NO entra:**
- Métricas de esfuerzo (horas) y de defectos (specs 12 y 11).
- Gráficos (Dashboard, spec 13) y reportes (spec 14).
- Pronósticos (por ejemplo, cuántos Sprints faltan para terminar el backlog).
- Pantallas del frontend (todavía no existe; se hacen cuando se arme la interfaz).

## Entradas

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |

Las métricas no se cargan: se calculan con los Sprints, sus historias y sus Story Points en el momento de la consulta.

## Salidas esperadas

Una sola respuesta con:

1. **Por cada Sprint** del proyecto (ordenados por fecha de inicio): su número, su estado (`Abierto` o `Cerrado`), si sus métricas son **parciales** (Sprint abierto), las historias planificadas y completadas, los Story Points planificados y completados, el porcentaje de historias completadas y la cantidad de historias sin estimar.
2. **Del proyecto:** la cantidad de Sprints cerrados, la **velocidad** del equipo (o que todavía no hay velocidad), y los totales de los Sprints cerrados: Story Points planificados y completados, historias planificadas y completadas, y porcentaje de historias completadas.

**Error:** si el proyecto no existe, se informa "El proyecto no existe".

## Reglas de negocio

- RN-1: Las **historias planificadas** de un Sprint son las asignadas a él (spec 05); las **completadas**, las que se completaron en él (spec 06).
- RN-2: Los **Story Points planificados** de un Sprint son la suma de los Story Points de sus historias planificadas que están estimadas; los **completados**, la suma de los de sus historias completadas que están estimadas. Las historias sin estimar no suman y se informan aparte.
- RN-3: En un Sprint **abierto** se usan los Story Points actuales de cada historia y sus métricas se marcan como **parciales**: el Sprint todavía no terminó.
- RN-4: Al **cerrar** un Sprint (spec 06) se **congelan** los Story Points de cada historia planificada tal como están en ese momento (o "sin estimar"). Las métricas de un Sprint cerrado usan esos valores congelados: si después se re-estima una historia que no se terminó, las métricas de ese Sprint no cambian.
- RN-5: El **porcentaje de historias completadas** de un Sprint es `completadas / planificadas × 100`, redondeado al entero más cercano con la mitad hacia arriba, igual que en la spec 03 (12,5 % → 13 %). Un Sprint sin historias tiene 0 %.
- RN-6: La **velocidad** del equipo es el promedio de los Story Points completados de **todos los Sprints cerrados** del proyecto, con **un decimal** (redondeo con la mitad hacia arriba: 12,25 → 12,3). Los Sprints abiertos no cuentan. Sin Sprints cerrados no hay velocidad.
- RN-7: Los **totales del proyecto** suman los Sprints cerrados: Story Points planificados y completados, historias planificadas y completadas; su porcentaje de historias completadas se calcula como en RN-5 sobre esos totales.
- RN-8: Las métricas se calculan en el momento de la consulta: siempre reflejan los datos actuales (y los congelados de los Sprints cerrados).

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- Los cálculos (RN-2, RN-5, RN-6, RN-7) son funciones del dominio sin base de datos y se desarrollan con TDD (consigna: pruebas unitarias para el cálculo de métricas).
- Los redondeos se hacen con números enteros, sin errores de coma flotante, como en la spec 03.

## Casos límite

- Proyecto sin Sprints: lista de Sprints vacía, velocidad "sin datos" y totales en 0.
- Solo hay un Sprint abierto: aparece con métricas parciales; la velocidad sigue sin datos.
- Sprint cerrado sin historias: 0 planificados, 0 completados, 0 %; cuenta en la velocidad con 0 Story Points completados.
- Sprint cerrado con todas sus historias sin estimar: 0 Story Points, pero sus historias sí cuentan para el porcentaje.
- Sprints cerrados con 10 y 15 Story Points completados: velocidad 12,5. Con 10, 15 y 12: velocidad 12,3 (37 / 3 = 12,33…).
- Una historia no completada en el Sprint 1 (cerrado) se re-estima y se completa en el Sprint 2: el Sprint 1 cuenta sus Story Points de antes; el Sprint 2, los nuevos.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Consultar las métricas de un proyecto que no existe | No encontrado: "El proyecto no existe" |

## Criterios de aceptación

- CA-1: Se calculan los Story Points planificados de cada Sprint (suma de sus historias estimadas) y el total de los Sprints cerrados.
- CA-2: Se calculan los Story Points completados de cada Sprint (suma de sus historias completadas en él y estimadas) y el total de los Sprints cerrados.
- CA-3: Se calcula la velocidad del equipo: promedio de los Story Points completados de todos los Sprints cerrados, con un decimal; sin Sprints cerrados, se informa que todavía no hay velocidad.
- CA-4: Se calcula el porcentaje de historias completadas de cada Sprint y del total de los Sprints cerrados, redondeado al entero.
- CA-5: Los cálculos se desarrollan con TDD y están cubiertos por tests unitarios Go.
- CA-6: Las métricas de un Sprint abierto se marcan como parciales; las de un Sprint cerrado usan los Story Points congelados al cerrarlo y no cambian si después se re-estiman sus historias.
- CA-7: Consultar las métricas de un proyecto que no existe informa que el proyecto no existe.

Relación con el Issue #14: Story Points planificados → CA-1; completados → CA-2; velocidad → CA-3; porcentaje de historias completadas → CA-4; TDD → CA-5. CA-6 sale de las decisiones de esta spec y CA-7 es el caso de error. El CA-5 se verifica con los tests unitarios del dominio (no tiene escenario BDD: no es un comportamiento visible).

## Escenarios BDD

```gherkin
# language: es
Característica: Métricas de avance del Sprint y del proyecto
  Para saber cuánto planificamos, cuánto completamos y a qué velocidad trabajamos
  Como integrante del equipo
  Quiero ver las métricas de avance de cada Sprint y del proyecto

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-31"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Y que el backlog tiene la historia "Pagar pedido" con prioridad 3
    Y que estimé la historia "Registrar usuario" con 5 Story Points
    Y que estimé la historia "Ver catálogo" con 8 Story Points

  Escenario: Métricas de un Sprint cerrado (CA-1, CA-2, CA-4)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Cuando consulto las métricas del proyecto
    Entonces el Sprint 1 tiene 13 Story Points planificados y 5 completados
    Y el Sprint 1 tiene 2 historias planificadas, 1 completada y 50 % completadas

  Escenario: Las historias sin estimar no suman Story Points pero sí cuentan como historias (CA-1, CA-4)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Pagar pedido"
    Y que completé en el Sprint la historia "Pagar pedido"
    Y que cerré el Sprint
    Cuando consulto las métricas del proyecto
    Entonces el Sprint 1 tiene 5 Story Points planificados y 0 completados
    Y el Sprint 1 tiene 1 historia sin estimar
    Y el Sprint 1 tiene 2 historias planificadas, 1 completada y 50 % completadas

  Escenario: Velocidad con un decimal (CA-3)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-15" al "2026-11-28" con la historia "Ver catálogo"
    Y que completé en el Sprint la historia "Ver catálogo"
    Y que cerré el Sprint
    Cuando consulto las métricas del proyecto
    Entonces la velocidad del equipo es 6.5 Story Points por Sprint
    Y el proyecto tiene 2 Sprints cerrados con 13 Story Points planificados y 13 completados

  Escenario: Sin Sprints cerrados no hay velocidad (CA-3)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando consulto las métricas del proyecto
    Entonces todavía no hay velocidad del equipo
    Y el proyecto tiene 0 Sprints cerrados con 0 Story Points planificados y 0 completados

  Escenario: Un Sprint cerrado sin completar nada baja la velocidad (CA-3)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-15" al "2026-11-28" con la historia "Ver catálogo"
    Y que cerré el Sprint
    Cuando consulto las métricas del proyecto
    Entonces la velocidad del equipo es 2.5 Story Points por Sprint
    Y el proyecto tiene 2 historias planificadas, 1 completada y 50 % completadas

  Escenario: Métricas parciales del Sprint abierto (CA-6)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Ver catálogo"
    Cuando consulto las métricas del proyecto
    Entonces el Sprint 1 está abierto y sus métricas son parciales
    Y el Sprint 1 tiene 13 Story Points planificados y 8 completados

  Escenario: Los Story Points de un Sprint cerrado quedan congelados (CA-6)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Ver catálogo"
    Y que cerré el Sprint
    Y que estimé la historia "Ver catálogo" con 3 Story Points
    Cuando consulto las métricas del proyecto
    Entonces el Sprint 1 tiene 8 Story Points planificados y 0 completados

  Escenario: Proyecto sin Sprints (CA-3, CA-6)
    Cuando consulto las métricas del proyecto
    Entonces no hay Sprints en las métricas
    Y todavía no hay velocidad del equipo

  Escenario: Métricas de un proyecto que no existe (CA-7)
    Cuando consulto las métricas del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"
```

Los pasos que no indican proyecto se refieren al de los Antecedentes, "Tienda Online". "Que estimé la historia … con N Story Points" es el paso de la spec 07; los de Sprints, los de las specs 05 y 06. La velocidad se escribe con punto decimal (6.5), como sale en la respuesta.

## Dependencias

- Spec 05 (Sprints), spec 06 (completar y cerrar) y spec 07 (Story Points): implementadas. Esta historia agrega a la spec 06 que, al cerrar, se congelen los Story Points (RN-4).
- Chequeo de independencia (INVEST): todo lo que usa ya está implementado; no depende de ninguna historia pendiente.
- La usan: spec 13 (Dashboard: grafica estas métricas) y spec 14 (reporte).

## Tareas

- [ ] `[P1]` Funciones de cálculo en `backend/dominio/` (métricas de un Sprint, velocidad, totales y redondeos) con TDD.
- [ ] `[P1]` Congelar los Story Points de las historias al cerrar un Sprint (dominio y base de datos).
- [ ] `[P1]` Caso de uso y ruta de la API para consultar las métricas de un proyecto.
- [ ] `[P1]` Escenarios de `backend/features/metricas-avance.feature` con sus pasos.
- [ ] `[P2]` Actualizar `backend/CLAUDE.md` (endpoint, JSON, specs relacionadas).

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented`.
