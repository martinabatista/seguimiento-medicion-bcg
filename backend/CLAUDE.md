# backend

> `CLAUDE.md` de la unidad backend. Autosuficiente para su carpeta; el raíz lo indexa.

## Qué es

Núcleo de la aplicación en Go: entidades del dominio (proyecto, integrante, historia, Sprint,
estimación, esfuerzo, defecto), reglas de negocio, cálculo de métricas y la API REST que las expone.

## Estructura

```
backend/
├── go.mod / go.sum          # módulo github.com/martinabatista/seguimiento-medicion-bcg/backend
├── cmd/api/                 # main: levanta la API (puerto en PORT, default 8080)
├── api/                     # rutas REST (net/http, patrones "GET /api/...") + tests
├── dominio/                 # entidades, reglas, casos de uso e interfaces de repositorio — Go puro
├── memoria/                 # repositorios en memoria para los tests (sin base real)
├── postgres/                # repositorios sobre PostgreSQL (database/sql + pgx) + tests de integración
│   └── esquema/             # tablas: NN_tabla.sql, las ejecuta Docker al crear la base
└── features/                # escenarios BDD *.feature (godog) + features_test.go con los pasos
```

Endpoints disponibles:

| Método y ruta | Qué hace | Respuestas |
|---|---|---|
| `GET /api/health` | Estado del servicio | 200 `{"status":"ok"}` |
| `POST /api/proyectos` | Crear proyecto (spec 01) | 201 proyecto · 400 |
| `GET /api/proyectos` | Listar en orden de creación | 200 lista (`[]` si no hay) |
| `GET /api/proyectos/{id}` | Ver un proyecto | 200 · 404 |
| `PUT /api/proyectos/{id}` | Modificar (reemplaza los 4 datos) | 200 · 400 · 404 |
| `POST /api/proyectos/{id}/backlog` | Crear historia en el backlog (spec 04) | 201 historia · 400 · 404 |
| `GET /api/proyectos/{id}/backlog` | Listar el backlog por prioridad, desempate por creación | 200 lista (`[]` si no hay) · 404 |
| `GET /api/proyectos/{id}/backlog/{historiaID}` | Ver una historia | 200 · 404 |
| `PUT /api/proyectos/{id}/backlog/{historiaID}` | Modificar (reemplaza datos y estado) | 200 · 400 · 404 |
| `DELETE /api/proyectos/{id}/backlog/{historiaID}` | Eliminar una historia | 204 · 404 |
| `PUT /api/proyectos/{id}/backlog/{historiaID}/story-points` | Estimar `{"story_points"}` (spec 07) | 200 historia · 400 · 404 · 409 completada |
| `DELETE /api/proyectos/{id}/backlog/{historiaID}/story-points` | Quitar la estimación (queda sin estimar) | 200 historia · 404 · 409 completada |
| `GET /api/story-points/escala` | Escala de Story Points `{"valores": [1, 2, 3, 5, 8, 13, 21]}` | 200 |
| `POST /api/proyectos/{id}/integrantes` | Agregar un integrante (spec 02) | 201 integrante · 400 · 404 · 409 email repetido |
| `GET /api/proyectos/{id}/integrantes` | Listar en orden de alta | 200 lista (`[]` si no hay) · 404 |
| `GET /api/proyectos/{id}/integrantes/{integranteID}` | Ver un integrante | 200 · 404 |
| `GET /api/proyectos/{id}/estado` | Estado del proyecto, resumen del backlog (spec 03) y Sprint en curso (spec 05) | 200 · 404 |
| `POST /api/proyectos/{id}/sprints` | Crear un Sprint (spec 05) | 201 Sprint · 400 · 404 · 409 superpuesto |
| `GET /api/proyectos/{id}/sprints` | Listar por fecha de inicio (sin historias) | 200 lista (`[]` si no hay) · 404 |
| `GET /api/proyectos/{id}/sprints/{sprintID}` | Ver un Sprint con sus historias | 200 · 404 |
| `PUT /api/proyectos/{id}/sprints/{sprintID}` | Modificar el Sprint Goal `{"sprint_goal"}` | 200 · 400 · 404 · 409 cerrado |
| `POST /api/proyectos/{id}/sprints/{sprintID}/historias` | Asignar `{"historia_id"}` | 200 Sprint · 404 · 409 |
| `DELETE /api/proyectos/{id}/sprints/{sprintID}/historias/{historiaID}` | Quitar una historia | 200 Sprint · 404 · 409 |
| `POST …/sprints/{sprintID}/historias/{historiaID}/completar` | Completar una historia en el Sprint (spec 06) | 200 Sprint · 404 · 409 |
| `DELETE …/sprints/{sprintID}/historias/{historiaID}/completar` | Deshacer (vuelve a `En progreso`) | 200 Sprint · 404 · 409 |
| `POST /api/proyectos/{id}/sprints/{sprintID}/cerrar` | Cerrar el Sprint (spec 06) | 200 Sprint · 404 · 409 ya cerrado |
| `GET /api/proyectos/{id}/sprints/anteriores` | Sprints cerrados con sus historias, del más reciente al más antiguo | 200 lista · 404 |
| `POST /api/proyectos/{id}/votaciones` | Abrir una votación de Planning Poker `{"historia_id"}` (spec 08) | 201 votación · 400 · 404 · 409 |
| `GET /api/proyectos/{id}/votaciones/{votacionID}` | Ver la votación (sin valores de la ronda sin revelar) | 200 · 404 |
| `POST …/votaciones/{votacionID}/votos` | Votar `{"integrante_id", "story_points"}` (solo Developers) | 200 votación · 400 · 404 · 409 |
| `POST …/votaciones/{votacionID}/revelar` | Revelar la ronda (votaron todos los Developers) | 200 votación · 404 · 409 |
| `POST …/votaciones/{votacionID}/rondas` | Empezar una nueva ronda | 200 votación · 404 · 409 |
| `POST …/votaciones/{votacionID}/acordar` | Registrar la acordada `{"story_points"}`: estima la historia y cierra | 200 votación · 400 · 404 · 409 |

