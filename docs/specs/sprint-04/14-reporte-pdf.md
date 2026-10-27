---
titulo: Reporte de proyecto o Sprint exportable a PDF
estado: implemented      # draft | in-review | approved | implemented
unidad: global
issue: "#21"
sprint: 04
autor: Luana Guajardo
aprobada_por: "Martina Batista"
fecha: 2026-10-27
specs_relacionadas: ["03-estado-proyecto", "05-crear-sprints", "06-cerrar-sprint", "07-story-points", "09-registro-esfuerzo", "10-metricas-avance", "11-defectos", "12-metricas-esfuerzo-calidad"]
---

# Spec 14 — Reporte de proyecto o Sprint exportable a PDF

> Describe **qué** y **para qué**, no el cómo. Spec en su rama `docs/` y PR propio; el código va después en otra rama `feat/`.

## Historia de usuario

Como integrante del equipo, quiero generar un reporte de un proyecto o de un Sprint y exportarlo a PDF, para compartir y archivar los resultados. ([#21](https://github.com/martinabatista/seguimiento-medicion-bcg/issues/21))

## Objetivo

Juntar en un solo informe lo que el sistema ya registra y calcula sobre un proyecto o un Sprint: las historias planificadas y completadas, sus estimaciones, el esfuerzo, las métricas y los defectos. El informe se puede ver en pantalla y **descargar como archivo PDF**, para entregarlo o guardarlo. El reporte **no inventa cálculos nuevos**: muestra los de las specs 03, 10, 11 y 12, con sus mismas reglas.

Origen: Consigna — Requerimiento 9 (Reportes) y *Sprint 4: exportar un informe del proyecto a PDF*.

## Alcance

**Entra:**
- **Reporte del proyecto:** los datos del proyecto, todo el backlog con estimaciones y esfuerzo, las métricas por Sprint y del proyecto, y los defectos.
- **Reporte de un Sprint:** los datos del Sprint, sus historias planificadas y cuáles se completaron, sus estimaciones y esfuerzo, sus métricas y sus defectos. Puede ser un Sprint abierto (el reporte es parcial) o cerrado.
- Una consulta del backend por cada reporte, que devuelve todos sus datos juntos y ya calculados (para la vista previa).
- La descarga de cada reporte como **archivo PDF**, armado en el backend.
- La pantalla "Reporte": se elige el proyecto completo o un Sprint, se ve el reporte y se descarga el PDF.

**NO entra:**
- Gráficos dentro del PDF (los gráficos están en el Dashboard, spec 13).
- Elegir qué secciones incluir, cambiar el diseño del PDF o exportar a otros formatos (Excel, Word).
- El detalle de cada registro de esfuerzo (fecha, integrante, actividad): el reporte muestra el **total por historia** (decisión de la dueña de la historia, para que el informe quede corto y legible). El detalle se consulta en la pantalla de Esfuerzo (spec 09).
- Guardar los reportes generados (un historial) o enviarlos por mail: el reporte se genera en el momento.
- Métricas nuevas: si un dato no lo calcula ya otra spec, no está en el reporte.

## Entradas

| Campo | Tipo | Obligatorio | Validación |
|---|---|---|---|
| Proyecto | identificador | Sí | Tiene que ser un proyecto existente (spec 01). |
| Sprint | identificador | Solo para el reporte de un Sprint | Tiene que ser un Sprint del proyecto (spec 05), abierto o cerrado. |
| Formato | vista previa o PDF | Sí | La vista previa devuelve los datos; el PDF, el archivo para descargar. |

Nada más: el reporte no se carga, se arma con los datos del momento y la fecha de hoy.

## Salidas esperadas

**Reporte del proyecto**, en este orden:
1. **Encabezado:** nombre y descripción del proyecto, fechas de inicio y fin, estado (`No iniciado`, `En curso` o `Finalizado`, spec 03), integrantes con su rol (spec 02) y la fecha en que se generó el reporte.
2. **Historias:** todas las historias del backlog, en el orden del backlog (spec 04 RN-8), cada una con su estado, sus Story Points (o "sin estimar"), sus horas estimadas (o "sin horas estimadas"), sus horas reales y su desviación (o "sin datos"), como en las specs 07 y 12.
3. **Sprints:** por cada Sprint, en orden de fecha de inicio: número, fechas, Sprint Goal, estado, historias planificadas y completadas con su porcentaje, Story Points planificados y completados, historias sin estimar (spec 10) y horas estimadas, reales y desviación (spec 12). El Sprint abierto va marcado como **parcial**.
4. **Métricas del proyecto:** la velocidad del equipo (o que todavía no hay), los totales de los Sprints cerrados (spec 10 RN-7) y las horas estimadas, reales y la desviación del proyecto, con las historias completadas sin horas estimadas (spec 12 RN-9, RN-10).
5. **Defectos:** cantidad de detectados y de resueltos (spec 11 CA-4) y la lista de defectos en el orden en que se registraron, cada uno con su descripción, severidad, estado, historia, Sprint de detección y Sprint de resolución (si tiene).

**Reporte de un Sprint**, en este orden:
1. **Encabezado:** nombre del proyecto, número del Sprint, fechas, Sprint Goal, estado (y si es **parcial**, porque está abierto) y la fecha en que se generó el reporte.
2. **Historias planificadas:** las asignadas al Sprint, en el orden del backlog, cada una con si se **completó en el Sprint**, sus Story Points (en un Sprint cerrado, los congelados al cerrarlo, spec 10 RN-4), sus horas estimadas, sus horas reales y su desviación (solo si se completó en el Sprint, spec 12 RN-7).
3. **Métricas del Sprint:** historias planificadas y completadas con su porcentaje, Story Points planificados y completados, historias sin estimar (spec 10), y horas estimadas, reales, desviación e historias completadas sin horas estimadas (spec 12 RN-8). Los mismos valores que la fila de ese Sprint en las specs 10 y 12.
4. **Defectos:** cantidad de detectados y de resueltos en el Sprint (spec 11 CA-4) y la lista de los defectos **detectados o resueltos en ese Sprint**, con los mismos datos que en el reporte del proyecto.

**PDF:** un archivo PDF con el mismo contenido que la vista previa, en ese mismo orden, con un título que dice de qué proyecto (y de qué Sprint) es. El nombre del archivo es `Reporte <nombre del proyecto>.pdf` o `Reporte <nombre del proyecto> - Sprint <número>.pdf`.

**Error:** no se genera nada y se informa qué no existe.

## Reglas de negocio

- RN-1: Solo se puede generar el reporte de un proyecto existente, y el de un Sprint de ese proyecto; un Sprint de otro proyecto se trata como inexistente.
- RN-2: El reporte **reutiliza** los cálculos de las specs 03 (estado), 10 (métricas de avance y velocidad), 11 (defectos detectados y resueltos) y 12 (horas y desviación), con sus mismas reglas: redondeos, Story Points congelados al cerrar un Sprint, métricas parciales del Sprint abierto, "sin estimar" y "sin datos". El mismo dato vale lo mismo en el reporte y en las otras pantallas.
- RN-3: Las **historias planificadas** de un Sprint son las asignadas a él y las **completadas**, las que se completaron en él (spec 10 RN-1). En el reporte de un Sprint cerrado, los Story Points de cada historia son los congelados al cerrarlo; en uno abierto, los actuales.
- RN-4: El **esfuerzo registrado** se muestra como el total de horas reales de cada historia (spec 12 RN-5), junto a sus horas estimadas y su desviación; y los totales del Sprint o del proyecto (spec 12 RN-8, RN-9). No se listan los registros uno por uno.
- RN-5: En el reporte de un Sprint, los **defectos** son los detectados en ese Sprint o resueltos en ese Sprint (spec 11 RN-10): un defecto detectado en el Sprint 1 y resuelto en el 2 aparece en los dos reportes. Las cantidades son las de la spec 11 para ese Sprint.
- RN-6: El reporte de un Sprint **abierto** se puede generar y va marcado como **parcial**: el Sprint todavía no terminó (spec 10 RN-3).
- RN-7: Un proyecto sin historias, sin Sprints, sin esfuerzo o sin defectos **igual tiene reporte**: cada sección dice que no hay datos ("sin historias", "sin Sprints", "sin defectos", "todavía no hay velocidad", desviación "sin datos"). Nunca es un error. Lo mismo un Sprint sin historias.
- RN-8: La **vista previa y el PDF** muestran los mismos datos, en el mismo orden, con los mismos textos para lo que falta. Los números se muestran como en las otras pantallas (horas con hasta 2 decimales, porcentajes enteros, velocidad con un decimal).
- RN-9: El reporte se arma en el momento, con los datos actuales y la fecha de hoy (la del servidor, como en la spec 03), que figura en el encabezado como fecha de generación.
- RN-10: El PDF muestra bien los textos en español (tildes, ñ, ¿, signos) y los textos largos (por ejemplo, una descripción de 1000 caracteres) siguen en la línea siguiente sin cortarse ni salirse de la hoja. Si el contenido no entra en una hoja, sigue en otra.

## Restricciones

- Las reglas de negocio y los cálculos se implementan en Go (backend); el frontend solo presenta y captura datos. **El PDF también se arma en Go** (decisión de la dueña de la historia): la pantalla solo lo descarga.
- El armado de cada reporte (qué historias, qué defectos y qué valores entran, RN-3 a RN-7) vive en el dominio y se desarrolla con TDD, reutilizando los cálculos ya implementados en las specs 03, 10, 11 y 12 (no los repite).
- Generar el PDF necesita una librería: se justifica en el PR que la agrega (regla de dependencias del repo). Tiene que ser de Go puro, sin programas externos instalados, para que corra igual en la CI y en las máquinas del equipo.
- Usa las specs 02, 03, 04, 05, 06, 07, 09, 10, 11 y 12, todas implementadas en Sprints anteriores: esta historia no espera a ninguna otra.

## Casos límite

- Proyecto recién creado (sin integrantes, historias, Sprints ni defectos): reporte con el encabezado y "sin …" en cada sección; velocidad "todavía no hay"; desviación "sin datos".
- Sprint sin historias asignadas: 0 planificadas, 0 completadas, 0 % y "sin historias".
- Sprint abierto: el reporte se genera y dice que es parcial.
- Sprint cerrado y una de sus historias re-estimada después de cerrarlo: el reporte del Sprint muestra los Story Points congelados; el del proyecto, en la lista de historias, los actuales.
- Historia sin Story Points ni horas estimadas: "sin estimar" y "sin horas estimadas"; su desviación "sin datos".
- Historia con horas reales pero no completada: muestra sus horas reales y la desviación "sin datos" (spec 12 RN-7).
- Defecto detectado en el Sprint 1 y resuelto en el 2: aparece en el reporte de los dos Sprints; en el del 1 cuenta como detectado y en el del 2 como resuelto.
- Defecto reabierto: en el reporte no tiene Sprint de resolución y no cuenta como resuelto (spec 11 RN-10).
- Nombre del proyecto con tildes (`Gestión de Turnos`): el PDF se llama `Reporte Gestión de Turnos.pdf` y el título lo muestra con tilde.
- Descripción de un defecto o de una historia muy larga: en el PDF sigue en varias líneas.

## Condiciones de error

| Situación | Respuesta esperada |
|---|---|
| El proyecto no existe | No encontrado: "El proyecto no existe" (spec 01) |
| El Sprint no existe o es de otro proyecto | No encontrado: "El Sprint no existe" (spec 05) |

No hay otros errores: lo que falta se muestra como "sin …" (RN-7).

## Criterios de aceptación

- CA-1: El reporte incluye las historias planificadas y completadas: en el del proyecto, todo el backlog con su estado y, por Sprint, las planificadas y completadas; en el de un Sprint, sus historias planificadas indicando cuáles se completaron en él.
- CA-2: El reporte incluye las estimaciones de cada historia: Story Points (congelados en un Sprint cerrado) y horas estimadas, o que no las tiene.
- CA-3: El reporte incluye el esfuerzo registrado: las horas reales de cada historia, su desviación y los totales del Sprint o del proyecto.
- CA-4: El reporte incluye las métricas: por Sprint y del proyecto (historias y Story Points completados, porcentaje, velocidad, horas y desviación), con los mismos valores que las specs 10 y 12, y el Sprint abierto marcado como parcial.
- CA-5: El reporte incluye los defectos: las cantidades de detectados y resueltos y la lista de defectos (en el de un Sprint, los detectados o resueltos en él).
- CA-6: El reporte del proyecto y el de un Sprint se pueden descargar como un archivo PDF con el mismo contenido, nombrado según el proyecto (y el Sprint).
- CA-7: Un proyecto o Sprint sin datos igual tiene reporte (con "sin …"); un proyecto o Sprint que no existe informa cuál no existe, sin generar nada.
- CA-8: La pantalla "Reporte" deja elegir el proyecto completo o un Sprint, muestra el reporte y descarga su PDF.

Relación con los criterios del Issue #21: CA-1 a CA-6 son los del Issue (CA-6 suma el reporte de un Sprint, que pide la historia de usuario). CA-7 cubre los casos sin datos y de error; CA-8, la pantalla (se verifica sobre la pantalla, no con escenarios).

## Escenarios BDD

```gherkin
# language: es
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
```

Los pasos que no indican proyecto se refieren al de los Antecedentes, "Tienda Online". "El Sprint 1" y "el Sprint 2" son los Sprints del escenario por su número. Las horas se escriben con punto decimal, como en la spec 12. Los pasos de proyectos, integrantes, backlog, Sprints, Story Points, horas estimadas, esfuerzo y defectos son los de las specs 01 a 12; "que hoy es …" es el de la spec 03. "El PDF incluye el texto …" verifica que el texto esté escrito en el archivo, no cómo se ve. CA-8 se verifica sobre la pantalla.

## Pantalla

`frontend/src/paginas/Reporte.tsx` (sección "Reporte" del menú): un desplegable para elegir "Proyecto completo" o uno de los Sprints; abajo, el reporte elegido con sus secciones (encabezado, historias, Sprints o métricas del Sprint, métricas y defectos) y un botón "Descargar PDF" que baja el archivo armado por el backend. No calcula nada: muestra lo que devuelve la consulta del reporte y los mensajes de error del backend.

## Dependencias

- Specs 02, 03, 04, 05, 06, 07, 09, 10, 11 y 12: implementadas en los Sprints 1 a 3. Esta historia no espera a ninguna otra.
- No cambia el comportamiento de ninguna otra spec: solo lee sus datos y reutiliza sus cálculos.
- Spec 13 (Dashboard): independiente; el reporte no muestra gráficos.

## Tareas

- [x] `[P1]` Armado del reporte del proyecto y del de un Sprint en `backend/dominio/` (RN-3 a RN-7), reutilizando los cálculos de las specs 03, 10, 11 y 12 (TDD).
- [x] `[P1]` Casos de uso de los dos reportes, con la verificación de proyecto y Sprint (RN-1).
- [x] `[P1]` Generación del PDF en el backend (RN-8, RN-10), con la librería justificada en el PR.
- [x] `[P1]` Rutas REST de la vista previa y de la descarga del PDF, con los escenarios de `backend/features/reporte.feature`.
- [x] `[P1]` Pantalla "Reporte" en el frontend (CA-8).
- [x] `[P2]` Actualizar `backend/CLAUDE.md` y `frontend/CLAUDE.md`.

## Definition of Done

- [x] Tests unitarios Go escritos primero (commits `test:` → `feat:` → `refactor:` visibles en el historial).
- [x] Escenarios BDD automatizados y en verde.
- [x] CI en verde (gofmt / vet / tests / lint / build).
- [x] Criterios de aceptación verificados.
- [x] `CLAUDE.md` de la unidad actualizado (specs relacionadas, trampas).
- [x] Spec en estado `implemented`.
