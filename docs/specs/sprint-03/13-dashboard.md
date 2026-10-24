---
titulo: Dashboard del proyecto
estado: in-review    # draft | in-review | approved | implemented
unidad: global
issue: "#19"
sprint: 03
autor: Florencia Cabeza
aprobada_por: ""
fecha: 2026-10-24
specs_relacionadas: ["03-estado-proyecto", "05-crear-sprints", "10-metricas-avance"]
---

# Spec 13 — Dashboard del proyecto

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero un Dashboard con el estado del proyecto y gráficos de las métricas, para entender de un vistazo cómo avanzamos. ([#19](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/19))

## Objetivo

Mostrar en una sola pantalla cómo va el proyecto: unos pocos **indicadores** (números que se leen de un vistazo) y tres **gráficos** de las métricas de avance. El Dashboard **no inventa métricas nuevas**: muestra las que ya calculan el estado del proyecto (spec 03) y las métricas de avance (spec 10), juntas y listas para graficar.

Origen: Consigna — Requerimiento 8 (Dashboard) y *Sprint 3: implementar la visualización*.

## Alcance

**Entra:**
- **Indicadores:** el estado del proyecto y los días que faltan; el Sprint en curso y los días que le quedan; el avance del backlog (historias y Story Points completados); la velocidad del equipo; las historias sin estimar.
- **Gráfico 1 — Story Points por Sprint:** planificados y completados en cada Sprint (barras).
- **Gráfico 2 — Velocidad:** Story Points completados en cada Sprint cerrado y su promedio (barras con una línea).
- **Gráfico 3 — Backlog por estado:** historias pendientes, en progreso y completadas (dona).
- Una consulta del backend que devuelve todo lo anterior junto y ya calculado.
- La pantalla del Dashboard que lo muestra.

Se implementa en **dos etapas**, con dos PRs de código:
1. **Backend** (ahora): la consulta del Dashboard, con TDD y BDD (CA-1 a CA-8).
2. **Pantalla** (cuando exista la interfaz web base, Issue #15): los indicadores y los tres gráficos (CA-9). Este PR cierra el Issue #19.

**NO entra:**
- Métricas de esfuerzo y de calidad (horas, desviación, defectos; specs 11 y 12): se pueden sumar en el Sprint 4 si queda tiempo.
- Filtros (por ejemplo, ver solo un Sprint), exportar o imprimir el Dashboard (el reporte en PDF es la spec 14).
- Pronósticos (por ejemplo, cuándo se terminaría el backlog).
- Elegir la librería de gráficos: se decide con el equipo en la etapa 2.

## Entradas

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |

Nada más: el Dashboard no se carga, se calcula en el momento de la consulta con la fecha de hoy.

## Salidas esperadas

Una sola respuesta con:

**Indicadores**
1. **Proyecto:** nombre, fechas, estado (`No iniciado`, `En curso` o `Finalizado`, spec 03) y los días que faltan (RN-3).
2. **Sprint en curso:** su número, su Sprint Goal, su fecha de fin, los días que le quedan y sus Story Points planificados y completados hasta ahora (parciales); o que no hay Sprint en curso.
3. **Avance del backlog:** historias completadas, total de historias y porcentaje; Story Points completados y totales; historias sin estimar.
4. **Velocidad del equipo**, o que todavía no hay (no hay Sprints cerrados).

**Datos de los gráficos**
5. **Story Points por Sprint:** por cada Sprint, en orden de fecha de inicio: su nombre ("Sprint 1"…), los Story Points planificados y completados, y si es parcial (Sprint abierto).
6. **Velocidad:** por cada Sprint cerrado, en orden de fecha de inicio: su nombre y sus Story Points completados; y el promedio (la velocidad).
7. **Backlog por estado:** cantidad de historias pendientes, en progreso y completadas, en ese orden.

**Error:** si el proyecto no existe, se informa "El proyecto no existe".

## Reglas de negocio

- RN-1: Solo se puede consultar el Dashboard de un proyecto existente.
- RN-2: El Dashboard **reutiliza** los cálculos de la spec 03 (estado, Sprint en curso, resumen del backlog) y de la spec 10 (métricas por Sprint, velocidad), con sus mismas reglas: redondeos, Story Points congelados al cerrar un Sprint, métricas parciales del Sprint abierto, historias sin estimar que no suman. Así, el mismo dato vale lo mismo en el Dashboard y en las otras consultas.
- RN-3: **Días que faltan del proyecto**, contados en días corridos desde hoy:
  - `No iniciado`: los días que faltan para empezar (fecha de inicio − hoy).
  - `En curso`: los días que faltan para terminar (fecha de finalización − hoy); 0 si termina hoy.
  - `Finalizado`: no se informan días.
- RN-4: **Días que le quedan al Sprint en curso:** fecha de fin del Sprint − hoy; 0 si termina hoy.
- RN-5: El gráfico de **Story Points por Sprint** incluye **todos** los Sprints del proyecto; el abierto va marcado como parcial (spec 10 RN-3).
- RN-6: El gráfico de **velocidad** incluye **solo los Sprints cerrados**, y su promedio es la velocidad de la spec 10 (RN-6 de esa spec). Sin Sprints cerrados, el gráfico está vacío y no hay velocidad.
- RN-7: El gráfico del **backlog por estado** cuenta todas las historias del backlog según su estado, como el resumen de la spec 03 (RN-4 de esa spec).
- RN-8: Un proyecto sin historias, sin Sprints o sin velocidad **igual tiene Dashboard**: los números van en 0, los gráficos de Sprints vacíos y se informa que no hay Sprint en curso ni velocidad. Nunca es un error.
- RN-9: Todo se calcula en el momento de la consulta, con la fecha de hoy: siempre refleja los datos actuales.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); **la pantalla solo dibuja** lo que devuelve la consulta del Dashboard: no suma, no promedia ni cuenta días.
- Los cálculos nuevos (días que faltan, armado de los datos de cada gráfico, RN-3 a RN-8) viven en el dominio y se desarrollan con TDD.
- La fecha de hoy es la del servidor, como en la spec 03; en los tests se fija (paso "que hoy es…").
- Los gráficos tienen que poder leerse sin colores: cada barra o porción lleva su número y su nombre.

## Casos límite

- El proyecto termina hoy: está `En curso` y faltan 0 días.
- El proyecto empieza mañana: está `No iniciado` y falta 1 día para empezar.
- El proyecto ya terminó: está `Finalizado` y no se informan días.
- El Sprint en curso termina hoy: le quedan 0 días.
- Hay Sprints pero ninguno incluye hoy (o el que lo incluye está cerrado): no hay Sprint en curso.
- Proyecto sin Sprints: los gráficos de Story Points por Sprint y de velocidad están vacíos, y no hay velocidad.
- Proyecto con un Sprint abierto y ninguno cerrado: el gráfico de Story Points por Sprint muestra el abierto (parcial); el de velocidad está vacío.
- Backlog vacío: 0 historias, 0 % completadas, 0 Story Points y el gráfico del backlog con 0, 0 y 0.
- Todas las historias sin estimar: 0 Story Points totales, y el indicador de sin estimar muestra la cantidad.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| El proyecto no existe | No encontrado: "El proyecto no existe" |

No hay otros errores: los datos que faltan se muestran como 0, vacíos o "no hay" (RN-8).

## Criterios de aceptación

- CA-1: El Dashboard muestra el estado del proyecto y los días que faltan para empezar o para terminar (ninguno si ya finalizó).
- CA-2: El Dashboard muestra el Sprint en curso con los días que le quedan y sus Story Points planificados y completados hasta ahora, o que no hay Sprint en curso.
- CA-3: El Dashboard muestra el avance del backlog: historias completadas sobre el total con su porcentaje, Story Points completados sobre los totales e historias sin estimar.
- CA-4: El Dashboard muestra la velocidad del equipo, o que todavía no hay porque no hay Sprints cerrados.
- CA-5: El Dashboard trae los datos del gráfico de Story Points planificados y completados por Sprint, con el Sprint abierto marcado como parcial.
- CA-6: El Dashboard trae los datos del gráfico de velocidad: los Story Points completados de cada Sprint cerrado y su promedio.
- CA-7: El Dashboard trae los datos del gráfico del backlog por estado: pendientes, en progreso y completadas.
- CA-8: Un proyecto sin historias ni Sprints igual tiene Dashboard (con 0, vacíos y "no hay"); un proyecto que no existe informa que no existe.
- CA-9: La pantalla del Dashboard muestra los indicadores y los tres gráficos con los datos de la consulta, cada barra o porción con su número (etapa 2, con la interfaz web base).

Relación con los criterios iniciales del Issue #19: "muestra información del estado del proyecto" → CA-1 a CA-4 y CA-9; "representaciones gráficas de algunas métricas de avance" → CA-5 a CA-7 y CA-9 (se eligieron Story Points planificados y completados, velocidad y porcentaje de historias completadas, de la spec 10).

## Escenarios BDD

Los escenarios cubren la consulta del Dashboard (etapa 1). CA-9 se verifica en la etapa 2, sobre la pantalla.

```gherkin
# language: es
Característica: Dashboard del proyecto
  Para entender de un vistazo cómo avanzamos
  Como integrante del equipo
  Quiero ver el estado del proyecto y los datos de los gráficos de sus métricas

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-31"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Y que el backlog tiene la historia "Pagar pedido" con prioridad 3
    Y que estimé la historia "Registrar usuario" con 5 Story Points
    Y que estimé la historia "Ver catálogo" con 8 Story Points

  Escenario: Proyecto en curso con los días que faltan para terminar (CA-1)
    Dado que hoy es "2026-11-20"
    Cuando consulto el Dashboard
    Entonces el proyecto está "En curso" y faltan 41 días para terminar

  Escenario: Proyecto que termina hoy (CA-1)
    Dado que hoy es "2026-12-31"
    Cuando consulto el Dashboard
    Entonces el proyecto está "En curso" y faltan 0 días para terminar

  Escenario: Proyecto que todavía no empezó (CA-1)
    Dado que hoy es "2026-10-25"
    Cuando consulto el Dashboard
    Entonces el proyecto está "No iniciado" y faltan 7 días para empezar

  Escenario: Proyecto finalizado (CA-1)
    Dado que hoy es "2027-01-05"
    Cuando consulto el Dashboard
    Entonces el proyecto está "Finalizado" sin días que falten

  Escenario: Sprint en curso con los días que le quedan (CA-2)
    Dado que existe un Sprint del "2026-11-15" al "2026-11-28" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que hoy es "2026-11-20"
    Cuando consulto el Dashboard
    Entonces el Sprint en curso es el Sprint 1 y le quedan 8 días
    Y el Sprint en curso tiene 13 Story Points planificados y 5 completados

  Escenario: No hay Sprint en curso (CA-2)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que hoy es "2026-11-20"
    Cuando consulto el Dashboard
    Entonces no hay Sprint en curso

  Escenario: Avance del backlog (CA-3)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que hoy es "2026-11-05"
    Cuando consulto el Dashboard
    Entonces el backlog tiene 1 de 3 historias completadas, el 33 %
    Y el backlog tiene 5 de 13 Story Points completados y 1 historia sin estimar

  Escenario: Velocidad del equipo con dos Sprints cerrados (CA-4)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-15" al "2026-11-28" con la historia "Ver catálogo"
    Y que completé en el Sprint la historia "Ver catálogo"
    Y que cerré el Sprint
    Cuando consulto el Dashboard
    Entonces la velocidad del equipo es 6.5

  Escenario: Todavía no hay velocidad (CA-4)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando consulto el Dashboard
    Entonces todavía no hay velocidad del equipo

  Escenario: Gráfico de Story Points por Sprint, con el abierto parcial (CA-5)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-15" al "2026-11-28" con la historia "Ver catálogo"
    Cuando consulto el Dashboard
    Entonces el gráfico de Story Points muestra el "Sprint 1" con 13 planificados y 5 completados
    Y el gráfico de Story Points muestra el "Sprint 2" con 8 planificados y 0 completados, parcial

  Escenario: El gráfico de velocidad solo muestra los Sprints cerrados (CA-6)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-15" al "2026-11-28" con la historia "Ver catálogo"
    Cuando consulto el Dashboard
    Entonces el gráfico de velocidad muestra solo el "Sprint 1" con 5 Story Points completados
    Y el promedio del gráfico de velocidad es 5

  Escenario: Gráfico del backlog por estado (CA-7)
    Dado que el backlog tiene la historia "Buscar producto" con prioridad 4 en estado "En progreso"
    Y que el backlog tiene la historia "Iniciar sesión" con prioridad 5 en estado "Completada"
    Cuando consulto el Dashboard
    Entonces el gráfico del backlog muestra 3 pendientes, 1 en progreso y 1 completada

  Escenario: Dashboard de un proyecto sin historias ni Sprints (CA-8)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Cuando consulto el Dashboard del proyecto "App Turnos"
    Entonces el gráfico del backlog muestra 0 pendientes, 0 en progreso y 0 completadas
    Y los gráficos de Sprints están vacíos
    Y todavía no hay velocidad del equipo
    Y no hay Sprint en curso

  Escenario: Dashboard de un proyecto que no existe (CA-8)
    Cuando consulto el Dashboard del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"
```

Los pasos que no indican proyecto se refieren al proyecto de los Antecedentes, "Tienda Online". Los pasos de proyectos, backlog, Story Points, Sprints y fecha de hoy (`que existe el proyecto…`, `que el backlog tiene la historia…`, `que estimé la historia…`, `que existe un Sprint…`, `que completé en el Sprint…`, `que cerré el Sprint`, `que hoy es…`) ya existen en las specs 01, 03, 04, 05, 06 y 07 y se reutilizan.

## Dependencias

- Spec 03 (estado del proyecto) y spec 10 (métricas de avance): implementadas. El Dashboard reutiliza sus cálculos.
- Specs 05 y 06 (Sprints) y 07 (Story Points): implementadas; se usan en los escenarios.
- Etapa 2: la interfaz web base (Issue #15, Sprint 2, tarea de las tres integrantes, incluida la dueña de esta spec).
- La usa: spec 14 (reporte en PDF), que puede reutilizar estos indicadores.
- Chequeo de independencia (INVEST): la etapa 1 solo usa historias implementadas. La etapa 2 depende de la #15, de un Sprint anterior y compartida con la dueña de esta spec.

## Tareas

**Etapa 1 — backend:**
- [ ] `[P1]` Días que faltan y armado de los datos de los tres gráficos (RN-3 a RN-8) en `backend/dominio/` (TDD), reutilizando el estado (spec 03) y las métricas (spec 10).
- [ ] `[P1]` Caso de uso consultar el Dashboard del proyecto.
- [ ] `[P1]` Ruta de la API del Dashboard en `backend/api/`.
- [ ] `[P1]` Escenarios de `backend/features/dashboard.feature` con sus pasos.
- [ ] `[P2]` Actualizar `backend/CLAUDE.md` (endpoint, JSON, specs relacionadas).

**Etapa 2 — pantalla (con la interfaz web base, #15):**
- [ ] `[P1]` Pantalla del Dashboard con los indicadores y los tres gráficos (barras, barras con línea y dona), cada uno con sus números.
- [ ] `[P2]` Elegir con el equipo la librería de gráficos y justificarla en esta spec.

## Definition of Done

- [ ] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [ ] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [ ] Criterios de aceptación verificados (CA-1 a CA-8 en la etapa 1; CA-9 en la etapa 2).
- [ ] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [ ] Spec en estado `implemented` (al terminar la etapa 2).