`DELETE /api/proyectos/{id}/backlog/{historiaID}` responde 409 si la historia está en un Sprint abierto (spec 05 CA-6).
`PUT /api/proyectos/{id}/backlog/{historiaID}` responde 409 si pone o saca `Completada` a una historia de un
Sprint abierto: se completa desde el Sprint (spec 06 CA-5).

Proyecto en JSON: `{"id", "nombre", "descripcion", "fecha_inicio", "fecha_fin"}`, fechas `AAAA-MM-DD`.

Historia en JSON: `{"id", "proyecto_id", "titulo", "descripcion", "prioridad", "criterios_aceptacion",
"estado", "story_points"}`. `estado` es `Pendiente` | `En progreso` | `Completada` (al crear se ignora:
siempre arranca `Pendiente`); `story_points` es `null` mientras está sin estimar. Los Story Points solo
cambian por las rutas `…/story-points` (spec 07): el `PUT` de la historia los conserva.

Integrante en JSON: `{"id", "proyecto_id", "nombre", "email", "rol"}`. `rol` es `Product Owner` |
`Scrum Master` | `Developer`; `email` se devuelve normalizado (sin espacios al borde y en minúsculas).

Estado en JSON: `{"proyecto": {…proyecto…}, "estado", "backlog": {"total", "pendientes", "en_progreso",
"completadas", "porcentaje_completadas", "story_points_totales", "story_points_completados", "sin_estimar"}}`.
`estado` es `No iniciado` | `En curso` | `Finalizado`, calculado con la fecha de hoy (no se guarda).
`"sprint_en_curso"` es el Sprint abierto que incluye hoy (sin historias), o `null` (spec 05 CA-5).

Sprint en JSON: `{"id", "proyecto_id", "numero", "fecha_inicio", "fecha_fin", "sprint_goal", "estado",
"historias"}`. `estado` es `Abierto` | `Cerrado`; `historias` (lista de historias como en el backlog,
ordenadas por prioridad, cada una con `"completada_en_sprint"`) solo aparece al crear, ver, modificar,
asignar, quitar, completar, deshacer, cerrar y en los anteriores; no al listar.

Votación en JSON: `{"id", "proyecto_id", "historia_id", "estado", "ronda", "estimacion_acordada", "votaron",
"faltan", "rondas"}`. `estado` es `Abierta` | `Cerrada`; `votaron` y `faltan` son `[{"id", "nombre"}]` de la
ronda actual; `rondas` trae **solo las reveladas**, cada una `{"numero", "votos": [{"integrante_id", "nombre",
"story_points"}], "hay_diferencias", "mas_bajo", "mas_alto"}` (extremos `{"story_points", "integrantes"}`).
Los valores de una ronda sin revelar no salen en ninguna respuesta (spec 08 RN-5).

