import { useEffect, useState, type FormEvent } from 'react'
import { mensajeDe, pedir, type DetalleSprint, type Historia, type Proyecto, type Sprint } from '../api'
import { Cargando, MensajeError } from '../componentes'

const sprintVacio = { fecha_inicio: '', fecha_fin: '', sprint_goal: '' }

// Sprints lista los Sprints del proyecto y permite crear uno (spec 05). Al elegir un Sprint se ve su
// detalle (DetalleDelSprint). Los cerrados se consultan igual, como Sprints anteriores (spec 06).
export default function Sprints({ proyecto }: { proyecto: Proyecto }) {
  const ruta = `/proyectos/${proyecto.id}/sprints`
  const [sprints, setSprints] = useState<Sprint[] | null>(null)
  const [nuevo, setNuevo] = useState(sprintVacio)
  const [elegido, setElegido] = useState<number | null>(null) // id del Sprint que se está viendo
  const [error, setError] = useState('')

  useEffect(() => {
    if (elegido === null) {
      pedir<Sprint[]>('GET', ruta)
        .then(setSprints)
        .catch((e) => setError(mensajeDe(e)))
    }
  }, [ruta, elegido])

  async function crear(evento: FormEvent) {
    evento.preventDefault()
    setError('')
    try {
      await pedir('POST', ruta, nuevo)
      setNuevo(sprintVacio)
      setSprints(await pedir<Sprint[]>('GET', ruta))
    } catch (e) {
      setError(mensajeDe(e)) // por ejemplo: "El Sprint se superpone con el Sprint 1"
    }
  }

  if (elegido !== null) {
    return <DetalleDelSprint proyecto={proyecto} sprintID={elegido} alVolver={() => setElegido(null)} />
  }

  return (
    <>
      <section>
        <h2>Sprints</h2>
        <MensajeError mensaje={error} />
        {sprints === null && !error && <Cargando />}
        {sprints?.length === 0 && <p className="suave">Todavía no hay Sprints: creá el primero.</p>}
        <ul className="lista">
          {sprints?.map((s) => (
            <li key={s.id}>
              <div>
                <strong>Sprint {s.numero}</strong>
                <span className="suave">
                  {' '}
                  · {s.fecha_inicio} al {s.fecha_fin} · {s.estado}
                </span>
                <p className="suave">{s.sprint_goal}</p>
              </div>
              <button onClick={() => setElegido(s.id)}>Ver</button>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Nuevo Sprint</h2>
        <form onSubmit={crear}>
          <label>
            Fecha de inicio
            <input type="date" value={nuevo.fecha_inicio} onChange={(e) => setNuevo({ ...nuevo, fecha_inicio: e.target.value })} />
          </label>
          <label>
            Fecha de fin
            <input type="date" value={nuevo.fecha_fin} onChange={(e) => setNuevo({ ...nuevo, fecha_fin: e.target.value })} />
          </label>
          <label>
            Sprint Goal
            <textarea value={nuevo.sprint_goal} onChange={(e) => setNuevo({ ...nuevo, sprint_goal: e.target.value })} />
          </label>
          <div className="acciones">
            <button type="submit">Crear Sprint</button>
          </div>
        </form>
      </section>
    </>
  )
}

type PropsDetalle = { proyecto: Proyecto; sprintID: number; alVolver: () => void }

// DetalleDelSprint muestra un Sprint con sus historias. Si está abierto se puede cambiar el Sprint
// Goal, asignar y quitar historias, completarlas o deshacerlo, y cerrarlo (specs 05 y 06). Si está
// cerrado, solo se consulta: qué se planificó y qué se completó en él.
function DetalleDelSprint({ proyecto, sprintID, alVolver }: PropsDetalle) {
  const ruta = `/proyectos/${proyecto.id}/sprints/${sprintID}`
  const [sprint, setSprint] = useState<DetalleSprint | null>(null)
  const [backlog, setBacklog] = useState<Historia[]>([])
  const [goal, setGoal] = useState('')
  const [aAsignar, setAAsignar] = useState('') // id de la historia elegida para asignar
  const [error, setError] = useState('')

  useEffect(() => {
    Promise.all([pedir<DetalleSprint>('GET', ruta), pedir<Historia[]>('GET', `/proyectos/${proyecto.id}/backlog`)])
      .then(([detalle, historias]) => {
        setSprint(detalle)
        setGoal(detalle.sprint_goal)
        setBacklog(historias)
      })
      .catch((e) => setError(mensajeDe(e)))
  }, [ruta, proyecto.id])

  // hacer ejecuta un cambio sobre el Sprint: el backend responde el Sprint actualizado, que reemplaza
  // al que se muestra. Devuelve si salió bien.
  async function hacer(metodo: string, subruta: string, cuerpo?: unknown): Promise<boolean> {
    setError('')
    try {
      setSprint(await pedir<DetalleSprint>(metodo, ruta + subruta, cuerpo))
      return true
    } catch (e) {
      setError(mensajeDe(e))
      return false
    }
  }

  async function guardarGoal(evento: FormEvent) {
    evento.preventDefault()
    await hacer('PUT', '', { sprint_goal: goal })
  }

  async function asignar(evento: FormEvent) {
    evento.preventDefault()
    if (await hacer('POST', '/historias', { historia_id: Number(aAsignar) })) {
      setAAsignar('')
    }
  }

  function cerrar() {
    if (window.confirm(`¿Cerrar el Sprint ${sprint?.numero}? Un Sprint cerrado ya no admite cambios.`)) {
      hacer('POST', '/cerrar')
    }
  }

  if (sprint === null) {
    return (
      <section>
        <button className="enlace" onClick={alVolver}>
          ← Volver a los Sprints
        </button>
        <MensajeError mensaje={error} />
        {!error && <Cargando />}
      </section>
    )
  }

  const abierto = sprint.estado === 'Abierto'
  // Para asignar se ofrecen las historias del backlog que no están en este Sprint.
  const disponibles = backlog.filter((h) => !sprint.historias.some((s) => s.id === h.id))

  return (
    <>
      <section>
        <button className="enlace" onClick={alVolver}>
          ← Volver a los Sprints
        </button>
        <h2>
          Sprint {sprint.numero} <span className="suave">· {sprint.estado}</span>
        </h2>
        <p className="suave">
          {sprint.fecha_inicio} al {sprint.fecha_fin}
        </p>
        <MensajeError mensaje={error} />
        {abierto ? (
          <form onSubmit={guardarGoal}>
            <label>
              Sprint Goal
              <textarea value={goal} onChange={(e) => setGoal(e.target.value)} />
            </label>
            <div className="acciones">
              <button type="submit">Guardar Sprint Goal</button>
            </div>
          </form>
        ) : (
          <p>
            <strong>Sprint Goal:</strong> {sprint.sprint_goal}
          </p>
        )}
      </section>

      <section>
        <h2>Historias del Sprint</h2>
        {sprint.historias.length === 0 && <p className="suave">El Sprint todavía no tiene historias.</p>}
        <ul className="lista">
          {sprint.historias.map((h) => (
            <li key={h.id}>
              <div>
                <strong>
                  {h.prioridad}. {h.titulo}
                </strong>
                <span className="suave">
                  {' '}
                  · {h.story_points === null ? 'Sin estimar' : `${h.story_points} SP`}
                </span>
                <p className={h.completada_en_sprint ? undefined : 'suave'}>
                  {h.completada_en_sprint ? '✓ Completada en este Sprint' : 'No completada'}
                </p>
              </div>
              {abierto && (
                <div className="acciones">
                  {h.completada_en_sprint ? (
                    <button className="secundario" onClick={() => hacer('DELETE', `/historias/${h.id}/completar`)}>
                      Deshacer
                    </button>
                  ) : (
                    <button onClick={() => hacer('POST', `/historias/${h.id}/completar`)}>Completar</button>
                  )}
                  <button className="secundario" onClick={() => hacer('DELETE', `/historias/${h.id}`)}>
                    Quitar
                  </button>
                </div>
              )}
            </li>
          ))}
        </ul>

        {abierto && (
          <form onSubmit={asignar}>
            <label>
              Asignar una historia del backlog
              <select value={aAsignar} onChange={(e) => setAAsignar(e.target.value)}>
                <option value="">Elegí una historia</option>
                {disponibles.map((h) => (
                  <option key={h.id} value={h.id}>
                    {h.prioridad}. {h.titulo} ({h.estado})
                  </option>
                ))}
              </select>
            </label>
            <div className="acciones">
              <button type="submit" disabled={aAsignar === ''}>
                Asignar
              </button>
            </div>
          </form>
        )}
      </section>

      {abierto && (
        <section>
          <h2>Cerrar el Sprint</h2>
          <p className="suave">Lo planificado y lo completado quedan registrados; las historias no completadas quedan libres para otro Sprint.</p>
          <button onClick={cerrar}>Cerrar Sprint {sprint.numero}</button>
        </section>
      )}
    </>
  )
}
