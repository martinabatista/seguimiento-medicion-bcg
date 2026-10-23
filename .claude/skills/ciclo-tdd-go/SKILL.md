---
name: "ciclo-tdd-go"
description: "Usar al implementar en Go una regla de negocio, métrica o validación del TP siguiendo RED → GREEN → REFACTOR, con una fase por commit y trazabilidad a la spec y sus criterios (CA-k)."
---

# Ciclo TDD en Go con evidencia en git

La cátedra evalúa el TDD mirando el historial de git: **cada fase es un commit separado**. Nunca
juntar test e implementación en el mismo commit. Pero **sin commits de relleno**: no un commit por
cada cosa que se hace, sino los necesarios para que el historial cuente lo que se hizo (ver
"Commits" al final).

## Antes de empezar

- Debe existir la spec en `docs/specs/sprint-NN/NN-slug.md` con `estado: approved` (la aprobó otra
  integrante y está mergeada en `main`). Si no, **frenar**: usar la skill `especificar-historia`.
- Trabajar en la rama `feat/slug`, creada desde `main` actualizado. Nunca en `main` (un hook lo
  bloquea).
- Trabajar **por capa**, no por criterio: un ciclo para el dominio (todas sus reglas) y otro para la
  integración (repositorio PostgreSQL + API + escenarios BDD). Dentro de cada ciclo se puede avanzar
  de a una regla corriendo los tests, pero **el commit es uno por fase y por capa**.
- Dónde va el código (ver `backend/CLAUDE.md`):
  - Reglas, cálculos, métricas y validaciones → `backend/dominio/` (Go puro, **sin** HTTP ni DB).
  - Acceso a datos → `backend/postgres/`. Rutas REST → `backend/api/`.
  - Los paquetes que no existen se crean con la primera historia que los necesite.
- Comandos desde `backend/`: `go test ./...`, `go vet ./...`, `gofmt -l .`.

## RED

1. Escribir el test en el paquete que corresponda, preferentemente table-driven:
   ```go
   // spec 10 CA-2: la velocidad es el promedio de Story Points completados por Sprint
   func TestVelocidad(t *testing.T) {
       casos := []struct{ nombre string; /* entradas */; esperado float64 }{ /* ... */ }
       for _, c := range casos {
           t.Run(c.nombre, func(t *testing.T) { /* ... */ })
       }
   }
   ```
2. El comentario del test cita la spec y el criterio (`spec NN CA-k` / `spec NN RN-k`): de ahí sale
   la trazabilidad. Casos de prueba: normal, alternativo, límite y error, igual que en la spec.
3. Correr `go test ./...` y confirmar que **falla por la razón correcta** (no por un error de
   compilación ajeno). Si hace falta, crear la firma devolviendo el valor cero para que compile.
4. Commit: `test(backend): <qué verifica> (spec NN RN-a a RN-b / CA-...)`, con todos los tests de
   la capa juntos.

## GREEN

1. Implementar lo **mínimo** para que pase. Nada extra, sin hardcodear valores del test.
2. `go test ./...` en verde.
3. Commit: `feat(backend): <qué implementa> (spec NN ...)`, con todo el código que hace pasar
   los tests de ese commit `test`.

## REFACTOR

1. Mejorar nombres, extraer funciones, quitar duplicación; errores explícitos e idiomáticos
   (`errors.New`, `fmt.Errorf("...: %w", err)`, errores centinela del dominio).
2. `gofmt -l .` (debe salir vacío), `go vet ./...` y `go test ./...` limpios.
3. Si no hay nada que mejorar, decirlo y **no hacer commit vacío**. Un `refactor` va solo si cambió
   algo de verdad (nombres, duplicación, estructura), nunca para formato o un typo.
4. Commit: `refactor(backend): <qué mejoró>`

## Escenarios BDD (godog)

- Cada CA de la spec debe quedar cubierto por un escenario en `backend/features/*.feature` o por un
  test. Escribir el escenario y sus pasos en `features_test.go` también en orden RED → GREEN.
- El escenario nombra su criterio: `Escenario: ... (CA-k)`.
- El `.feature` y sus pasos en godog van **en el mismo commit `test`** que los tests de la API y del
  repositorio (ciclo de integración): no hay commits aparte para el `.feature` ni para los pasos.

## Al cerrar el ciclo

- Mostrar la cobertura: `go test ./... -cover`. No hay % mínimo fijado todavía (pendiente en la
  spec 00); las métricas, cálculos y validaciones deben tener tests unitarios sí o sí.
- Verificar el gate antes del PR: `gofmt`, `go vet`, `go test ./...` (la CI corre con `-race`).
  **Sin `--no-verify`.**
- Actualizar `backend/CLAUDE.md` (specs relacionadas, trampas) y pasar la spec a `implemented`
  en el mismo PR. El PR usa `.github/pull_request_template.md`, con `Closes #N`.
- Explicar brevemente cada línea nueva para que la integrante pueda justificarla en la defensa.

## Commits

Criterio: **cada commit es un paso con sentido propio** que ayuda a entender lo que se hizo. Si un
commit se puede describir solo como "y además…", probablemente va junto con otro; si un commit mezcla
dos cosas que conviene leer por separado, se separa. Estructura de referencia de un PR de historia
(`feat/slug`), que suele quedar en **4 a 6 commits** (una historia con dos entidades puede necesitar
un ciclo más de dominio, y está bien):

```
test(backend): reglas de <entidad> (spec NN RN-1 a RN-k)                    ← RED dominio
feat(backend): <entidad> y sus casos de uso en el dominio (spec NN)         ← GREEN dominio
refactor(backend): <qué mejoró>                                             ← solo si hace falta
test(backend): API, repositorio y escenarios BDD de <x> (spec NN CA-1 a CA-k) ← RED integración
feat(backend): tabla, repositorio PostgreSQL y rutas REST de <x> (spec NN)   ← GREEN integración
docs: <x> implementada (spec NN, backend/CLAUDE.md, README)                 ← cierre
```

- El paso de la spec a `approved` y a `implemented` va en el commit `docs` final, **no** en commits
  sueltos. La aprobación ya queda registrada en el PR de la spec.
- El PR de la spec (`docs/spec-NN-slug`) lleva **un solo commit**: `docs(docs): spec NN <slug>`.
- **No son commits:** correcciones de formato, typos o "me faltó un archivo". Se arreglan antes de
  pushear sumándolos al commit que corresponde (`git commit --amend` si es el último y no se pusheó).
- Nada de commits "WIP", "arreglo", "cambios" ni mensajes que no digan qué se hizo.
- Conventional Commits con scope = unidad (`backend`, `frontend`, `docs`, `ci`).
- Commitear y pushear **solo cuando la integrante lo pida**. Los commits son de la integrante:
  sin líneas `Co-Authored-By` ni menciones a herramientas de IA en el mensaje.
