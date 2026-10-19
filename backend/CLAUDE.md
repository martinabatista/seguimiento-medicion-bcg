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
- `postgres/`: acceso a PostgreSQL (`DATABASE_URL`, ver `.env.example`).

Endpoint disponible: `GET /api/health` → `{"status":"ok"}`.

## Límites

- No renderiza interfaz: el frontend consume la API.
- El paquete de dominio no conoce HTTP ni PostgreSQL; solo la capa de API y la de postgres.

## Contratos públicos

La **API REST** que consume `frontend/`. Cambiar una ruta, un payload o un código de error rompe
al frontend: actualizar la spec correspondiente y `frontend/` en el mismo PR.

## Dependencias

- **Internas:** ninguna.
- **Externas:** PostgreSQL (`docker compose up -d` en la raíz, puerto 5434); godog v0.16 (BDD).
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
