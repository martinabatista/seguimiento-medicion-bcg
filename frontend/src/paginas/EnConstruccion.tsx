// EnConstruccion ocupa el lugar de una sección que todavía no tiene pantalla.
export default function EnConstruccion({ titulo }: { titulo: string }) {
  return (
    <section>
      <h2>{titulo}</h2>
      <p className="suave">Esta sección todavía está en construcción.</p>
    </section>
  )
}
