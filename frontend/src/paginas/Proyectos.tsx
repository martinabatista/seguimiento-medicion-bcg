import { useEffect, useState, type FormEvent } from 'react'
import { mensajeDe, pedir, type Proyecto } from '../api'
import { Cargando, MensajeError } from '../componentes'

const proyectoVacio = { nombre: '', descripcion: '', fecha_inicio: '', fecha_fin: '' }

// Proyectos lista los proyectos, permite crear uno nuevo, modificarlo y elegir con cuál trabajar (spec 01).
// Es la pantalla de entrada: las demás secciones son de un proyecto elegido.
export default function Proyectos({ alElegir }: { alElegir: (p: Proyecto) => void }) {
  const [proyectos, setProyectos] = useState<Proyecto[] | null>(null)
  const [formulario, setFormulario] = useState(proyectoVacio)
  const [editando, setEditando] = useState<number | null>(null) // id del proyecto que se modifica
  const [error, setError] = useState('')

  useEffect(() => {
    pedir<Proyecto[]>('GET', '/proyectos')
      .then(setProyectos)
      .catch((e) => setError(mensajeDe(e)))
  }, [])

  // guardar crea el proyecto o, si se está modificando uno, reemplaza sus datos (CA-1, CA-2).
  async function guardar(evento: FormEvent) {
    evento.preventDefault()
    setError('')
    try {
      if (editando === null) {
        const creado = await pedir<Proyecto>('POST', '/proyectos', formulario)
        setProyectos([...(proyectos ?? []), creado])
      } else {
        const modificado = await pedir<Proyecto>('PUT', `/proyectos/${editando}`, formulario)
        setProyectos((proyectos ?? []).map((p) => (p.id === modificado.id ? modificado : p)))
      }
      cancelar()
    } catch (e) {
      setError(mensajeDe(e)) // por ejemplo: "La fecha de finalización debe ser posterior a la de inicio"
    }
  }

  function modificar(p: Proyecto) {
    setError('')
    setEditando(p.id)
    setFormulario({ nombre: p.nombre, descripcion: p.descripcion, fecha_inicio: p.fecha_inicio, fecha_fin: p.fecha_fin })
  }

  function cancelar() {
    setEditando(null)
    setFormulario(proyectoVacio)
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
              <div className="acciones">
                <button onClick={() => alElegir(p)}>Elegir</button>
                <button className="secundario" onClick={() => modificar(p)}>
                  Modificar
                </button>
              </div>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>{editando === null ? 'Nuevo proyecto' : 'Modificar proyecto'}</h2>
        <form onSubmit={guardar}>
          <label>
            Nombre
            <input value={formulario.nombre} onChange={(e) => setFormulario({ ...formulario, nombre: e.target.value })} />
          </label>
          <label>
            Descripción
            <textarea
              value={formulario.descripcion}
              onChange={(e) => setFormulario({ ...formulario, descripcion: e.target.value })}
            />
          </label>
          <label>
            Fecha de inicio
            <input
              type="date"
              value={formulario.fecha_inicio}
              onChange={(e) => setFormulario({ ...formulario, fecha_inicio: e.target.value })}
            />
          </label>
          <label>
            Fecha de finalización
            <input
              type="date"
              value={formulario.fecha_fin}
              onChange={(e) => setFormulario({ ...formulario, fecha_fin: e.target.value })}
            />
          </label>
          <div className="acciones">
            <button type="submit">{editando === null ? 'Crear proyecto' : 'Guardar cambios'}</button>
            {editando !== null && (
              <button type="button" className="secundario" onClick={cancelar}>
                Cancelar
              </button>
            )}
          </div>
        </form>
      </section>
    </main>
  )
}
