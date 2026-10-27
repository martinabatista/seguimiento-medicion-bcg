---
titulo: Métricas de esfuerzo y calidad
estado: implemented  # draft | in-review | approved | implemented
unidad: global
issue: "#18"
sprint: 03
autor: Martina Batista
aprobada_por: "Florencia Cabeza"
fecha: 2026-10-26
specs_relacionadas: ["04-product-backlog", "06-cerrar-sprint", "07-story-points", "09-registro-esfuerzo", "10-metricas-avance"]
---

# Spec 12 — Métricas de esfuerzo y calidad

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero ver las métricas de esfuerzo, para evaluar la precisión de nuestras estimaciones comparando las horas que estimamos con las que trabajamos. ([#18](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/18))

## Objetivo

Comparar el esfuerzo estimado con el esfuerzo real, como pide la consigna: cada historia tiene sus **horas estimadas**, el registro de esfuerzo (spec 09) da las **horas reales**, y la **desviación** muestra cuánto nos equivocamos al estimar, por historia, por Sprint y para el proyecto. Con eso el equipo puede ajustar sus próximas estimaciones. Es la base de la parte de esfuerzo del reporte (spec 14).

Origen: Consigna — Requerimiento 5 (registro de esfuerzo: "comparar posteriormente el esfuerzo estimado con el esfuerzo real") y Requerimiento 7 (Métricas: horas estimadas, horas reales y desviación).

## Alcance

**Entra:**
- Horas estimadas de cada historia del backlog: cargarlas, cambiarlas y quitarlas.
- Métricas de esfuerzo por historia, por Sprint y del proyecto: horas estimadas, horas reales y desviación (en horas y en porcentaje).
- La pantalla: todo en la pantalla de Métricas (spec 10), en una parte de esfuerzo con una fila por historia, donde se cargan sus horas estimadas y se ven sus horas reales y su desviación, más los totales por Sprint y del proyecto. Así no se toca la pantalla del Backlog.

**NO entra:**
- Métricas de defectos: cantidad de defectos detectados y resueltos (spec 11, CA-4). Se sacaron de esta spec para que no dependa de la HU-11, del mismo Sprint.
- Horas estimadas por tarea o por integrante: la estimación es de la historia completa.
- Convertir Story Points a horas: son dos estimaciones distintas y no se relacionan.
- Gráficos (Dashboard, spec 13) y reportes (spec 14).

## Entradas

**Cargar o cambiar las horas estimadas de una historia:**

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |
| Historia | identificador | Sí | Tiene que ser una historia del backlog de ese proyecto (spec 04), no `Completada`. |
| Horas estimadas | número | Sí | Mayor a 0, con hasta 2 decimales, sin máximo. |

**Quitar las horas estimadas de una historia:** proyecto e historia, con las mismas validaciones. La historia vuelve a quedar sin horas estimadas.

**Consultar las métricas de esfuerzo:** solo el proyecto. Las métricas no se cargan: se calculan en el momento de la consulta.

## Salidas esperadas

- Al cargar, cambiar o quitar: la historia con sus horas estimadas (o sin horas estimadas).
- Al consultar, una sola respuesta con:
    1. **Por cada historia** del backlog (en el orden del backlog, spec 04): su estado, sus horas estimadas (o "sin horas estimadas"), sus horas reales y, si entra en la desviación (RN-7), su desviación en horas y en porcentaje.
    2. **Por cada Sprint** del proyecto (por fecha de inicio): su número, su estado, si es **parcial** (Sprint abierto), las horas estimadas y reales de las historias completadas en él que tienen horas estimadas, su desviación en horas y en porcentaje (o "sin datos"), y cuántas historias completadas en él no tienen horas estimadas.
    3. **Del proyecto:** las horas estimadas de todo el backlog, las horas reales de todo el esfuerzo registrado, la desviación de las historias completadas con horas estimadas (o "sin datos") y cuántas historias completadas no tienen horas estimadas.

**Errores:** ver Condiciones de error.

## Reglas de negocio

- RN-1: Una historia arranca **sin horas estimadas**. Se le pueden cargar, cambiar por otro valor o quitar (vuelve a sin horas estimadas).
- RN-2: Las horas estimadas son un número **mayor a 0, con hasta 2 decimales y sin máximo**: una historia puede llevar varios días de trabajo. Se guardan y se suman exactas, sin redondear (como las horas de la spec 09).
- RN-3: No se pueden cargar, cambiar ni quitar las horas estimadas de una historia `Completada`, igual que sus Story Points (spec 07 RN-4): la estimación se compara con lo trabajado, y cambiarla después de terminar arreglaría la desviación a mano. Si se deshace la completada (spec 06), se puede volver a estimar.
- RN-4: Cargar horas estimadas no cambia ningún otro dato de la historia (título, prioridad, estado, Story Points). Modificar la historia desde el backlog (spec 04) o estimar sus Story Points (spec 07) no cambia sus horas estimadas.
- RN-5: Las **horas reales** de una historia son el total de sus registros de esfuerzo (spec 09 RN-6), sin importar la fecha ni el Sprint en que se cargaron. Se pueden seguir cargando después de completarla (spec 09 RN-1), y las métricas reflejan siempre lo cargado.
- RN-6: La **desviación** compara lo real con lo estimado:
    - en horas: `horas reales − horas estimadas`. Positiva: se trabajó **más** de lo estimado; negativa: **menos**.
    - en porcentaje: `desviación en horas / horas estimadas × 100`, redondeado al entero más cercano con la mitad alejándose de cero (+12,5 % → +13 %; −12,5 % → −13 %).
- RN-7: Entran en la desviación **solo las historias completadas que tienen horas estimadas**. Una historia sin terminar todavía no tiene todas sus horas: si entrara, parecería que se estimó de más. Las historias no completadas igual muestran sus horas estimadas y sus horas reales hasta ahora.
- RN-8: **Por Sprint:** se usan las historias **completadas en ese Sprint** (spec 06), con sus horas estimadas y sus horas reales (RN-5). La desviación del Sprint se calcula sobre las sumas (RN-6), no como promedio de las desviaciones de cada historia. En un Sprint **abierto** las métricas son **parciales**.
- RN-9: **Del proyecto:** las horas estimadas son la suma de todas las historias del backlog que tienen horas estimadas; las horas reales, la suma de todo el esfuerzo registrado en el proyecto. La desviación del proyecto se calcula sobre las sumas de las historias completadas con horas estimadas.
- RN-10: Si ninguna historia entra en una desviación (RN-7), esa desviación es **"sin datos"**, no 0: 0 significaría que se estimó perfecto. Las historias completadas sin horas estimadas se informan aparte (por Sprint y del proyecto), para saber cuántas quedaron afuera.
- RN-11: La historia tiene que pertenecer al proyecto indicado: una historia de otro proyecto se trata como inexistente (como en spec 04 RN-10).

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- Los cálculos (RN-5 a RN-10) son funciones del dominio sin base de datos y se desarrollan con TDD (consigna: pruebas unitarias para el cálculo de métricas).
- Las horas se manejan en centésimas enteras y el porcentaje se redondea con enteros, sin errores de coma flotante (como en las specs 09 y 10).
- Las horas estimadas se guardan en la historia: hace falta una columna nueva en la tabla de historias, y cada integrante recrea su base local.

## Casos límite

- Horas estimadas 0,01 (el mínimo con 2 decimales): se acepta. 0: se rechaza. Negativas: se rechazan.
- Horas estimadas con 3 decimales (1,125): se rechaza.
- Horas estimadas 500: se acepta (no hay máximo).
- Historia completada con 8 h estimadas y 10 h reales: desviación +2 h y +25 %.
- Historia completada con 8 h estimadas y 6 h reales: desviación −2 h y −25 %.
- Historia completada con 8 h estimadas y sin esfuerzo registrado: horas reales 0, desviación −8 h y −100 % (se terminó sin cargar horas: el dato lo muestra).
- Historia completada con 3 h estimadas y 3,5 h reales: +0,5 h y +17 % (16,67 % se redondea a 17).
- Historia en progreso con 8 h estimadas y 10 h reales: muestra sus horas, pero no tiene desviación.
- Historia completada sin horas estimadas: no tiene desviación y cuenta entre las "completadas sin horas estimadas".
- Proyecto sin historias ni esfuerzo: todo en 0 y desviación "sin datos".
- Sprint con dos historias completadas (8 h estimadas y 10 h reales; 4 h estimadas y 2 h reales): estimadas 12 h, reales 12 h, desviación 0 h y 0 % (se compensan: la desviación es sobre las sumas, RN-8).
- Historia completada y deshecha (vuelve a `En progreso`): se pueden volver a cargar sus horas estimadas y sale de la desviación hasta que se complete otra vez.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Faltan las horas estimadas | Dato inválido: "Las horas estimadas son obligatorias" |
| Horas estimadas que no son un número mayor a 0 o tienen más de 2 decimales | Dato inválido: "Las horas estimadas deben ser un número mayor a 0, con hasta 2 decimales" |
| Cargar, cambiar o quitar las horas estimadas de una historia completada | Conflicto: "La historia ya está completada" (el mismo mensaje de las specs 05 y 07) |
| El proyecto no existe | No encontrado: "El proyecto no existe" |
| La historia no existe o es de otro proyecto | No encontrado: "La historia no existe" |

## Criterios de aceptación

- CA-1: Se pueden cargar, cambiar y quitar las horas estimadas de una historia, con un número mayor a 0 y hasta 2 decimales.
- CA-2: No se pueden cargar, cambiar ni quitar las horas estimadas de una historia completada.
- CA-3: Se calculan las horas estimadas y las horas reales de cada historia, de las historias completadas en cada Sprint y del proyecto.
- CA-4: Se calcula la desviación entre las horas reales y las estimadas, en horas y en porcentaje, de cada historia completada con horas estimadas, de cada Sprint y del proyecto; sin historias que entren, se informa "sin datos".
- CA-5: Las métricas de un Sprint abierto se marcan como parciales, y las historias completadas sin horas estimadas se informan aparte.
- CA-6: Los cálculos se desarrollan con TDD y están cubiertos por tests unitarios Go.
- CA-7: Cargar horas estimadas o consultar las métricas de un proyecto o una historia que no existe informa que no existe.

Relación con el Issue #18: horas estimadas → CA-1, CA-2 y CA-3; horas reales → CA-3; desviación → CA-4; TDD → CA-6. CA-5 sale de las decisiones de esta spec y CA-7 son los casos de error. El CA-6 se verifica con los tests unitarios del dominio (no tiene escenario BDD: no es un comportamiento visible).

## Escenarios BDD

```gherkin
# language: es
Característica: Métricas de esfuerzo
  Para evaluar la precisión de nuestras estimaciones
  Como integrante del equipo
  Quiero comparar las horas que estimamos con las que trabajamos

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-31"
    Y que hoy es "2026-11-20"
    Y que el proyecto tiene la integrante "Ana" con rol "Developer"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2

  Escenario: Cargar las horas estimadas de una historia (CA-1)
    Cuando estimo la historia "Registrar usuario" en 8 horas
    Entonces la historia "Registrar usuario" tiene 8 horas estimadas

  Escenario: Cambiar las horas estimadas, con decimales (CA-1)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Cuando estimo la historia "Registrar usuario" en 6.5 horas
    Entonces la historia "Registrar usuario" tiene 6.5 horas estimadas

  Escenario: Quitar las horas estimadas (CA-1)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Cuando quito las horas estimadas de la historia "Registrar usuario"
    Entonces la historia "Registrar usuario" no tiene horas estimadas

  Esquema del escenario: Validar las horas estimadas (CA-1)
    Cuando estimo la historia "Registrar usuario" en <horas> horas
    Entonces <resultado>

    Ejemplos:
      | horas | resultado                                                                                           |
      | 0.01  | la historia "Registrar usuario" tiene 0.01 horas estimadas                                          |
      | 500   | la historia "Registrar usuario" tiene 500 horas estimadas                                           |
      | 0     | se informa el error "Las horas estimadas deben ser un número mayor a 0, con hasta 2 decimales"      |
      | -3    | se informa el error "Las horas estimadas deben ser un número mayor a 0, con hasta 2 decimales"      |
      | 1.125 | se informa el error "Las horas estimadas deben ser un número mayor a 0, con hasta 2 decimales"      |

  Escenario: Faltan las horas estimadas (CA-1)
    Cuando estimo la historia "Registrar usuario" sin indicar las horas
    Entonces se informa el error "Las horas estimadas son obligatorias"

  Escenario: No se pueden estimar las horas de una historia completada (CA-2)
    Dado que la historia "Registrar usuario" pasa a estado "Completada"
    Cuando estimo la historia "Registrar usuario" en 8 horas
    Entonces se informa el error "La historia ya está completada"

  Escenario: No se pueden quitar las horas estimadas de una historia completada (CA-2)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que la historia "Registrar usuario" pasa a estado "Completada"
    Cuando quito las horas estimadas de la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está completada"

  Escenario: Horas estimadas y reales de cada historia (CA-3)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que Ana registró 3 horas en la historia "Registrar usuario" el "2026-11-10"
    Y que Ana registró 2.5 horas en la historia "Ver catálogo" el "2026-11-11"
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces la historia "Registrar usuario" tiene 8 horas estimadas y 3 horas reales
    Y la historia "Ver catálogo" no tiene horas estimadas y tiene 2.5 horas reales
    Y el proyecto tiene 8 horas estimadas y 5.5 horas reales

  Escenario: Una historia sin terminar no tiene desviación (CA-4)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que Ana registró 10 horas en la historia "Registrar usuario" el "2026-11-10"
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces la historia "Registrar usuario" no tiene desviación
    Y la desviación del proyecto es "sin datos"

  Esquema del escenario: Desviación de una historia completada (CA-4)
    Dado que estimé la historia "Registrar usuario" en <estimadas> horas
    Y que Ana registró <reales> horas en la historia "Registrar usuario" el "2026-11-10"
    Y que la historia "Registrar usuario" pasa a estado "Completada"
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces la historia "Registrar usuario" tiene una desviación de <horas> horas y <porcentaje> %

    Ejemplos:
      | estimadas | reales | horas | porcentaje |
      | 8         | 10     | 2     | 25         |
      | 8         | 6      | -2    | -25        |
      | 3         | 3.5    | 0.5   | 17         |
      | 8         | 8      | 0     | 0          |

  Escenario: Historia completada sin esfuerzo registrado (CA-4)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que la historia "Registrar usuario" pasa a estado "Completada"
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces la historia "Registrar usuario" tiene una desviación de -8 horas y -100 %

  Escenario: La desviación de un Sprint se calcula sobre las sumas (CA-3, CA-4)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que estimé la historia "Ver catálogo" en 4 horas
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que Ana registró 10 horas en la historia "Registrar usuario" el "2026-11-10"
    Y que Ana registró 2 horas en la historia "Ver catálogo" el "2026-11-11"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Ver catálogo"
    Y que cerré el Sprint
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces el Sprint 1 tiene 12 horas estimadas, 12 horas reales y una desviación de 0 horas y 0 %

  Escenario: Las métricas de un Sprint abierto son parciales (CA-5)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que existe un Sprint del "2026-11-15" al "2026-11-28" con la historia "Registrar usuario"
    Y que Ana registró 10 horas en la historia "Registrar usuario" el "2026-11-16"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces las métricas de esfuerzo del Sprint 1 son parciales
    Y el Sprint 1 tiene 8 horas estimadas, 10 horas reales y una desviación de 2 horas y 25 %

  Escenario: Las historias completadas sin horas estimadas se informan aparte (CA-5)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Ver catálogo"
    Y que completé en el Sprint la historia "Ver catálogo"
    Y que cerré el Sprint
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces el Sprint 1 tiene 1 historia completada sin horas estimadas
    Y la desviación del Sprint 1 es "sin datos"
    Y el proyecto tiene 1 historia completada sin horas estimadas

  Escenario: Consultar las métricas de esfuerzo de un proyecto que no existe (CA-7)
    Cuando consulto las métricas de esfuerzo del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Estimar las horas de una historia que no existe (CA-7)
    Cuando estimo la historia con identificador 999 en 8 horas
    Entonces se informa el error "La historia no existe"
```

Los pasos que ya existen (crear proyecto, integrantes, historias, Sprints, completar, cerrar, registrar esfuerzo, `que hoy es`) se reutilizan al implementar; los textos exactos pueden ajustarse a los pasos del `.feature`.

## Trazabilidad

| Criterio | Escenarios BDD (`backend/features/metricas-esfuerzo.feature`) | Tests Go | Código |
|---|---|---|---|
| CA-1 Cargar, cambiar y quitar | Cargar · Cambiar, con decimales · Quitar · Validar las horas estimadas · Faltan las horas estimadas | `TestParsearHorasEstimadas`, `TestEstimarHoras` (dominio y API), `TestHorasEstimadasDeHistoria` | `dominio/horas_estimadas.go`, `api/metricas_esfuerzo.go` (`estimarHoras`, `quitarHorasEstimadas`), `postgres/historias.go` |
| CA-2 Historia completada | No se pueden estimar · No se pueden quitar | `TestEstimarHorasErrores` (dominio y API) | `ServicioBacklog.cambiarSiNoEstaCompletada` |
| CA-3 Horas estimadas y reales | Horas estimadas y reales de cada historia · La desviación de un Sprint se calcula sobre las sumas | `TestCalcularMetricasEsfuerzo`, `TestMetricasEsfuerzoConsultar`, `TestConsultarMetricasEsfuerzo` | `dominio/metricas_esfuerzo.go` (`CalcularMetricasEsfuerzo`), `dominio/servicio_metricas_esfuerzo.go` |
| CA-4 Desviación | Una historia sin terminar no tiene desviación · Desviación de una historia completada · Historia completada sin esfuerzo registrado | `TestCalcularDesviacion`, `TestMetricasEsfuerzoNegativaYSinDatos` | `CalcularDesviacion` |
| CA-5 Sprint parcial y completadas sin horas | Las métricas de un Sprint abierto son parciales · Las historias completadas sin horas estimadas se informan aparte | `TestCalcularMetricasEsfuerzo`, `TestCalcularMetricasEsfuerzoVacio` | `CalcularMetricasEsfuerzo` |
| CA-6 TDD | — (no es un comportamiento visible) | los tests unitarios del dominio, escritos antes del código | — |
| CA-7 No existe | Consultar las métricas de un proyecto que no existe · Estimar las horas de una historia que no existe | `TestEstimarHorasErrores`, `TestMetricasEsfuerzoProyectoInexistente` | `verificarProyecto`, `responderError` |

## Pantalla

`frontend/src/paginas/Metricas.tsx` (sección "Métricas" del menú), debajo de las métricas de avance: las horas estimadas y reales del proyecto con su desviación, una tabla por historia donde se cargan, cambian o quitan sus horas estimadas (una historia completada solo se ve) y una tabla por Sprint. No valida ni calcula nada: muestra lo que manda el backend y sus mensajes de error.

## Dependencias

- Spec 04 (backlog) y spec 07 (Story Points): de Sprints anteriores. Las horas estimadas se cargan en la historia, igual que los Story Points, y se suma la columna nueva a la tabla de historias.
- Spec 06 (completar y cerrar Sprints): del Sprint 1, para saber qué se completó en cada Sprint.
- Spec 09 (registro de esfuerzo): del Sprint 2, para las horas reales.
- Spec 10 (métricas de avance): misma dueña y ya implementada; esta spec sigue su formato (Sprints parciales, "sin datos", redondeo con enteros).
- La usa: spec 14 (reporte), que muestra estas métricas.

Ninguna dependencia es de otra integrante del mismo Sprint o de uno posterior.

## Tareas

- [x] `[P1]` Horas estimadas en la historia: validación (RN-2), regla de la historia completada (RN-3) y caso de uso para cargarlas y quitarlas, con TDD.
- [x] `[P1]` Columna de horas estimadas en la tabla de historias y su lectura y escritura en `backend/postgres/`.
- [x] `[P1]` Funciones de cálculo en `backend/dominio/` (horas reales, desviación en horas y en %, por historia, por Sprint y del proyecto) con TDD.
- [x] `[P1]` Rutas de la API para cargar y quitar las horas estimadas y para consultar las métricas de esfuerzo.
- [x] `[P1]` Escenarios de `backend/features/metricas-esfuerzo.feature` con sus pasos.
- [x] `[P2]` Pantalla: la parte de esfuerzo en la pantalla de Métricas, con la carga de las horas estimadas de cada historia.
- [x] `[P2]` Actualizar `backend/CLAUDE.md` y `frontend/CLAUDE.md` (endpoints, JSON, specs relacionadas, recrear la base).

## Definition of Done

- [x] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [x] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [x] Criterios de aceptación verificados.
- [x] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [x] Spec en estado `implemented`.
