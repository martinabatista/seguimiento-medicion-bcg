test → feat → test → feat → feat (pantalla)
| Indicador | Valor |
|---|---|
| Cobertura del backend Go (todos los paquetes juntos) | **88,7 %** de las sentencias |
| Cobertura del dominio (reglas de negocio y cálculos) | **93,6 %** |
| Tests unitarios y de integración en Go | **227** funciones de test |
| Escenarios BDD automatizados (godog) | **359** escenarios y **2309** pasos, todos en verde |
| Criterios de aceptación cubiertos por escenarios BDD | **94 de 98** (los 4 restantes se verifican con tests unitarios o a mano, ver sección 4) |
| Commits de TDD | **52** `test:` (rojo), **54** `feat:` (verde) y **8** `refactor:` |
| Historias de usuario implementadas | **14 de 14** |
| Pull Requests mergeados | **52**, todos aprobados por otra integrante |
| Ejecuciones de la CI | **111**, todas en verde |

## 2. Cobertura de pruebas

Medida con `go test -cover`. Cada paquete muestra la cobertura que le dan sus propios tests:

| Paquete | Qué tiene | Cobertura |
|---|---|---|
| `dominio` | Entidades, reglas de negocio, validaciones y cálculos de métricas | **93,6 %** |
| `api` | Rutas REST: leen el pedido, llaman al dominio y responden JSON | **93,8 %** |
| `pdf` | Armado del reporte en PDF (spec 14) | **96,9 %** |
| `postgres` | Acceso a datos con SQL (tests de integración contra PostgreSQL real) | **82,7 %** |
| `memoria` | Repositorios en memoria que usan los tests | se prueba a través de los otros paquetes |
| `cmd/api` | `main`: conecta la base y arma los servicios | sin tests propios |

**Mínimo del equipo (spec 00): 90 % en el dominio y 80 % en el total del backend. Se cumplen los dos.**

**Cobertura total del backend: 88,7 %** (3139 de 3539 sentencias), medida con `-coverpkg=./...`. Así cuenta también el código que un paquete ejecuta desde los tests de otro, por ejemplo los repositorios en memoria que usan los tests del dominio. Sin esa opción, `go test` informa un total de 82,6 %, porque `memoria` y `cmd/api` no tienen tests propios y figuran con 0 %.

**Lo que queda sin cubrir** son sobre todo errores que no se pueden provocar sin romper la base de datos: una conexión que se cae o una fila que no se puede leer. La lógica de negocio, que es lo que la cátedra exige probar (métricas, estimación, reglas y validaciones), está en el dominio y supera el 93 %.

**Frontend:** 8 tests con Vitest sobre el cliente de la API, que es la única pieza con lógica propia: armar el pedido y mostrar el mensaje de error del backend. Las pantallas no tienen reglas de negocio, solo muestran lo que calcula Go. La CI además corre `lint`, `typecheck` y `build` en cada PR.

## 3. Pruebas automatizadas

| Tipo | Dónde | Cantidad |
|---|---|---|
| Unitarios del dominio | `backend/dominio/*_test.go` | 133 tests en 28 archivos |
| De la API | `backend/api/*_test.go` | 72 tests en 15 archivos |
| De integración con PostgreSQL | `backend/postgres/*_test.go` | 17 tests en 7 archivos |
| Del PDF | `backend/pdf/*_test.go` | 5 tests |
| Escenarios BDD (Gherkin en español) | `backend/features/*.feature` | 359 escenarios en 15 archivos |
| Frontend | `frontend/src/api.test.ts` | 8 tests |

**Escenarios BDD por funcionalidad:**

| Archivo | Spec | Escenarios |
|---|---|---|
| `proyectos.feature` | 01 | 15 |
| `integrantes.feature` | 02 | 21 |
| `estado-proyecto.feature` | 03 | 10 |
| `product-backlog.feature` | 04 | 18 |
| `sprints.feature` | 05 | 18 |
| `cerrar-sprint.feature` | 06 | 15 |
| `story-points.feature` | 07 | 14 |
| `planning-poker.feature` | 08 | 29 |
| `registro-esfuerzo.feature` | 09 | 21 |
| `metricas-avance.feature` | 10 | 9 |
| `defectos.feature` | 11 | 23 |
| `metricas-esfuerzo.feature` | 12 | 16 |
| `dashboard.feature` | 13 | 14 |
| `reporte.feature` | 14 | 22 |
| `estado-servicio.feature` | 00 | 1 |

