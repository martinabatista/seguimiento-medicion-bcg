-- Spec 11: defectos encontrados en las historias del backlog. Las validaciones de negocio están en
-- el dominio (Go); acá solo se reflejan como restricciones de la tabla. VARCHAR(n) cuenta
-- caracteres, no bytes. sprint_resolucion_id es NULL mientras el defecto no está Resuelto (RN-9).
-- La referencia a la historia NO tiene ON DELETE CASCADE: una historia con defectos registrados no
-- se puede eliminar (RN-14), así los defectos no se pierden.
CREATE TABLE defectos (
    id                   BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    proyecto_id          BIGINT        NOT NULL REFERENCES proyectos (id),
    historia_id          BIGINT        NOT NULL REFERENCES historias (id),
    descripcion          VARCHAR(1000) NOT NULL,
    severidad            VARCHAR(10)   NOT NULL CHECK (severidad IN ('Baja', 'Media', 'Alta', 'Crítica')),
    estado               VARCHAR(20)   NOT NULL CHECK (estado IN ('Abierto', 'En corrección', 'Resuelto')),
    sprint_deteccion_id  BIGINT        NOT NULL REFERENCES sprints (id),
    sprint_resolucion_id BIGINT        REFERENCES sprints (id),
    CHECK ((estado = 'Resuelto') = (sprint_resolucion_id IS NOT NULL))
);
