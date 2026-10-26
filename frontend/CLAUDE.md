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
        └── EnConstruccion.tsx  # lugar de las secciones que todavía no tienen pantalla
```

## Cómo agregar una pantalla

1. Crear `src/paginas/MiPantalla.tsx` (tomar `Proyectos.tsx` como modelo). Recibe el proyecto elegido.
2. Sumar en `api.ts` los tipos de los datos que usa (copiando el JSON de `backend/CLAUDE.md`).
3. En `App.tsx`, cambiar su línea de `secciones`: `pantalla: (proyecto) => <MiPantalla proyecto={proyecto} />`.
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
- **Estilos en un solo archivo CSS**, sin librerías. La librería de gráficos se elige con el Dashboard
  (spec 13, etapa 2).
- **Dependencias:** de la app, `react` y `react-dom`; de desarrollo, Vite, TypeScript, ESLint (mínimo) y
  Vitest. Cada dependencia nueva se justifica en el PR que la agrega.

## Límites

- No valida datos ni calcula métricas: eso es del backend.
- Se usa con `npm run dev` (para la demo alcanza). Publicarla en un servidor no está en el alcance.

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

Origen: interfaz web base (Issue #15). Última revisión: 2026-10-26.
