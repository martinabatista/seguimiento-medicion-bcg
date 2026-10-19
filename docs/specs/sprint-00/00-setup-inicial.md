---
titulo: Setup inicial del proyecto
estado: draft
unidad: global
issue: "#1"
sprint: 00
autor: Martina Batista
aprobada_por: ""
fecha: 2026-10-17
specs_relacionadas: []
---

# Spec 00 — Setup inicial del proyecto

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como equipo de desarrollo, quiero tener el repositorio, el tablero y el entorno listos, para
empezar el Sprint 1 trabajando en paralelo sin fricción.

## Objetivo

Cumplir el **Sprint 0 (La Preparación)** de la consigna: entorno configurado, tablero Scrum armado
y estructura base del código lista para aplicar SDD, BDD y TDD desde la primera historia.

## Alcance

**Entra:**
- Esqueleto del backend Go que compila y tiene al menos un test pasando.
- Esqueleto del frontend React + Vite + TS que compila.
- PostgreSQL levantable localmente de forma reproducible para las tres integrantes.
- Estructura para escenarios BDD con godog, con un escenario de ejemplo en verde.
- CI en GitHub Actions corriendo sobre PRs a `main`.
- Protección de la rama `main` en GitHub.
- GitHub Project (tablero), Milestones `Sprint 0` a `Sprint 4` y labels básicos.
- Product Backlog inicial cargado como Issues.

**NO entra:**
- Funcionalidades del producto (gestión de proyectos, backlog, etc.): van en sus propias specs.

## Entradas

No aplica (configuración del proyecto).

## Salidas esperadas

- `go test ./...` y `npm run build` pasan en una máquina recién clonada.
- El tablero muestra el Product Backlog inicial.

## Reglas de negocio

No aplica.

## Restricciones

- El núcleo y las reglas de negocio en Go (consigna).
- Cualquier integrante debe poder levantar el entorno siguiendo instrucciones escritas.

## Casos límite

- Integrante con Windows (rutas, finales de línea): el entorno debe funcionar igual.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| Un PR rompe el formato o los tests | La CI falla y bloquea el merge |
| Commit directo a `main` | Lo bloquea la protección de ramas de GitHub |

## Criterios de aceptación

- CA-1: El backend compila y `go test ./...` pasa en CI.
- CA-2: El frontend compila (`npm run build`) en CI.
- CA-3: Existe un escenario BDD de ejemplo ejecutado por godog en CI.
- CA-4: `main` está protegida (requiere PR + CI verde).
- CA-5: El tablero de GitHub Projects tiene los Milestones de los Sprints y el backlog inicial.
- CA-6: El README explica cómo levantar el entorno.

## Escenarios BDD

No aplica (configuración). El escenario de ejemplo de CA-3 solo valida la integración de godog.

## Tareas

- [ ] `[P1]` Instalar Go, Node y Docker en las máquinas del equipo (hecho en la de Martina).
- [x] `[P1]` Inicializar el módulo Go y la estructura de `backend/`.
- [ ] `[P1]` Inicializar el frontend con Vite + React + TS.
- [x] `[P1]` Definir cómo se levanta PostgreSQL localmente (docker compose, puerto 5434).
- [x] `[P1]` Integrar godog con un escenario de ejemplo.
- [x] `[P1]` Configurar la protección de la rama `main`.
- [x] `[P1]` Crear Project, Milestones y labels en GitHub.
- [ ] `[P2]` Definir % mínimo de cobertura.
- [x] `[P2]` Actualizar el README con instrucciones de setup.

## Definition of Done

- [ ] CI en verde.
- [ ] Criterios de aceptación verificados.
- [ ] `CLAUDE.md` de cada unidad actualizado.
- [ ] Spec en estado `implemented`.
