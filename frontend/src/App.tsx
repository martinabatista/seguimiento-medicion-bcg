import { useState, type ReactNode } from 'react'
import type { Proyecto } from './api'
import Backlog from './paginas/Backlog'
import Dashboard from './paginas/Dashboard'
import EnConstruccion from './paginas/EnConstruccion'
import Integrantes from './paginas/Integrantes'
import Metricas from './paginas/Metricas'
import PlanningPoker from './paginas/PlanningPoker'
import Proyectos from './paginas/Proyectos'
import Sprints from './paginas/Sprints'

// Las secciones del menú de un proyecto. Para sumar una pantalla: crear el archivo en paginas/ y
// cambiar su línea acá (pantalla recibe el proyecto elegido).
const secciones: { titulo: string; pantalla: (proyecto: Proyecto) => ReactNode }[] = [
  { titulo: 'Estado', pantalla: () => <EnConstruccion titulo="Estado" /> },
  { titulo: 'Backlog', pantalla: (proyecto) => <Backlog proyecto={proyecto} /> },
  { titulo: 'Integrantes', pantalla: (proyecto) => <Integrantes proyecto={proyecto} /> },
  { titulo: 'Sprints', pantalla: (proyecto) => <Sprints proyecto={proyecto} /> },
  { titulo: 'Planning Poker', pantalla: (proyecto) => <PlanningPoker proyecto={proyecto} /> },
  { titulo: 'Esfuerzo', pantalla: () => <EnConstruccion titulo="Esfuerzo" /> },
  { titulo: 'Métricas', pantalla: (proyecto) => <Metricas proyecto={proyecto} /> },
  { titulo: 'Dashboard', pantalla: (proyecto) => <Dashboard proyecto={proyecto} /> },
]

// App es el layout: el encabezado y, según haya un proyecto elegido o no, la lista de proyectos o
// el menú con la sección elegida. La navegación es un estado: no usa librerías de rutas.
export default function App() {
  const [proyecto, setProyecto] = useState<Proyecto | null>(null)
  const [seccion, setSeccion] = useState(0)

  function elegir(p: Proyecto) {
    setProyecto(p)
    setSeccion(0)
  }

  return (
    <>
      <header>
        <h1>Software Metrics &amp; Estimation</h1>
        {proyecto && (
          <p>
            Proyecto: <strong>{proyecto.nombre}</strong>{' '}
            <button className="enlace" onClick={() => setProyecto(null)}>
              Cambiar de proyecto
            </button>
          </p>
        )}
      </header>

      {proyecto === null ? (
        <Proyectos alElegir={elegir} />
      ) : (
        <>
          <nav>
            {secciones.map((s, i) => (
              <button key={s.titulo} aria-current={i === seccion ? 'page' : undefined} onClick={() => setSeccion(i)}>
                {s.titulo}
              </button>
            ))}
          </nav>
          <main>{secciones[seccion].pantalla(proyecto)}</main>
        </>
      )}
    </>
  )
}
