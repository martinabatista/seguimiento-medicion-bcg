import { useEffect, useState } from 'react'
import { mensajeDe, pedir, type Dashboard as DatosDashboard, type Proyecto } from '../api'
import { Cargando, MensajeError } from '../componentes'

// Dashboard: el estado del proyecto de un vistazo (spec 13). Los números llegan calculados por el
// backend; la pantalla solo los dibuja. Los gráficos son HTML y CSS, sin librería: el alto de cada
// barra es su valor sobre el mayor del gráfico.
export default function Dashboard({ proyecto }: { proyecto: Proyecto }) {
  const [datos, setDatos] = useState<DatosDashboard | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    pedir<DatosDashboard>('GET', `/proyectos/${proyecto.id}/dashboard`)
      .then(setDatos)
      .catch((e) => setError(mensajeDe(e)))
  }, [proyecto.id])

  if (error) {
    return <MensajeError mensaje={error} />
  }
  if (datos === null) {
    return <Cargando />
  }

  const { backlog, sprint_en_curso: sprint, graficos } = datos
  return (
    <>
      <section>
        <h2>Dashboard</h2>
        <div className="tarjetas">
          <Tarjeta numero={datos.estado} etiqueta={textoDeDias(datos)} />
          <Tarjeta
            numero={sprint ? `Sprint ${sprint.numero}` : 'Sin Sprint'}
            etiqueta={
              sprint
                ? `en curso · le quedan ${sprint.dias_restantes} días · ${sprint.story_points_completados} de ${sprint.story_points_planificados} SP`
                : 'no hay un Sprint en curso'
            }
          />
          <Tarjeta numero={`${backlog.completadas} de ${backlog.total}`} etiqueta={`historias completadas · ${backlog.porcentaje_completadas} %`} />
          <Tarjeta numero={`${backlog.story_points_completados} de ${backlog.story_points_totales}`} etiqueta="Story Points completados" />
          <Tarjeta
            numero={datos.velocidad === null ? '—' : decimal(datos.velocidad)}
            etiqueta={datos.velocidad === null ? 'velocidad: todavía no hay Sprints cerrados' : 'velocidad (SP por Sprint)'}
          />
          <Tarjeta numero={String(backlog.sin_estimar)} etiqueta="historias sin estimar" />
        </div>
      </section>

      <section>
        <h2>Story Points por Sprint</h2>
        <Leyenda items={[{ nombre: 'Planificados', clase: 'serie-1' }, { nombre: 'Completados', clase: 'serie-2' }]} />
        <Barras
          grupos={graficos.story_points_por_sprint.map((s) => ({
            nombre: s.parcial ? `${s.sprint} (parcial)` : s.sprint,
            barras: [
              { valor: s.planificados, clase: 'serie-1', detalle: `${s.sprint}: ${s.planificados} planificados` },
              { valor: s.completados, clase: 'serie-2', detalle: `${s.sprint}: ${s.completados} completados` },
            ],
          }))}
        />
      </section>

      <section>
        <h2>Velocidad del equipo</h2>
        <p className="suave">Story Points completados en cada Sprint cerrado y su promedio.</p>
        <Barras
          grupos={graficos.velocidad.sprints.map((s) => ({
            nombre: s.sprint,
            barras: [{ valor: s.completados, clase: 'serie-1', detalle: `${s.sprint}: ${s.completados} completados` }],
          }))}
          promedio={graficos.velocidad.promedio}
        />
      </section>

      <section>
        <h2>Backlog por estado</h2>
        <Dona porciones={graficos.backlog_por_estado} />
      </section>
    </>
  )
}

// ---------- piezas ----------

function Tarjeta({ numero, etiqueta }: { numero: string; etiqueta: string }) {
  return (
    <div className="tarjeta">
      <strong>{numero}</strong>
      <span className="suave">{etiqueta}</span>
    </div>
  )
}

