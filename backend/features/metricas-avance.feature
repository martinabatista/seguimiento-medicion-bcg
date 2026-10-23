# language: es
# Spec 10 — Métricas de avance del Sprint y del proyecto (Issue #14).
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
