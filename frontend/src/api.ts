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

// Tipos de los datos que devuelve el backend (los JSON están en backend/CLAUDE.md).
// Cada integrante suma acá los de sus pantallas.

export interface Proyecto {
  id: number
  nombre: string
  descripcion: string
  fecha_inicio: string // AAAA-MM-DD
  fecha_fin: string
}

// Historia del Product Backlog (specs 04 y 07).
export interface Historia {
  id: number
  proyecto_id: number
  titulo: string
  descripcion: string
  prioridad: number // 1 es la más importante
  criterios_aceptacion: string
  estado: 'Pendiente' | 'En progreso' | 'Completada'
  story_points: number | null // null = sin estimar
}

// Integrante del proyecto (spec 02).
export interface Integrante {
  id: number
  proyecto_id: number
  nombre: string
  email: string
  rol: 'Product Owner' | 'Scrum Master' | 'Developer'
}

// Votación de Planning Poker (spec 08). Los valores de los votos solo aparecen en "rondas", que trae
// las rondas ya reveladas: de la ronda sin revelar solo se sabe quién votó y quién falta.
export interface Votacion {
  id: number
  proyecto_id: number
  historia_id: number
  estado: 'Abierta' | 'Cerrada'
  ronda: number // la ronda actual
  estimacion_acordada: number | null
  votaron: Persona[]
  faltan: Persona[]
  rondas: RondaRevelada[]
}

export interface Persona {
  id: number
  nombre: string
}

export interface RondaRevelada {
  numero: number
  votos: { integrante_id: number; nombre: string; story_points: number }[]
  hay_diferencias: boolean
  mas_bajo: { story_points: number; integrantes: Persona[] }
  mas_alto: { story_points: number; integrantes: Persona[] }
}

// mensajeDe devuelve el texto para mostrar de cualquier error.
export function mensajeDe(error: unknown): string {
  return error instanceof Error ? error.message : 'Ocurrió un error inesperado'
}
