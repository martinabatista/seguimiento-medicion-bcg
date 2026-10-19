---
name: "chequear-trazabilidad"
description: "Usar para verificar o generar la matriz de trazabilidad del TP: HU (Issue) → spec SDD → criterios (CA) → escenarios BDD → tests → código Go, y detectar huecos."
---

# Chequear trazabilidad

La trazabilidad se demuestra en vivo en la presentación final:
**Historia (Issue) → Spec SDD → Criterios de aceptación → Escenarios BDD → Tests → Código Go.**
Esta skill arma la matriz y marca lo que falta. Es de **solo lectura** salvo el archivo de
resultado.

## Fuentes

- Specs: `docs/specs/sprint-NN/NN-slug.md`. Frontmatter: `issue`, `estado`, `sprint`. Los
  criterios son `CA-k` y las reglas `RN-k`, **locales a cada spec**: la clave única es
  `spec NN + CA-k`.
- Escenarios: `backend/features/*.feature`. El nombre de cada escenario indica el criterio:
  `Escenario: ... (CA-k)`. Los pasos están en `backend/features/features_test.go`.
- Tests: `backend/**/*_test.go`. Un comentario cita `spec NN CA-k` / `spec NN RN-k`.
- Código: funciones que ejercitan esos tests (`backend/dominio/`, `backend/api/`,
  `backend/postgres/`).
- Commits: `git log --oneline` (Conventional Commits `test(backend)` / `feat(backend)` /
  `refactor(backend)`). El número de spec va en el mensaje o en el cuerpo, y el PR cierra el
  Issue con `Closes #N`.
- Issues y PRs: `gh issue view N` / `gh pr list` si `gh` está disponible; si no, solo lo local.

## Pasos

1. Listar las specs y sus CA/RN. Ignorar `_template.md`.
   Ej.: `grep -rnE '^- (CA|RN)-[0-9]+' docs/specs --include=*.md`
2. Para cada `spec NN CA-k` buscar: escenario(s) `.feature` cuyo nombre lo cita, test(s) `_test.go`
   que lo citan, archivo y función Go que lo implementa, y los commits RED/GREEN/REFACTOR (en
   `git log` de la rama o de `main`). Para cada RN, qué CA o test la cubre.
3. Generar/actualizar `docs/trazabilidad.md` con una tabla por spec:
   `| Spec (Issue) | Estado | CA/RN | Escenario BDD | Test Go | Código | Commits TDD |`
4. Reportar huecos, ordenados por gravedad:
   - Spec en `implemented` con CA sin escenario ni test
   - CA sin escenario BDD
   - CA/RN sin test
   - Escenario sin CA en el nombre, o que cita un CA que no existe
   - Test que no cita ningún `spec NN CA-k` / `RN-k`
   - Código de dominio sin commit `test(backend)` previo (RED)
   - Escenario sin pasos definidos en godog (pendiente)
   - Spec con `_TODO_` pendientes pero con código ya escrito
   - Spec `implemented` sin PR/Issue vinculado, o spec `draft` con código
5. Correr desde `backend/`: `go test ./... -cover` (incluye los escenarios godog) y agregar el
   resumen al final de `docs/trazabilidad.md`.
6. **No inventar vínculos**: si algo no se encuentra, marcarlo como faltante. Una spec todavía
   en `draft` sin código no es un hueco: figura como "pendiente", no como error.

## Cuándo correrla

Al cerrar cada Sprint (antes de la Review), antes de la presentación final y antes de abrir un PR
de código, para confirmar que cada CA de la spec quedó cubierto.

## Nota

`docs/trazabilidad.md` es un archivo generado: se actualiza en un commit propio
(`docs(docs): actualizar matriz de trazabilidad`) y **solo se commitea cuando la integrante lo
pida**.
