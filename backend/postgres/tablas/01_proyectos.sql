-- Spec 01: proyectos. Las validaciones de negocio están en el dominio (Go);
-- acá solo se reflejan como restricciones de la tabla.
CREATE TABLE proyectos (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nombre       TEXT         NOT NULL,
    descripcion  VARCHAR(500) NOT NULL,
    fecha_inicio DATE         NOT NULL,
    fecha_fin    DATE         NOT NULL,
    CHECK (fecha_fin > fecha_inicio)
);
