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
└── features/                # escenarios BDD *.feature (godog) + features_test.go con los pasos
```

Se crean con la primera historia que las necesite:
- `dominio/`: entidades, reglas y cálculos — Go puro, sin HTTP ni DB.
- `postgres/`: acceso a PostgreSQL con `pgx` (`DATABASE_URL`, ver `.env.example`) y sus
  migraciones en `postgres/migraciones/` (ver "Decisiones técnicas").

Endpoint disponible: `GET /api/health` → `{"status":"ok"}`.

## Decisiones técnicas

Acordadas por el equipo el 2026-10-19. Cambiarlas requiere acuerdo de las tres y actualizar este archivo.

**Acceso a datos: `pgx` v5, SQL escrito a mano, sin ORM.**
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

**Migraciones: `goose`, un archivo SQL por cambio.**
- Van en `postgres/migraciones/` con número correlativo: `00001_crear_proyectos.sql`, con las
  secciones `-- +goose Up` y `-- +goose Down`.
- Nunca se edita una migración ya mergeada: se agrega una nueva. Si no, las bases de las tres se
  desincronizan.
- Cada historia crea las suyas con el número siguiente. Avisar en el chat antes de tomar un número
  para no chocar entre ramas.

Las tres dependencias (`pgx`, `goose`, y `godog` ya existente) se justifican en la spec de la historia
que las incorpora (la 01 para `pgx` y `goose`).

## Límites

- No renderiza interfaz: el frontend consume la API.
- El paquete de dominio no conoce HTTP ni PostgreSQL; solo la capa de API y la de postgres.

## Contratos públicos

La **API REST** que consume `frontend/`. Cambiar una ruta, un payload o un código de error rompe
al frontend: actualizar la spec correspondiente y `frontend/` en el mismo PR.

## Dependencias

- **Internas:** ninguna.
- **Externas:** PostgreSQL (`docker compose up -d` en la raíz, puerto 5434); godog v0.16 (BDD).
  Decididas, se agregan con la spec 01: `pgx` v5 (acceso a datos) y `goose` (migraciones).
  Cada dependencia nueva se justifica en la spec que la agrega.

## Specs relacionadas

- `docs/specs/sprint-00/00-setup-inicial.md` — estructura inicial y CI.

## Cómo testear

```bash
cd backend
gofmt -l .          # debe salir vacío
go vet ./...
go test ./... -cover   # incluye los escenarios BDD de features/
go run ./cmd/api       # levantar la API
```

Escenarios BDD: se escriben en `features/<funcionalidad>.feature` con `# language: es`
(Característica / Escenario / Dado / Cuando / Entonces) y sus pasos se registran en
`features/features_test.go`. Cada escenario indica en su nombre el criterio que cubre: `(CA-3)`.

Umbrales de calidad: tests unitarios obligatorios para métricas, cálculos de estimación, reglas de
negocio y validaciones (exigencia de la cátedra). TODO: fijar % mínimo de cobertura en la spec de setup.

## Trampas conocidas

- _(vacío)_

## Estado

Origen: bootstrap del flujo (spec 00). Última revisión: 2026-10-17.
