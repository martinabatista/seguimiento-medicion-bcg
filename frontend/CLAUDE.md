# frontend

> `CLAUDE.md` de la unidad frontend. Autosuficiente para su carpeta; el raíz lo indexa.

## Qué es

La interfaz web: presenta los datos y captura lo que carga la usuaria. **No tiene reglas de negocio**:
todas las validaciones y los cálculos están en el backend Go (regla de la cátedra). Si algo está mal,
la pantalla muestra el mensaje que manda la API, tal cual.

Pensada para ser lo más simple posible: React + Vite + TypeScript, sin librería de rutas ni de estilos.

## Estructura

```
frontend/
├── package.json · package-lock.json   # dependencias y comandos (dev, build, lint, typecheck, test)
├── index.html · vite.config.ts        # la página; Vite reenvía /api al backend (puerto 8080)
├── tsconfig.json · eslint.config.js   # TypeScript y ESLint (configuración mínima)
└── src/
    ├── main.tsx         # arranca React
    ├── App.tsx          # layout: encabezado, lista de proyectos o menú de secciones del proyecto elegido
    ├── api.ts           # pedir(): todos los pedidos a la API; ErrorAPI; tipos de los datos
    ├── api.test.ts      # tests del cliente de la API (Vitest)
    ├── componentes.tsx  # Cargando y MensajeError, comunes a todas las pantallas
    ├── estilos.css      # todos los estilos
    └── paginas/         # una pantalla por archivo
        ├── Proyectos.tsx       # listar, crear y elegir proyecto (spec 01)
        ├── Estado.tsx          # estado por fechas, Sprint en curso y avance del backlog (spec 03)
        ├── Backlog.tsx         # backlog con Story Points (specs 04 y 07)
        ├── Integrantes.tsx     # listar y agregar integrantes con su rol (spec 02)
        ├── Sprints.tsx         # crear Sprints, asignar y completar historias, cerrar y consultar (specs 05 y 06)
        ├── PlanningPoker.tsx   # votar en secreto, revelar, nueva ronda y acordar (spec 08)
        ├── Esfuerzo.tsx        # registrar, consultar y eliminar horas por historia (spec 09)
        ├── Defectos.tsx        # métricas de defectos, lista con cambio de estado y registro (spec 11)
        ├── Metricas.tsx        # avance (spec 10) y esfuerzo, con la carga de horas estimadas (spec 12)
        └── Dashboard.tsx       # indicadores y 3 gráficos en HTML y CSS (spec 13)
```

Todas las secciones del menú tienen su pantalla (la de "En construcción" se quitó al completar la #15).

## Cómo agregar una pantalla

1. Crear `src/paginas/MiPantalla.tsx` (tomar `Proyectos.tsx` como modelo). Recibe el proyecto elegido.
2. Sumar en `api.ts` los tipos de los datos que usa (copiando el JSON de `backend/CLAUDE.md`).
3. En `App.tsx`, agregar su línea en `secciones`: `{ titulo: 'Mi sección', pantalla: (proyecto) => <MiPantalla proyecto={proyecto} /> }`.
4. Para los pedidos, `pedir<Tipo>('GET', '/proyectos/' + proyecto.id + '/…')`. Mostrar la carga con
   `<Cargando />` y los errores con `<MensajeError mensaje={mensajeDe(e)} />`.

## Decisiones técnicas

Acordadas al crear la unidad (Issue #15, 2026-10-26):

- **Navegación sin librería:** un estado con la sección elegida y un menú. La dirección del navegador no
  cambia; alcanza para el TP y no suma conceptos.
- **Un solo cliente de la API** (`pedir`): arma el pedido JSON, devuelve los datos y, si hay error,
  lanza `ErrorAPI` con el mensaje del backend (`{"error": "…"}`) o uno claro si el servidor no responde.
- **Proxy de Vite** en desarrollo (`/api` → `http://localhost:8080`): el frontend y la API se ven como
  un solo sitio, así no hace falta configurar CORS en el backend.
- **Estilos en un solo archivo CSS**, sin librerías. Los gráficos del Dashboard también son HTML y CSS
  (barras con su alto en % y la dona con `conic-gradient`), sin librería de gráficos (spec 13).
- **Dependencias:** de la app, `react` y `react-dom`; de desarrollo, Vite, TypeScript, ESLint (mínimo) y
  Vitest. Cada dependencia nueva se justifica en el PR que la agrega.

## Límites

- No valida datos ni calcula métricas: eso es del backend.
- Se usa con `npm run dev` (para la demo alcanza). Publicarla en un servidor no está en el alcance.

## Trampas conocidas

- Si se cambia de rama con `npm run dev` andando y la rama no tiene `frontend/` (por ejemplo, una vieja), Vite
  se reinicia sin su configuración y **pierde el proxy**: los pedidos a `/api` devuelven la página. Solución:
  frenar `npm run dev` (Ctrl+C) y volver a lanzarlo.
- Algunos campos del navegador **validan solos** y no dejan enviar el formulario, así nunca se ve el mensaje
  del backend: `type="email"` (formato), `step` en `type="number"` (decimales), `required`, `min`/`max`.
  Para que la validación sea la de Go: email como texto común (`inputMode="email"`) y números con
  `type="number" step="any"` (Integrantes y Esfuerzo).

## Contratos públicos

Consume la API REST del backend: rutas y JSON en `backend/CLAUDE.md`. Si cambia un JSON, actualizar
sus tipos en `api.ts`.

## Cómo testear

```bash
cd frontend
npm test            # Vitest
npm run lint        # ESLint
npm run typecheck   # TypeScript
npm run build       # la CI corre lint, typecheck y build
```

## Estado

Origen: interfaz web base (Issue #15). Última revisión: 2026-10-27 (pantalla de Defectos, spec 11).