function Leyenda({ items }: { items: { nombre: string; clase: string }[] }) {
  return (
    <p className="leyenda">
      {items.map((i) => (
        <span key={i.nombre}>
          <span className={`muestra ${i.clase}`} /> {i.nombre}
        </span>
      ))}
    </p>
  )
}

type Grupo = { nombre: string; barras: { valor: number; clase: string; detalle: string }[] }

// Barras dibuja un grupo de barras por Sprint, cada una con su número arriba (se lee sin colores).
// Si llega promedio, lo marca con una línea punteada a su altura.
function Barras({ grupos, promedio = null }: { grupos: Grupo[]; promedio?: number | null }) {
  if (grupos.length === 0) {
    return <p className="suave">Todavía no hay Sprints.</p>
  }
  const mayor = Math.max(1, promedio ?? 0, ...grupos.flatMap((g) => g.barras.map((b) => b.valor)))
  const alto = (valor: number) => `${(valor / mayor) * 100}%`
  return (
    <div className="grafico">
      <div className="area">
        {promedio !== null && (
          <div className="promedio" style={{ bottom: alto(promedio) }} title={`Promedio: ${decimal(promedio)}`}>
            <span>Promedio: {decimal(promedio)}</span>
          </div>
        )}
        {grupos.map((g) => (
          <div key={g.nombre} className="grupo">
            {g.barras.map((b) => (
              <div key={b.detalle} className="columna" title={b.detalle}>
                <span className="valor">{b.valor}</span>
                <div className={`barra ${b.clase}`} style={{ height: alto(b.valor) }} />
              </div>
            ))}
          </div>
        ))}
      </div>
      <div className="ejes">
        {grupos.map((g) => (
          <span key={g.nombre}>{g.nombre}</span>
        ))}
      </div>
    </div>
  )
}

const clasesPorEstado: Record<string, string> = { Pendiente: 'serie-1', 'En progreso': 'serie-2', Completada: 'serie-3' }

// Dona dibuja las historias por estado con un degradado circular de CSS; al lado, la leyenda con
// el nombre y la cantidad de cada estado (se lee sin colores).
function Dona({ porciones }: { porciones: { estado: string; cantidad: number }[] }) {
  const total = porciones.reduce((suma, p) => suma + p.cantidad, 0)
  let acumulado = 0
  const tramos = porciones.map((p) => {
    const desde = (acumulado / total) * 100
    acumulado += p.cantidad
    return `var(--${clasesPorEstado[p.estado]}) ${desde}% ${(acumulado / total) * 100}%`
  })
  const fondo = total === 0 ? 'var(--borde)' : `conic-gradient(${tramos.join(', ')})`
  const detalle = porciones.map((p) => `${p.estado}: ${p.cantidad}`).join(' · ')
  return (
    <div className="dona-y-leyenda">
      <div className="dona" style={{ background: fondo }} title={detalle} role="img" aria-label={`Backlog por estado: ${detalle}`}>
        <span>
          <strong>{total}</strong>
          <br />
          historias
        </span>
      </div>
      <ul className="leyenda-dona">
        {porciones.map((p) => (
          <li key={p.estado}>
            <span className={`muestra ${clasesPorEstado[p.estado]}`} /> {p.estado}: <strong>{p.cantidad}</strong>
          </li>
        ))}
      </ul>
    </div>
  )
}

// ---------- textos ----------

function textoDeDias(d: DatosDashboard): string {
  if (d.dias_para_empezar !== null) {
    return `faltan ${d.dias_para_empezar} días para empezar`
  }
  if (d.dias_para_terminar !== null) {
    return d.dias_para_terminar === 0 ? 'termina hoy' : `faltan ${d.dias_para_terminar} días para terminar`
  }
  return 'el proyecto ya terminó'
}

// decimal muestra 6.5 como "6,5" (con coma, como se escribe en castellano).
function decimal(n: number): string {
  return n.toLocaleString('es-AR', { maximumFractionDigits: 1 })
}
