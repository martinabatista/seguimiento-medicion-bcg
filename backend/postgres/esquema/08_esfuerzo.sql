-- Spec 09: registros de esfuerzo sobre las historias del backlog. Las validaciones de negocio están
-- en el dominio (Go); acá solo se reflejan como restricciones de la tabla. VARCHAR(n) cuenta
-- caracteres, no bytes. Las horas se guardan en centésimas enteras (1,5 h = 150) para que las
-- sumas sean exactas.
-- La referencia a la historia NO tiene ON DELETE CASCADE: una historia con esfuerzo registrado no
-- se puede eliminar (RN-10), así las horas reales no se pierden.
CREATE TABLE registros_esfuerzo (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    proyecto_id   BIGINT       NOT NULL REFERENCES proyectos (id),
    historia_id   BIGINT       NOT NULL REFERENCES historias (id),
    integrante_id BIGINT       NOT NULL REFERENCES integrantes (id),
    fecha         DATE         NOT NULL,
    actividad     VARCHAR(300) NOT NULL,
    centesimas    INTEGER      NOT NULL CHECK (centesimas > 0 AND centesimas <= 2400)
);
