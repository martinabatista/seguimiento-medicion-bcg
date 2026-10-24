# language: es
# Spec 09 — Registrar esfuerzo (Issue #13).
# En las tablas, "(vacío)" representa un texto vacío: Gherkin no puede escribirlo. Las horas se
# escriben con coma decimal; los pasos las mandan a la API como número.
Característica: Registrar esfuerzo
  Para comparar después el esfuerzo estimado con el real
  Como integrante del equipo
  Quiero registrar las horas que trabajé en una historia

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Y que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que hoy es "2026-11-20"

  Escenario: Registrar esfuerzo con todos sus datos (CA-1)
    Cuando registro 2,5 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Diseño del formulario"
    Entonces el registro de esfuerzo queda guardado con un identificador
    Y el registro tiene integrante "Ana Pérez", fecha "2026-11-18", actividad "Diseño del formulario" y 2,5 horas

  Escenario: Registrar esfuerzo en una historia completada (CA-2)
    Dado que la historia "Registrar usuario" pasa a estado "Completada"
    Cuando registro 1 hora de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Pruebas finales"
    Entonces el registro de esfuerzo queda guardado con un identificador

  Escenario: Cualquier rol puede registrar esfuerzo (CA-2)
    Dado que el proyecto tiene al integrante "Bruno Díaz" con email "bruno@mail.com" y rol "Scrum Master"
    Cuando registro 1 hora de "Bruno Díaz" el "2026-11-18" en "Registrar usuario" con la actividad "Refinamiento"
    Entonces el registro de esfuerzo queda guardado con un identificador

  Esquema del escenario: Validar las horas (CA-3)
    Cuando registro <horas> horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Desarrollo"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | horas   | resultado                                                                 |
      | 0,01    | registrado                                                                |
      | 1,5     | registrado                                                                |
      | 24      | registrado                                                                |
      | 24,01   | Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales |
      | 0       | Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales |
      | -2      | Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales |
      | 1,125   | Las horas deben ser un número mayor a 0 y hasta 24, con hasta 2 decimales |
      | (vacío) | Las horas son obligatorias                                                |

  Escenario: Consultar el esfuerzo de una historia ordenado por fecha (CA-4)
    Dado que registré 2 horas de "Ana Pérez" el "2026-11-19" en "Registrar usuario" con la actividad "Backend"
    Y que registré 1,5 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Diseño"
    Y que registré 0,25 horas de "Ana Pérez" el "2026-11-19" en "Registrar usuario" con la actividad "Revisión"
    Cuando consulto el esfuerzo de "Registrar usuario"
    Entonces veo las actividades "Diseño", "Backend" y "Revisión" en ese orden
    Y el total de horas es 3,75

  Escenario: Consultar una historia sin esfuerzo registrado (CA-4)
    Cuando consulto el esfuerzo de "Registrar usuario"
    Entonces no hay registros de esfuerzo
    Y el total de horas es 0

  Escenario: El total no pierde precisión con decimales (CA-4)
    Dado que registré 0,1 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "A"
    Y que registré 0,2 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "B"
    Cuando consulto el esfuerzo de "Registrar usuario"
    Entonces el total de horas es 0,3

  Esquema del escenario: Validar la fecha (CA-5)
    Cuando registro 1 hora de "Ana Pérez" el "<fecha>" en "Registrar usuario" con la actividad "Desarrollo"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | fecha      | resultado                                             |
      | 2026-11-01 | registrado                                            |
      | 2026-11-20 | registrado                                            |
      | 2026-11-21 | La fecha no puede ser posterior a hoy                 |
      | 2026-10-31 | La fecha debe estar dentro de las fechas del proyecto |
      | 2026-02-30 | La fecha no es válida                                 |
      | (vacío)    | La fecha es obligatoria                               |

  Escenario: Registrar el último día del proyecto (CA-5)
    Dado que hoy es "2026-12-20"
    Cuando registro 1 hora de "Ana Pérez" el "2026-12-15" en "Registrar usuario" con la actividad "Cierre"
    Entonces el registro de esfuerzo queda guardado con un identificador

  Escenario: Rechazar una fecha posterior al fin del proyecto (CA-5)
    Dado que hoy es "2026-12-20"
    Cuando registro 1 hora de "Ana Pérez" el "2026-12-16" en "Registrar usuario" con la actividad "Cierre"
    Entonces se informa el error "La fecha debe estar dentro de las fechas del proyecto"

  Esquema del escenario: Validar la actividad (CA-5)
    Cuando registro 1 hora de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con una actividad de <largo> caracteres
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                        |
      | 300   | registrado                                       |
      | 301   | La actividad no puede superar los 300 caracteres |
      | 0     | La actividad es obligatoria                      |

  Escenario: Rechazar una actividad con solo espacios (CA-5)
    Cuando registro 1 hora de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "   "
    Entonces se informa el error "La actividad es obligatoria"
    Y no hay esfuerzo registrado en "Registrar usuario"

  Escenario: Si hay varios datos inválidos se informa el primero (CA-5)
    Cuando registro 0 horas de "Ana Pérez" el "2026-11-21" en "Registrar usuario" con la actividad ""
    Entonces se informa el error "La fecha no puede ser posterior a hoy"

  Escenario: Eliminar un registro de esfuerzo (CA-6)
    Dado que registré 2 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Backend"
    Y que registré 1 hora de "Ana Pérez" el "2026-11-19" en "Registrar usuario" con la actividad "Pruebas"
    Cuando elimino el registro de la actividad "Backend"
    Y consulto el esfuerzo de "Registrar usuario"
    Entonces veo solo la actividad "Pruebas"
    Y el total de horas es 1

  Escenario: Eliminar dos veces el mismo registro (CA-6)
    Dado que registré 2 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Backend"
    Cuando elimino el registro de la actividad "Backend"
    Y elimino otra vez el registro de la actividad "Backend"
    Entonces se informa el error "El registro de esfuerzo no existe"

  Escenario: No se puede eliminar del backlog una historia con esfuerzo (CA-7)
    Dado que registré 2 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Backend"
    Cuando elimino la historia "Registrar usuario"
    Entonces se informa el error "La historia tiene esfuerzo registrado"
    Y el backlog tiene solo la historia "Registrar usuario"

  Escenario: Sin registros, la historia se puede eliminar del backlog (CA-7)
    Dado que registré 2 horas de "Ana Pérez" el "2026-11-18" en "Registrar usuario" con la actividad "Backend"
    Y que eliminé el registro de la actividad "Backend"
    Cuando elimino la historia "Registrar usuario"
    Entonces el backlog está vacío

  Escenario: Registrar en una historia que no existe (CA-8)
    Cuando registro 1 hora de "Ana Pérez" el "2026-11-18" en la historia con identificador 999 con la actividad "Desarrollo"
    Entonces se informa el error "La historia no existe"

  Escenario: Registrar con un integrante de otro proyecto (CA-8)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto "App Turnos" tiene al integrante "Carla Gómez" con email "carla@mail.com" y rol "Developer"
    Cuando registro 1 hora de "Carla Gómez" el "2026-11-18" en "Registrar usuario" con la actividad "Desarrollo"
    Entonces se informa el error "El integrante no existe"
    Y no hay esfuerzo registrado en "Registrar usuario"

  Escenario: Registrar sin indicar el integrante (CA-8)
    Cuando registro 1 hora sin integrante el "2026-11-18" en "Registrar usuario" con la actividad "Desarrollo"
    Entonces se informa el error "El integrante es obligatorio"

  Escenario: Consultar el esfuerzo en un proyecto que no existe (CA-8)
    Cuando consulto el esfuerzo de la historia 1 del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"
