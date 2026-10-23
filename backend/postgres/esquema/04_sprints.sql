-- Spec 05: Sprints de un proyecto y las historias asignadas a cada uno. Las validaciones de negocio
-- (fechas dentro del proyecto, sin superposición, una historia en un solo Sprint abierto) están en
-- el dominio (Go); acá solo se reflejan como restricciones de las tablas.
CREATE TABLE sprints (
    id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    proyecto_id  BIGINT       NOT NULL REFERENCES proyectos (id),
    numero       INTEGER      NOT NULL CHECK (numero >= 1),
    fecha_inicio DATE         NOT NULL,
    fecha_fin    DATE         NOT NULL,
    sprint_goal  VARCHAR(300) NOT NULL,
    estado       TEXT         NOT NULL CHECK (estado IN ('Abierto', 'Cerrado')),
    CHECK (fecha_fin > fecha_inicio),
    CONSTRAINT sprints_numero_unico UNIQUE (proyecto_id, numero)
);

-- Si una historia se elimina del backlog (solo se puede si no está en un Sprint abierto, spec 05
-- RN-11), sale también de los Sprints cerrados donde estuvo.
CREATE TABLE sprint_historias (
    sprint_id   BIGINT NOT NULL REFERENCES sprints (id) ON DELETE CASCADE,
    historia_id BIGINT NOT NULL REFERENCES historias (id) ON DELETE CASCADE,
    PRIMARY KEY (sprint_id, historia_id)
);
