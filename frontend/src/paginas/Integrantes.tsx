import { useEffect, useState, type FormEvent } from 'react'
import { mensajeDe, pedir, type Integrante, type Proyecto } from '../api'
import { Cargando, MensajeError } from '../componentes'

const roles: Integrante['rol'][] = ['Product Owner', 'Scrum Master', 'Developer']
const integranteVacio = { nombre: '', email: '', rol: '' }

// Integrantes lista los integrantes del proyecto en el orden en que se agregaron y permite agregar
// uno nuevo (spec 02). Las reglas están en el backend: el email llega normalizado y los errores
// (dato inválido o email repetido) se muestran con el mensaje de la API.
export default function Integrantes({ proyecto }: { proyecto: Proyecto }) {
  const ruta = `/proyectos/${proyecto.id}/integrantes`
  const [integrantes, setIntegrantes] = useState<Integrante[] | null>(null)
  const [nuevo, setNuevo] = useState(integranteVacio)
  const [error, setError] = useState('')

  useEffect(() => {
    pedir<Integrante[]>('GET', ruta)
      .then(setIntegrantes)
      .catch((e) => setError(mensajeDe(e)))
  }, [ruta])

  async function agregar(evento: FormEvent) {
    evento.preventDefault()
    setError('')
    try {
      const agregado = await pedir<Integrante>('POST', ruta, nuevo)
      setIntegrantes([...(integrantes ?? []), agregado])
      setNuevo(integranteVacio)
    } catch (e) {
      setError(mensajeDe(e)) // por ejemplo: "El integrante ya está registrado en el proyecto"
    }
  }

  return (
    <>
      <section>
        <h2>Integrantes</h2>
        <MensajeError mensaje={error} />
        {integrantes === null && !error && <Cargando />}
        {integrantes?.length === 0 && <p className="suave">El proyecto todavía no tiene integrantes: agregá el primero.</p>}
        <ul className="lista">
          {integrantes?.map((i) => (
            <li key={i.id}>
              <div>
                <strong>{i.nombre}</strong>
                <span className="suave"> · {i.rol}</span>
                <p className="suave">{i.email}</p>
              </div>
            </li>
          ))}
        </ul>
      </section>

      <section>
        <h2>Nuevo integrante</h2>
        <form onSubmit={agregar}>
          <label>
            Nombre
            <input value={nuevo.nombre} onChange={(e) => setNuevo({ ...nuevo, nombre: e.target.value })} />
          </label>
          <label>
            Email
            {/* Sin type="email": así el formato lo valida el backend y se ve su mensaje, no el del navegador. */}
            <input inputMode="email" value={nuevo.email} onChange={(e) => setNuevo({ ...nuevo, email: e.target.value })} />
          </label>
          <label>
            Rol
            <select value={nuevo.rol} onChange={(e) => setNuevo({ ...nuevo, rol: e.target.value })}>
              <option value="">Elegí un rol</option>
              {roles.map((rol) => (
                <option key={rol}>{rol}</option>
              ))}
            </select>
          </label>
          <button type="submit">Agregar integrante</button>
        </form>
      </section>
    </>
  )
}
