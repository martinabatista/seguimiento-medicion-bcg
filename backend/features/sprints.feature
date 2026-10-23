# language: es
# Spec 05 — Crear Sprints, definir Sprint Goal y asignar historias (Issue #8).
# En las tablas, "(vacío)" representa un texto vacío: Gherkin no puede escribirlo.
Característica: Crear Sprints, definir Sprint Goal y asignar historias
  Para planificar el trabajo de cada iteración
  Como integrante del equipo
  Quiero crear Sprints con su objetivo y asignarles historias del backlog

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-31"

  Escenario: Crear un Sprint (CA-1)
    Cuando creo un Sprint del "2026-11-01" al "2026-11-14" con objetivo "Registrar clientes"
    Entonces el Sprint queda registrado como el Sprint 1, abierto y sin historias

  Escenario: Los Sprints se numeran en orden (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando creo un Sprint del "2026-11-15" al "2026-11-28" con objetivo "Catálogo"
    Entonces el Sprint queda registrado como el Sprint 2, abierto y sin historias

  Esquema del escenario: Validar las fechas del Sprint (CA-1)
    Dado que existe un Sprint del "2026-11-15" al "2026-11-28"
    Cuando creo un Sprint del "<inicio>" al "<fin>" con objetivo "Objetivo"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | inicio     | fin        | resultado                                                  |
      | 2026-11-01 | 2026-11-14 | registrado                                                 |
      | 2026-11-29 | 2026-12-31 | registrado                                                 |
      | 2026-11-29 | 2026-11-30 | registrado                                                 |
      | 2026-10-31 | 2026-11-10 | El Sprint debe estar dentro de las fechas del proyecto     |
      | 2026-12-20 | 2027-01-02 | El Sprint debe estar dentro de las fechas del proyecto     |
      | 2026-11-10 | 2026-11-15 | El Sprint se superpone con otro Sprint del proyecto        |
      | 2026-11-28 | 2026-12-05 | El Sprint se superpone con otro Sprint del proyecto        |
      | 2026-12-01 | 2026-12-01 | La fecha de finalización debe ser posterior a la de inicio |
      | (vacío)    | 2026-12-01 | Las fechas de inicio y finalización son obligatorias       |

  Esquema del escenario: Validar el Sprint Goal (CA-2)
    Cuando creo un Sprint del "2026-11-01" al "2026-11-14" con un objetivo de <largo> caracteres
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                          |
      | 300   | registrado                                         |
      | 301   | El Sprint Goal no puede superar los 300 caracteres |
      | 0     | El Sprint Goal es obligatorio                      |

  Escenario: Modificar el Sprint Goal (CA-2)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando cambio su objetivo a "Registrar y loguear clientes"
    Entonces el Sprint tiene objetivo "Registrar y loguear clientes"

  Escenario: Asignar historias a un Sprint (CA-3)
    Dado que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando asigno al Sprint las historias "Ver catálogo" y "Registrar usuario"
    Entonces el Sprint tiene las historias "Registrar usuario" y "Ver catálogo" en ese orden

  Escenario: Quitar una historia de un Sprint (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando quito del Sprint la historia "Registrar usuario"
    Entonces el Sprint no tiene historias

  Escenario: No se puede asignar una historia completada (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1 en estado "Completada"
    Y que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando asigno al Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está completada"

  Escenario: No se puede quitar una historia que no está en el Sprint (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando quito del Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia no está asignada a este Sprint"

  Escenario: Listar los Sprints por fecha de inicio (CA-3)
    Dado que existe un Sprint del "2026-11-15" al "2026-11-28"
    Y que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando listo los Sprints del proyecto
    Entonces veo los Sprints que empiezan el "2026-11-01" y el "2026-11-15" en ese orden

  Escenario: Una historia no puede estar en dos Sprints abiertos (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Y que existe un Sprint del "2026-11-15" al "2026-11-28"
    Cuando asigno al segundo Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está asignada a un Sprint abierto"

  Escenario: Asignar dos veces la misma historia al mismo Sprint (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando asigno al Sprint la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está asignada a un Sprint abierto"

  Esquema del escenario: Sprint en curso en el estado del proyecto (CA-5)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Y que existe un Sprint del "2026-11-16" al "2026-11-28"
    Y que hoy es "<hoy>"
    Cuando consulto el estado del proyecto
    Entonces el Sprint en curso es "<sprint>"

    Ejemplos:
      | hoy        | sprint          |
      | 2026-11-01 | Sprint 1        |
      | 2026-11-14 | Sprint 1        |
      | 2026-11-15 | ninguno         |
      | 2026-11-20 | Sprint 2        |
      | 2026-12-10 | ninguno         |

  Escenario: No se puede eliminar del backlog una historia de un Sprint abierto (CA-6)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando elimino la historia "Registrar usuario"
    Entonces se informa el error "La historia está asignada a un Sprint abierto"

  Escenario: Se puede eliminar una historia después de quitarla del Sprint (CA-6)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Registrar usuario"
    Cuando quito del Sprint la historia "Registrar usuario"
    Y elimino la historia "Registrar usuario"
    Entonces el backlog está vacío

  Escenario: Crear un Sprint en un proyecto que no existe (CA-7)
    Cuando creo un Sprint en el proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Ver un Sprint que no existe (CA-7)
    Cuando veo el Sprint con identificador 999
    Entonces se informa el error "El Sprint no existe"

  Escenario: Asignar una historia que no existe (CA-7)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando asigno al Sprint la historia con identificador 999
    Entonces se informa el error "La historia no existe"
