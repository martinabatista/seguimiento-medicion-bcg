# language: es
# Spec 07 — Estimar historias con Story Points (Issue #10).
# En las tablas, "(vacío)" representa un valor que falta: Gherkin no puede escribirlo.
# Los pasos que no nombran un proyecto se refieren al proyecto de los Antecedentes.
Característica: Estimar historias con Story Points
  Para conocer el tamaño relativo del trabajo
  Como integrante del equipo
  Quiero asignar Story Points a las historias del backlog

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Asignar Story Points a una historia sin estimar (CA-1)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando estimo la historia "Registrar usuario" con 5 Story Points
    Entonces la historia "Registrar usuario" tiene 5 Story Points
    Y la historia tiene prioridad 1 y estado "Pendiente"

  Escenario: Estimar una historia en progreso (CA-1)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1 en estado "En progreso"
    Cuando estimo la historia "Registrar usuario" con 8 Story Points
    Entonces la historia "Registrar usuario" tiene 8 Story Points

  Escenario: Cambiar la estimación de una historia (CA-2)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que estimé la historia "Registrar usuario" con 3 Story Points
    Cuando estimo la historia "Registrar usuario" con 13 Story Points
    Entonces la historia "Registrar usuario" tiene 13 Story Points

  Escenario: Quitar la estimación de una historia (CA-2)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que estimé la historia "Registrar usuario" con 3 Story Points
    Cuando quito la estimación de la historia "Registrar usuario"
    Entonces la historia "Registrar usuario" está sin estimar

  Escenario: Quitar la estimación de una historia que ya está sin estimar (CA-2)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando quito la estimación de la historia "Registrar usuario"
    Entonces la historia "Registrar usuario" está sin estimar

  Esquema del escenario: Solo se aceptan valores de la escala (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando estimo la historia "Registrar usuario" con "<valor>" Story Points
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | valor   | resultado                                                                   |
      | 1       | estimada                                                                    |
      | 21      | estimada                                                                    |
      | 0       | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | 4       | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | 22      | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | -3      | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | 2.5     | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | (vacío) | Los Story Points son obligatorios                                           |

  Escenario: Un valor fuera de la escala no cambia la estimación (CA-3)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que estimé la historia "Registrar usuario" con 3 Story Points
    Cuando estimo la historia "Registrar usuario" con 4 Story Points
    Entonces se informa el error "Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21"
    Y la historia "Registrar usuario" tiene 3 Story Points

  Escenario: No se puede estimar una historia completada (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1 en estado "Completada"
    Cuando estimo la historia "Registrar usuario" con 5 Story Points
    Entonces se informa el error "La historia ya está completada"
    Y la historia "Registrar usuario" está sin estimar

  Escenario: No se puede quitar la estimación de una historia completada (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que estimé la historia "Registrar usuario" con 3 Story Points
    Y que existe un Sprint del "2026-11-02" al "2026-11-15" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando quito la estimación de la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está completada"
    Y la historia "Registrar usuario" tiene 3 Story Points

  Escenario: Una historia completada y deshecha se puede volver a estimar (CA-4)
    Dado que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-02" al "2026-11-15" con la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando deshago en el Sprint la historia "Registrar usuario"
    Y estimo la historia "Registrar usuario" con 2 Story Points
    Entonces la historia "Registrar usuario" tiene 2 Story Points

  Escenario: Consultar la escala de Story Points (CA-5)
    Cuando consulto la escala de Story Points
    Entonces la escala es 1, 2, 3, 5, 8, 13 y 21

  Escenario: Estimar en un proyecto que no existe (CA-6)
    Cuando estimo la historia con identificador 1 del proyecto con identificador 999 con 3 Story Points
    Entonces se informa el error "El proyecto no existe"

  Escenario: Estimar una historia que no existe (CA-6)
    Cuando estimo la historia con identificador 999 con 3 Story Points
    Entonces se informa el error "La historia no existe"

  Escenario: Una historia de otro proyecto se trata como inexistente (CA-6)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Cuando estimo la historia "Registrar usuario" desde el proyecto "App Turnos" con 3 Story Points
    Entonces se informa el error "La historia no existe"
    Y la historia "Registrar usuario" está sin estimar
