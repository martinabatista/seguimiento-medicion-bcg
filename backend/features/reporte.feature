# language: es
# Spec 14 — Reporte de proyecto o Sprint exportable a PDF (Issue #21).
# "El Sprint 1" y "el Sprint 2" son los Sprints del escenario por su número. Las horas usan punto decimal.
Característica: Reporte de proyecto o Sprint exportable a PDF
  Para compartir y archivar los resultados
  Como integrante del equipo
  Quiero generar el reporte de un proyecto o de un Sprint y descargarlo en PDF

  Antecedentes:
    Dado que existe el proyecto "Tienda Online" con inicio "2026-11-01" y fin "2026-12-31"
    Y que hoy es "2026-11-20"
    Y que el proyecto tiene al integrante "Ana" con email "ana@mail.com" y rol "Developer"
    Y que el backlog tiene la historia "Registrar usuario" con prioridad 1
    Y que el backlog tiene la historia "Ver catálogo" con prioridad 2
    Y que el backlog tiene la historia "Pagar pedido" con prioridad 3

  Escenario: El reporte del proyecto tiene el encabezado y todo el backlog (CA-1)
    Cuando genero el reporte del proyecto
    Entonces el reporte es del proyecto "Tienda Online", "En curso", generado el "2026-11-20"
    Y el reporte tiene al integrante "Ana" con rol "Developer"
    Y el reporte muestra las historias "Registrar usuario", "Ver catálogo" y "Pagar pedido" en ese orden

  Escenario: El reporte del proyecto muestra lo planificado y completado en cada Sprint (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Cuando genero el reporte del proyecto
    Entonces la historia "Registrar usuario" del reporte está "Completada"
    Y el Sprint 1 del reporte tiene 2 historias planificadas, 1 completada y 50 % completadas

  Escenario: El reporte de un Sprint indica qué historias planificadas se completaron (CA-1)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando genero el reporte del Sprint 1
    Entonces el reporte del Sprint muestra las historias "Registrar usuario" y "Ver catálogo"
    Y en el reporte del Sprint la historia "Registrar usuario" está completada y "Ver catálogo" no
    Y el reporte no incluye la historia "Pagar pedido"

  Escenario: El reporte incluye las estimaciones de cada historia (CA-2)
    Dado que estimé la historia "Registrar usuario" con 5 Story Points
    Y que estimé la historia "Registrar usuario" en 8 horas
    Cuando genero el reporte del proyecto
    Entonces la historia "Registrar usuario" del reporte tiene 5 Story Points y 8 horas estimadas
    Y la historia "Ver catálogo" del reporte está sin estimar y sin horas estimadas

  Escenario: El reporte de un Sprint cerrado usa los Story Points congelados al cerrarlo (CA-2)
    Dado que estimé la historia "Ver catálogo" con 5 Story Points
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que cerré el Sprint
    Y que estimé la historia "Ver catálogo" con 8 Story Points
    Cuando genero el reporte del Sprint 1
    Entonces la historia "Ver catálogo" del reporte tiene 5 Story Points

  Escenario: El reporte incluye el esfuerzo registrado y la desviación (CA-3)
    Dado que estimé la historia "Registrar usuario" en 8 horas
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que registré 6 horas de "Ana" el "2026-11-05" en "Registrar usuario" con la actividad "Backend"
    Y que registré 4 horas de "Ana" el "2026-11-06" en "Registrar usuario" con la actividad "Pantalla"
    Y que completé en el Sprint la historia "Registrar usuario"
    Cuando genero el reporte del Sprint 1
    Entonces la historia "Registrar usuario" del reporte tiene 10 horas reales y una desviación de 2 horas y 25 %
    Y el Sprint del reporte tiene 8 horas estimadas, 10 horas reales y una desviación de 2 horas y 25 %

  Escenario: Una historia sin terminar muestra sus horas reales sin desviación (CA-3)
    Dado que estimé la historia "Ver catálogo" en 5 horas
    Y que registré 2 horas de "Ana" el "2026-11-05" en "Ver catálogo" con la actividad "Diseño"
    Cuando genero el reporte del proyecto
    Entonces la historia "Ver catálogo" del reporte tiene 2 horas reales y desviación "sin datos"

  Escenario: El reporte del proyecto incluye las métricas por Sprint y la velocidad (CA-4)
    Dado que estimé la historia "Registrar usuario" con 5 Story Points
    Y que estimé la historia "Ver catálogo" con 3 Story Points
    Y que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Y que completé en el Sprint la historia "Registrar usuario"
    Y que cerré el Sprint
    Cuando genero el reporte del proyecto
    Entonces el Sprint 1 del reporte tiene 8 Story Points planificados y 5 completados
    Y la velocidad del reporte es 5.0
    Y el proyecto del reporte tiene 1 Sprint cerrado con 8 Story Points planificados y 5 completados

  Escenario: El reporte de un Sprint abierto es parcial (CA-4)
    Dado que existe un Sprint del "2026-11-16" al "2026-11-28" con las historias "Registrar usuario" y "Ver catálogo"
    Cuando genero el reporte del Sprint 1
    Entonces el reporte del Sprint dice que es parcial

  Escenario: Sin Sprints cerrados no hay velocidad (CA-4)
    Dado que existe un Sprint del "2026-11-16" al "2026-11-28" con las historias "Registrar usuario" y "Ver catálogo"
    Cuando genero el reporte del proyecto
    Entonces el reporte dice que todavía no hay velocidad
    Y el Sprint 1 del reporte es parcial

  Escenario: El reporte del proyecto incluye todos los defectos (CA-5)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Y que existe un Sprint del "2026-11-16" al "2026-11-28"
    Y que registré el defecto "El alta acepta emails sin @" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que registré el defecto "Fotos que no cargan" de severidad "Crítica" en "Ver catálogo" detectado en el Sprint 2
    Y que resolví el defecto "El alta acepta emails sin @" en el Sprint 2
    Cuando genero el reporte del proyecto
    Entonces el reporte tiene 2 defectos detectados y 1 resuelto
    Y el reporte muestra los defectos "El alta acepta emails sin @" y "Fotos que no cargan" en ese orden
    Y el defecto "El alta acepta emails sin @" del reporte está "Resuelto" en el Sprint 2

  Escenario: El reporte de un Sprint incluye los defectos detectados o resueltos en él (CA-5)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Y que existe un Sprint del "2026-11-16" al "2026-11-28"
    Y que registré el defecto "A" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que registré el defecto "B" de severidad "Baja" en "Registrar usuario" detectado en el Sprint 2
    Y que registré el defecto "C" de severidad "Media" en "Ver catálogo" detectado en el Sprint 1
    Y que resolví el defecto "A" en el Sprint 2
    Cuando genero el reporte del Sprint 2
    Entonces el reporte tiene 1 defecto detectado y 1 resuelto
    Y el reporte muestra los defectos "A" y "B" en ese orden

  Escenario: Un defecto reabierto no figura como resuelto en el reporte (CA-5)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Y que registré el defecto "A" de severidad "Alta" en "Registrar usuario" detectado en el Sprint 1
    Y que resolví el defecto "A" en el Sprint 1
    Y que cambié el estado del defecto "A" a "Abierto"
    Cuando genero el reporte del proyecto
    Entonces el reporte tiene 1 defecto detectado y 0 resueltos
    Y el defecto "A" del reporte está "Abierto" sin Sprint de resolución

  Escenario: Descargar el reporte del proyecto en PDF (CA-6)
    Cuando descargo el reporte del proyecto en PDF
    Entonces recibo un archivo PDF llamado "Reporte Tienda Online.pdf"
    Y el PDF incluye el texto "Tienda Online"
    Y el PDF incluye las secciones "Historias", "Sprints", "Métricas" y "Defectos"

  Escenario: Descargar el reporte de un Sprint en PDF (CA-6)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14" con las historias "Registrar usuario" y "Ver catálogo"
    Cuando descargo el reporte del Sprint 1 en PDF
    Entonces recibo un archivo PDF llamado "Reporte Tienda Online - Sprint 1.pdf"
    Y el PDF incluye el texto "Registrar usuario"

  Escenario: El PDF muestra bien las tildes y la ñ (CA-6)
    Dado que existe el proyecto "Gestión de Turnos" con inicio "2026-11-01" y fin "2026-12-31"
    Cuando descargo el reporte del proyecto "Gestión de Turnos" en PDF
    Entonces recibo un archivo PDF llamado "Reporte Gestión de Turnos.pdf"
    Y el PDF incluye el texto "Gestión de Turnos"

  Escenario: Un proyecto sin datos igual tiene reporte (CA-7)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Cuando genero el reporte del proyecto "App Turnos"
    Entonces el reporte no tiene historias, Sprints ni defectos
    Y el reporte dice que todavía no hay velocidad
    Y la desviación del proyecto en el reporte es "sin datos"

  Escenario: Un Sprint sin historias igual tiene reporte (CA-7)
    Dado que existe un Sprint del "2026-11-01" al "2026-11-14"
    Cuando genero el reporte del Sprint 1
    Entonces el reporte del Sprint no tiene historias y tiene 0 % completadas

  Escenario: Reporte de un proyecto que no existe (CA-7)
    Cuando genero el reporte del proyecto con identificador 999
    Entonces se informa el error "El proyecto no existe"

  Escenario: Reporte de un Sprint que no existe (CA-7)
    Cuando genero el reporte del Sprint con identificador 999
    Entonces se informa el error "El Sprint no existe"

  Escenario: Reporte de un Sprint de otro proyecto (CA-7)
    Dado que existe el proyecto "App Turnos" con inicio "2026-11-05" y fin "2027-01-20"
    Y que el proyecto "App Turnos" tiene un Sprint del "2026-11-05" al "2026-11-20"
    Cuando genero el reporte del Sprint de "App Turnos" desde el proyecto "Tienda Online"
    Entonces se informa el error "El Sprint no existe"

  Escenario: No se descarga el PDF de un proyecto que no existe (CA-7)
    Cuando descargo el reporte del proyecto con identificador 999 en PDF
    Entonces se informa el error "El proyecto no existe"
