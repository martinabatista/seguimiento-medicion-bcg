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

Proyecto en JSON: `{"id", "nombre", "descripcion", "fecha_inicio", "fecha_fin"}`, fechas `AAAA-MM-DD`.

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
  `02_integrantes.sql`… Cada uno con su `CREATE TABLE`.
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

- _(vacío)_

## Estado

Origen: bootstrap del flujo (spec 00). Última revisión: 2026-10-20 (spec 01).
