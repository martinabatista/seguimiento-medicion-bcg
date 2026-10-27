import { useEffect, useState } from 'react'
import { mensajeDe, pedir, type Metricas as TipoMetricas, type MetricasEsfuerzo, type Proyecto } from '../api'
import { Cargando, MensajeError } from '../componentes'
import { decimal, desviacion, horas } from '../formato'

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

      <Esfuerzo proyecto={proyecto} />
    </>
  )
}

// Esfuerzo muestra las métricas de esfuerzo (spec 12) y permite cargar las horas estimadas de cada
// historia. Las validaciones y los cálculos (horas reales, desviación) son del backend.
function Esfuerzo({ proyecto }: { proyecto: Proyecto }) {
  const ruta = `/proyectos/${proyecto.id}/metricas/esfuerzo`
  const [esfuerzo, setEsfuerzo] = useState<MetricasEsfuerzo | null>(null)
  const [cargadas, setCargadas] = useState<Record<number, string>>({}) // lo que se escribió, por historia
  const [error, setError] = useState('')

  useEffect(() => {
    pedir<MetricasEsfuerzo>('GET', ruta)
      .then(setEsfuerzo)
      .catch((e) => setError(mensajeDe(e)))
  }, [ruta])

  // cambiarHoras guarda (o quita, con horas null) las horas estimadas y vuelve a pedir las métricas.
  async function cambiarHoras(historiaID: number, horas: string | null) {
    setError('')
    const rutaHoras = `/proyectos/${proyecto.id}/backlog/${historiaID}/horas-estimadas`
    try {
      if (horas === null) {
        await pedir('DELETE', rutaHoras)
      } else {
        // Vacías no se mandan, para que el backend informe "Las horas estimadas son obligatorias".
        await pedir('PUT', rutaHoras, horas === '' ? {} : { horas_estimadas: Number(horas) })
      }
      setCargadas((anteriores) => {
        const resto = { ...anteriores }
        delete resto[historiaID] // ya se guardó: se vuelve a mostrar lo que dice el backend
        return resto
      })
      setEsfuerzo(await pedir<MetricasEsfuerzo>('GET', ruta))
    } catch (e) {
      setError(mensajeDe(e))
    }
  }

  if (esfuerzo === null) {
    return (
      <section>
        <h2>Métricas de esfuerzo</h2>
        <MensajeError mensaje={error} />
        {!error && <Cargando />}
      </section>
    )
  }

  return (
    <>
      <section>
        <h2>Métricas de esfuerzo del proyecto</h2>
        <p className="suave">
          La desviación compara las horas reales con las estimadas, solo en las historias completadas que tienen
          horas estimadas. Positiva: se trabajó más de lo estimado.
        </p>
        <div className="tarjetas">
          <Tarjeta numero={`${horas(esfuerzo.horas_estimadas)} h`} etiqueta="horas estimadas del backlog" />
          <Tarjeta numero={`${horas(esfuerzo.horas_reales)} h`} etiqueta="horas reales registradas" />
          <Tarjeta numero={desviacion(esfuerzo.desviacion)} etiqueta="desviación de las completadas" />
          <Tarjeta numero={String(esfuerzo.completadas_sin_horas_estimadas)} etiqueta="completadas sin horas estimadas" />
        </div>
      </section>

      <section>
        <h2>Esfuerzo por historia</h2>
        <MensajeError mensaje={error} />
        {esfuerzo.historias.length === 0 ? (
          <p className="suave">El backlog está vacío.</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Historia</th>
                <th>Horas estimadas</th>
                <th>Horas reales</th>
                <th>Desviación</th>
              </tr>
            </thead>
            <tbody>
              {esfuerzo.historias.map((h) => (
                <tr key={h.historia_id}>
                  <td>
                    {h.titulo}
                    <span className="suave"> · {h.estado}</span>
                  </td>
                  <td>
                    {h.estado === 'Completada' ? (
                      // Una historia completada no admite cambiar su estimación (lo controla el backend).
                      h.horas_estimadas === null ? (
                        <span className="suave">Sin horas estimadas</span>
                      ) : (
                        `${horas(h.horas_estimadas)} h`
                      )
                    ) : (
                      <div className="acciones">
                        <input
                          type="number"
                          step="0.01"
                          aria-label={`Horas estimadas de ${h.titulo}`}
                          placeholder="Sin estimar"
                          value={cargadas[h.historia_id] ?? h.horas_estimadas ?? ''}
                          onChange={(e) => setCargadas({ ...cargadas, [h.historia_id]: e.target.value })}
                        />
                        <button onClick={() => cambiarHoras(h.historia_id, cargadas[h.historia_id] ?? String(h.horas_estimadas ?? ''))}>
                          Guardar
                        </button>
                        {h.horas_estimadas !== null && (
                          <button className="secundario" onClick={() => cambiarHoras(h.historia_id, null)}>
                            Quitar
                          </button>
                        )}
                      </div>
                    )}
                  </td>
                  <td>{horas(h.horas_reales)} h</td>
                  <td>{h.desviacion === null ? <span className="suave">—</span> : desviacion(h.desviacion)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>

      <section>
        <h2>Esfuerzo por Sprint</h2>
        {esfuerzo.sprints.length === 0 ? (
          <p className="suave">Todavía no hay Sprints.</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>Sprint</th>
                <th>Horas estimadas</th>
                <th>Horas reales</th>
                <th>Desviación</th>
                <th>Completadas sin horas estimadas</th>
              </tr>
            </thead>
            <tbody>
              {esfuerzo.sprints.map((s) => (
                <tr key={s.sprint_id}>
                  <td>
                    Sprint {s.numero}
                    <span className="suave"> · {s.parcial ? 'abierto (parcial)' : 'cerrado'}</span>
                  </td>
                  <td>{horas(s.horas_estimadas)} h</td>
                  <td>{horas(s.horas_reales)} h</td>
                  <td>{desviacion(s.desviacion)}</td>
                  <td>{s.completadas_sin_horas_estimadas}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
        <p className="suave">Cada Sprint cuenta las historias que se completaron en él.</p>
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
