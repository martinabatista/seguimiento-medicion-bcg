---
description: Ponerse al día al arrancar una sesión — traer lo último, ver qué hicieron las compañeras y reportar dónde estamos antes de tocar nada.
---

Ponete al día antes de tocar nada. No codees, no edites archivos ni commitees hasta que te dé el OK.

## 1. Traer y leer

- `git switch main` y `git pull` (si hay cambios sin commitear, avisame y no cambies de rama).
- Leé `CLAUDE.md` (reglas del equipo).
- Qué se hizo últimamente: `git log --format="%h %an %ad %s" --date=short -15`.
- PRs: `gh pr list --state merged --limit 8` y `gh pr list --state open`. Marcá los abiertos que
  esperan MI aprobación (no los abrí yo y no los aprobé).
- Sprint actual: `gh api repos/martinabatista/seguimiento-medicion-bcg/milestones` → el Milestone
  abierto cuya fecha de cierre (`due_on`) es la más próxima a hoy.
- Issues de ese Sprint: `gh issue list --milestone "<título del Milestone>" --state all
  --json number,title,state,assignees`.
- Mis Issues abiertos: `gh issue list --assignee @me --state open`.
- Specs del Sprint actual: el campo `estado` de cada archivo en `docs/specs/sprint-NN/`.

## 2. Recordatorios

- Nunca commit directo a `main`: rama desde `main` actualizado → PR → aprobación de otra integrante.
- Cada historia: primero la spec (rama `docs/spec-NN-slug`, PR para aprobarla), después el código
  (rama `feat/slug`, PR con `Closes #N`).
- Crear la rama y commitear en comandos SEPARADOS: el hook revisa la rama antes de ejecutar.
- Una historia no se programa sin su spec en `approved`.

## 3. Reporte (pocas líneas, en español)

- **DÓNDE ESTAMOS:** Sprint actual (nombre y fechas) y estado de sus historias y specs.
- **QUÉ HICIERON LAS DEMÁS:** lo último de cada compañera (commits, PRs, specs).
- **QUÉ SIGUE PARA MÍ:** mis Issues abiertos y el próximo paso de cada uno.
- **PENDIENTES:** PRs esperando mi aprobación, bloqueos o decisiones abiertas.

Después PARÁ y esperá el OK.
