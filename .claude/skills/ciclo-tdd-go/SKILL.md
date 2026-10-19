---
name: "ciclo-tdd-go"
description: "Usar al implementar en Go una regla de negocio, métrica o validación del TP siguiendo RED → GREEN → REFACTOR, con una fase por commit y trazabilidad a la spec y sus criterios (CA-k)."
---

# Ciclo TDD en Go con evidencia en git

La cátedra evalúa el TDD mirando el historial de git: **cada fase es un commit separado**. Nunca
juntar test e implementación en el mismo commit.

## Antes de empezar

- Debe existir la spec en `docs/specs/sprint-NN/NN-slug.md` con `estado: approved` (la aprobó otra
  integrante y está mergeada en `main`). Si no, **frenar**: usar la skill `especificar-historia`.
- Trabajar en la rama `feat/slug`, creada desde `main` actualizado. Nunca en `main` (un hook lo
  bloquea).
- Elegir **UN** criterio `CA-k` o regla `RN-k` por ciclo. Ciclos chicos.
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
4. Commit: `test(backend): <qué verifica> (spec NN CA-k)`

## GREEN

1. Implementar lo **mínimo** para que pase. Nada extra, sin hardcodear valores del test.
2. `go test ./...` en verde.
3. Commit: `feat(backend): <qué implementa> (spec NN CA-k)`

## REFACTOR

1. Mejorar nombres, extraer funciones, quitar duplicación; errores explícitos e idiomáticos
   (`errors.New`, `fmt.Errorf("...: %w", err)`, errores centinela del dominio).
2. `gofmt -l .` (debe salir vacío), `go vet ./...` y `go test ./...` limpios.
3. Si no hay nada que mejorar, decirlo y **no hacer commit vacío**.
4. Commit: `refactor(backend): <qué mejoró>`

## Escenarios BDD (godog)

- Cada CA de la spec debe quedar cubierto por un escenario en `backend/features/*.feature` o por un
  test. Escribir el escenario y sus pasos en `features_test.go` también en orden RED → GREEN.
- El escenario nombra su criterio: `Escenario: ... (CA-k)`.
- Si el escenario es lo único que cubre un CA, el commit RED es el `.feature` + pasos en rojo
  (`test(backend): escenario <nombre> (spec NN CA-k)`).

## Al cerrar el ciclo

- Mostrar la cobertura: `go test ./... -cover`. No hay % mínimo fijado todavía (pendiente en la
  spec 00); las métricas, cálculos y validaciones deben tener tests unitarios sí o sí.
- Verificar el gate antes del PR: `gofmt`, `go vet`, `go test ./...` (la CI corre con `-race`).
  **Sin `--no-verify`.**
- Actualizar `backend/CLAUDE.md` (specs relacionadas, trampas) y pasar la spec a `implemented`
  en el mismo PR. El PR usa `.github/pull_request_template.md`, con `Closes #N`.
- Explicar brevemente cada línea nueva para que la integrante pueda justificarla en la defensa.

## Commits

- Conventional Commits con scope = unidad (`backend`, `frontend`, `docs`, `ci`).
- Commitear y pushear **solo cuando la integrante lo pida**. Los commits son de la integrante:
  sin líneas `Co-Authored-By` ni menciones a herramientas de IA en el mensaje.
