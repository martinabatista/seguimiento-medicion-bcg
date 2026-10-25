-- Spec 02: integrantes de un proyecto. Las validaciones de negocio están en el dominio (Go);
-- acá solo se reflejan como restricciones de la tabla. VARCHAR(n) cuenta caracteres, no bytes.
-- El email llega normalizado (sin espacios al borde y en minúsculas, RN-3), así que la
-- restricción UNIQUE alcanza para que no se repita en el proyecto (RN-6).
CREATE TABLE integrantes (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    proyecto_id BIGINT       NOT NULL REFERENCES proyectos (id),
    nombre      VARCHAR(100) NOT NULL,
    email       VARCHAR(254) NOT NULL,
    rol         TEXT         NOT NULL CHECK (rol IN ('Product Owner', 'Scrum Master', 'Developer')),
    CONSTRAINT integrantes_email_unico UNIQUE (proyecto_id, email)
);
