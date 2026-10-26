import { useEffect, useState, type FormEvent } from 'react'
import { mensajeDe, pedir, type Proyecto } from '../api'
import { Cargando, MensajeError } from '../componentes'

const proyectoVacio = { nombre: '', descripcion: '', fecha_inicio: '', fecha_fin: '' }

// Proyectos lista los proyectos, permite crear uno nuevo y elegir con cuál trabajar (spec 01).
// Es la pantalla de entrada: las demás secciones son de un proyecto elegido.
export default function Proyectos({ alElegir }: { alElegir: (p: Proyecto) => void }) {
  const [proyectos, setProyectos] = useState<Proyecto[] | null>(null)
  const [nuevo, setNuevo] = useState(proyectoVacio)
  const [error, setError] = useState('')

  useEffect(() => {
    pedir<Proyecto[]>('GET', '/proyectos')
      .then(setProyectos)
      .catch((e) => setError(mensajeDe(e)))
  }, [])

  async function crear(evento: FormEvent) {
    evento.preventDefault()
    setError('')
    try {
      const creado = await pedir<Proyecto>('POST', '/proyectos', nuevo)
      setProyectos([...(proyectos ?? []), creado])
      setNuevo(proyectoVacio)
    } catch (e) {
      setError(mensajeDe(e)) // por ejemplo: "La fecha de finalización debe ser posterior a la de inicio"
    }
  }

  return (
    <main>
      <section>
        <h2>Proyectos</h2>
        <MensajeError mensaje={error} />
        {proyectos === null && !error && <Cargando />}
        {proyectos?.length === 0 && <p className="suave">Todavía no hay proyectos: creá el primero.</p>}
        <ul className="lista">
          {proyectos?.map((p) => (
            <li key={p.id}>
              <div>
                <strong>{p.nombre}</strong>
                <span className="suave">
                  {' '}
                  · {p.fecha_inicio} al {p.fecha_fin}
                </span>
                <p className="suave">{p.descripcion}</p>
              </div>
              <button onClick={() => alElegir(p)}>Elegir</button>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Nuevo proyecto</h2>
        <form onSubmit={crear}>
          <label>
            Nombre
            <input value={nuevo.nombre} onChange={(e) => setNuevo({ ...nuevo, nombre: e.target.value })} />
          </label>
          <label>
            Descripción
            <textarea value={nuevo.descripcion} onChange={(e) => setNuevo({ ...nuevo, descripcion: e.target.value })} />
          </label>
          <label>
            Fecha de inicio
            <input type="date" value={nuevo.fecha_inicio} onChange={(e) => setNuevo({ ...nuevo, fecha_inicio: e.target.value })} />
          </label>
          <label>
            Fecha de finalización
            <input type="date" value={nuevo.fecha_fin} onChange={(e) => setNuevo({ ...nuevo, fecha_fin: e.target.value })} />
          </label>
          <button type="submit">Crear proyecto</button>
        </form>
      </section>
    </main>
  )
}
