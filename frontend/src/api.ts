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
  horas_estimadas: number | null // null = sin horas estimadas (spec 12)
}

// Integrante del proyecto (spec 02).
export interface Integrante {
  id: number
  proyecto_id: number
  nombre: string
  email: string
  rol: 'Product Owner' | 'Scrum Master' | 'Developer'
}

// Sprint del proyecto (spec 05). Al listarlos llega sin historias.
export interface Sprint {
  id: number
  proyecto_id: number
  numero: number
  fecha_inicio: string // AAAA-MM-DD
  fecha_fin: string
  sprint_goal: string
  estado: 'Abierto' | 'Cerrado'
}

// Estado de un proyecto (spec 03), ya calculado por el backend con la fecha de hoy. El Sprint en curso
// es el abierto cuyo período incluye hoy (spec 05 CA-5), o null si no hay.
export interface EstadoProyecto {
  proyecto: Proyecto
  estado: 'No iniciado' | 'En curso' | 'Finalizado'
  backlog: {
    total: number
    pendientes: number
    en_progreso: number
    completadas: number
    porcentaje_completadas: number
    story_points_totales: number
    story_points_completados: number
    sin_estimar: number
  }
  sprint_en_curso: Sprint | null
}

// Un Sprint con sus historias (specs 05 y 06): llega al verlo y como respuesta de cada cambio.
// completada_en_sprint dice si se completó en ese Sprint; en uno cerrado ya no cambia.
export interface DetalleSprint extends Sprint {
  historias: (Historia & { completada_en_sprint: boolean })[]
}

// Métricas de avance (spec 10), ya calculadas por el backend. Los Sprints van por fecha de inicio;
// "parcial" marca el abierto. La velocidad y los totales cuentan solo los Sprints cerrados.
export interface Metricas {
  sprints: {
    sprint_id: number
    numero: number
    estado: Sprint['estado']
    parcial: boolean
    historias_planificadas: number
    historias_completadas: number
    porcentaje_completadas: number
    sin_estimar: number
    story_points_planificados: number
    story_points_completados: number
  }[]
  sprints_cerrados: number
  velocidad: number | null // null = todavía no hay Sprints cerrados
  totales: {
    historias_planificadas: number
    historias_completadas: number
    porcentaje_completadas: number
    story_points_planificados: number
    story_points_completados: number
  }
}

// Métricas de esfuerzo (spec 12), ya calculadas por el backend. Las horas son números con hasta 2
// decimales; la desviación es positiva si se trabajó más de lo estimado y null cuando no hay datos.
export interface Desviacion {
  horas: number
  porcentaje: number
}

export interface MetricasEsfuerzo {
  historias: {
    historia_id: number
    titulo: string
    estado: Historia['estado']
    horas_estimadas: number | null
    horas_reales: number
    desviacion: Desviacion | null // solo las completadas con horas estimadas
  }[]
  sprints: {
    sprint_id: number
    numero: number
    estado: Sprint['estado']
    parcial: boolean
    horas_estimadas: number
    horas_reales: number
    desviacion: Desviacion | null
    completadas_sin_horas_estimadas: number
  }[]
  horas_estimadas: number
  horas_reales: number
  desviacion: Desviacion | null
  completadas_sin_horas_estimadas: number
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

// Dashboard del proyecto (spec 13): indicadores y datos de los gráficos, ya calculados por el backend.
// Donde no hay dato llega null.
export interface Dashboard {
  proyecto: Proyecto
  estado: 'No iniciado' | 'En curso' | 'Finalizado'
  dias_para_empezar: number | null
  dias_para_terminar: number | null
  sprint_en_curso: {
    id: number
    numero: number
    sprint_goal: string
    fecha_fin: string
    dias_restantes: number
    story_points_planificados: number
    story_points_completados: number
  } | null
  backlog: {
    total: number
    completadas: number
    porcentaje_completadas: number
    story_points_totales: number
    story_points_completados: number
    sin_estimar: number
  }
  velocidad: number | null
  graficos: {
    story_points_por_sprint: { sprint: string; numero: number; planificados: number; completados: number; parcial: boolean }[]
    velocidad: { sprints: { sprint: string; numero: number; completados: number }[]; promedio: number | null }
    backlog_por_estado: { estado: Historia['estado']; cantidad: number }[]
  }
}

// mensajeDe devuelve el texto para mostrar de cualquier error.
export function mensajeDe(error: unknown): string {
  return error instanceof Error ? error.message : 'Ocurrió un error inesperado'
}
