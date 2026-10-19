---
name: "especificar-historia"
description: "Usar al convertir una historia de usuario del TP en spec SDD (docs/specs/sprint-NN/NN-slug.md), con criterios de aceptación y escenarios BDD en Gherkin en español, trazables."
---

# Especificar una historia (SDD + BDD)

Se usa ANTES de escribir código para cualquier funcionalidad no trivial. La spec describe **qué** y
**para qué**, nunca código. Va en su propia rama `docs/spec-NN-slug` y su propio PR; el código
viene después en otra rama `feat/slug` (ver `CLAUDE.md` raíz).

## Convenciones del repo

- Spec: `docs/specs/sprint-NN/NN-slug.md` (`sprint-NN` = Milestone; `NN` = número de la HU, p. ej.
  `sprint-01/07-story-points.md`). Ya hay un borrador por cada HU del backlog: **completarlo**, no
  crear otro archivo.
- Plantilla: `docs/specs/_template.md`. Respetar su frontmatter y sus secciones.
- Issue de GitHub: campo `issue: "#N"` y link en la sección Historia de usuario.
- IDs **locales a la spec**: `RN-k` (reglas de negocio) y `CA-k` (criterios de aceptación). Para
  referirse a uno desde afuera se cita con la spec: `spec 07 CA-2`.
- Estados: `draft` → `in-review` (PR abierto) → `approved` (la aprueba **otra** integrante; recién
  ahí se programa) → `implemented` (se mergeó el PR de código). Esta skill deja la spec en `draft`
  o `in-review`; **nunca** la pasa a `approved` ni completa `aprobada_por`.
- Los `.feature` van en `backend/features/<funcionalidad>.feature` (ver más abajo).

## Pasos

1. Leer la HU y su spec borrador. Si falta rol, acción o beneficio, preguntar a la integrante.
   Revisar `Dependencias` y las specs relacionadas (`specs_relacionadas`) para no contradecirlas.
2. Completar **todas** las secciones de la plantilla, reemplazando cada `_TODO_`:
   Historia de usuario · Objetivo · Alcance (Entra / NO entra) · Entradas (campo, tipo,
   obligatoriedad, validación) · Salidas esperadas · Reglas de negocio (`RN-k`, una por línea,
   verificables) · Restricciones · Casos límite · Condiciones de error (tabla situación →
   respuesta) · Criterios de aceptación (`CA-k`, observables y testeables) · Escenarios BDD ·
   Tareas · Definition of Done.
3. Reglas de contenido:
   - **Sin código** en la spec (ni SQL, ni structs, ni rutas con handlers). Solo comportamiento.
   - Las reglas de negocio y los cálculos se implementan en Go; el frontend solo presenta y
     captura datos. Dejarlo en Restricciones.
   - Nada de valores inventados: si la consigna o el Issue no definen algo (escala, límites,
     fórmulas), **preguntar** y registrar la decisión en la spec.
4. Escenarios BDD en Gherkin español dentro de la spec, y copiados al `.feature` al implementar.
   Palabras clave: `Característica`, `Antecedentes`, `Escenario`, `Esquema del escenario`,
   `Ejemplos`, `Dado`, `Cuando`, `Entonces`, `Y`, `Pero`. Por cada CA, como mínimo un caso
   **normal**, uno **alternativo**, uno **límite** y uno de **error**. Usar `Esquema del escenario`
   para variantes de datos. Pasos en lenguaje de negocio (sin HTTP, SQL ni nombres de funciones).
5. **Cada escenario indica en su nombre el criterio que cubre**: `Escenario: ... (CA-3)`. Es la
   convención del repo y de la skill `chequear-trazabilidad`.
6. Verificar: todo CA tiene al menos un escenario y todo escenario nombra un CA; las RN están
   cubiertas por algún CA; no quedan `_TODO_`.
7. Si la spec cambia el contrato de la API o afecta a otra spec, actualizar esa otra spec en el
   mismo PR y anotarla en `specs_relacionadas`.

## Archivo `.feature` (lo crea la rama `feat/`, no la de la spec)

- `backend/features/<funcionalidad>.feature`, con `# language: es` en la primera línea.
- Los pasos se registran en `backend/features/features_test.go` (godog). Reutilizar pasos
  existentes antes de crear nuevos.
- Antecedentes (`Antecedentes:`) para el estado común; un escenario = un comportamiento.

## Cierre

- Frontmatter: `estado: draft` mientras se escribe; `in-review` al abrir el PR. `fecha` en
  formato `YYYY-MM-DD`; `autor` = quien escribe la spec.
- PR de la spec: usar `.github/pull_request_template.md`, referenciar el Issue **sin cerrarlo**
  ("Spec de #N") y pedir aprobación a otra integrante.
- Commit sugerido (Conventional Commits, scope = unidad): `docs(docs): spec NN <slug>`.
  Commitear y pushear **solo cuando la integrante lo pida**.

## Para entender, no solo generar

Al terminar, explicar en 3-4 líneas cuál es la regla de negocio más delicada y por qué, para que
el equipo pueda defenderla en la presentación. Lo generado con IA se declara en la sección
"Uso de IA" del PR.
