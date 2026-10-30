# Manual de usuario — Software Metrics & Estimation

Guía breve para usar la aplicación: administrar un proyecto con Scrum y ver cómo avanza.
Cómo está construida: [documentación técnica](documentacion-tecnica.md).

## Antes de empezar

La aplicación se abre en el navegador, en **http://localhost:5173**. Para que funcione tienen que
estar andando la base de datos, el backend y el frontend (ver [Levantar el proyecto](../README.md#levantar-el-proyecto)).

Si algo no se puede hacer, la pantalla lo avisa con un mensaje en rojo que explica por qué (por
ejemplo, "El título es obligatorio"). Corregí el dato y volvé a intentar.

## Recorrido recomendado

```
1. Proyectos      → crear el proyecto y elegirlo
2. Integrantes    → cargar al equipo con su rol
3. Backlog        → cargar las historias de usuario
4. Planning Poker → estimar las historias en equipo
5. Sprints        → crear el Sprint y asignarle historias
6. Esfuerzo       → registrar las horas trabajadas
7. Defectos       → registrar los errores encontrados
8. Sprints        → completar historias y cerrar el Sprint
9. Estado · Métricas · Dashboard · Reporte → ver cómo va el proyecto
```

## 1. Proyectos

Es la primera pantalla.

- **Crear:** completá nombre, descripción, fecha de inicio y fecha de finalización, y tocá **Crear proyecto**.
- **Modificar:** tocá **Modificar** en el proyecto, cambiá los datos y tocá **Guardar cambios** (o **Cancelar**).
- **Elegir:** tocá **Elegir** para entrar al proyecto. Arriba aparece su nombre y el menú de secciones.
  Con **Cambiar de proyecto** volvés a la lista.

## 2. Integrantes

Completá nombre, email y rol (**Product Owner**, **Scrum Master** o **Developer**) y tocá **Agregar integrante**.
El email no se puede repetir dentro del proyecto. Solo los Developers votan en el Planning Poker.

## 3. Backlog

El Product Backlog: las historias de usuario, ordenadas por prioridad (1 es la más importante).

- **Agregar:** completá título, descripción, prioridad y, si querés, los criterios de aceptación. Tocá
  **Agregar historia**. Toda historia arranca en estado **Pendiente**.
- **Modificar:** tocá **Modificar**. Además de los datos podés cambiar el estado. Una historia que está
  en un Sprint abierto se completa desde la sección Sprints.
- **Estimar a mano:** elegí los Story Points en la lista de la historia (1, 2, 3, 5, 8, 13 o 21) o
  **Sin estimar** para quitarlos. Una historia completada ya no cambia su estimación.
- **Eliminar:** tocá **Eliminar** y confirmá. No se puede si la historia está en un Sprint abierto o
  tiene horas o defectos registrados.

## 4. Planning Poker

Para que el equipo estime una historia votando en secreto.

1. En **Historia para estimar**, elegí la historia y tocá **Abrir votación**. El proyecto tiene que
   tener al menos un Developer.
2. Cada Developer elige su carta y toca **Votar**. Los votos no se muestran: solo se ve quién **ya votó**.
   Si alguien vota de nuevo antes de revelar, vale el último voto.
3. Cuando votaron todos, tocá **Revelar votos**. Se ven los valores y, si hay diferencias, quién votó
   el más bajo y el más alto, para que lo expliquen.
4. Si no hay acuerdo, tocá **Nueva ronda** y se vuelve a votar.
5. Cuando se ponen de acuerdo, elegí la **Estimación acordada** y tocá **Acordar**. La historia queda
   estimada con ese valor y la votación se cierra.

Una votación que quedó sin terminar aparece en **Votaciones abiertas**: tocá **Continuar** para retomarla.

## 5. Sprints

- **Crear:** completá fecha de inicio, fecha de fin y Sprint Goal (el objetivo), y tocá **Crear Sprint**.
  Dos Sprints no pueden superponerse en fechas.
- **Ver:** tocá **Ver** en un Sprint para entrar a su detalle. Con **← Volver a los Sprints** volvés a la lista.

Dentro de un Sprint abierto:

- **Sprint Goal:** cambialo y tocá **Guardar Sprint Goal**.
- **Asignar historias:** elegí una historia del backlog y tocá **Asignar**. **Quitar** la saca del Sprint.
- **Completar:** tocá **Completar** en cada historia terminada. **Deshacer** la vuelve a "En progreso".
- **Cerrar:** al terminar el Sprint, tocá **Cerrar Sprint** y confirmá. Lo planificado y lo completado
  quedan registrados y ya no se pueden cambiar. Las historias que no se completaron quedan libres para otro Sprint.

## 6. Esfuerzo

Para registrar las horas que trabajó cada integrante.

1. Elegí la **Historia**. Abajo se ven sus registros por fecha y el total de horas.
2. Completá integrante, fecha, actividad realizada y horas trabajadas (con coma, por ejemplo **1,5**), y tocá **Registrar**.
3. Si te equivocaste, tocá **Eliminar** en el registro y confirmá.

## 7. Defectos

- **Registrar:** completá la descripción, elegí la severidad (**Baja**, **Media**, **Alta** o
  **Crítica**), la historia relacionada y el Sprint en que se detectó. Tocá **Registrar**. Hace falta
  tener al menos un Sprint creado.
- **Seguimiento:** en cada defecto elegí el estado nuevo (**Abierto**, **En corrección** o **Resuelto**)
  y tocá **Guardar**. Al pasarlo a Resuelto, elegí también en qué Sprint se resolvió.
- **Métricas de defectos:** cuántos se detectaron y cuántos se resolvieron, en total y por Sprint.

## 8. Ver cómo va el proyecto

| Sección | Qué muestra |
|---|---|
| **Estado** | Si el proyecto está No iniciado, En curso o Finalizado (según las fechas), el Sprint en curso y el avance del backlog: historias por estado y Story Points |
| **Métricas** | Por Sprint: historias y Story Points planificados y completados. Del proyecto: la **velocidad** del equipo (promedio de Story Points completados por Sprint cerrado). Y las métricas de esfuerzo: horas estimadas, horas reales y desviación |
| **Dashboard** | Todo de un vistazo: tarjetas con los indicadores principales y tres gráficos (Story Points por Sprint, velocidad del equipo y backlog por estado) |
| **Reporte** | El resumen del proyecto completo o de un Sprint, para descargar en PDF |

**Horas estimadas:** se cargan en **Métricas**, en la tabla de esfuerzo por historia: escribí las horas
y tocá **Guardar** (**Quitar** las borra). Con eso la app calcula la desviación entre lo estimado y lo
real. Una desviación positiva quiere decir que se trabajó más de lo estimado.

**Reporte en PDF:** en **Reporte de**, elegí **Proyecto completo** o un Sprint. Se ve una vista previa;
tocá **Descargar PDF** para guardarlo.

## Preguntas frecuentes

- **No puedo eliminar una historia.** Está en un Sprint abierto o tiene horas o defectos registrados.
- **No puedo revelar los votos.** Todavía falta votar algún Developer: la lista muestra quién ya votó.
- **No puedo registrar un defecto.** Primero hay que crear un Sprint en la sección Sprints.
- **La velocidad aparece vacía.** Se calcula con Sprints cerrados: aparece al cerrar el primero.
- **Un Sprint cerrado no me deja cambiar nada.** Es así a propósito: lo cerrado queda como registro.
