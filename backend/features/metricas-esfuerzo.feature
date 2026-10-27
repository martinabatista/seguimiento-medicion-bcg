# language: es
# Spec 12 — Métricas de esfuerzo (Issue #18).
Característica: Métricas de esfuerzo
  Para evaluar la precisión de nuestras estimaciones
  Como integrante del equipo
  Quiero comparar las horas que estimamos con las que trabajamos

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-31"
    Y que hoy es "2026-11-20"
    Y que el proyecto tiene al integrante "Ana" con email "ana@mail.com" y rol "Developer"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2

  Escenario: Cargar las horas estimadas de una historia (CA-1)
    Cuando estimo la historia "Registrar usuario" en 8 horas
    Entonces la historia "Registrar usuario" tiene 8 horas estimadas

  Escenario: Cambiar las horas estimadas, con decimales (CA-1)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Cuando estimo la historia "Registrar usuario" en 6.5 horas
    Entonces la historia "Registrar usuario" tiene 6.5 horas estimadas

  Escenario: Quitar las horas estimadas (CA-1)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Cuando quito las horas estimadas de la historia "Registrar usuario"
    Entonces la historia "Registrar usuario" no tiene horas estimadas

  Esquema del escenario: Validar las horas estimadas (CA-1)
    Cuando estimo la historia "Registrar usuario" en <horas> horas
    Entonces <resultado>

    Ejemplos:
      | horas | resultado                                                                                      |
      | 0.01  | la historia "Registrar usuario" tiene 0.01 horas estimadas                                     |
      | 500   | la historia "Registrar usuario" tiene 500 horas estimadas                                      |
      | 0     | se informa el error "Las horas estimadas deben ser un número mayor a 0, con hasta 2 decimales" |
      | -3    | se informa el error "Las horas estimadas deben ser un número mayor a 0, con hasta 2 decimales" |
      | 1.125 | se informa el error "Las horas estimadas deben ser un número mayor a 0, con hasta 2 decimales" |

  Escenario: Faltan las horas estimadas (CA-1)
    Cuando estimo la historia "Registrar usuario" sin indicar las horas
    Entonces se informa el error "Las horas estimadas son obligatorias"

  Escenario: No se pueden estimar las horas de una historia completada (CA-2)
    Dado que la historia "Registrar usuario" pasa a estado "Completada"
    Cuando estimo la historia "Registrar usuario" en 8 horas
    Entonces se informa el error "La historia ya está completada"

  Escenario: No se pueden quitar las horas estimadas de una historia completada (CA-2)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que la historia "Registrar usuario" pasa a estado "Completada"
    Cuando quito las horas estimadas de la historia "Registrar usuario"
    Entonces se informa el error "La historia ya está completada"

  Escenario: Horas estimadas y reales de cada historia (CA-3)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que registré 3 horas de "Ana" el "2026-11-10" en "Registrar usuario" con la actividad "Programar"
    Y que registré 2.5 horas de "Ana" el "2026-11-11" en "Ver catálogo" con la actividad "Diseñar"
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces la historia "Registrar usuario" tiene 8 horas estimadas y 3 horas reales
    Y la historia "Ver catálogo" no tiene horas estimadas y tiene 2.5 horas reales
    Y el proyecto tiene 8 horas estimadas y 5.5 horas reales

  Escenario: Una historia sin terminar no tiene desviación (CA-4)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que registré 10 horas de "Ana" el "2026-11-10" en "Registrar usuario" con la actividad "Programar"
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces la historia "Registrar usuario" no tiene desviación
    Y la desviación del proyecto es "sin datos"

  Esquema del escenario: Desviación de una historia completada (CA-4)
    Dado que estimé la historia "Registrar usuario" en <estimadas> horas
    Y que registré <reales> horas de "Ana" el "2026-11-10" en "Registrar usuario" con la actividad "Programar"
    Y que la historia "Registrar usuario" pasa a estado "Completada"
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces la historia "Registrar usuario" tiene una desviación de <horas> horas y <porcentaje> %

    Ejemplos:
      | estimadas | reales | horas | porcentaje |
      | 8         | 10     | 2     | 25         |
      | 8         | 6      | -2    | -25        |
      | 3         | 3.5    | 0.5   | 17         |
      | 8         | 8      | 0     | 0          |

  Escenario: Historia completada sin esfuerzo registrado (CA-4)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que la historia "Registrar usuario" pasa a estado "Completada"
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces la historia "Registrar usuario" tiene una desviación de -8 horas y -100 %

  Escenario: La desviación de un Sprint se calcula sobre las sumas (CA-3, CA-4)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que estimé la historia "Ver catálogo" en 4 horas
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que registré 10 horas de "Ana" el "2026-11-10" en "Registrar usuario" con la actividad "Programar"
    Y que registré 2 horas de "Ana" el "2026-11-11" en "Ver catálogo" con la actividad "Diseñar"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que completé en el Sprint la historia "Ver catálogo"
    Y que cerré el Sprint
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces el Sprint 1 tiene 12 horas estimadas, 12 horas reales y una desviación de 0 horas y 0 %

  Escenario: Las métricas de un Sprint abierto son parciales (CA-5)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que existe un Sprint del "2026-11-15" al "2026-11-28" con la historia "Registrar usuario"
    Y que registré 10 horas de "Ana" el "2026-11-16" en "Registrar usuario" con la actividad "Programar"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces las métricas de esfuerzo del Sprint 1 son parciales
    Y el Sprint 1 tiene 8 horas estimadas, 10 horas reales y una desviación de 2 horas y 25 %

  Escenario: Las historias completadas sin horas estimadas se informan aparte (CA-5)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con la historia "Ver catálogo"
    Y que completé en el Sprint la historia "Ver catálogo"
    Y que cerré el Sprint
    Cuando consulto las métricas de esfuerzo del proyecto
    Entonces el Sprint 1 tiene 1 historia completada sin horas estimadas
    Y la desviación del Sprint 1 es "sin datos"
    Y el proyecto tiene 1 historia completada sin horas estimadas

  Escenario: Consultar las métricas de esfuerzo de un proyecto que no existe (CA-7)
    Cuando consulto las métricas de esfuerzo del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Estimar las horas de una historia que no existe (CA-7)
    Cuando estimo la historia con identificador 999 en 8 horas
    Entonces se informa el error "La historia no existe"
