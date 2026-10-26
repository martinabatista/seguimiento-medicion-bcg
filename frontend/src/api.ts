// Cliente de la API REST del backend: todos los pedidos de la app pasan por pedir().
// El frontend no valida nada: las reglas están en Go, y si algo está mal se muestra el mensaje
// que manda el backend ({"error": "..."}) tal cual.

// ErrorAPI es un error del pedido, con el mensaje para mostrar y el código HTTP (0 si no hubo respuesta).
export class ErrorAPI extends Error {
  estado: number

  constructor(mensaje: string, estado: number) {
    super(mensaje)
    this.estado = estado
  }
}

// pedir hace un pedido a /api + ruta y devuelve los datos de la respuesta (undefined si es 204).
export async function pedir<T>(metodo: string, ruta: string, cuerpo?: unknown): Promise<T> {
  let respuesta: Response
  try {
    respuesta = await fetch(`/api${ruta}`, {
      method: metodo,
      headers: cuerpo === undefined ? undefined : { 'Content-Type': 'application/json' },
      body: cuerpo === undefined ? undefined : JSON.stringify(cuerpo),
    })
  } catch {
    throw new ErrorAPI('No se pudo conectar con el servidor', 0)
  }
  if (respuesta.status === 204) {
    return undefined as T
  }
  const datos = await respuesta.json().catch(() => null)
  if (!respuesta.ok) {
    throw new ErrorAPI(datos?.error ?? `Error del servidor (${respuesta.status})`, respuesta.status)
  }
  return datos as T
}

// mensajeDe devuelve el texto para mostrar de cualquier error.
export function mensajeDe(error: unknown): string {
  return error instanceof Error ? error.message : 'Ocurrió un error inesperado'
}
