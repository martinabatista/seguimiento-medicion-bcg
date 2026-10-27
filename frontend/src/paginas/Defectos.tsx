import { useCallback, useEffect, useState, type FormEvent } from 'react'
import { mensajeDe, pedir, type Defecto, type Historia, type MetricasDefectos, type Proyecto, type Sprint } from '../api'
import { Cargando, MensajeError } from '../componentes'

const severidades: Defecto['severidad'][] = ['Baja', 'Media', 'Alta', 'Crítica']
const estados: Defecto['estado'][] = ['Abierto', 'En corrección', 'Resuelto']
const defectoVacio = { descripcion: '', severidad: '', historia_id: '', sprint_deteccion_id: '' }

// Defectos registra los defectos encontrados y sigue su estado hasta resolverlos (spec 11): arriba
// las métricas de detectados y resueltos (del proyecto y por Sprint), después la lista con el cambio
// de estado y al final el formulario. Las validaciones (Sprint de resolución obligatorio y no
// anterior al de detección, etc.) y el conteo son del backend: la pantalla muestra sus mensajes.
export default function Defectos({ proyecto }: { proyecto: Proyecto }) {
  const rutaProyecto = `/proyectos/${proyecto.id}`
  const [defectos, setDefectos] = useState<Defecto[] | null>(null)
  const [metricas, setMetricas] = useState<MetricasDefectos | null>(null)
  const [historias, setHistorias] = useState<Historia[]>([])
  const [sprints, setSprints] = useState<Sprint[]>([])
  const [nuevo, setNuevo] = useState(defectoVacio)
  const [error, setError] = useState('')

  // cargarDefectos pide la lista y las métricas juntas: cambian con cada registro o cambio de estado.
  const cargarDefectos = useCallback(async () => {
    const [lista, conteo] = await Promise.all([
      pedir<Defecto[]>('GET', `${rutaProyecto}/defectos`),
      pedir<MetricasDefectos>('GET', `${rutaProyecto}/metricas/defectos`),
    ])
    setDefectos(lista)
    setMetricas(conteo)
  }, [rutaProyecto])

  useEffect(() => {
    Promise.all([pedir<Historia[]>('GET', `${rutaProyecto}/backlog`), pedir<Sprint[]>('GET', `${rutaProyecto}/sprints`), cargarDefectos()])
      .then(([listaHistorias, listaSprints]) => {
        setHistorias(listaHistorias)
        setSprints(listaSprints)
      })
      .catch((e) => setError(mensajeDe(e)))
  }, [rutaProyecto, cargarDefectos])

  // hacer ejecuta un pedido que cambia los defectos y, si sale bien, vuelve a pedir la lista y las
  // métricas (así los números siempre son los del backend). Devuelve si salió bien.
  async function hacer(pedido: () => Promise<unknown>): Promise<boolean> {
    setError('')
    try {
      await pedido()
      await cargarDefectos()
      return true
    } catch (e) {
      setError(mensajeDe(e))
      return false
    }
  }

  async function registrar(evento: FormEvent) {
    evento.preventDefault()
    // Los identificadores vacíos no se mandan, para que el backend informe cuál es obligatorio.
    const datos = {
      ...nuevo,
      historia_id: nuevo.historia_id === '' ? undefined : Number(nuevo.historia_id),
      sprint_deteccion_id: nuevo.sprint_deteccion_id === '' ? undefined : Number(nuevo.sprint_deteccion_id),
    }
    if (await hacer(() => pedir('POST', `${rutaProyecto}/defectos`, datos))) {
      setNuevo(defectoVacio)
    }
  }

  function cambiarEstado(id: number, estado: string, sprintResolucionID: string) {
    hacer(() =>
      pedir('PUT', `${rutaProyecto}/defectos/${id}/estado`, {
        estado,
        sprint_resolucion_id: sprintResolucionID === '' ? undefined : Number(sprintResolucionID),
      }),
    )
  }

  if (defectos === null || metricas === null) {
    return error ? <MensajeError mensaje={error} /> : <Cargando />
  }

  return (
    <>
      <section>
        <h2>Métricas de defectos</h2>
        <div className="tarjetas">
          <div className="tarjeta">
            <strong>{metricas.detectados}</strong>
            <span className="suave">defectos detectados</span>
          </div>
          <div className="tarjeta">
            <strong>{metricas.resueltos}</strong>
            <span className="suave">defectos resueltos</span>
          </div>
        </div>
        {metricas.sprints.length === 0 ? (
          <p className="suave">Todavía no hay Sprints: se crean en la sección Sprints.</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Sprint</th>
                <th>Detectados</th>
                <th>Resueltos</th>
              </tr>
            </thead>
            <tbody>
              {metricas.sprints.map((s) => (
                <tr key={s.sprint_id}>
                  <td>Sprint {s.numero}</td>
                  <td>{s.detectados}</td>
                  <td>{s.resueltos}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>

      <section>
        <h2>Defectos</h2>
        <MensajeError mensaje={error} />
        {defectos.length === 0 && <p className="suave">Todavía no hay defectos registrados.</p>}
        <ul className="lista">
          {defectos.map((d) => (
            <FilaDefecto key={`${d.id}-${d.estado}-${d.sprint_resolucion?.id}`} defecto={d} sprints={sprints} alGuardar={cambiarEstado} />
          ))}
        </ul>
      </section>

      <section>
        <h2>Registrar defecto</h2>
        {sprints.length === 0 && <p className="suave">El proyecto todavía no tiene Sprints: el defecto se registra en el Sprint en que se detectó.</p>}
        <form onSubmit={registrar}>
          <label>
            Descripción
            <textarea value={nuevo.descripcion} onChange={(e) => setNuevo({ ...nuevo, descripcion: e.target.value })} />
          </label>
          <label>
            Severidad
            <select value={nuevo.severidad} onChange={(e) => setNuevo({ ...nuevo, severidad: e.target.value })}>
              <option value="">Elegí la severidad</option>
              {severidades.map((s) => (
                <option key={s}>{s}</option>
              ))}
            </select>
          </label>
          <label>
            Historia relacionada
            <select value={nuevo.historia_id} onChange={(e) => setNuevo({ ...nuevo, historia_id: e.target.value })}>
              <option value="">Elegí una historia</option>
              {historias.map((h) => (
                <option key={h.id} value={h.id}>
                  {h.prioridad}. {h.titulo} ({h.estado})
                </option>
              ))}
            </select>
          </label>
          <label>
            Sprint de detección
            <SelectorSprint sprints={sprints} valor={nuevo.sprint_deteccion_id} alCambiar={(v) => setNuevo({ ...nuevo, sprint_deteccion_id: v })} />
          </label>
          <button type="submit">Registrar</button>
        </form>
      </section>
    </>
  )
}

// FilaDefecto muestra un defecto y deja elegir el estado nuevo; el Sprint de resolución se elige
// solo al pasar a Resuelto. Arranca con lo que tiene guardado (la key cambia cuando el backend lo cambia).
function FilaDefecto({
  defecto,
  sprints,
  alGuardar,
}: {
  defecto: Defecto
  sprints: Sprint[]
  alGuardar: (id: number, estado: string, sprintResolucionID: string) => void
}) {
  const [estado, setEstado] = useState<string>(defecto.estado)
  const [sprintResolucion, setSprintResolucion] = useState(defecto.sprint_resolucion ? String(defecto.sprint_resolucion.id) : '')

  return (
    <li>
      <div>
        <strong>
          {defecto.severidad} · {defecto.descripcion}
        </strong>
        <p className="suave">
          {defecto.historia_titulo} · detectado en el Sprint {defecto.sprint_deteccion.numero}
          {defecto.sprint_resolucion && ` · resuelto en el Sprint ${defecto.sprint_resolucion.numero}`}
        </p>
        <p className="suave">Estado: {defecto.estado}</p>
      </div>
      <div className="acciones">
        <select aria-label="Estado nuevo" value={estado} onChange={(e) => setEstado(e.target.value)}>
          {estados.map((e) => (
            <option key={e}>{e}</option>
          ))}
        </select>
        {estado === 'Resuelto' && (
          <SelectorSprint sprints={sprints} valor={sprintResolucion} alCambiar={setSprintResolucion} etiqueta="Sprint de resolución" />
        )}
        <button className="secundario" onClick={() => alGuardar(defecto.id, estado, estado === 'Resuelto' ? sprintResolucion : '')}>
          Guardar
        </button>
      </div>
    </li>
  )
}

function SelectorSprint({
  sprints,
  valor,
  alCambiar,
  etiqueta,
}: {
  sprints: Sprint[]
  valor: string
  alCambiar: (valor: string) => void
  etiqueta?: string
}) {
  return (
    <select aria-label={etiqueta} value={valor} onChange={(e) => alCambiar(e.target.value)}>
      <option value="">Elegí un Sprint</option>
      {sprints.map((s) => (
        <option key={s.id} value={s.id}>
          Sprint {s.numero} ({s.fecha_inicio} al {s.fecha_fin}, {s.estado.toLowerCase()})
        </option>
      ))}
    </select>
  )
}
