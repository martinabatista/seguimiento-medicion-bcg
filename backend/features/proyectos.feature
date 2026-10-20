# language: es
# Spec 01 — Crear y modificar proyectos (Issue #4).
# En las tablas, "(vacío)" y "(solo espacios)" representan esos textos: Gherkin no puede escribirlos.
Característica: Crear y modificar proyectos
  Para administrar el seguimiento de un proyecto de software
  Como integrante del equipo
  Quiero crear, modificar y consultar proyectos

  Escenario: Crear un proyecto con todos sus datos (CA-1)
    Cuando creo un proyecto con nombre "Tienda Online", descripción "Venta de ropa", inicio "2026-11-01" y fin "2026-12-15"
    Entonces el proyecto queda registrado con un identificador
    Y el proyecto tiene nombre "Tienda Online", descripción "Venta de ropa", inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Crear un proyecto que dura un solo día (CA-1)
    Cuando creo un proyecto con nombre "Hackatón", descripción "Evento de un día", inicio "2026-11-01" y fin "2026-11-02"
    Entonces el proyecto queda registrado con un identificador

  Escenario: Modificar el nombre de un proyecto (CA-2)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Cuando cambio su nombre a "Tienda Online v2"
    Entonces el proyecto tiene nombre "Tienda Online v2"

  Escenario: Modificar las fechas de un proyecto (CA-2)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Cuando cambio sus fechas a inicio "2026-11-10" y fin "2027-01-31"
    Entonces el proyecto tiene inicio "2026-11-10" y fin "2027-01-31"

  Escenario: Una modificación inválida no cambia el proyecto (CA-2)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Cuando cambio sus fechas a inicio "2026-11-10" y fin "2026-11-05"
    Entonces se informa el error "La fecha de finalización debe ser posterior a la de inicio"
    Y el proyecto tiene inicio "2026-11-01" y fin "2026-12-15"

  Esquema del escenario: Rechazar fechas inválidas (CA-3)
    Cuando creo un proyecto con nombre "Tienda Online", descripción "Venta de ropa", inicio "<inicio>" y fin "<fin>"
    Entonces se informa el error "<mensaje>"
    Y no se registra ningún proyecto

    Ejemplos:
      | inicio     | fin        | mensaje                                                    |
      | 2026-11-10 | 2026-11-10 | La fecha de finalización debe ser posterior a la de inicio |
      | 2026-11-10 | 2026-11-01 | La fecha de finalización debe ser posterior a la de inicio |
      | (vacío)    | 2026-11-01 | Las fechas de inicio y finalización son obligatorias       |
      | 2026-11-10 | (vacío)    | Las fechas de inicio y finalización son obligatorias       |
      | 2026-02-30 | 2026-03-10 | La fecha no es válida                                      |

  Esquema del escenario: Rechazar un proyecto sin nombre (CA-4)
    Cuando creo un proyecto con nombre "<nombre>", descripción "Venta de ropa", inicio "2026-11-01" y fin "2026-12-15"
    Entonces se informa el error "El nombre es obligatorio"
    Y no se registra ningún proyecto

    Ejemplos:
      | nombre           |
      | (vacío)          |
      | (solo espacios)  |

  Escenario: Aceptar un nombre muy largo (CA-4)
    Cuando creo un proyecto con un nombre de 300 caracteres, descripción "Venta de ropa", inicio "2026-11-01" y fin "2026-12-15"
    Entonces el proyecto queda registrado con un identificador

  Escenario: Aceptar un nombre repetido (CA-4)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Cuando creo un proyecto con nombre "Tienda Online", descripción "Otra tienda", inicio "2027-01-01" y fin "2027-02-01"
    Entonces hay 2 proyectos registrados

  Esquema del escenario: Validar el largo de la descripción (CA-5)
    Cuando creo un proyecto con nombre "Tienda Online", una descripción de <largo> caracteres, inicio "2026-11-01" y fin "2026-12-15"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                          |
      | 500   | registrado                                         |
      | 501   | La descripción no puede superar los 500 caracteres |
      | 0     | La descripción es obligatoria                      |

  Escenario: Listar proyectos en orden de creación (CA-6)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Y que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Cuando listo los proyectos
    Entonces veo los proyectos "Tienda Online" y "App Turnos" en ese orden

  Escenario: Listar sin proyectos cargados (CA-6)
    Cuando listo los proyectos
    Entonces la lista de proyectos está vacía

  Escenario: Ver un proyecto por su identificador (CA-6)
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Cuando veo ese proyecto por su identificador
    Entonces el proyecto tiene nombre "Tienda Online", descripción "Venta de ropa", inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Ver un proyecto que no existe (CA-7)
    Cuando veo el proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Modificar un proyecto que no existe (CA-7)
    Cuando cambio el nombre del proyecto con identificador 999 a "Fantasma"
    Entonces se informa el error "El proyecto no existe"
