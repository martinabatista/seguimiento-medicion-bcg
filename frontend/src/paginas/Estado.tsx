import { useEffect, useState } from 'react'
import { mensajeDe, pedir, type EstadoProyecto, type Proyecto } from '../api'
import { Cargando, MensajeError } from '../componentes'

// Estado muestra la situación actual del proyecto (spec 03): sus datos con el estado calculado por
// fechas, el Sprint en curso (spec 05 CA-5) y el resumen del avance del backlog. Todo llega calculado
// por el backend en el momento de la consulta; la pantalla solo lo presenta.
export default function Estado({ proyecto }: { proyecto: Proyecto }) {
  const [datos, setDatos] = useState<EstadoProyecto | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    pedir<EstadoProyecto>('GET', `/proyectos/${proyecto.id}/estado`)
      .then(setDatos)
      .catch((e) => setError(mensajeDe(e)))
  }, [proyecto.id])

  if (error) {
    return <MensajeError mensaje={error} />
  }
  if (datos === null) {
    return <Cargando />
  }

  const { proyecto: p, backlog, sprint_en_curso: sprint } = datos
  return (
    <>
      <section>
        <h2>Estado del proyecto</h2>
        <div className="tarjetas">
          <div className="tarjeta">
            <span className="suave">Estado</span>
            <strong>{datos.estado}</strong>
            <span className="suave">
              del {p.fecha_inicio} al {p.fecha_fin}
            </span>
          </div>
          <div className="tarjeta">
            <span className="suave">Sprint en curso</span>
            {sprint ? (
              <>
                <strong>Sprint {sprint.numero}</strong>
                <span className="suave">
                  del {sprint.fecha_inicio} al {sprint.fecha_fin}
                </span>
              </>
            ) : (
              <>
                <strong>—</strong>
                <span className="suave">No hay un Sprint en curso</span>
              </>
            )}
          </div>
        </div>
        <p className="suave">{p.descripcion}</p>
        {sprint && (
          <p>
            <strong>Sprint Goal:</strong> {sprint.sprint_goal}
          </p>
        )}
      </section>

      <section>
        <h2>Avance del backlog</h2>
        {backlog.total === 0 ? (
          <p className="suave">El backlog todavía no tiene historias.</p>
        ) : (
          <table>
            <tbody>
              <tr>
                <th>Historias</th>
                <td>{backlog.total}</td>
              </tr>
              <tr>
                <th>Pendientes</th>
                <td>{backlog.pendientes}</td>
              </tr>
              <tr>
                <th>En progreso</th>
                <td>{backlog.en_progreso}</td>
              </tr>
              <tr>
                <th>Completadas</th>
                <td>
                  {backlog.completadas} ({backlog.porcentaje_completadas} %)
                </td>
              </tr>
              <tr>
                <th>Story Points</th>
                <td>
                  {backlog.story_points_completados} completados de {backlog.story_points_totales}
                </td>
              </tr>
              <tr>
                <th>Sin estimar</th>
                <td>{backlog.sin_estimar}</td>
              </tr>
            </tbody>
          </table>
        )}
      </section>
    </>
  )
}
