# language: es
# Spec 13 — Dashboard del proyecto (Issue #19), etapa 1: la consulta del Dashboard.
# Los pasos que no nombran un proyecto se refieren al proyecto de los Antecedentes.
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
