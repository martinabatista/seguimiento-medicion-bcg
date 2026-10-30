# Documentación técnica — Software Metrics & Estimation

Cómo está construida la aplicación, cómo se corre y cómo se prueba. Trabajo Práctico Integrador de
Ingeniería y Calidad de Software 2026 — UTN FRSR. Equipo: Martina Batista, Florencia Cabeza y Luana Guajardo.

Para usar la aplicación, ver el [manual de usuario](manual-usuario.md). Los números de pruebas y cobertura
están en el [informe de métricas](informe-metricas.md).

## 1. Qué es

Una aplicación web para administrar proyectos de software con Scrum y medir su estimación,
planificación, seguimiento y calidad. El recorrido completo:

```
crear proyecto e integrantes → cargar el Product Backlog → estimar con Planning Poker (Story Points)
→ planificar el Sprint → registrar horas y defectos → cerrar el Sprint
→ consultar métricas, Dashboard y reporte en PDF
```

| Parte | Tecnología |
|---|---|
| Backend: reglas de negocio, cálculos y API REST | Go (versión en `backend/go.mod`), solo la librería estándar `net/http` |
| Base de datos | PostgreSQL 17 en Docker; acceso con `database/sql` + driver `pgx` v5, SQL escrito a mano |
| Reporte en PDF | `go-pdf/fpdf` (Go puro) |
| Frontend | React 19 + Vite + TypeScript, sin librerías de rutas, estilos ni gráficos |
| Pruebas | `go test` (TDD), godog (BDD con Gherkin en español), Vitest (frontend) |
| Integración continua | GitHub Actions |

**Regla de la cátedra:** todas las reglas de negocio, validaciones y cálculos están en Go. El frontend
solo muestra datos y captura lo que se carga; si algo está mal, muestra el mensaje que manda el backend.

## 2. Arquitectura

```
 Navegador
    │
    ▼
 frontend/  (React)  ── pide /api/... ──►  backend/api/      rutas REST: lee el pedido y responde JSON
                         (proxy de Vite)        │
                                                ▼
                                          backend/dominio/  entidades, reglas y métricas (Go puro)
                                                │  usa interfaces de repositorio
                                     ┌──────────┴──────────┐
                                     ▼                     ▼
                              backend/postgres/      backend/memoria/
                              (la app real)          (solo para los tests)
                                     │
                                     ▼
                                PostgreSQL
```

