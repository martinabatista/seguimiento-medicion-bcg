import { afterEach, describe, expect, it, vi } from 'vitest'
import { ErrorAPI, pedir } from './api'

// Reemplaza fetch del navegador por uno de prueba que devuelve la respuesta indicada.
function responder(respuesta: Response | Error) {
  const fetchDePrueba = vi.fn(() => (respuesta instanceof Error ? Promise.reject(respuesta) : Promise.resolve(respuesta)))
  vi.stubGlobal('fetch', fetchDePrueba)
  return fetchDePrueba
}

function json(cuerpo: unknown, estado: number) {
  return new Response(JSON.stringify(cuerpo), { status: estado, headers: { 'Content-Type': 'application/json' } })
}

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('pedir', () => {
  it('devuelve los datos de una respuesta correcta y manda el cuerpo como JSON', async () => {
    const fetchDePrueba = responder(json({ id: 1, nombre: 'Tienda Online' }, 201))

    const datos = await pedir<{ id: number; nombre: string }>('POST', '/proyectos', { nombre: 'Tienda Online' })

    expect(datos).toEqual({ id: 1, nombre: 'Tienda Online' })
    expect(fetchDePrueba).toHaveBeenCalledWith('/api/proyectos', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ nombre: 'Tienda Online' }),
    })
  })

  it('usa el mensaje de error del backend tal cual, con su código', async () => {
    responder(json({ error: 'El nombre es obligatorio' }, 400))

    const error = await pedir('POST', '/proyectos', {}).catch((e: unknown) => e)

    expect(error).toBeInstanceOf(ErrorAPI)
    expect((error as ErrorAPI).message).toBe('El nombre es obligatorio')
    expect((error as ErrorAPI).estado).toBe(400)
  })

  it('acepta una respuesta 204 sin contenido', async () => {
    responder(new Response(null, { status: 204 }))

    await expect(pedir('DELETE', '/proyectos/1/backlog/1')).resolves.toBeUndefined()
  })

  it('informa un error claro si el servidor no responde', async () => {
    responder(new TypeError('Failed to fetch'))

    await expect(pedir('GET', '/proyectos')).rejects.toThrow('No se pudo conectar con el servidor')
  })

  it('informa el código si el servidor responde un error sin mensaje', async () => {
    responder(new Response('<html>Bad Gateway</html>', { status: 502 }))

    await expect(pedir('GET', '/proyectos')).rejects.toThrow('Error del servidor (502)')
  })
})
