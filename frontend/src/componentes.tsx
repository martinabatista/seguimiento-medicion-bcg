// Componentes que comparten todas las pantallas, para mostrar la carga y los errores igual en todas.

export function Cargando() {
  return <p className="suave">Cargando…</p>
}

// MensajeError muestra el mensaje tal cual llega (normalmente, el del backend). Sin mensaje, no muestra nada.
export function MensajeError({ mensaje }: { mensaje: string }) {
  if (!mensaje) {
    return null
  }
  return (
    <p className="error" role="alert">
      {mensaje}
    </p>
  )
}