## Decisiones técnicas

Acordadas por el equipo el 2026-10-19 (acceso a datos y esquema revisados el 2026-10-20).
Cambiarlas requiere acuerdo de las tres y actualizar este archivo.

**Acceso a datos: `database/sql` (librería estándar de Go) con el driver `pgx` v5, SQL escrito a
mano, sin ORM.**
- Se programa contra `database/sql` (`db.QueryRow`, `db.Exec`, `rows.Scan`); `pgx` solo se registra
  como driver con `import _ "github.com/jackc/pgx/v5/stdlib"` y `sql.Open("pgx", DATABASE_URL)`.
- `postgres/` implementa interfaces que define el dominio; el dominio no conoce la base. Las reglas y
  los cálculos se testean con repositorios en memoria, sin base real.
- Solo `postgres/` necesita una base real (`docker compose up -d`) para sus tests de integración.
- La conexión usa `DATABASE_URL` (ver `.env.example`). Nada de credenciales en el código.
- Por qué: el SQL se ve y se puede defender en la presentación, y un ORM mezcla tags y tipos de la
  base en las entidades del dominio.

**Identificadores y rutas REST**
- El id es un número autoincremental generado por PostgreSQL (`BIGINT GENERATED ALWAYS AS IDENTITY`),
  no se genera en Go.
- Lo que pertenece a un proyecto cuelga de `/api/proyectos/{id}/...`, en español:

  ```
  POST /api/proyectos                    GET  /api/proyectos/{id}/integrantes
  GET  /api/proyectos                    GET  /api/proyectos/{id}/backlog
  GET  /api/proyectos/{id}               POST /api/proyectos/{id}/sprints
  PUT  /api/proyectos/{id}
  ```
- Rutas con `net/http` (patrones `"GET /api/..."` y `r.PathValue("id")`), sin librerías extra.
- Errores con el mismo formato en todas las rutas: `{"error": "mensaje"}` y el código HTTP que
  corresponda (400 dato inválido, 404 no existe, 409 conflicto).
- Campos JSON en español y `snake_case` (`fecha_inicio`, `fecha_fin`).

**Tablas: archivos SQL simples en `postgres/esquema/`, sin herramienta de migraciones.**
- Un archivo por tabla, numerado en el orden en que se crean: `01_proyectos.sql`,
  `02_historias.sql`, `03_integrantes.sql`… Cada uno con su `CREATE TABLE`.
- Docker los ejecuta en orden **solo cuando crea la base por primera vez** (los monta en
  `/docker-entrypoint-initdb.d`, ver `docker-compose.yml`).
- Si se agrega o cambia un archivo, cada integrante recrea su base (los datos son de prueba):
  `docker compose down -v && docker compose up -d`.
- Se puede editar un archivo existente mientras su historia esté abierta; avisar en el PR cuando haga
  falta recrear la base.
- Por qué: un TP no necesita conservar datos entre cambios, y así no se suma una herramienta más.

Las dependencias (`pgx`, y `godog` ya existente) se justifican en la spec de la historia que las
incorpora (la 01 para `pgx`).

## Límites

- No renderiza interfaz: el frontend consume la API.
- El paquete de dominio no conoce HTTP ni PostgreSQL; solo la capa de API y la de postgres.

## Contratos públicos

La **API REST** que consume `frontend/`. Cambiar una ruta, un payload o un código de error rompe
al frontend: actualizar la spec correspondiente y `frontend/` en el mismo PR.

## Dependencias

- **Internas:** ninguna.
- **Externas:** PostgreSQL (`docker compose up -d` en la raíz, puerto 5434); godog v0.16 (BDD).
  Decidida, se agrega con la spec 01: `pgx` v5 (driver de PostgreSQL para `database/sql`).
  Cada dependencia nueva se justifica en la spec que la agrega.

## Specs relacionadas

- `docs/specs/sprint-00/00-setup-inicial.md` — estructura inicial y CI.
- `docs/specs/sprint-01/01-proyectos.md` — proyectos: `dominio/proyecto.go`, `dominio/servicio_proyectos.go`,
  `postgres/proyectos.go`, `api/proyectos.go`, `features/proyectos.feature`. Si cambia el JSON de
  proyecto o sus errores, revisar las specs 02 a 05, que cuelgan de `/api/proyectos/{id}`.
