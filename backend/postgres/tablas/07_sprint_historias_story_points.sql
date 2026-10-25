-- Spec 10: al cerrar un Sprint se congelan los Story Points de cada historia planificada, para que sus
-- métricas no cambien si después se re-estima (RN-4). NULL = sin estimar (o Sprint todavía abierto).
ALTER TABLE sprint_historias ADD COLUMN story_points_al_cerrar INTEGER;