Los escenarios de tipo "Esquema del escenario" se ejecutan una vez por fila de ejemplos: por eso los 359 escenarios ejecutados son más que los escenarios escritos.

**Tamaño del código:** el backend tiene unas 7000 líneas de código Go de producción y unas 12 000 líneas de tests en Go, más 1872 líneas de escenarios Gherkin. Hay más código de pruebas que de producción: 1,7 líneas de test por cada línea de código.

## 4. Trazabilidad: criterios de aceptación cubiertos

Cada escenario BDD indica en su nombre el criterio que cubre, por ejemplo `(CA-4)`. Sobre las 14 specs de historias de usuario:

| Spec | Criterios | Con escenario BDD |
|---|---|---|
| 01 Proyectos | 7 | 7 |
| 02 Integrantes | 6 | 6 |
| 03 Estado del proyecto | 4 | 4 |
| 04 Product Backlog | 7 | 7 |
| 05 Crear Sprints | 7 | 7 |
| 06 Cerrar Sprint | 6 | 6 |
| 07 Story Points | 6 | 6 |
| 08 Planning Poker | 9 | 9 |
| 09 Registrar esfuerzo | 8 | 8 |
| 10 Métricas de avance | 7 | 6 |
| 11 Defectos | 7 | 7 |
| 12 Métricas de esfuerzo | 7 | 6 |
| 13 Dashboard | 9 | 8 |
| 14 Reporte PDF | 8 | 7 |
| **Total** | **98** | **94** |

Los 4 criterios sin escenario BDD no son comportamientos de la API, y cada uno se verifica de otra forma:

- **Spec 10 CA-5 y spec 12 CA-6** ("los cálculos se desarrollan con TDD"): se verifican con los tests unitarios del dominio y con el orden de los commits (sección 5).
- **Spec 13 CA-9 y spec 14 CA-8** (las pantallas del Dashboard y del Reporte): se probaron a mano en la interfaz web, como dice cada PR.

La spec 00 (setup inicial) tiene criterios de infraestructura, como la CI en verde, la rama `main` protegida o el README. Se verifican con la configuración del repositorio, no con escenarios.

Cada spec tiene su tabla de trazabilidad (criterio → escenarios → tests → código) y el recorrido completo de una funcionalidad está en la guía del proyecto, con el ejemplo del CA-4 de la HU-01.

## 5. Evidencias de aplicación de TDD

La regla del equipo es que cada fase de TDD sea un commit separado: `test:` (el test que falla), `feat:` (el código que lo hace pasar) y, si hace falta, `refactor:` (ordenar sin cambiar el comportamiento).

**Commits por tipo y por integrante** (sin contar los merges):

| Integrante | `test:` | `feat:` | `refactor:` | Total de commits |
|---|---|---|---|---|
| Martina Batista | 16 | 17 | 2 | 65 |
| Florencia Cabeza | 18 | 19 | 5 | 63 |
| Luana Guajardo | 18 | 18 | 1 | 52 |
| **Total** | **52** | **54** | **8** | **180** |

El resto de los commits son de documentación (`docs:`, 55) y mantenimiento (`chore:`, 7). Hay más `feat:` que `test:` porque las pantallas del frontend se commitean como `feat:` sin un `test:` antes: no tienen reglas de negocio que probar.

**Orden de los commits en cada historia del backend** (PR de código, en el orden en que se hicieron):

