## Spec implementada

<!-- Link a la spec en docs/specs/. PR de spec (rama docs/): pedir aprobación. PR de código (rama feat/): la spec ya debe estar aprobada y mergeada. -->

Closes #   <!-- solo en el PR de código; el PR de la spec referencia el Issue sin cerrarlo: "Spec de #N" -->

## Qué resuelve

<!-- El para qué. Una o dos líneas. -->

## Cambios realizados

<!-- Bullets de qué se tocó. -->

## Gate (verde antes de abrir el PR — sin `--no-verify`)

- [ ] backend: `gofmt` / `go vet` / `go test` ✅
- [ ] frontend: lint / typecheck / build ✅ (si aplica)
- [ ] escenarios BDD (godog) ✅ (si aplica)
- [ ] sin credenciales en el diff (`.env` en `.gitignore`)
- [ ] commits convencionales (scope = unidad)

## TDD y trazabilidad

- [ ] Los tests se escribieron antes del código (commits `test:` → `feat:` → `refactor:` en el historial).
- [ ] Cada criterio de aceptación (CA-x) de la spec está cubierto por un escenario BDD o un test.

## Doc viva

- [ ] `CLAUDE.md` de la unidad afectada actualizado (specs relacionadas, trampas).
- [ ] Spec pasada a `implemented`. Si esto afecta otra spec, la actualicé en el mismo PR.

## Uso de IA

<!-- Qué se generó con asistencia de IA y cómo se revisó/validó. -->

## Notas

<!-- Decisiones, trade-offs, deuda técnica. -->
