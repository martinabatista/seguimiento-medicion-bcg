# seguimiento-medicion-bcg — Software Metrics & Estimation

Aplicación web para administrar proyectos de software y medir su estimación, planificación,
seguimiento y calidad. Trabajo Práctico Integrador de Ingeniería y Calidad de Software 2026 — UTN FRSR.

Integrantes: Batista Martina, Cabeza Florencia, Guajardo Luana

- Reglas del equipo y arquitectura: [`CLAUDE.md`](CLAUDE.md)
- Especificaciones (SDD) por Sprint: [`docs/specs/`](docs/specs/)
- Documentación técnica: [`docs/documentacion-tecnica.md`](docs/documentacion-tecnica.md)
- Manual de usuario: [`docs/manual-usuario.md`](docs/manual-usuario.md)
- Tablero Scrum: [GitHub Projects](https://github.com/users/martinabatista/projects/4)

## Requisitos

- [Git](https://git-scm.com/)
- [Go](https://go.dev/dl/) 1.27 (la versión de `backend/go.mod`)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/), para la base de datos PostgreSQL
- [Node.js](https://nodejs.org/) 24, para el frontend

## Levantar el proyecto

```bash
git clone https://github.com/martinabatista/seguimiento-medicion-bcg.git
cd seguimiento-medicion-bcg
cp .env.example .env
```

Base de datos (con Docker Desktop abierto):

```bash
docker compose up -d        # PostgreSQL 17 en localhost:5434; crea las tablas de backend/postgres/tablas/
docker compose down         # para apagarla
docker compose down -v && docker compose up -d   # recrear la base cuando cambian las tablas
```

Backend: necesita la variable de entorno `DATABASE_URL` (dónde está la base). Se define **una sola
vez** por compu, en PowerShell, y después se cierra y se vuelve a abrir VS Code:

```powershell
setx DATABASE_URL "postgres://sme:sme_dev@localhost:5434/sme?sslmode=disable"
```

```bash
cd backend
go run ./cmd/api            # API en http://localhost:8080 — probar: /api/health, /api/proyectos
```

Frontend (con el backend andando, en otra terminal):

```bash
cd frontend
npm install                 # solo la primera vez, o cuando cambian las dependencias
npm run dev                 # la app en http://localhost:5173 (los pedidos a /api van al backend)
```

## Tests

```bash
cd backend
go test ./... -cover        # tests unitarios y escenarios BDD (backend/features/*.feature)
                            # con DATABASE_URL definida, también los de integración con PostgreSQL
go vet ./...
gofmt -l .                  # no debe listar nada

cd frontend
npm test                    # tests del frontend (Vitest)
npm run lint                # ESLint
npm run typecheck           # TypeScript
npm run build               # compila la app
```

La CI corre estos mismos controles en cada Pull Request (en el frontend: lint, typecheck y build).

## Cómo trabajamos

Cada historia pasa por dos ramas, siempre creadas desde `main` actualizado:

1. `docs/spec-NN-slug`: se completa la spec → Pull Request que otra integrante revisa y aprueba.
2. `feat/slug`: con la spec ya en `main`, se implementa con BDD + TDD → Pull Request con `Closes #N` (frase que se escribe en la descripción de un Pull Request para conectarlo con un Issue. N es el número del Issue).

Todo entra a `main` por Pull Request, con la CI en verde y la aprobación de otra integrante.
