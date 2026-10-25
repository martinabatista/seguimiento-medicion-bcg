-- Spec 08: votaciones de Planning Poker, sus rondas y sus votos. Las reglas de negocio están en el
-- dominio (Go); acá solo se reflejan como restricciones de las tablas.
-- Si se elimina una historia del backlog, se eliminan sus votaciones (y con ellas, rondas y votos).
CREATE TABLE votaciones (
    id                  BIGINT  GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    proyecto_id         BIGINT  NOT NULL REFERENCES proyectos (id),
    historia_id         BIGINT  NOT NULL REFERENCES historias (id) ON DELETE CASCADE,
    estado              TEXT    NOT NULL CHECK (estado IN ('Abierta', 'Cerrada')),
    estimacion_acordada INTEGER -- NULL mientras la votación está abierta (RN-9)
);

-- RN-1: una historia no puede tener dos votaciones abiertas a la vez.
CREATE UNIQUE INDEX votaciones_una_abierta_por_historia ON votaciones (historia_id) WHERE estado = 'Abierta';

CREATE TABLE rondas (
    votacion_id BIGINT  NOT NULL REFERENCES votaciones (id) ON DELETE CASCADE,
    numero      INTEGER NOT NULL CHECK (numero >= 1),
    revelada    BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (votacion_id, numero)
);

-- RN-4: un voto por Developer en cada ronda.
CREATE TABLE votos (
    votacion_id   BIGINT  NOT NULL,
    ronda         INTEGER NOT NULL,
    integrante_id BIGINT  NOT NULL REFERENCES integrantes (id),
    valor         INTEGER NOT NULL,
    PRIMARY KEY (votacion_id, ronda, integrante_id),
    FOREIGN KEY (votacion_id, ronda) REFERENCES rondas (votacion_id, numero) ON DELETE CASCADE
);
