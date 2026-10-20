# language: es
# Spec 04 — Administrar el Product Backlog (Issue #7).
# En las tablas, "(vacío)" representa un texto vacío: Gherkin no puede escribirlo.
# Los pasos que no nombran un proyecto ("el backlog") se refieren al proyecto de los Antecedentes.
Característica: Administrar el Product Backlog
  Para mantener priorizado el trabajo pendiente
  Como integrante del equipo
  Quiero crear, modificar, eliminar y listar las historias del backlog de un proyecto

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Crear una historia con todos sus datos (CA-1)
    Cuando creo en el backlog la historia "Registrar usuario" con descripción "Alta de clientes", prioridad 1 y criterios "Valida el email"
    Entonces la historia queda registrada con un identificador
    Y la historia tiene título "Registrar usuario", descripción "Alta de clientes", prioridad 1 y criterios "Valida el email"
    Y la historia está en estado "Pendiente" y sin estimar

  Escenario: Crear una historia sin criterios de aceptación (CA-1)
    Cuando creo en el backlog la historia "Ver catálogo" con descripción "Listado de productos", prioridad 2 y sin criterios
    Entonces la historia queda registrada con un identificador
    Y la historia no tiene criterios de aceptación

  Escenario: Modificar la prioridad y el estado de una historia (CA-2)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 3
    Cuando cambio su prioridad a 1 y su estado a "En progreso"
    Entonces la historia tiene prioridad 1 y estado "En progreso"

  Escenario: Una modificación inválida no cambia la historia (CA-2)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 3
    Cuando cambio su prioridad a 0
    Entonces se informa el error "La prioridad debe ser un número entero mayor o igual a 1"
    Y la historia tiene prioridad 3 y estado "Pendiente"

  Escenario: Eliminar una historia del backlog (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Cuando elimino la historia "Registrar usuario"
    Entonces el backlog tiene solo la historia "Ver catálogo"

  Escenario: Eliminar dos veces la misma historia (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando elimino la historia "Registrar usuario"
    Y elimino otra vez la historia "Registrar usuario"
    Entonces se informa el error "La historia no existe"

  Escenario: Listar el backlog ordenado por prioridad (CA-4)
    Dado que el backlog tiene la historia "Pagar pedido" con prioridad 3
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Cuando listo el backlog
    Entonces veo las historias "Registrar usuario", "Ver catálogo" y "Pagar pedido" en ese orden

  Escenario: A igual prioridad se respeta el orden de creación (CA-4)
    Dado que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Y que el backlog tiene la historia "Buscar producto" con prioridad 2
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando listo el backlog
    Entonces veo las historias "Registrar usuario", "Ver catálogo" y "Buscar producto" en ese orden

  Escenario: Listar un backlog vacío (CA-4)
    Cuando listo el backlog
    Entonces el backlog está vacío

  Escenario: Cada proyecto ve solo su backlog (CA-4)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando listo el backlog del proyecto "App Turnos"
    Entonces el backlog está vacío

  Escenario: Ver una historia por su identificador (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando veo esa historia por su identificador
    Entonces la historia tiene título "Registrar usuario" y prioridad 1

  Esquema del escenario: Validar título, descripción y criterios (CA-5)
    Cuando creo en el backlog una historia con un título de <titulo> caracteres, una descripción de <descripcion> caracteres y criterios de <criterios> caracteres
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | titulo | descripcion | criterios | resultado                                                          |
      | 150    | 1000        | 2000      | registrada                                                         |
      | 151    | 10          | 0         | El título no puede superar los 150 caracteres                      |
      | 0      | 10          | 0         | El título es obligatorio                                           |
      | 10     | 1001        | 0         | La descripción no puede superar los 1000 caracteres                |
      | 10     | 0           | 0         | La descripción es obligatoria                                      |
      | 10     | 10          | 2001      | Los criterios de aceptación no pueden superar los 2000 caracteres  |

  Escenario: Rechazar un título con solo espacios (CA-5)
    Cuando creo en el backlog la historia "   " con descripción "Alta de clientes", prioridad 1 y sin criterios
    Entonces se informa el error "El título es obligatorio"
    Y el backlog está vacío

  Esquema del escenario: Validar la prioridad (CA-6)
    Cuando creo en el backlog la historia "Registrar usuario" con descripción "Alta de clientes", prioridad "<prioridad>" y sin criterios
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | prioridad | resultado                                                  |
      | 1         | registrada                                                 |
      | 1000      | registrada                                                 |
      | 0         | La prioridad debe ser un número entero mayor o igual a 1   |
      | -2        | La prioridad debe ser un número entero mayor o igual a 1   |
      | 2.5       | La prioridad debe ser un número entero mayor o igual a 1   |
      | (vacío)   | La prioridad es obligatoria                                |

  Esquema del escenario: Validar el estado al modificar (CA-6)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando cambio su estado a "<estado>"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | estado      | resultado              |
      | Pendiente   | modificada             |
      | En progreso | modificada             |
      | Completada  | modificada             |
      | Terminada   | El estado no es válido |
      | pendiente   | El estado no es válido |

  Escenario: Crear una historia en un proyecto que no existe (CA-7)
    Cuando creo una historia en el backlog del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Ver una historia que no existe (CA-7)
    Cuando veo la historia con identificador 999
    Entonces se informa el error "La historia no existe"

  Escenario: Una historia de otro proyecto se trata como inexistente (CA-7)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando veo esa historia desde el proyecto "App Turnos"
    Entonces se informa el error "La historia no existe"
