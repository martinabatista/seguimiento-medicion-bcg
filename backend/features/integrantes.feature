# language: es
# Spec 02 — Registrar integrantes del proyecto (Issue #5).
# En las tablas, "(vacío)" representa un texto vacío y "(solo espacios)" uno con solo espacios:
# Gherkin no puede escribirlos.
Característica: Registrar integrantes del proyecto
  Para poder asignarles estimaciones, esfuerzo y trabajo
  Como integrante del equipo
  Quiero registrar a los integrantes de un proyecto

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"

  Escenario: Agregar un integrante con todos sus datos (CA-1)
    Cuando agrego al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Entonces el integrante queda registrado con un identificador
    Y el integrante tiene nombre "Ana Pérez", email "ana@mail.com" y rol "Developer"

  Esquema del escenario: Agregar integrantes con cada rol (CA-1)
    Cuando agrego al integrante "Ana Pérez" con email "ana@mail.com" y rol "<rol>"
    Entonces el integrante queda registrado con un identificador
    Y el integrante tiene rol "<rol>"

    Ejemplos:
      | rol           |
      | Product Owner |
      | Scrum Master  |
      | Developer     |

  Escenario: El email se guarda normalizado (CA-1)
    Cuando agrego al integrante "Ana Pérez" con email " Ana@Mail.COM " y rol "Developer"
    Entonces el integrante queda registrado con un identificador
    Y el integrante tiene email "ana@mail.com"

  Escenario: Listar integrantes en el orden en que se agregaron (CA-2)
    Dado que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Y que el proyecto tiene al integrante "Bruno Díaz" con email "bruno@mail.com" y rol "Scrum Master"
    Cuando listo los integrantes
    Entonces veo los integrantes "Ana Pérez" y "Bruno Díaz" en ese orden

  Escenario: Listar un proyecto sin integrantes (CA-2)
    Cuando listo los integrantes
    Entonces la lista de integrantes está vacía

  Escenario: Cada proyecto ve solo sus integrantes (CA-2)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando listo los integrantes del proyecto "App Turnos"
    Entonces la lista de integrantes está vacía

  Escenario: Ver un integrante por su identificador (CA-2)
    Dado que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando veo ese integrante por su identificador
    Entonces el integrante tiene nombre "Ana Pérez", email "ana@mail.com" y rol "Developer"

  Esquema del escenario: No registrar dos veces el mismo email en un proyecto (CA-3)
    Dado que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando agrego al integrante "Ana P." con email "<email>" y rol "Scrum Master"
    Entonces se informa el error "El integrante ya está registrado en el proyecto"
    Y el proyecto tiene 1 integrante

    Ejemplos:
      | email        |
      | ana@mail.com |
      | ANA@Mail.com |

  Escenario: El mismo email con un espacio adelante también está repetido (CA-3)
    Dado que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando agrego al integrante "Ana P." con email " ana@mail.com" y rol "Scrum Master"
    Entonces se informa el error "El integrante ya está registrado en el proyecto"
    Y el proyecto tiene 1 integrante

  Escenario: El mismo email puede estar en otro proyecto (CA-3)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando agrego al proyecto "App Turnos" al integrante "Ana Pérez" con email "ana@mail.com" y rol "Product Owner"
    Entonces el integrante queda registrado con un identificador

  Escenario: Dos integrantes con el mismo nombre y distinto email (CA-3)
    Dado que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando agrego al integrante "Ana Pérez" con email "ana.perez@otro.com" y rol "Developer"
    Entonces el proyecto tiene 2 integrantes

  Esquema del escenario: Validar el nombre (CA-4)
    Cuando agrego un integrante con un nombre de <largo> caracteres, email "ana@mail.com" y rol "Developer"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                     |
      | 100   | registrado                                    |
      | 101   | El nombre no puede superar los 100 caracteres |
      | 0     | El nombre es obligatorio                      |

  Escenario: Rechazar un nombre con solo espacios (CA-4)
    Cuando agrego al integrante "   " con email "ana@mail.com" y rol "Developer"
    Entonces se informa el error "El nombre es obligatorio"
    Y la lista de integrantes está vacía

  Esquema del escenario: Validar el email (CA-4)
    Cuando agrego al integrante "Ana Pérez" con email "<email>" y rol "Developer"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | email              | resultado               |
      | ana@mail.com       | registrado              |
      | ana.p@mail.com.ar  | registrado              |
      | (vacío)            | El email es obligatorio |
      | (solo espacios)    | El email es obligatorio |
      | ana.mail.com       | El email no es válido   |
      | ana@@mail.com      | El email no es válido   |
      | @mail.com          | El email no es válido   |
      | ana@mail           | El email no es válido   |
      | ana@.com           | El email no es válido   |
      | ana@mail.          | El email no es válido   |
      | ana pérez@mail.com | El email no es válido   |

  Esquema del escenario: Validar el largo del email (CA-4)
    Cuando agrego al integrante "Ana Pérez" con un email de <largo> caracteres y rol "Developer"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                    |
      | 254   | registrado                                   |
      | 255   | El email no puede superar los 254 caracteres |

  Escenario: Si hay varios datos inválidos se informa el primero (CA-4)
    Cuando agrego al integrante "" con email "ana.mail.com" y rol "Tester"
    Entonces se informa el error "El nombre es obligatorio"

  Esquema del escenario: Validar el rol (CA-5)
    Cuando agrego al integrante "Ana Pérez" con email "ana@mail.com" y rol "<rol>"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | rol          | resultado             |
      | Scrum Master | registrado            |
      | (vacío)      | El rol es obligatorio |
      | Tester       | El rol no es válido   |
      | developer    | El rol no es válido   |

  Escenario: Agregar un integrante a un proyecto que no existe (CA-6)
    Cuando agrego un integrante al proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Listar los integrantes de un proyecto que no existe (CA-6)
    Cuando listo los integrantes del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Ver un integrante que no existe (CA-6)
    Cuando veo el integrante con identificador 999
    Entonces se informa el error "El integrante no existe"

  Escenario: Un integrante de otro proyecto se trata como inexistente (CA-6)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto tiene al integrante "Ana Pérez" con email "ana@mail.com" y rol "Developer"
    Cuando veo ese integrante desde el proyecto "App Turnos"
    Entonces se informa el error "El integrante no existe"
