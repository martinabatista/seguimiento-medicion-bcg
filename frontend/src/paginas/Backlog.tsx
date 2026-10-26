import { useEffect, useState, type FormEvent } from 'react'
import { mensajeDe, pedir, type Historia, type Proyecto } from '../api'
import { Cargando, MensajeError } from '../componentes'

const estados: Historia['estado'][] = ['Pendiente', 'En progreso', 'Completada']
const formularioVacio = { titulo: '', descripcion: '', prioridad: '', criterios_aceptacion: '', estado: 'Pendiente' }

// Backlog muestra el Product Backlog del proyecto (spec 04) y permite estimar cada historia con
// Story Points (spec 07). Las reglas están en el backend: la lista llega ya ordenada por prioridad y
// los errores se muestran con el mensaje de la API.
export default function Backlog({ proyecto }: { proyecto: Proyecto }) {
  const ruta = `/proyectos/${proyecto.id}/backlog`
  const [historias, setHistorias] = useState<Historia[] | null>(null)
  const [escala, setEscala] = useState<number[]>([])
  const [formulario, setFormulario] = useState(formularioVacio)
  const [editando, setEditando] = useState<number | null>(null) // id de la historia que se modifica
  const [error, setError] = useState('')

  useEffect(() => {
    Promise.all([pedir<Historia[]>('GET', ruta), pedir<{ valores: number[] }>('GET', '/story-points/escala')])
      .then(([lista, respuesta]) => {
        setHistorias(lista)
        setEscala(respuesta.valores)
      })
      .catch((e) => setError(mensajeDe(e)))
  }, [ruta])

  // hacer ejecuta un pedido que cambia el backlog y, si sale bien, vuelve a pedir la lista (así el
  // orden siempre es el del backend). Devuelve si salió bien.
  async function hacer(pedido: () => Promise<unknown>): Promise<boolean> {
    setError('')
    try {
      await pedido()
      setHistorias(await pedir<Historia[]>('GET', ruta))
      return true
    } catch (e) {
      setError(mensajeDe(e))
      return false
    }
  }

  async function guardar(evento: FormEvent) {
    evento.preventDefault()
    const datos = {
      ...formulario,
      // Vacía no se manda, para que el backend informe "La prioridad es obligatoria".
      prioridad: formulario.prioridad === '' ? undefined : Number(formulario.prioridad),
    }
    const salioBien = await hacer(() =>
      editando === null ? pedir('POST', ruta, datos) : pedir('PUT', `${ruta}/${editando}`, datos),
    )
    if (salioBien) {
      cancelar()
    }
  }

  function modificar(h: Historia) {
    setEditando(h.id)
    setFormulario({
      titulo: h.titulo,
      descripcion: h.descripcion,
      prioridad: String(h.prioridad),
      criterios_aceptacion: h.criterios_aceptacion,
      estado: h.estado,
    })
  }

  function cancelar() {
    setEditando(null)
    setFormulario(formularioVacio)
  }

  function eliminar(h: Historia) {
    if (window.confirm(`¿Eliminar la historia "${h.titulo}"?`)) {
      hacer(() => pedir('DELETE', `${ruta}/${h.id}`))
    }
  }

  // estimar asigna los Story Points elegidos, o quita la estimación si se eligió "Sin estimar".
  function estimar(h: Historia, valor: string) {
    hacer(() =>
      valor === ''
        ? pedir('DELETE', `${ruta}/${h.id}/story-points`)
        : pedir('PUT', `${ruta}/${h.id}/story-points`, { story_points: Number(valor) }),
    )
  }

  return (
    <>
      <section>
        <h2>Backlog</h2>
        <MensajeError mensaje={error} />
        {historias === null && !error && <Cargando />}
        {historias?.length === 0 && <p className="suave">El backlog está vacío: agregá la primera historia.</p>}
        <ul className="lista">
          {historias?.map((h) => (
            <li key={h.id}>
              <div>
                <strong>
                  {h.prioridad}. {h.titulo}
                </strong>
                <span className="suave"> · {h.estado}</span>
                <p className="suave">{h.descripcion}</p>
                {h.criterios_aceptacion && <p className="suave">Criterios: {h.criterios_aceptacion}</p>}
              </div>
              <div className="acciones">
                <label>
                  Story Points
                  <select value={h.story_points ?? ''} onChange={(e) => estimar(h, e.target.value)}>
                    <option value="">Sin estimar</option>
                    {escala.map((v) => (
                      <option key={v} value={v}>
                        {v}
                      </option>
                    ))}
                  </select>
                </label>
                <button onClick={() => modificar(h)}>Modificar</button>
                <button className="secundario" onClick={() => eliminar(h)}>
                  Eliminar
                </button>
              </div>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>{editando === null ? 'Nueva historia' : 'Modificar historia'}</h2>
        <form onSubmit={guardar}>
          <label>
            Título
            <input value={formulario.titulo} onChange={(e) => setFormulario({ ...formulario, titulo: e.target.value })} />
          </label>
          <label>
            Descripción
            <textarea value={formulario.descripcion} onChange={(e) => setFormulario({ ...formulario, descripcion: e.target.value })} />
          </label>
          <label>
            Prioridad (1 es la más importante)
            <input type="number" value={formulario.prioridad} onChange={(e) => setFormulario({ ...formulario, prioridad: e.target.value })} />
          </label>
          <label>
            Criterios de aceptación (opcional)
            <textarea
              value={formulario.criterios_aceptacion}
              onChange={(e) => setFormulario({ ...formulario, criterios_aceptacion: e.target.value })}
            />
          </label>
          {editando !== null && (
            <label>
              Estado
              <select value={formulario.estado} onChange={(e) => setFormulario({ ...formulario, estado: e.target.value })}>
                {estados.map((estado) => (
                  <option key={estado}>{estado}</option>
                ))}
              </select>
            </label>
          )}
          <div className="acciones">
            <button type="submit">{editando === null ? 'Agregar historia' : 'Guardar cambios'}</button>
            {editando !== null && (
              <button type="button" className="secundario" onClick={cancelar}>
                Cancelar
              </button>
            )}
          </div>
        </form>
      </section>
    </>
  )
}
