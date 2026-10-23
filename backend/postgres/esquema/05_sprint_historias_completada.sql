-- Spec 06: cada historia asignada a un Sprint registra si se completó en ese Sprint. El registro no
-- cambia aunque después se modifique la historia en el backlog (spec 06 RN-8).
ALTER TABLE sprint_historias ADD COLUMN completada BOOLEAN NOT NULL DEFAULT false;