- `docs/specs/sprint-01/02-integrantes.md` — integrantes: `dominio/integrante.go`,
  `dominio/servicio_integrantes.go`, `memoria/integrantes.go`, `postgres/integrantes.go`,
  `postgres/esquema/03_integrantes.sql`, `api/integrantes.go`, `features/integrantes.feature` (pasos en
  `features/integrantes_test.go`). Las specs 08 (Planning Poker) y 09 (esfuerzo) usan el id del
  integrante: si cambia su JSON o su identidad (el email), actualizar la spec 02.
- `docs/specs/sprint-01/03-estado-proyecto.md` — estado del proyecto: `dominio/estado_proyecto.go`
  (`CalcularEstado`, `ResumirBacklog`), `dominio/servicio_estado.go`, `api/estado.go`,
  `features/estado-proyecto.feature` (pasos en `features/estado_proyecto_test.go`). El Sprint en curso
  lo agregó la spec 05 (CA-5); los escenarios con Story Points esperan a la spec 07. Si cambian los
  estados o los Story Points de la historia (spec 04), revisar `ResumirBacklog`. La spec 13 (Dashboard)
  va a usar este estado.
- `docs/specs/sprint-01/04-product-backlog.md` — backlog: `dominio/historia.go`, `dominio/servicio_backlog.go`,
  `memoria/historias.go`, `postgres/historias.go`, `postgres/esquema/02_historias.sql`, `api/backlog.go`,
  `features/product-backlog.feature` (pasos en `features/backlog_test.go`). Las specs 05 a 08 usan la
  historia: si cambian su estado o sus Story Points, actualizar la spec 04.
- `docs/specs/sprint-01/05-crear-sprints.md` — Sprints: `dominio/sprint.go` (`NuevoSprint`, `SprintEnCurso`),
  `dominio/servicio_sprints.go`, `memoria/sprints.go`, `postgres/sprints.go`, `postgres/esquema/04_sprints.sql`,
  `api/sprints.go`, `features/sprints.feature` (pasos en `features/sprints_test.go`). También toca el
  estado (`ServicioEstado`, CA-5) y el backlog (`ServicioBacklog.Eliminar`, CA-6). La spec 06 cierra los
  Sprints (estado `Cerrado`) y la 10 calcula sus métricas.
- `docs/specs/sprint-01/06-cerrar-sprint.md` — cierre de Sprints: `ServicioSprints` (`CompletarHistoria`,
  `DeshacerHistoria`, `Cerrar`, `Anteriores`), `Sprint.Completadas`, regla del backlog en
  `ServicioBacklog.Modificar` (`cambiaCompletada`), `postgres/esquema/05_sprint_historias_completada.sql`,
  `features/cerrar-sprint.feature` (pasos en `features/cerrar_sprint_test.go`). La spec 10 (velocidad) usa
  `Sprint.Completadas` de los Sprints cerrados, no el estado actual de las historias (RN-8).
- `docs/specs/sprint-01/07-story-points.md` — Story Points: `dominio/story_points.go` (`EscalaStoryPoints`,
  `EnEscala`, `ParsearStoryPoints`), `dominio/servicio_story_points.go` (`ServicioBacklog.Estimar`,
  `QuitarEstimacion`), `api/story_points.go`, `features/story-points.feature` (pasos en
  `features/story_points_test.go`). Sin tabla nueva: usa `historias.story_points`. La spec 08 (Planning Poker)
  registra la estimación acordada con `Estimar` y usa la misma escala; la 10 suma estos Story Points.
- `docs/specs/sprint-02/08-planning-poker.md` — Planning Poker: `dominio/planning_poker.go` (`Votacion`,
  `NuevaVotacion`, `Votar`, `Revelar`, `NuevaRonda`, `Acordar`, `Vista` con los votos ocultos y las
  diferencias), `dominio/servicio_planning_poker.go`, `memoria/votaciones.go`, `postgres/votaciones.go`,
  `postgres/esquema/06_votaciones.sql`, `api/planning_poker.go`, `features/planning-poker.feature` (pasos en
  `features/planning_poker_test.go`). Usa los roles de la spec 02 (solo `Developer` vota) y `Estimar` de la 07.

## Cómo testear

```bash
cd backend
gofmt -l .          # debe salir vacío
go vet ./...
go test ./... -cover   # incluye los escenarios BDD de features/
go run ./cmd/api       # levantar la API (necesita DATABASE_URL, ver el README)
```

