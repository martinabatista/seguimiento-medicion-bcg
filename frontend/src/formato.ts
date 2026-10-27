// Formatos de los números que llegan del backend, comunes a las pantallas de Métricas y Reporte (y
// a los textos del PDF del reporte, spec 14 RN-8): solo muestran, no calculan.
import type { Desviacion } from './api'

// decimal muestra la velocidad con coma y un decimal, como llega del backend (6.5 → "6,5").
export function decimal(n: number): string {
  return n.toLocaleString('es-AR', { maximumFractionDigits: 1 })
}

// horas muestra horas con coma y hasta 2 decimales, como llegan del backend (2.25 → "2,25").
export function horas(n: number): string {
  return n.toLocaleString('es-AR', { maximumFractionDigits: 2 })
}

// desviacion muestra la desviación con su signo: "+2 h (+25 %)". Sin datos, el texto indicado.
export function desviacion(d: Desviacion | null, sinDatos = '—'): string {
  if (d === null) {
    return sinDatos
  }
  const signo = (n: number) => (n > 0 ? '+' : '')
  return `${signo(d.horas)}${horas(d.horas)} h (${signo(d.porcentaje)}${d.porcentaje} %)`
}