| PR | Historia | Integrante | Secuencia |
|---|---|---|---|
| #34 | HU-01 Proyectos | Martina | test → feat → test → feat → test → feat → test → feat → test → test |
| #37 | HU-04 Product Backlog | Florencia | test → test → feat → refactor → test → feat → test → feat → test → feat → refactor → test → feat → test → feat → test |
| #39 | HU-02 Integrantes | Luana | test → test → feat → test → feat → refactor → test → feat → test → feat → test |
| #41 | HU-03 Estado del proyecto | Luana | test → test → feat → test → feat → test → feat → test |
| #44 | HU-05 Crear Sprints | Martina | test → feat → test → feat |
| #46 | HU-06 Cerrar Sprint | Martina | test → feat → refactor → test → feat |
| #48 | HU-07 Story Points | Florencia | test → feat → refactor → test → feat |
| #50 | HU-08 Planning Poker | Florencia | test → feat → test → feat |
| #53 | HU-10 Métricas de avance | Martina | test → feat → test → feat |
| #56 | HU-09 Registrar esfuerzo | Luana | test → feat → test → feat |
| #58 | HU-13 Dashboard | Florencia | test → feat → refactor → test → feat → refactor |
| #69 | HU-12 Métricas de esfuerzo | Martina | test → feat → refactor → test → feat → feat (pantalla) |
| #72 | HU-11 Defectos | Luana | test → feat → test → feat → feat (pantalla) |
| #74 | HU-14 Reporte PDF | Luana | test → feat → test → feat → feat (pantalla) |

**En todas las historias del backend, cada `feat:` tiene antes su `test:`.** Las primeras historias hicieron un ciclo por capa (dominio, base de datos, API y escenarios BDD); las siguientes, dos ciclos: uno para el dominio (reglas y cálculos) y otro para la integración (API, base de datos y escenarios BDD). El último `feat` de las historias más recientes es su pantalla, que no lleva test porque no tiene reglas de negocio. Los refactors aparecen cuando el código nuevo repetía algo existente. Por ejemplo, en la HU-12 se unificó la lectura de horas con la del registro de esfuerzo, y la regla de "historia completada" con la de los Story Points.

**Un ejemplo de lo que aporta TDD:** en la HU-12, un test escrito antes del código mostró que modificar una historia desde el backlog perdía sus horas estimadas. Se corrigió antes de llegar a `main`.

## 6. Gestión del proyecto

**Cumplimiento de los Sprints** (Issues de cada Milestone):

| Sprint | Historias | Issues cerrados |
|---|---|---|
| Sprint 0 — La Preparación | — | 3 de 3 |
| Sprint 1 — El MVP | HU-01 a HU-07 | 8 de 8 |
| Sprint 2 — La Interfaz | HU-08, HU-09, HU-10 e interfaz web | 5 de 5 |
| Sprint 3 — Funcionalidad y Calidad | HU-11, HU-12, HU-13 | 4 de 4 |
| Sprint 4 — Cierre y Entrega Final | HU-14 | 1 de 5 (en curso) |

Las historias de los Sprints 3 y 4 se terminaron antes de lo planificado, en paralelo con el Sprint 2.

**Pull Requests:** 52 mergeados a `main`, todos con la CI en verde y aprobados por una integrante distinta de la autora (`main` está protegida y no acepta commits directos).

| Integrante | PRs |
|---|---|
| Martina Batista | 23 |
| Florencia Cabeza | 17 |
| Luana Guajardo | 12 |

Cada historia pasó por dos PRs: primero la spec (rama `docs/spec-NN-…`), y recién con la spec aprobada el código (rama `feat/…`, con `Closes #N`).

**Integración continua:** 111 ejecuciones de GitHub Actions, todas en verde. En cada PR corren `gofmt`, `go vet`, los tests con `-race` y los escenarios BDD del backend, y `lint`, `typecheck` y `build` del frontend.

## 7. Cómo volver a medir

```bash
# Cobertura por paquete (con la base levantada: docker compose up -d)
cd backend
go test ./... -cover

# Cobertura total, contando lo que cada paquete ejecuta desde los tests de otro
go test -p 1 -coverpkg=./... -coverprofile=cobertura.out ./...
go tool cover -func=cobertura.out | tail -n 1
go tool cover -html=cobertura.out          # informe navegable por archivo

# Escenarios BDD
go test -v ./features/ | tail -n 3

# Frontend
cd ../frontend
npm test

# Commits por tipo (evidencia de TDD)
git log main --no-merges --format=%s | sed -E 's/^([a-z]+).*/\1/' | sort | uniq -c
```