Escenarios BDD: se escriben en `features/<funcionalidad>.feature` con `# language: es`
(Característica / Escenario / Dado / Cuando / Entonces) y sus pasos se registran en
`features/features_test.go`. Cada escenario indica en su nombre el criterio que cubre: `(CA-3)`.

Umbrales de calidad: tests unitarios obligatorios para métricas, cálculos de estimación, reglas de
negocio y validaciones (exigencia de la cátedra). TODO: fijar % mínimo de cobertura en la spec de setup.

## Trampas conocidas

- Eliminar una historia del backlog borra sus votaciones de Planning Poker (`ON DELETE CASCADE` en
  `06_votaciones.sql`). `votos` referencia a `integrantes`, así que el `TRUNCATE … CASCADE` también lo vacía.
- `VistaVotacion` (dominio) es lo único que se muestra de una votación: nunca devolver `Votacion` directo
  por la API, porque tiene los valores de la ronda sin revelar (spec 08 RN-5).
- `historias` e `integrantes` referencian a `proyectos`: los tests de integración vacían con
  `TRUNCATE proyectos RESTART IDENTITY CASCADE` (sin `CASCADE`, PostgreSQL no deja truncar `proyectos`).
- Cada tabla nueva en `postgres/esquema/` obliga a recrear la base local
  (`docker compose down -v && docker compose up -d`), o a correr solo ese archivo:
  `docker compose exec -T db psql -U sme -d sme < backend/postgres/esquema/NN_tabla.sql`.
- El orden del backlog (spec 04 RN-8) lo aplica el dominio (`OrdenarBacklog`), no el `ORDER BY` del SQL:
  el repositorio puede devolver las historias en cualquier orden.
- Email repetido de un integrante (spec 02 RN-6): lo controla el servicio y **también** la restricción
  `integrantes_email_unico` de la tabla, para dos altas simultáneas. `postgres/integrantes.go` traduce esa
  violación (código `23505`) a `ErrIntegranteYaRegistrado`, que la API responde como 409. Si se renombra
  la restricción en el SQL, hay que cambiarla también en el repositorio. La comparación funciona porque
  el email se guarda ya normalizado (minúsculas, sin espacios al borde).
- La verificación de que el proyecto existe (`verificarProyecto` en `dominio/servicio_proyectos.go`) la
  comparten todos los servicios de lo que cuelga de un proyecto: usarla en los nuevos (Sprints, etc.).
- La fecha de hoy (spec 03) **no** se lee con `time.Now()` dentro del dominio: `ServicioEstado` recibe un
  reloj (`func() time.Time`). `cmd/api` le pasa `time.Now`, los tests una fecha fija y los escenarios BDD
  la del paso `que hoy es "AAAA-MM-DD"` (sin ese paso, la del sistema). Hacer lo mismo en cualquier
  cálculo que dependa del día (por ejemplo, el Sprint en curso).
- `go test -race` necesita cgo: en Windows sin gcc no corre localmente; la CI (Linux) sí lo corre.
- Los tests de integración (`postgres/*_test.go`) usan la misma base que `go run ./cmd/api` y la dejan
  con sus datos de prueba (por ejemplo, un Sprint cerrado). Antes de probar la API a mano, vaciarla:
  `docker compose exec -T db psql -U sme -d sme -c "TRUNCATE proyectos RESTART IDENTITY CASCADE"`.
- Conflictos (409): el dominio los devuelve como `ErrConflicto{Mensaje}` (como `ErrDatoInvalido` para
  los 400) y `responderError` los traduce solo. Usarlo para las reglas nuevas que choquen con lo ya
  registrado, en vez de agregar casos sueltos.
- `sprint_historias` borra en cascada: si se elimina una historia del backlog (solo se puede fuera de un
  Sprint abierto, spec 05 RN-11), sale también de los Sprints cerrados donde estuvo.
- Lo completado **en un Sprint** se guarda en `sprint_historias.completada` (spec 06), aparte del estado
  de la historia en el backlog: un Sprint cerrado conserva lo que se completó en él aunque después la
  historia cambie. Para métricas por Sprint (spec 10), leer eso y no el estado actual.
- `docker compose up -d --wait` puede decir "Healthy" mientras PostgreSQL todavía se reinicia después de
  ejecutar los `.sql` del esquema: si un `psql` justo después falla con "shutting down", esperar unos segundos.

## Estado

Origen: bootstrap del flujo (spec 00). Última revisión: 2026-10-23 (spec 08).
