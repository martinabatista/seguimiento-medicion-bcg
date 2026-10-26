import { useEffect, useState } from 'react'
import { mensajeDe, pedir, type Historia, type Integrante, type Proyecto, type RondaRevelada, type Votacion } from '../api'
import { Cargando, MensajeError } from '../componentes'

// PlanningPoker: el equipo estima una historia votando en secreto (spec 08). Todas las reglas están
// en el backend: quién vota, que los votos estén ocultos, cuándo se revela y las diferencias.
export default function PlanningPoker({ proyecto }: { proyecto: Proyecto }) {
  const rutaProyecto = `/proyectos/${proyecto.id}`
  const [historias, setHistorias] = useState<Historia[] | null>(null)
  const [developers, setDevelopers] = useState<Integrante[]>([])
  const [escala, setEscala] = useState<number[]>([])
  const [abiertas, setAbiertas] = useState<Votacion[]>([])
  const [actual, setActual] = useState<Votacion | null>(null) // la votación que se está mirando
  const [historiaElegida, setHistoriaElegida] = useState('')
  const [cartas, setCartas] = useState<Record<number, string>>({}) // carta elegida por cada Developer, antes de votar
  const [acordada, setAcordada] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    Promise.all([
      pedir<Historia[]>('GET', `${rutaProyecto}/backlog`),
      pedir<Integrante[]>('GET', `${rutaProyecto}/integrantes`),
      pedir<{ valores: number[] }>('GET', '/story-points/escala'),
      pedir<Votacion[]>('GET', `${rutaProyecto}/votaciones`),
    ])
      .then(([listaHistorias, integrantes, respuestaEscala, votaciones]) => {
        setHistorias(listaHistorias)
        setDevelopers(integrantes.filter((i) => i.rol === 'Developer'))
        setEscala(respuestaEscala.valores)
        setAbiertas(votaciones)
      })
      .catch((e) => setError(mensajeDe(e)))
  }, [rutaProyecto])

  // hacer ejecuta un pedido sobre una votación, muestra cómo quedó y vuelve a pedir las votaciones
  // abiertas y el backlog (al acordar, cambian los Story Points de la historia).
  async function hacer(pedido: () => Promise<Votacion>) {
    setError('')
    try {
      setActual(await pedido())
      setAbiertas(await pedir<Votacion[]>('GET', `${rutaProyecto}/votaciones`))
      setHistorias(await pedir<Historia[]>('GET', `${rutaProyecto}/backlog`))
    } catch (e) {
      setError(mensajeDe(e))
    }
  }

  const rutaVotacion = `${rutaProyecto}/votaciones/${actual?.id}`
  const titulo = (historiaID: number) => historias?.find((h) => h.id === historiaID)?.titulo ?? `Historia ${historiaID}`

  function abrir() {
    hacer(() => pedir<Votacion>('POST', `${rutaProyecto}/votaciones`, { historia_id: Number(historiaElegida) }))
    setHistoriaElegida('')
  }

  function votar(developer: Integrante) {
    const carta = cartas[developer.id] ?? ''
    hacer(() =>
      pedir<Votacion>('POST', `${rutaVotacion}/votos`, {
        integrante_id: developer.id,
        story_points: carta === '' ? undefined : Number(carta), // vacía: el backend avisa que falta
      }),
    )
    setCartas({ ...cartas, [developer.id]: '' }) // la carta deja de verse: el voto es secreto
  }

  function acordar() {
    hacer(() => pedir<Votacion>('POST', `${rutaVotacion}/acordar`, { story_points: acordada === '' ? undefined : Number(acordada) }))
    setAcordada('')
  }

  if (historias === null && !error) {
    return <Cargando />
  }

  const reveladaActual = actual?.rondas.find((r) => r.numero === actual.ronda)
  const anteriores = actual?.rondas.filter((r) => r.numero !== actual.ronda) ?? []
  const sinVotacion = (historias ?? []).filter((h) => h.estado !== 'Completada' && !abiertas.some((v) => v.historia_id === h.id))

  return (
    <>
      <section>
        <h2>Planning Poker</h2>
        <MensajeError mensaje={error} />
        <div className="acciones">
          <label>
            Historia para estimar
            <select value={historiaElegida} onChange={(e) => setHistoriaElegida(e.target.value)}>
              <option value="">Elegí una historia</option>
              {sinVotacion.map((h) => (
                <option key={h.id} value={h.id}>
                  {h.prioridad}. {h.titulo}
                </option>
              ))}
            </select>
          </label>
          <button onClick={abrir} disabled={historiaElegida === ''}>
            Abrir votación
          </button>
        </div>

        {abiertas.length > 0 && (
          <>
            <h3>Votaciones abiertas</h3>
            <ul className="lista">
              {abiertas.map((v) => (
                <li key={v.id}>
                  <span>
                    <strong>{titulo(v.historia_id)}</strong> <span className="suave">· ronda {v.ronda}</span>
                  </span>
                  <button className="secundario" onClick={() => setActual(v)}>
                    Continuar
                  </button>
                </li>
              ))}
            </ul>
          </>
        )}
      </section>

      {actual && (
        <section>
          <h2>
            {titulo(actual.historia_id)} <span className="suave">· ronda {actual.ronda}</span>
          </h2>

          {actual.estado === 'Cerrada' ? (
            <p>
              Votación cerrada: se acordaron <strong>{actual.estimacion_acordada} Story Points</strong>.
            </p>
          ) : reveladaActual ? (
            <>
              <Ronda ronda={reveladaActual} />
              <div className="acciones">
                <button className="secundario" onClick={() => hacer(() => pedir<Votacion>('POST', `${rutaVotacion}/rondas`))}>
                  Nueva ronda
                </button>
                <label>
                  Estimación acordada
                  <select value={acordada} onChange={(e) => setAcordada(e.target.value)}>
                    <option value="">Elegí un valor</option>
                    {escala.map((valor) => (
                      <option key={valor}>{valor}</option>
                    ))}
                  </select>
                </label>
                <button onClick={acordar}>Acordar</button>
              </div>
            </>
          ) : (
            <>
              <ul className="lista">
                {developers.map((d) => (
                  <li key={d.id}>
                    <span>
                      <strong>{d.nombre}</strong>{' '}
                      {actual.votaron.some((p) => p.id === d.id) && <span className="suave">· Ya votó</span>}
                    </span>
                    <div className="acciones">
                      <select
                        aria-label={`Carta de ${d.nombre}`}
                        value={cartas[d.id] ?? ''}
                        onChange={(e) => setCartas({ ...cartas, [d.id]: e.target.value })}
                      >
                        <option value="">Carta</option>
                        {escala.map((valor) => (
                          <option key={valor}>{valor}</option>
                        ))}
                      </select>
                      <button onClick={() => votar(d)}>Votar</button>
                    </div>
                  </li>
                ))}
              </ul>
              <button onClick={() => hacer(() => pedir<Votacion>('POST', `${rutaVotacion}/revelar`))}>Revelar votos</button>
            </>
          )}

          {anteriores.length > 0 && (
            <>
              <h3>Rondas anteriores</h3>
              {anteriores.map((r) => (
                <Ronda key={r.numero} ronda={r} />
              ))}
            </>
          )}
        </section>
      )}
    </>
  )
}

// Ronda muestra una ronda revelada: los votos y, si los hay, las diferencias (spec 08 CA-4, CA-5).
function Ronda({ ronda }: { ronda: RondaRevelada }) {
  const nombres = (personas: { nombre: string }[]) => personas.map((p) => p.nombre).join(' y ')
  return (
    <div>
      <h4>Ronda {ronda.numero}</h4>
      <ul>
        {ronda.votos.map((v) => (
          <li key={v.integrante_id}>
            {v.nombre}: <strong>{v.story_points}</strong>
          </li>
        ))}
      </ul>
      {ronda.hay_diferencias ? (
        <p className="error">
          ¡Hay diferencias! El voto más bajo es {ronda.mas_bajo.story_points} ({nombres(ronda.mas_bajo.integrantes)}) y el más
          alto, {ronda.mas_alto.story_points} ({nombres(ronda.mas_alto.integrantes)}).
        </p>
      ) : (
        <p className="suave">Sin diferencias: todos votaron {ronda.mas_bajo.story_points}.</p>
      )}
    </div>
  )
}
