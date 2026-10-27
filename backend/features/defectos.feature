# language: es
# Spec 11 — Gestionar defectos (Issue #17).
# En las tablas, "(vacío)" representa un texto vacío y "(solo espacios)" uno con solo espacios.
# "El Sprint 1" y "el Sprint 2" son los Sprints de los Antecedentes por su número.
Característica: Gestionar defectos
  Para hacer seguimiento de la calidad del producto
  Como integrante del equipo
  Quiero registrar los defectos encontrados y seguir su estado

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-15"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que existe un Sprint del "2026-11-01" al "2026-11-14"
    Y que existe un Sprint del "2026-11-16" al "2026-11-28"

  Escenario: Registrar un defecto con todos sus datos (CA-1)
    Cuando registro el defecto "El alta acepta emails sin @" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Entonces el defecto queda registrado con un identificador
    Y el defecto tiene descripción "El alta acepta emails sin @", severidad "Alta", historia "Registrar usuario" y Sprint de detección 1
    Y el defecto está "Abierto" y sin Sprint de resolución

  Esquema del escenario: Registrar defectos con cada severidad (CA-1)
    Cuando registro el defecto "Falla" de severidad "<severidad>" en "Registrar usuario" detectado en el Sprint 1
    Entonces el defecto queda registrado con un identificador

    Ejemplos:
      | severidad |
      | Baja      |
      | Media     |
      | Alta      |
      | Crítica   |

  Escenario: Registrar un defecto en una historia completada y en un Sprint cerrado (CA-1)
    Dado que la historia "Registrar usuario" pasa a estado "Completada"
    Y que cerré el Sprint 1
    Cuando registro el defecto "Falla tardía" de severidad "Media" en "Registrar usuario" detectado en el Sprint 1
    Entonces el defecto queda registrado con un identificador

  Escenario: Pasar un defecto a En corrección (CA-2)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando cambio el estado del defecto "Falla" a "En corrección"
    Entonces el defecto está "En corrección" y sin Sprint de resolución

  Escenario: Reabrir un defecto resuelto (CA-2)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que resolví el defecto "Falla" en el Sprint 2
    Cuando cambio el estado del defecto "Falla" a "Abierto"
    Entonces el defecto está "Abierto" y sin Sprint de resolución

  Esquema del escenario: Validar el estado (CA-2)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando cambio el estado del defecto "Falla" a "<estado>"
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | estado        | resultado              |
      | En corrección | modificado             |
      | Abierto       | modificado             |
      | Cerrado       | El estado no es válido |
      | resuelto      | El estado no es válido |

  Escenario: Resolver un defecto registra el Sprint de resolución (CA-3)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando resuelvo el defecto "Falla" en el Sprint 2
    Entonces el defecto está "Resuelto" con Sprint de resolución 2

  Escenario: Resolver en el mismo Sprint de detección (CA-3)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando resuelvo el defecto "Falla" en el Sprint 1
    Entonces el defecto está "Resuelto" con Sprint de resolución 1

  Escenario: No se puede resolver en un Sprint anterior al de detección (CA-3)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 2
    Cuando resuelvo el defecto "Falla" en el Sprint 1
    Entonces se informa el error "El Sprint de resolución no puede ser anterior al de detección"
    Y el defecto está "Abierto" y sin Sprint de resolución

  Escenario: No se puede resolver sin indicar el Sprint de resolución (CA-3)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando cambio el estado del defecto "Falla" a "Resuelto"
    Entonces se informa el error "Para resolver el defecto hay que indicar el Sprint de resolución"

  Escenario: Contar defectos detectados y resueltos (CA-4)
    Dado que registré el defecto "A" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que registré el defecto "B" de severidad "Baja" en "Registrar usuario" detectado en el Sprint 1
    Y que registré el defecto "C" de severidad "Media" en "Registrar usuario" detectado en el Sprint 2
    Y que resolví el defecto "A" en el Sprint 2
    Y que resolví el defecto "C" en el Sprint 2
    Cuando consulto las métricas de defectos
    Entonces el proyecto tiene 3 defectos detectados y 2 resueltos
    Y el Sprint 1 tiene 2 defectos detectados y 0 resueltos
    Y el Sprint 2 tiene 1 defecto detectado y 2 resueltos

  Escenario: Un defecto reabierto deja de contar como resuelto (CA-4)
    Dado que registré el defecto "A" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que resolví el defecto "A" en el Sprint 1
    Y que cambié el estado del defecto "A" a "En corrección"
    Cuando consulto las métricas de defectos
    Entonces el proyecto tiene 1 defecto detectado y 0 resueltos
    Y el Sprint 1 tiene 1 defecto detectado y 0 resueltos

  Escenario: Métricas sin defectos (CA-4)
    Cuando consulto las métricas de defectos
    Entonces el proyecto tiene 0 defectos detectados y 0 resueltos
    Y el Sprint 1 tiene 0 defectos detectados y 0 resueltos
    Y el Sprint 2 tiene 0 defectos detectados y 0 resueltos

  Esquema del escenario: Validar los datos al registrar (CA-5)
    Cuando registro el defecto "<descripcion>" de severidad "<severidad>" en "Registrar usuario" detectado en el Sprint 1
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | descripcion     | severidad | resultado                     |
      | (vacío)         | Alta      | La descripción es obligatoria |
      | (solo espacios) | Alta      | La descripción es obligatoria |
      | Falla           | (vacío)   | La severidad es obligatoria   |
      | Falla           | alta      | La severidad no es válida     |
      | Falla           | Critica   | La severidad no es válida     |

  Esquema del escenario: Validar el largo de la descripción (CA-5)
    Cuando registro un defecto con una descripción de <largo> caracteres de severidad "Baja" en "Registrar usuario" detectado en el Sprint 1
    Entonces el resultado es "<resultado>"

    Ejemplos:
      | largo | resultado                                           |
      | 1000  | registrado                                          |
      | 1001  | La descripción no puede superar los 1000 caracteres |

  Escenario: Registrar sin historia relacionada (CA-5)
    Cuando registro el defecto "Falla" de severidad "Alta" sin historia detectado en el Sprint 1
    Entonces se informa el error "La historia relacionada es obligatoria"

  Escenario: Registrar sin Sprint de detección (CA-5)
    Cuando registro el defecto "Falla" de severidad "Alta" en "Registrar usuario" sin Sprint de detección
    Entonces se informa el error "El Sprint de detección es obligatorio"
    Y no hay defectos registrados

  Escenario: Listar los defectos en el orden en que se registraron (CA-5)
    Dado que registré el defecto "Primero" de severidad "Baja" en "Registrar usuario" detectado en el Sprint 2
    Y que registré el defecto "Segundo" de severidad "Crítica" en "Registrar usuario" detectado en el Sprint 1
    Cuando listo los defectos
    Entonces veo los defectos "Primero" y "Segundo" en ese orden

  Escenario: No se puede eliminar del backlog una historia con defectos (CA-6)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Cuando elimino la historia "Registrar usuario"
    Entonces se informa el error "La historia tiene defectos registrados"
    Y el backlog tiene solo la historia "Registrar usuario"

  Escenario: Registrar en una historia que no existe (CA-7)
    Cuando registro el defecto "Falla" de severidad "Alta" en la historia con identificador 999 detectado en el Sprint 1
    Entonces se informa el error "La historia no existe"
    Y no hay defectos registrados

  Escenario: Registrar en un Sprint que no existe (CA-7)
    Cuando registro el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint con identificador 999
    Entonces se informa el error "El Sprint no existe"

  Escenario: Resolver en un Sprint de otro proyecto (CA-7)
    Dado que registré el defecto "Falla" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto "App Turnos" tiene un Sprint del "2026-11-05" al "2026-11-20"
    Cuando resuelvo el defecto "Falla" en el Sprint de "App Turnos"
    Entonces se informa el error "El Sprint no existe"

  Escenario: Ver un defecto que no existe (CA-7)
    Cuando veo el defecto con identificador 999
    Entonces se informa el error "El defecto no existe"
