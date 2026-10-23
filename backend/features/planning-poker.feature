# language: es
# Spec 08 — Planning Poker (Issue #12).
# En las tablas, "(vacío)" representa un valor que falta: Gherkin no puede escribirlo.
# Los pasos que no nombran un proyecto se refieren al proyecto de los Antecedentes.
Característica: Planning Poker
  Para llegar a una estimación consensuada sin que unos influyan a otros
  Como integrante del equipo
  Quiero estimar las historias del backlog votando en secreto con la escala de Story Points

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Y que el proyecto tiene al integrante "Ana" con email "ana@mail.com" y rol "Developer"
    Y que el proyecto tiene al integrante "Beto" con email "beto@mail.com" y rol "Developer"
    Y que el proyecto tiene al integrante "Caro" con email "caro@mail.com" y rol "Product Owner"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1

  Escenario: Abrir una votación para una historia (CA-1)
    Cuando abro una votación para la historia "Registrar usuario"
    Entonces la votación está abierta en la ronda 1
    Y faltan votar "Ana" y "Beto"

  Escenario: No se pueden abrir dos votaciones para la misma historia (CA-1)
    Dado que abrí una votación para la historia "Registrar usuario"
    Cuando abro una votación para la historia "Registrar usuario"
    Entonces se informa el error "La historia ya tiene una votación abierta"

  Escenario: No se puede abrir una votación para una historia completada (CA-1)
    Dado que el backlog tiene la historia "Pagar pedido" con prioridad 2 en estado "Completada"
    Cuando abro una votación para la historia "Pagar pedido"
    Entonces se informa el error "La historia ya está completada"

  Escenario: No se puede abrir una votación en un proyecto sin Developers (CA-1)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el backlog del proyecto "App Turnos" tiene la historia "Pedir turno" en estado "Pendiente" sin estimar
    Cuando abro una votación para la historia "Pedir turno" del proyecto "App Turnos"
    Entonces se informa el error "El proyecto no tiene Developers para votar"

  Escenario: Un Developer vota (CA-2)
    Dado que abrí una votación para la historia "Registrar usuario"
    Cuando "Ana" vota 5
    Entonces ya votó "Ana"
    Y falta votar "Beto"

  Escenario: Un Developer cambia su voto antes de revelar (CA-2)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 3
    Cuando "Ana" vota 8
    Y "Beto" vota 8
    Y revelo la ronda
    Entonces veo los votos: "Ana" 8 y "Beto" 8

  Escenario: Un integrante que no es Developer no puede votar (CA-2)
    Dado que abrí una votación para la historia "Registrar usuario"
    Cuando "Caro" vota 5
    Entonces se informa el error "Solo los Developers pueden votar"

  Esquema del escenario: Solo se vota con valores de la escala (CA-2)
    Dado que abrí una votación para la historia "Registrar usuario"
    Cuando "Ana" vota "<valor>"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | valor   | resultado                                                                   |
      | 1       | votado                                                                      |
      | 21      | votado                                                                      |
      | 4       | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | 0       | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | 2.5     | Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21    |
      | (vacío) | Los Story Points son obligatorios                                           |

  Escenario: Los votos están ocultos hasta revelar la ronda (CA-3)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 5
    Cuando veo la votación
    Entonces ya votó "Ana"
    Y falta votar "Beto"
    Y no se ve ningún valor de los votos

  Escenario: Revelar la ronda cuando votaron todos los Developers (CA-4)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 3
    Y que "Beto" votó 13
    Cuando revelo la ronda
    Entonces veo los votos: "Ana" 3 y "Beto" 13

  Escenario: No se puede revelar si falta un voto (CA-4)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 5
    Cuando revelo la ronda
    Entonces se informa el error "Faltan votos: todavía no votaron todos los Developers"

  Escenario: Una ronda revelada no admite más votos (CA-4)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 5
    Y que "Beto" votó 5
    Y que revelé la ronda
    Cuando "Ana" vota 8
    Entonces se informa el error "La ronda ya fue revelada"

  Escenario: Detectar diferencias entre los votos (CA-5)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 3
    Y que "Beto" votó 13
    Cuando revelo la ronda
    Entonces hay diferencias entre los votos
    Y el voto más bajo es 3 de "Ana"
    Y el voto más alto es 13 de "Beto"

  Escenario: Sin diferencias cuando todos votan lo mismo (CA-5)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 5
    Y que "Beto" votó 5
    Cuando revelo la ronda
    Entonces no hay diferencias entre los votos

  Escenario: Dos Developers empatan en el voto más bajo (CA-5)
    Dado que el proyecto tiene al integrante "Dani" con email "dani@mail.com" y rol "Developer"
    Y que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 3
    Y que "Beto" votó 8
    Y que "Dani" votó 3
    Cuando revelo la ronda
    Entonces el voto más bajo es 3 de "Ana" y "Dani"
    Y el voto más alto es 8 de "Beto"

  Escenario: Empezar una nueva ronda (CA-6)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 3
    Y que "Beto" votó 13
    Y que revelé la ronda
    Cuando empiezo una nueva ronda
    Entonces la votación está abierta en la ronda 2
    Y faltan votar "Ana" y "Beto"
    Y en la ronda 1 se ven los votos: "Ana" 3 y "Beto" 13

  Escenario: No se puede empezar una nueva ronda sin revelar la actual (CA-6)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 3
    Cuando empiezo una nueva ronda
    Entonces se informa el error "Primero hay que revelar la ronda actual"

  Escenario: Registrar la estimación acordada (CA-7)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 5
    Y que "Beto" votó 5
    Y que revelé la ronda
    Cuando registro la estimación acordada de 5
    Entonces la historia "Registrar usuario" tiene 5 Story Points
    Y la votación está cerrada con la estimación acordada de 5

  Escenario: Acordar un valor que nadie votó (CA-7)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 3
    Y que "Beto" votó 8
    Y que revelé la ronda
    Cuando registro la estimación acordada de 5
    Entonces la historia "Registrar usuario" tiene 5 Story Points

  Escenario: No se puede acordar sin revelar la ronda (CA-7)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 5
    Cuando registro la estimación acordada de 5
    Entonces se informa el error "Primero hay que revelar la ronda actual"
    Y la historia "Registrar usuario" está sin estimar

  Escenario: La estimación acordada tiene que ser de la escala (CA-7)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 5
    Y que "Beto" votó 5
    Y que revelé la ronda
    Cuando registro la estimación acordada de 4
    Entonces se informa el error "Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21"
    Y la historia "Registrar usuario" está sin estimar

  Escenario: Una votación cerrada no admite más votos (CA-7)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 5
    Y que "Beto" votó 5
    Y que revelé la ronda
    Y que registré la estimación acordada de 5
    Cuando "Ana" vota 8
    Entonces se informa el error "La votación está cerrada"

  Escenario: Después de cerrar una votación se puede abrir otra para la misma historia (CA-7)
    Dado que abrí una votación para la historia "Registrar usuario"
    Y que "Ana" votó 5
    Y que "Beto" votó 5
    Y que revelé la ronda
    Y que registré la estimación acordada de 5
    Cuando abro una votación para la historia "Registrar usuario"
    Entonces la votación está abierta en la ronda 1

  Escenario: Abrir una votación para una historia que no existe (CA-8)
    Cuando abro una votación para la historia con identificador 999
    Entonces se informa el error "La historia no existe"

  Escenario: Ver una votación que no existe (CA-8)
    Cuando veo la votación con identificador 999
    Entonces se informa el error "La votación no existe"

  Escenario: Un integrante de otro proyecto no puede votar (CA-8)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y agrego al proyecto "App Turnos" al integrante "Eli" con email "eli@mail.com" y rol "Developer"
    Y que abrí una votación para la historia "Registrar usuario"
    Cuando "Eli" vota 5
    Entonces se informa el error "El integrante no existe"
