# language: es
# Spec 06 — Registrar historias completadas, cerrar Sprint y consultar Sprints anteriores (Issue #9).
Característica: Registrar historias completadas, cerrar Sprint y consultar Sprints anteriores
  Para dejar constancia de lo logrado en cada iteración
  Como integrante del equipo
  Quiero completar historias, cerrar Sprints y consultar los anteriores

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-31"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2

  Escenario: Completar una historia del Sprint (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Cuando completo en el Sprint la historia "Registrar usuario"
    Entonces en el Sprint la historia "Registrar usuario" está completada y "Ver catálogo" no
    Y la historia "Registrar usuario" del backlog está en estado "Completada"

  Escenario: Deshacer una historia completada (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando deshago en el Sprint la historia "Registrar usuario"
    Entonces en el Sprint ninguna historia está completada
    Y la historia "Registrar usuario" del backlog está en estado "En progreso"

  Escenario: Completar una historia que no está en el Sprint (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Ver catálogo"
    Cuando completo en el Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia no está asignada a este Sprint"

  Escenario: Completar dos veces la misma historia (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando completo en el Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está completada"

  Escenario: Deshacer una historia que no está completada (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando deshago en el Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia no está completada"

  Escenario: Cerrar un Sprint con historias sin terminar (CA-2)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando cierro el Sprint
    Entonces el Sprint está cerrado con las historias "Registrar usuario" y "Ver catálogo"
    Y en el Sprint la historia "Registrar usuario" está completada y "Ver catálogo" no

  Escenario: Las historias no completadas quedan libres para otro Sprint (CA-2)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Ver catálogo"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-15" al "2026-11-28"
    Cuando asigno al Sprint la historia "Ver catálogo"
    Entonces el Sprint tiene la historia "Ver catálogo"

  Escenario: Cerrar un Sprint vacío antes de su fecha de fin (CA-2)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Y que hoy es "2026-11-05"
    Cuando cierro el Sprint
    Entonces el Sprint está cerrado sin historias

  Esquema del escenario: Un Sprint cerrado no admite cambios (CA-3)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que cerré el Sprint
    Cuando <accion>
    Entonces se informa el error "El Sprint está cerrado"

    Ejemplos:
      | accion                                             |
      | cambio su objetivo a "Otro objetivo"               |
      | asigno al Sprint la historia "Ver catálogo"        |
      | quito del Sprint la historia "Registrar usuario"   |
      | completo en el Sprint la historia "Registrar usuario" |
      | cierro el Sprint                                   |

  Escenario: Consultar los Sprints anteriores (CA-4)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-15" al "2026-11-28" con la historia "Ver catálogo"
    Y que cerré el Sprint
    Y que existe un Sprint del "2026-11-29" al "2026-12-12"
    Cuando consulto los Sprints anteriores
    Entonces veo los Sprints anteriores que empiezan el "2026-11-15" y el "2026-11-01" en ese orden
    Y en el Sprint anterior que empieza el "2026-11-01" la historia "Registrar usuario" está completada

  Escenario: Consultar Sprints anteriores sin Sprints cerrados (CA-4)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando consulto los Sprints anteriores
    Entonces no hay Sprints anteriores

  Escenario: El registro del Sprint cerrado no cambia (CA-4)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Cuando la historia "Registrar usuario" pasa a estado "Pendiente"
    Entonces en el Sprint la historia "Registrar usuario" está completada

  Esquema del escenario: El backlog no cambia "Completada" de una historia en un Sprint abierto (CA-5)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando la historia "Registrar usuario" pasa a estado "<estado>"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | estado      | resultado                                                        |
      | En progreso | modificada                                                       |
      | Completada  | La historia está en un Sprint abierto: se completa desde el Sprint |

  Escenario: Cerrar un Sprint que no existe (CA-6)
    Cuando cierro el Sprint con identificador 999
    Entonces se informa el error "El Sprint no existe"

  Escenario: Completar una historia que no existe (CA-6)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando completo en el Sprint la historia con identificador 999
    Entonces se informa el error "La historia no existe"
