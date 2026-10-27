import { useEffect, useState, type FormEvent } from 'react'
import { mensajeDe, pedir, type EsfuerzoDeHistoria, type Historia, type Integrante, type Proyecto } from '../api'
import { Cargando, MensajeError } from '../componentes'

const registroVacio = { integrante_id: '', fecha: '', actividad: '', horas: '' }

// Esfuerzo registra las horas trabajadas en cada historia (spec 09): se elige una historia del
// backlog, se ven sus registros por fecha con el total y se cargan o eliminan registros. El orden,
// el total y todas las validaciones (fecha dentro del proyecto y no futura, horas hasta 24 con 2
// decimales) son del backend: la pantalla muestra sus mensajes de error.
export default function Esfuerzo({ proyecto }: { proyecto: Proyecto }) {
  const rutaProyecto = `/proyectos/${proyecto.id}`
  const [historias, setHistorias] = useState<Historia[] | null>(null)
  const [integrantes, setIntegrantes] = useState<Integrante[]>([])
  const [historiaID, setHistoriaID] = useState('')
  const [esfuerzo, setEsfuerzo] = useState<EsfuerzoDeHistoria | null>(null)
  const [nuevo, setNuevo] = useState(registroVacio)
  const [error, setError] = useState('')

  const rutaEsfuerzo = `${rutaProyecto}/backlog/${historiaID}/esfuerzo`

  useEffect(() => {
    Promise.all([pedir<Historia[]>('GET', `${rutaProyecto}/backlog`), pedir<Integrante[]>('GET', `${rutaProyecto}/integrantes`)])
      .then(([listaHistorias, listaIntegrantes]) => {
        setHistorias(listaHistorias)
        setIntegrantes(listaIntegrantes)
      })
      .catch((e) => setError(mensajeDe(e)))
  }, [rutaProyecto])

  // elegir cambia de historia y pide su esfuerzo.
  async function elegir(id: string) {
    setHistoriaID(id)
    setEsfuerzo(null)
    setError('')
    if (id === '') {
      return
    }
    try {
      setEsfuerzo(await pedir<EsfuerzoDeHistoria>('GET', `${rutaProyecto}/backlog/${id}/esfuerzo`))
    } catch (e) {
      setError(mensajeDe(e))
    }
  }

  // hacer ejecuta un pedido que cambia el esfuerzo y, si sale bien, lo vuelve a pedir (así el orden
  // y el total siempre son los del backend). Devuelve si salió bien.
  async function hacer(pedido: () => Promise<unknown>): Promise<boolean> {
    setError('')
    try {
      await pedido()
      setEsfuerzo(await pedir<EsfuerzoDeHistoria>('GET', rutaEsfuerzo))
      return true
    } catch (e) {
      setError(mensajeDe(e))
      return false
    }
  }

  async function registrar(evento: FormEvent) {
    evento.preventDefault()
    const datos = {
      ...nuevo,
      // Vacíos no se mandan, para que el backend informe "El integrante es obligatorio" o
      // "Las horas son obligatorias".
      integrante_id: nuevo.integrante_id === '' ? undefined : Number(nuevo.integrante_id),
      horas: nuevo.horas === '' ? undefined : Number(nuevo.horas),
    }
    if (await hacer(() => pedir('POST', rutaEsfuerzo, datos))) {
      setNuevo(registroVacio)
    }
  }

  function eliminar(id: number, actividad: string) {
    if (window.confirm(`¿Eliminar el registro "${actividad}"?`)) {
      hacer(() => pedir('DELETE', `${rutaEsfuerzo}/${id}`))
    }
  }

  if (historias === null) {
    return error ? <MensajeError mensaje={error} /> : <Cargando />
  }

  return (
    <>
      <section>
        <h2>Esfuerzo</h2>
        {historias.length === 0 ? (
          <p className="suave">El backlog está vacío: primero agregá una historia.</p>
        ) : (
          <label>
            Historia
            <select value={historiaID} onChange={(e) => elegir(e.target.value)}>
              <option value="">Elegí una historia</option>
              {historias.map((h) => (
                <option key={h.id} value={h.id}>
                  {h.prioridad}. {h.titulo} ({h.estado})
                </option>
              ))}
            </select>
          </label>
        )}
        <MensajeError mensaje={error} />
        {historiaID !== '' && esfuerzo === null && !error && <Cargando />}
        {esfuerzo && (
          <>
            <p>
              <strong>Total: {esfuerzo.total_horas} horas</strong>
            </p>
            {esfuerzo.registros.length === 0 && <p className="suave">Esta historia todavía no tiene horas registradas.</p>}
            <ul className="lista">
              {esfuerzo.registros.map((r) => (
                <li key={r.id}>
                  <div>
                    <strong>
                      {r.horas} h · {r.actividad}
                    </strong>
                    <p className="suave">
                      {r.fecha} · {r.integrante_nombre}
                    </p>
                  </div>
                  <button className="secundario" onClick={() => eliminar(r.id, r.actividad)}>
                    Eliminar
                  </button>
                </li>
              ))}
            </ul>
          </>
        )}
      </section>

      {historiaID !== '' && (
        <section>
          <h2>Registrar horas</h2>
          {integrantes.length === 0 && <p className="suave">El proyecto todavía no tiene integrantes: agregalos en "Integrantes".</p>}
          <form onSubmit={registrar}>
            <label>
              Integrante
              <select value={nuevo.integrante_id} onChange={(e) => setNuevo({ ...nuevo, integrante_id: e.target.value })}>
                <option value="">Elegí un integrante</option>
                {integrantes.map((i) => (
                  <option key={i.id} value={i.id}>
                    {i.nombre} ({i.rol})
                  </option>
                ))}
              </select>
            </label>
            <label>
              Fecha
              <input type="date" value={nuevo.fecha} onChange={(e) => setNuevo({ ...nuevo, fecha: e.target.value })} />
            </label>
            <label>
              Actividad realizada
              <input value={nuevo.actividad} onChange={(e) => setNuevo({ ...nuevo, actividad: e.target.value })} />
            </label>
            <label>
              Horas trabajadas (por ejemplo, 1,5)
              {/* step="any": el navegador no limita los decimales ni el rango; eso lo valida el backend. */}
              <input type="number" step="any" value={nuevo.horas} onChange={(e) => setNuevo({ ...nuevo, horas: e.target.value })} />
            </label>
            <button type="submit">Registrar</button>
          </form>
        </section>
      )}
    </>
  )
}
