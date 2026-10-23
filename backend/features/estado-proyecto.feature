# language: es
# Spec 03 — Consultar el estado de un proyecto (Issue #6): CA-1, CA-3 y CA-4.
# El CA-2 (Sprint en curso) lo implementa la spec 05 (CA-5): sus escenarios están en sprints.feature.
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

  Escenario: Resumen del backlog con historias en todos los estados (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" en estado "Completada" con 5 Story Points
    Y que el backlog tiene la historia "Ver catálogo" en estado "Completada" con 3 Story Points
    Y que el backlog tiene la historia "Buscar producto" en estado "En progreso" con 8 Story Points
    Y que el backlog tiene la historia "Pagar pedido" en estado "Pendiente" sin estimar
    Cuando consulto el estado del proyecto
    Entonces el resumen del backlog tiene 4 historias: 1 pendiente, 1 en progreso y 2 completadas
    Y el porcentaje de historias completadas es 50
    Y los Story Points son 16 en total y 8 completados, con 1 historia sin estimar

  Escenario: Una historia estimada que no está completada suma solo al total (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" en estado "Pendiente" con 13 Story Points
    Y que el backlog tiene la historia "Ver catálogo" en estado "Completada" sin estimar
    Cuando consulto el estado del proyecto
    Entonces los Story Points son 13 en total y 0 completados, con 1 historia sin estimar

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
