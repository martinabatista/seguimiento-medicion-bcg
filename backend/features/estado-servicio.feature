# language: es
Característica: Estado del servicio
  Para saber si la API está disponible
  Como integrante del equipo
  Quiero consultar el estado del servicio

  Escenario: El servicio informa que está activo (CA-3)
    Cuando consulto el estado del servicio
    Entonces la respuesta tiene código 200
    Y el estado informado es "ok"
