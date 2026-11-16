import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'
import Backlog from './Backlog'

// El formulario se dibuja como HTML (sin navegador): alcanza para revisar los atributos de los campos.
// Los pedidos a la API no se hacen porque useEffect no corre al dibujar así.
const html = renderToStaticMarkup(
  <Backlog proyecto={{ id: 1, nombre: 'Tienda Online', descripcion: 'd', fecha_inicio: '2026-11-01', fecha_fin: '2026-12-15' }} />,
)

describe('formulario de historia (spec 04 RN-4: prioridad entera mayor o igual a 1)', () => {
  it('las flechas del campo de prioridad no bajan de 1', () => {
    const campo = html.match(/Prioridad \(1 es la más importante\)<input[^>]*>/)?.[0] ?? ''

    expect(campo).toContain('type="number"')
    expect(campo).toContain('min="1"')
    expect(campo).toContain('step="1"')
  })

  it('el navegador no bloquea el envío: el error de una prioridad escrita a mano lo informa el backend', () => {
    const formulario = html.match(/<form[^>]*>/)?.[0] ?? ''

    expect(formulario).toMatch(/novalidate/i) // en HTML los atributos no distinguen mayúsculas
  })
})
