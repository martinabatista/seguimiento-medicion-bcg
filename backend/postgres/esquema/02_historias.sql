-- Spec 04: historias del Product Backlog. Las validaciones de negocio están en el dominio (Go);
-- acá solo se reflejan como restricciones de la tabla. VARCHAR(n) cuenta caracteres, no bytes.
CREATE TABLE historias (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    proyecto_id  BIGINT        NOT NULL REFERENCES proyectos (id),
    titulo       VARCHAR(150)  NOT NULL,
    descripcion  VARCHAR(1000) NOT NULL,
    prioridad    INTEGER       NOT NULL CHECK (prioridad >= 1),
    criterios    VARCHAR(2000) NOT NULL DEFAULT '',
    estado       TEXT          NOT NULL CHECK (estado IN ('Pendiente', 'En progreso', 'Completada')),
    story_points INTEGER       -- NULL = sin estimar (se asignan en las specs 07 y 08)
);