**Regla de capas:** el dominio no conoce HTTP ni la base de datos. Define interfaces (por ejemplo, "un
repositorio de historias") y las capas de afuera las implementan. Así:

- las reglas y los cálculos se prueban rápido, con repositorios en memoria y sin base real;
- cambiar la base o la forma de exponer la API no toca las reglas.

### Un pedido de punta a punta

Ejemplo: registrar un defecto.

1. La pantalla **Defectos** llama a `pedir('POST', '/proyectos/1/defectos', {...})` (`frontend/src/api.ts`).
2. Vite reenvía `/api/...` al backend en el puerto 8080.
3. `backend/api/defectos.go` lee el JSON y llama al servicio del dominio.
4. `backend/dominio/servicio_defectos.go` valida (descripción obligatoria, severidad válida, la historia
   y el Sprint existen) y guarda con el repositorio.
5. `backend/postgres/defectos.go` hace el `INSERT` en la tabla `defectos`.
6. La API responde `201` con el defecto en JSON, o un error `{"error": "mensaje"}`.

### Errores

Todas las rutas responden los errores igual: `{"error": "mensaje en español"}` con este código:

| Código | Cuándo | Ejemplo |
|---|---|---|
| 400 | Un dato es inválido | "La descripción es obligatoria" |
| 404 | Lo que se pide no existe | un proyecto o una historia que no está |
| 409 | La acción choca con el estado actual | estimar una historia ya completada, cerrar un Sprint cerrado |

## 3. Backend

```
backend/
├── cmd/api/      # main: conecta la base, arma los servicios y levanta la API (puerto 8080)
├── api/          # rutas REST + tests de cada ruta
├── dominio/      # entidades, reglas, cálculos, servicios e interfaces de repositorio + tests unitarios
├── memoria/      # repositorios en memoria (los usan los tests)
├── postgres/     # repositorios sobre PostgreSQL + tests de integración
│   └── tablas/   # un archivo NN_tabla.sql por tabla
├── pdf/          # arma los PDF del reporte (solo presenta, no calcula)
└── features/     # escenarios BDD (*.feature) y sus pasos en Go
```

### Funcionalidades y rutas

Cada funcionalidad es una historia de usuario con su spec SDD y sus escenarios BDD. Las rutas cuelgan
de `/api/proyectos/{id}/...`. El contrato completo (rutas, JSON y códigos) está en
[`backend/CLAUDE.md`](../backend/CLAUDE.md).

| Spec | Funcionalidad | Rutas principales | Escenarios BDD |
|---|---|---|---|
| [01](specs/sprint-01/01-proyectos.md) | Crear y modificar proyectos | `/api/proyectos` | `proyectos.feature` |
| [02](specs/sprint-01/02-integrantes.md) | Integrantes con su rol Scrum | `…/integrantes` | `integrantes.feature` |
| [03](specs/sprint-01/03-estado-proyecto.md) | Estado del proyecto y avance del backlog | `…/estado` | `estado-proyecto.feature` |
| [04](specs/sprint-01/04-product-backlog.md) | Product Backlog | `…/backlog` | `product-backlog.feature` |
| [05](specs/sprint-01/05-crear-sprints.md) | Crear Sprints, Sprint Goal y asignar historias | `…/sprints` | `sprints.feature` |
| [06](specs/sprint-01/06-cerrar-sprint.md) | Completar historias y cerrar el Sprint | `…/completar`, `…/cerrar` | `cerrar-sprint.feature` |
| [07](specs/sprint-01/07-story-points.md) | Estimar con Story Points | `…/story-points` | `story-points.feature` |
| [08](specs/sprint-02/08-planning-poker.md) | Planning Poker | `…/votaciones` | `planning-poker.feature` |
| [09](specs/sprint-02/09-registro-esfuerzo.md) | Registrar horas trabajadas | `…/esfuerzo` | `registro-esfuerzo.feature` |
| [10](specs/sprint-02/10-metricas-avance.md) | Métricas de avance y velocidad | `…/metricas` | `metricas-avance.feature` |
| [11](specs/sprint-03/11-defectos.md) | Defectos | `…/defectos` | `defectos.feature` |
| [12](specs/sprint-03/12-metricas-esfuerzo-calidad.md) | Horas estimadas y métricas de esfuerzo | `…/horas-estimadas`, `…/metricas/esfuerzo` | `metricas-esfuerzo.feature` |
| [13](specs/sprint-03/13-dashboard.md) | Dashboard | `…/dashboard` | `dashboard.feature` |
| [14](specs/sprint-04/14-reporte-pdf.md) | Reporte exportable a PDF | `…/reporte`, `…/reporte/pdf` | `reporte.feature` |

### Base de datos

Las tablas se crean con archivos SQL simples, sin herramienta de migraciones. Docker los ejecuta en
orden la primera vez que crea la base:

| Archivo | Qué guarda |
|---|---|
| `01_proyectos.sql` | proyectos |
| `02_historias.sql` | historias del backlog |
| `03_integrantes.sql` | integrantes y su rol |
| `04_sprints.sql` | Sprints y qué historias tiene cada uno |
| `05_sprint_historias_completada.sql` | en qué Sprint se completó cada historia |
| `06_votaciones.sql` | votaciones, rondas y votos del Planning Poker |
| `07_sprint_historias_story_points.sql` | los Story Points de cada historia al cerrar el Sprint |
| `08_esfuerzo.sql` | horas trabajadas |
| `09_historias_horas_estimadas.sql` | horas estimadas de cada historia |
| `10_defectos.sql` | defectos |

Los ids los genera PostgreSQL (`GENERATED ALWAYS AS IDENTITY`). La conexión se lee de la variable
`DATABASE_URL`: no hay credenciales en el código.

## 4. Frontend

```
frontend/src/
├── main.tsx · App.tsx   # arranque y layout: lista de proyectos o menú de secciones
├── api.ts               # pedir(): el único lugar que habla con la API; tipos de los datos
├── componentes.tsx      # Cargando y MensajeError
├── formato.ts           # formato de números y fechas
├── estilos.css          # todos los estilos
└── paginas/             # una pantalla por sección del menú
```

- **Navegación sin librería:** la sección elegida es un estado de React; la dirección del navegador no cambia.
- **Un solo cliente de la API:** `pedir()` arma el pedido, devuelve los datos o lanza el mensaje de error del backend.
- **Gráficos en HTML y CSS** (barras con su alto en %, dona con `conic-gradient`), con colores revisados
  para daltonismo y el número escrito en cada barra.
- Detalle y "cómo agregar una pantalla" en [`frontend/CLAUDE.md`](../frontend/CLAUDE.md).

## 5. Cómo correr

Los pasos completos están en el [README](../README.md#levantar-el-proyecto). En resumen:

| Qué | Requisito | Comando | Queda en |
|---|---|---|---|
| Base de datos | Docker Desktop abierto | `docker compose up -d` | `localhost:5434` |
| Backend | Go y la variable `DATABASE_URL` | `cd backend && go run ./cmd/api` | `http://localhost:8080` |
| Frontend | Node.js 24 | `cd frontend && npm install && npm run dev` | `http://localhost:5173` |

Si cambian las tablas, hay que recrear la base: `docker compose down -v && docker compose up -d` (borra los datos de prueba).

## 6. Cómo testear

Los comandos están en el [README](../README.md#tests). Qué prueba cada tipo:

| Tipo | Dónde | Qué prueba |
|---|---|---|
| Unitarios (TDD) | `backend/dominio/*_test.go` | reglas, validaciones y cálculos, con repositorios en memoria |
| De la API | `backend/api/*_test.go` | cada ruta: códigos HTTP y JSON |
| De integración | `backend/postgres/*_test.go` | el SQL contra PostgreSQL real (corren si está definida `DATABASE_URL`) |
| Escenarios BDD | `backend/features/*.feature` | los criterios de aceptación de cada spec, en Gherkin, con godog |
| Del PDF | `backend/pdf/pdf_test.go` | que el reporte se arme con sus datos |
| Frontend | `frontend/src/api.test.ts` | el cliente de la API (Vitest) |

Para correr un solo escenario BDD: `cd backend && go test ./features/ -run "TestFeatures/<nombre del escenario>"`.

**Integración continua:** en cada Pull Request, GitHub Actions corre `gofmt`, `go vet` y `go test` con
cobertura (unitarios, de la API y BDD; la CI no tiene base, así que los de integración se corren en cada
compu), y en el frontend `lint`, `typecheck` y `build`. No se mergea sin la CI en verde.

## 7. Cómo se desarrolló

Cada funcionalidad siguió el mismo camino, y cada paso queda enlazado con el siguiente:

```
Historia (Issue) → Spec SDD → Criterios de aceptación (CA-n) → Escenarios BDD (CA-n) → Tests → Código Go
```

1. **SDD:** la spec (`docs/specs/sprint-NN/NN-slug.md`) dice qué y para qué, con reglas (RN-n) y
   criterios (CA-n). Otra integrante la aprueba antes de programar.
2. **BDD:** cada criterio se escribe como escenario Given–When–Then en un `.feature`, con el CA en el nombre.
3. **TDD:** RED → GREEN → REFACTOR, con un commit por fase (`test:` → `feat:` → `refactor:`), visible en el historial.
4. **Pull Request:** con la CI en verde y la aprobación de otra integrante; el PR cierra el Issue (`Closes #N`).

Cada spec termina con una tabla de trazabilidad que une cada criterio con su escenario y su test. Las
reglas del equipo están en [`CLAUDE.md`](../CLAUDE.md).

## 8. Decisiones técnicas

| Decisión | Por qué |
|---|---|
| SQL a mano con `database/sql`, sin ORM | El SQL se ve y se puede explicar; el dominio no se mezcla con la base |
| Tablas en archivos SQL, sin migraciones | Los datos son de prueba: recrear la base alcanza y no suma otra herramienta |
| Rutas con `net/http`, sin framework | La librería estándar de Go ya tiene patrones de rutas |
| Métricas calculadas al pedirlas, no guardadas | Siempre coinciden con los datos; el Dashboard y el reporte reutilizan los mismos cálculos |
| Proxy de Vite en desarrollo | El frontend y la API se ven como un solo sitio: no hace falta configurar CORS |
| Sin librería de gráficos | Tres gráficos simples no justifican una dependencia |

## 9. Límites

- Se usa en una compu de desarrollo (`npm run dev`); publicarla en un servidor no está en el alcance.
- No tiene usuarios ni inicio de sesión: cualquiera que abre la app ve todos los proyectos.
- En el Planning Poker los votos se cargan desde una misma pantalla, uno por Developer.
