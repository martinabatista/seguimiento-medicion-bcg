-- Spec 12: horas estimadas de cada historia, en centésimas enteras (8 h = 800) para que las sumas sean
-- exactas, como las horas del esfuerzo (spec 09). NULL = sin horas estimadas.
ALTER TABLE historias ADD COLUMN horas_estimadas_centesimas INTEGER CHECK (horas_estimadas_centesimas > 0);
