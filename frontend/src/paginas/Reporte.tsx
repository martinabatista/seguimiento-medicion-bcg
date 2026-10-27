import { useEffect, useState } from 'react'
import {
  descargar,
  mensajeDe,
  pedir,
  type DefectosDeReporte,
  type HistoriaDeReporte,
  type Proyecto,
  type ReporteProyecto,
  type ReporteSprint,
  type Sprint,
  type SprintDeReporte,
} from '../api'
import { Cargando, MensajeError } from '../componentes'
import { decimal, desviacion, horas } from '../formato'

// Reporte muestra el reporte del proyecto completo o de un Sprint (spec 14) y descarga su PDF, que
// arma el backend. La pantalla no calcula nada: muestra lo que devuelve la consulta del reporte, con
// los mismos textos que el PDF para lo que falta ("sin estimar", "sin datos"…).
export default function Reporte({ proyecto }: { proyecto: Proyecto }) {
  const rutaProyecto = `/proyectos/${proyecto.id}`
  const [sprints, setSprints] = useState<Sprint[]>([])
  const [eleccion, setEleccion] = useState('') // '' = proyecto completo; si no, el id del Sprint
  const [reporte, setReporte] = useState<ReporteProyecto | ReporteSprint | null>(null)
  const [error, setError] = useState('')
  const [descargando, setDescargando] = useState(false)

  const rutaReporte = eleccion === '' ? `${rutaProyecto}/reporte` : `${rutaProyecto}/sprints/${eleccion}/reporte`

  useEffect(() => {
    pedir<Sprint[]>('GET', `${rutaProyecto}/sprints`)
      .then(setSprints)
      .catch((e) => setError(mensajeDe(e)))
  }, [rutaProyecto])

  useEffect(() => {
    setReporte(null)
    setError('')
    pedir<ReporteProyecto | ReporteSprint>('GET', rutaReporte)
      .then(setReporte)
      .catch((e) => setError(mensajeDe(e)))
  }, [rutaReporte])

  async function descargarPDF() {
    setError('')
    setDescargando(true)
    try {
      await descargar(`${rutaReporte}/pdf`)
    } catch (e) {
      setError(mensajeDe(e))
    } finally {
      setDescargando(false)
    }
  }

  return (
    <>
      <section>
        <h2>Reporte</h2>
        <div className="acciones">
          <label>
            Reporte de
            <select value={eleccion} onChange={(e) => setEleccion(e.target.value)}>
              <option value="">Proyecto completo</option>
              {sprints.map((s) => (
                <option key={s.id} value={s.id}>
                  Sprint {s.numero} ({s.fecha_inicio} al {s.fecha_fin}, {s.estado.toLowerCase()})
                </option>
              ))}
            </select>
          </label>
          <button onClick={descargarPDF} disabled={reporte === null || descargando}>
            {descargando ? 'Descargando…' : 'Descargar PDF'}
          </button>
        </div>
        <MensajeError mensaje={error} />
        {reporte === null && !error && <Cargando />}
      </section>

      {reporte && ('sprint' in reporte ? <VistaSprint reporte={reporte} /> : <VistaProyecto reporte={reporte} />)}
    </>
  )
}

