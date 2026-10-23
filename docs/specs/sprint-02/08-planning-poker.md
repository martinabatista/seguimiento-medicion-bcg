---
titulo: Planning Poker
estado: implemented  # draft | in-review | approved | implemented
unidad: backend
issue: "#12"
sprint: 02
autor: Florencia Cabeza
aprobada_por: "Martina Batista"
fecha: 2026-10-23
specs_relacionadas: ["02-integrantes", "04-product-backlog", "07-story-points"]
---

# Spec 08 — Planning Poker

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero estimar historias con Planning Poker, para llegar a una estimación consensuada sin que unos influyan a otros. ([#12](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/12))

## Objetivo

Estimar una historia **en equipo**: cada Developer vota en secreto con una carta de la escala de Story Points (spec 07), los votos se revelan todos juntos, se detectan las diferencias para que el equipo las discuta, se repite en nuevas rondas si hace falta y, al final, se registra la estimación acordada como Story Points de la historia. Votar en secreto evita el **efecto ancla**: que el primer número que se dice en voz alta arrastre a los demás.

Origen: Consigna — Requerimiento 4 (mecanismo de Planning Poker: registrar estimaciones individuales, mantenerlas ocultas hasta finalizar la votación, mostrarlas, detectar diferencias, realizar nuevas rondas y registrar la estimación acordada).

## Alcance

**Entra:**
- Abrir una votación de Planning Poker para una historia del backlog.
- Votar: cada Developer del proyecto elige un valor de la escala en la ronda actual, y puede cambiarlo antes de revelar.
- Ver la votación: mientras la ronda no se revela, quién votó y quién falta, sin los valores.
- Revelar la ronda cuando votaron todos los Developers: se muestran todos los votos y si hay diferencias.
- Empezar una nueva ronda, conservando las anteriores con sus votos.
- Registrar la estimación acordada: queda como Story Points de la historia y la votación se cierra.

**NO entra:**
- La escala y las reglas para estimar una historia: spec 07 (esta historia las reutiliza).
- Votos de integrantes que no son Developers (el Product Owner explica la historia y el Scrum Master facilita: no estiman).
- Un tiempo límite para votar, o votaciones de varias historias a la vez en una misma sesión.
- Listar todas las votaciones de una historia: se puede sumar cuando se arme la interfaz.
- Pantallas del frontend (Sprint 2, tarea #15).

## Entradas

| Operación | Datos | Validación |
|---|---|---|
| Abrir una votación | Proyecto e historia del backlog | El proyecto y la historia tienen que existir; la historia no puede estar `Completada`; no puede tener otra votación abierta; el proyecto tiene que tener al menos un Developer. |
| Votar | Votación, integrante y valor | El integrante tiene que ser del proyecto y tener rol `Developer`; el valor tiene que ser de la escala de la spec 07 (1, 2, 3, 5, 8, 13, 21); la votación tiene que estar abierta y la ronda actual sin revelar. |
| Ver la votación | Votación | Tiene que existir y ser del proyecto. |
| Revelar la ronda | Votación | La ronda actual no puede estar revelada y tienen que haber votado todos los Developers del proyecto. |
| Nueva ronda | Votación | La ronda actual tiene que estar revelada y la votación, abierta. |
| Registrar la estimación acordada | Votación y valor | La ronda actual tiene que estar revelada; el valor tiene que ser de la escala (no hace falta que alguien lo haya votado). |

## Salidas esperadas

- **Abrir:** la votación queda abierta, en la ronda 1, sin votos, con la lista de Developers que faltan votar.
- **Votar:** se devuelve la votación con el integrante entre los que ya votaron, **sin mostrar su valor**.
- **Ver (ronda sin revelar):** el estado de la votación, el número de ronda, quiénes votaron y quiénes faltan, sin ningún valor; y las rondas anteriores, ya reveladas, con sus votos.
- **Revelar:** todos los votos de la ronda (integrante y valor, en el orden en que se agregaron los integrantes al proyecto), si hay diferencias, el voto más bajo y el más alto, y quiénes los dieron.
- **Nueva ronda:** la votación pasa a la ronda siguiente, sin votos; las rondas anteriores se siguen viendo con sus votos.
- **Registrar la acordada:** la historia queda con esos Story Points (spec 07) y la votación queda cerrada con la estimación acordada.
- **Error:** no cambia nada y se informa un mensaje que explica qué está mal.

## Reglas de negocio

- RN-1: Una votación es de **una** historia del backlog de un proyecto. No se puede abrir para una historia `Completada` (spec 07 RN-4) y una historia no puede tener **dos votaciones abiertas** a la vez.
- RN-2: Solo votan los integrantes del proyecto con rol **`Developer`** (spec 02 RN-5). Para abrir una votación, el proyecto tiene que tener al menos un Developer.
- RN-3: Cada voto es un valor de la escala de Story Points de la spec 07 (RN-1 y RN-2 de esa spec).
- RN-4: Cada Developer tiene **un voto por ronda**. Si vuelve a votar antes de revelar, el voto nuevo **reemplaza** al anterior.
- RN-5: Mientras la ronda no se revela, los valores de los votos **no se muestran a nadie**: solo quién votó y quién falta.
- RN-6: La ronda se puede revelar solo cuando **votaron todos los Developers** del proyecto. Una vez revelada no admite más votos ni cambios.
- RN-7: Al revelar, **hay diferencias** si no todos votaron el mismo valor. Se informan el voto más bajo y el más alto, cada uno con los integrantes que lo dieron (puede haber más de uno).
- RN-8: Después de revelar se puede empezar una **nueva ronda**: arranca sin votos y con el número siguiente. Las rondas anteriores se conservan, con sus votos, y se pueden consultar.
- RN-9: Después de revelar se puede **registrar la estimación acordada**: cualquier valor de la escala, aunque nadie lo haya votado. Se guarda como Story Points de la historia con la operación de la spec 07 (y con sus reglas: si la historia se completó mientras tanto, se rechaza) y la votación se **cierra**.
- RN-10: Una votación cerrada no admite votos, revelaciones, rondas nuevas ni otra estimación acordada. Para volver a estimar la historia se abre una votación nueva.
- RN-11: La votación, la historia y el integrante tienen que pertenecer al proyecto indicado; si no, se tratan como inexistentes.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos.
- Las reglas RN-1 a RN-10 (incluido el cálculo de diferencias de RN-7) viven en el dominio y se desarrollan con TDD.
- La escala y el registro de la estimación se reutilizan de la spec 07: no se repiten.
- El valor de un voto no se devuelve en ninguna respuesta mientras su ronda no está revelada (RN-5): la regla vale para la API, no solo para la pantalla.

## Casos límite

- Proyecto con un solo Developer: puede abrir, votar y revelar (no hay diferencias posibles).
- Todos votan lo mismo: no hay diferencias; el más bajo y el más alto son el mismo valor.
- Dos Developers empatan en el voto más bajo: se informan los dos.
- Un Developer cambia su voto antes de revelar: cuenta solo el último.
- Votar 1 o 21 (los extremos de la escala): se acepta. Votar 4, 0, 22 o 2.5: se rechaza.
- Revelar cuando falta un solo voto: se rechaza.
- Acordar un valor que nadie votó (votaron 3 y 8, se acuerda 5): se acepta.
- Después de cerrar una votación, abrir otra para la misma historia: se acepta (por ejemplo, para re-estimar).
- La historia se completa en un Sprint mientras la votación está abierta: al registrar la acordada se rechaza ("La historia ya está completada").
- Un integrante de otro proyecto intenta votar: se informa que el integrante no existe.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Abrir para una historia completada, o registrar la acordada si se completó | Conflicto: "La historia ya está completada" |
| Abrir para una historia que ya tiene una votación abierta | Conflicto: "La historia ya tiene una votación abierta" |
| Abrir en un proyecto sin Developers | Conflicto: "El proyecto no tiene Developers para votar" |
| Vota un integrante que no es Developer | Conflicto: "Solo los Developers pueden votar" |
| Falta el valor del voto o de la acordada | Dato inválido: "Los Story Points son obligatorios" (spec 07) |
| Valor fuera de la escala | Dato inválido: "Los Story Points deben ser un valor de la escala: 1, 2, 3, 5, 8, 13 o 21" (spec 07) |
| Votar o revelar una ronda ya revelada | Conflicto: "La ronda ya fue revelada" |
| Revelar cuando falta algún voto | Conflicto: "Faltan votos: todavía no votaron todos los Developers" |
| Nueva ronda o acordada sin revelar la ronda actual | Conflicto: "Primero hay que revelar la ronda actual" |
| Cualquier operación sobre una votación cerrada (salvo verla) | Conflicto: "La votación está cerrada" |
| El proyecto no existe | No encontrado: "El proyecto no existe" |
| La historia no existe (o es de otro proyecto) | No encontrado: "La historia no existe" |
| El integrante no existe (o es de otro proyecto) | No encontrado: "El integrante no existe" |
| La votación no existe (o es de otro proyecto) | No encontrado: "La votación no existe" |

Si hay más de un problema, primero se informa un pedido mal armado o un valor fuera de la escala (400), después lo que no existe (404) y por último el conflicto (409): el mismo orden que el resto de la API (specs 01, 04 y 07).

## Criterios de aceptación

- CA-1: Se puede abrir una votación para una historia del backlog que no está completada; arranca en la ronda 1 sin votos. Una historia no puede tener dos votaciones abiertas, y el proyecto tiene que tener al menos un Developer.
- CA-2: Cada Developer del proyecto registra su voto con un valor de la escala en la ronda actual, y lo puede cambiar antes de revelar; un integrante que no es Developer no puede votar.
- CA-3: Mientras la ronda no se revela, los votos permanecen ocultos: se ve quién votó y quién falta, pero no los valores.
- CA-4: La ronda se revela solo cuando votaron todos los Developers; al revelarla se muestran todos los votos y no admite más votos.
- CA-5: Al revelar se detectan las diferencias entre los votos, informando el voto más bajo y el más alto con quiénes los dieron.
- CA-6: Después de revelar se puede empezar una nueva ronda sin votos; las rondas anteriores se conservan con sus votos.
- CA-7: Después de revelar se registra la estimación acordada (cualquier valor de la escala): queda como Story Points de la historia y la votación se cierra y no admite más cambios.
- CA-8: Operar con un proyecto, una historia, un integrante o una votación que no existe (o es de otro proyecto) informa cuál no existe, sin cambiar nada.

Relación con los criterios iniciales del Issue #12: "registrar las estimaciones individuales" → CA-2; "ocultas hasta finalizar la votación" → CA-3; "al finalizar se muestran todas" → CA-4; "detectar diferencias" → CA-5; "nuevas rondas" → CA-6; "registrar la estimación acordada" → CA-7. CA-1 y CA-8 agregan cómo se abre una votación y los errores.

## Escenarios BDD

```gherkin
# language: es
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
```

Los pasos que no indican proyecto se refieren al proyecto de los Antecedentes, "Tienda Online". En las tablas, `(vacío)` representa un valor que falta (Gherkin no puede escribirlo), igual que en las specs 01, 04 y 07. Los pasos de proyectos, integrantes, backlog y Story Points (`que existe el proyecto…`, `que el proyecto tiene al integrante…`, `que el backlog tiene la historia…`, `la historia "…" tiene N Story Points`, `la historia "…" está sin estimar`) ya existen en las specs 01, 02, 04 y 07 y se reutilizan.

## Trazabilidad

Historia [#12](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/12) → esta spec → criterios → escenarios de `backend/features/planning-poker.feature` → tests → código Go.

| Criterio | Escenarios BDD | Tests Go | Código Go |
|---|---|---|---|
| CA-1 Abrir | Abrir una votación · No se pueden abrir dos · Historia completada · Proyecto sin Developers | `TestNuevaVotacion`, `TestPlanningPokerAbrirErrores`, `TestAbrirVotacion` | `dominio/planning_poker.go` (`NuevaVotacion`), `dominio/servicio_planning_poker.go` (`Abrir`), `api/planning_poker.go` (`abrir`), `postgres/esquema/06_votaciones.sql` |
| CA-2 Votar | Un Developer vota · Cambia su voto · No es Developer · Solo valores de la escala | `TestVotar`, `TestPlanningPokerVotarErrores`, `TestPlanningPokerErrores` | `dominio/planning_poker.go` (`Votar`), `api/planning_poker.go` (`votar`) |
| CA-3 Votos ocultos | Los votos están ocultos hasta revelar | `TestVistaDeVotacion`, `TestVotosOcultosEnLaAPI` | `dominio/planning_poker.go` (`Vista`), `api/planning_poker.go` (`votacionAJSON`) |
| CA-4 Revelar | Revelar cuando votaron todos · Falta un voto · Ronda revelada no admite votos | `TestRevelar`, `TestRevelarVotacion` | `dominio/planning_poker.go` (`Revelar`) |
| CA-5 Diferencias | Detectar diferencias · Sin diferencias · Empate en el más bajo | `TestResultadoDeRonda`, `TestRevelarVotacion` | `dominio/planning_poker.go` (`revelar`, `quienesVotaron`) |
| CA-6 Nueva ronda | Empezar una nueva ronda · Sin revelar la actual | `TestNuevaRonda`, `TestNuevaRondaYAcordar`, `TestRepositorioVotaciones` | `dominio/planning_poker.go` (`NuevaRonda`), `postgres/votaciones.go` |
| CA-7 Acordar | Registrar la acordada · Un valor que nadie votó · Sin revelar · Fuera de la escala · Votación cerrada · Abrir otra después de cerrar | `TestAcordar`, `TestPlanningPokerRecorridoCompleto`, `TestPlanningPokerAcordarHistoriaCompletada`, `TestNuevaRondaYAcordar` | `dominio/planning_poker.go` (`Acordar`), `dominio/servicio_planning_poker.go` (`Acordar` → `ServicioBacklog.Estimar`, spec 07) |
| CA-8 Inexistentes | Historia que no existe · Votación que no existe · Integrante de otro proyecto | `TestPlanningPokerAbrirErrores`, `TestPlanningPokerVotarErrores`, `TestPlanningPokerErrores`, `TestRepositorioVotaciones` | `dominio/servicio_planning_poker.go` (`votacionConEquipo`), `api/planning_poker.go` (`leerIDsVotacion`) |

## Dependencias

- Spec 02 (integrantes y sus roles): implementada.
- Spec 04 (backlog): implementada.
- Spec 07 (escala de Story Points y estimar una historia): misma dueña, Sprint 1; **su PR de código tiene que estar mergeado antes de implementar esta spec**.
- La usan: spec 10 (métricas de avance: suma los Story Points acordados) y spec 13 (Dashboard).
- Chequeo de independencia (INVEST): todos los criterios usan historias implementadas (02, 04) o de la misma dueña que se hacen antes (07); ninguno depende de una historia pendiente de otra integrante.

## Tareas

- [x] `[P1]` Votación, rondas, votos y cálculo de diferencias (RN-1 a RN-10) en `backend/dominio/` (TDD).
- [x] `[P1]` Casos de uso abrir, votar, ver, revelar, nueva ronda y registrar la acordada, usando la operación estimar de la spec 07.
- [x] `[P1]` Tablas de votaciones, rondas y votos en `backend/postgres/esquema/06_votaciones.sql` y su repositorio en `backend/postgres/votaciones.go`.
- [x] `[P1]` Rutas de la API de Planning Poker en `backend/api/`.
- [x] `[P1]` Escenarios de `backend/features/planning-poker.feature` con sus pasos en `features/planning_poker_test.go`.
- [x] `[P2]` Actualizar `backend/CLAUDE.md` (endpoints disponibles, specs relacionadas, trampas).

## Definition of Done

- [x] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [x] Escenarios BDD automatizados y en verde.
- [ ] CI en verde (gofmt / vet / tests / lint / build).
- [x] Criterios de aceptación verificados.
- [x] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [x] Spec en estado `implemented`.
