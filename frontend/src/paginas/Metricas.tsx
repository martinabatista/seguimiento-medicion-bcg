import { useEffect, useState } from 'react'
import { mensajeDe, pedir, type Metricas as TipoMetricas, type Proyecto } from '../api'
import { Cargando, MensajeError } from '../componentes'

// Metricas muestra las métricas de avance del proyecto y de cada Sprint (spec 10). Pide la consulta
// una sola vez y muestra lo que llega: todos los cálculos (porcentajes, velocidad, totales) son del backend.
export default function Metricas({ proyecto }: { proyecto: Proyecto }) {
  const [metricas, setMetricas] = useState<TipoMetricas | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    pedir<TipoMetricas>('GET', `/proyectos/${proyecto.id}/metricas`)
      .then(setMetricas)
      .catch((e) => setError(mensajeDe(e)))
  }, [proyecto.id])

  if (metricas === null) {
    return (
      <section>
        <h2>Métricas</h2>
        <MensajeError mensaje={error} />
        {!error && <Cargando />}
      </section>
    )
  }

  const { totales } = metricas
  return (
    <>
      <section>
        <h2>Métricas del proyecto</h2>
        <p className="suave">Cuentan solo los Sprints cerrados: {metricas.sprints_cerrados}.</p>
        <div className="tarjetas">
          <Tarjeta
            numero={metricas.velocidad === null ? '—' : decimal(metricas.velocidad)}
            etiqueta={metricas.velocidad === null ? 'velocidad: todavía no hay Sprints cerrados' : 'velocidad (SP por Sprint)'}
          />
          <Tarjeta
            numero={`${totales.story_points_completados} de ${totales.story_points_planificados}`}
            etiqueta="Story Points completados"
          />
          <Tarjeta
            numero={`${totales.historias_completadas} de ${totales.historias_planificadas}`}
            etiqueta={`historias completadas (${totales.porcentaje_completadas} %)`}
          />
        </div>
      </section>

      <section>
        <h2>Métricas por Sprint</h2>
        {metricas.sprints.length === 0 ? (
          <p className="suave">Todavía no hay Sprints: se crean en la sección Sprints.</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Sprint</th>
                <th>Story Points completados</th>
                <th>Historias completadas</th>
                <th>Sin estimar</th>
              </tr>
            </thead>
            <tbody>
              {metricas.sprints.map((s) => (
                <tr key={s.sprint_id}>
                  <td>
                    Sprint {s.numero}
                    <span className="suave"> · {s.parcial ? 'abierto (parcial)' : 'cerrado'}</span>
                  </td>
                  <td>
                    {s.story_points_completados} de {s.story_points_planificados}
                  </td>
                  <td>
                    {s.historias_completadas} de {s.historias_planificadas} ({s.porcentaje_completadas} %)
                  </td>
                  <td>{s.sin_estimar}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        <p className="suave">
          Un Sprint abierto es parcial: sus números cambian hasta que se cierra. Un Sprint cerrado usa los Story
          Points que tenían sus historias al cerrarlo.
        </p>
      </section>
    </>
  )
}

function Tarjeta({ numero, etiqueta }: { numero: string; etiqueta: string }) {
  return (
    <div className="tarjeta">
      <strong>{numero}</strong>
      <span className="suave">{etiqueta}</span>
    </div>
  )
}

// decimal muestra la velocidad con coma y un decimal, como llega del backend (6.5 → "6,5").
function decimal(n: number): string {
  return n.toLocaleString('es-AR', { maximumFractionDigits: 1 })
}