function VistaProyecto({ reporte }: { reporte: ReporteProyecto }) {
  const { proyecto, metricas } = reporte
  const t = metricas.totales
  return (
    <>
      <section>
        <h2>Reporte del proyecto {proyecto.nombre}</h2>
        {proyecto.descripcion && <p>{proyecto.descripcion}</p>}
        <p className="suave">
          {proyecto.fecha_inicio} al {proyecto.fecha_fin} · {reporte.estado} · Generado el {reporte.generado_el}
        </p>
        <h3>Integrantes</h3>
        {reporte.integrantes.length === 0 ? (
          <p className="suave">sin integrantes</p>
        ) : (
          <p>{reporte.integrantes.map((i) => `${i.nombre} (${i.rol})`).join(' · ')}</p>
        )}
      </section>

      <section>
        <h2>Historias</h2>
        <TablaHistorias historias={reporte.historias} deSprint={false} />
      </section>

      <section>
        <h2>Sprints</h2>
        {reporte.sprints.length === 0 ? (
          <p className="suave">sin Sprints</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Sprint</th>
                <th>Historias completadas</th>
                <th>Story Points completados</th>
                <th>Horas estimadas / reales</th>
                <th>Desviación</th>
              </tr>
            </thead>
            <tbody>
              {reporte.sprints.map((s) => (
                <tr key={s.sprint_id}>
                  <td>
                    {nombreSprint(s)}
                    <br />
                    <span className="suave">{s.sprint_goal}</span>
                  </td>
                  <td>
                    {s.avance.historias_completadas} de {s.avance.historias_planificadas} ({s.avance.porcentaje_completadas} %)
                  </td>
                  <td>
                    {s.avance.story_points_completados} de {s.avance.story_points_planificados}
                  </td>
                  <td>
                    {horas(s.esfuerzo.horas_estimadas)} / {horas(s.esfuerzo.horas_reales)} h
                  </td>
                  <td>{desviacion(s.esfuerzo.desviacion, 'sin datos')}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>

      <section>
        <h2>Métricas del proyecto</h2>
        <ul>
          <li>
            Velocidad del equipo:{' '}
            {metricas.velocidad === null ? 'todavía no hay velocidad (no hay Sprints cerrados)' : `${decimal(metricas.velocidad)} Story Points por Sprint`}
          </li>
          <li>
            Sprints cerrados: {metricas.sprints_cerrados} · Historias: {t.historias_completadas} de {t.historias_planificadas} completadas (
            {t.porcentaje_completadas} %) · Story Points: {t.story_points_completados} de {t.story_points_planificados} completados
          </li>
          <li>
            Horas estimadas: {horas(metricas.horas_estimadas)} · Horas reales: {horas(metricas.horas_reales)}
          </li>
          <li>Desviación: {desviacion(metricas.desviacion, 'sin datos')}</li>
          <li>Historias completadas sin horas estimadas: {metricas.completadas_sin_horas_estimadas}</li>
        </ul>
      </section>

      <SeccionDefectos defectos={reporte.defectos} />
    </>
  )
}

function VistaSprint({ reporte }: { reporte: ReporteSprint }) {
  const s = reporte.sprint
  return (
    <>
      <section>
        <h2>Reporte del Sprint {s.numero}</h2>
        <p>
          Proyecto {reporte.proyecto.nombre} · Sprint Goal: {s.sprint_goal}
        </p>
        <p className="suave">
          {s.fecha_inicio} al {s.fecha_fin} · {estadoSprint(s)} · Generado el {reporte.generado_el}
        </p>
      </section>

      <section>
        <h2>Historias planificadas</h2>
        <TablaHistorias historias={reporte.historias} deSprint />
      </section>

      <section>
        <h2>Métricas del Sprint</h2>
        <ul>
          <li>
            Historias: {s.avance.historias_completadas} de {s.avance.historias_planificadas} completadas ({s.avance.porcentaje_completadas} %) ·
            Sin estimar: {s.avance.sin_estimar}
          </li>
          <li>
            Story Points: {s.avance.story_points_completados} de {s.avance.story_points_planificados} completados
          </li>
          <li>
            Horas estimadas: {horas(s.esfuerzo.horas_estimadas)} · Horas reales: {horas(s.esfuerzo.horas_reales)} · Desviación:{' '}
            {desviacion(s.esfuerzo.desviacion, 'sin datos')}
          </li>
          <li>Historias completadas sin horas estimadas: {s.esfuerzo.completadas_sin_horas_estimadas}</li>
        </ul>
      </section>

      <SeccionDefectos defectos={reporte.defectos} />
    </>
  )
}

// TablaHistorias muestra las estimaciones y el esfuerzo de cada historia (CA-2, CA-3). En el
// reporte de un Sprint, en vez del estado muestra si se completó en ese Sprint.
function TablaHistorias({ historias, deSprint }: { historias: HistoriaDeReporte[]; deSprint: boolean }) {
  if (historias.length === 0) {
    return <p className="suave">sin historias</p>
  }
  return (
    <table>
      <thead>
        <tr>
          <th>Historia</th>
          <th>{deSprint ? 'En el Sprint' : 'Estado'}</th>
          <th>Story Points</th>
          <th>Horas estimadas</th>
          <th>Horas reales</th>
          <th>Desviación</th>
        </tr>
      </thead>
      <tbody>
        {historias.map((h) => (
          <tr key={h.historia_id}>
            <td>{h.titulo}</td>
            <td>{deSprint ? (h.completada_en_sprint ? 'completada' : 'no completada') : h.estado}</td>
            <td>{h.story_points ?? 'sin estimar'}</td>
            <td>{h.horas_estimadas === null ? 'sin horas estimadas' : `${horas(h.horas_estimadas)} h`}</td>
            <td>{horas(h.horas_reales)} h</td>
            <td>{desviacion(h.desviacion, 'sin datos')}</td>
          </tr>
        ))}
      </tbody>
    </table>
  )
}

function SeccionDefectos({ defectos }: { defectos: DefectosDeReporte }) {
  return (
    <section>
      <h2>Defectos</h2>
      <p>
        Detectados: {defectos.detectados} · Resueltos: {defectos.resueltos}
      </p>
      {defectos.lista.length === 0 ? (
        <p className="suave">sin defectos</p>
      ) : (
        <ul className="lista">
          {defectos.lista.map((d) => (
            <li key={d.id}>
              <div>
                <strong>
                  {d.severidad} · {d.descripcion}
                </strong>
                <p className="suave">
                  {d.estado} · {d.historia_titulo} · detectado en el Sprint {d.sprint_deteccion.numero}
                  {d.sprint_resolucion ? ` · resuelto en el Sprint ${d.sprint_resolucion.numero}` : ' · sin Sprint de resolución'}
                </p>
              </div>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}

function nombreSprint(s: SprintDeReporte): string {
  return `Sprint ${s.numero} (${s.fecha_inicio} al ${s.fecha_fin}, ${estadoSprint(s)})`
}

// estadoSprint marca el Sprint abierto como parcial (spec 14 RN-6), como el PDF.
function estadoSprint(s: SprintDeReporte): string {
  return s.parcial ? `${s.estado}, parcial` : s.estado
}
