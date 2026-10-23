package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioSprints guarda los Sprints en la tabla sprints y sus historias en sprint_historias.
type RepositorioSprints struct {
	db *sql.DB
}

func NuevoRepositorioSprints(db *sql.DB) *RepositorioSprints {
	return &RepositorioSprints{db: db}
}

// Guardar inserta el Sprint; el ID lo genera PostgreSQL (spec 05 RN-13). Un Sprint nuevo no tiene
// historias (RN-6).
func (r *RepositorioSprints) Guardar(ctx context.Context, s dominio.Sprint) (dominio.Sprint, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO sprints (proyecto_id, numero, fecha_inicio, fecha_fin, sprint_goal, estado)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id`,
		s.ProyectoID, s.Numero, s.FechaInicio, s.FechaFin, s.Goal, string(s.Estado),
	).Scan(&s.ID)
	if err != nil {
		return dominio.Sprint{}, fmt.Errorf("guardar Sprint: %w", err)
	}
	if s.Historias == nil {
		s.Historias = []int64{}
	}
	if s.Completadas == nil {
		s.Completadas = []int64{}
	}
	return s, nil
}

// Actualizar reemplaza el Sprint Goal, el estado y las historias asignadas (con cuáles se completaron
// en el Sprint, spec 06), todo junto en una transacción. Devuelve ErrSprintNoExiste si el Sprint no es
// de ese proyecto.
func (r *RepositorioSprints) Actualizar(ctx context.Context, s dominio.Sprint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("actualizar Sprint %d: %w", s.ID, err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx,
		`UPDATE sprints SET sprint_goal = $1, estado = $2 WHERE id = $3 AND proyecto_id = $4`,
		s.Goal, string(s.Estado), s.ID, s.ProyectoID)
	if err != nil {
		return fmt.Errorf("actualizar Sprint %d: %w", s.ID, err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("actualizar Sprint %d: %w", s.ID, err)
	} else if n == 0 {
		return dominio.ErrSprintNoExiste
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM sprint_historias WHERE sprint_id = $1`, s.ID); err != nil {
		return fmt.Errorf("actualizar historias del Sprint %d: %w", s.ID, err)
	}
	for _, historiaID := range s.Historias {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO sprint_historias (sprint_id, historia_id, completada, story_points_al_cerrar) VALUES ($1, $2, $3, $4)`,
			s.ID, historiaID, s.EstaCompletada(historiaID), s.StoryPointsAlCerrar[historiaID]); err != nil {
			return fmt.Errorf("asignar historia %d al Sprint %d: %w", historiaID, s.ID, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("actualizar Sprint %d: %w", s.ID, err)
	}
	return nil
}

// Listar devuelve los Sprints del proyecto con sus historias. El orden por fecha de inicio lo
// aplica el dominio.
func (r *RepositorioSprints) Listar(ctx context.Context, proyectoID int64) ([]dominio.Sprint, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, proyecto_id, numero, fecha_inicio, fecha_fin, sprint_goal, estado
		 FROM sprints
		 WHERE proyecto_id = $1
		 ORDER BY id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar Sprints: %w", err)
	}
	defer rows.Close()

	sprints := []dominio.Sprint{}
	for rows.Next() {
		s, err := escanearSprint(rows)
		if err != nil {
			return nil, err
		}
		sprints = append(sprints, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listar Sprints: %w", err)
	}
	for i := range sprints {
		if err := r.cargarHistorias(ctx, &sprints[i]); err != nil {
			return nil, err
		}
	}
	return sprints, nil
}

// Obtener devuelve el Sprint de ese proyecto con sus historias, o ErrSprintNoExiste.
func (r *RepositorioSprints) Obtener(ctx context.Context, proyectoID, id int64) (dominio.Sprint, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, proyecto_id, numero, fecha_inicio, fecha_fin, sprint_goal, estado
		 FROM sprints
		 WHERE id = $1 AND proyecto_id = $2`, id, proyectoID)
	s, err := escanearSprint(row)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Sprint{}, dominio.ErrSprintNoExiste
	}
	if err != nil {
		return dominio.Sprint{}, err
	}
	if err := r.cargarHistorias(ctx, &s); err != nil {
		return dominio.Sprint{}, err
	}
	return s, nil
}

// cargarHistorias completa el Sprint con sus historias asignadas, las que se completaron en él
// (spec 06) y, si está cerrado, los Story Points congelados al cerrarlo (spec 10 RN-4).
func (r *RepositorioSprints) cargarHistorias(ctx context.Context, s *dominio.Sprint) error {
	rows, err := r.db.QueryContext(ctx,
		`SELECT historia_id, completada, story_points_al_cerrar
		 FROM sprint_historias WHERE sprint_id = $1 ORDER BY historia_id`, s.ID)
	if err != nil {
		return fmt.Errorf("historias del Sprint %d: %w", s.ID, err)
	}
	defer rows.Close()
	s.Historias, s.Completadas = []int64{}, []int64{}
	congelados := map[int64]*int{}
	for rows.Next() {
		var id int64
		var completada bool
		var puntos sql.NullInt64
		if err := rows.Scan(&id, &completada, &puntos); err != nil {
			return fmt.Errorf("historias del Sprint %d: %w", s.ID, err)
		}
		s.Historias = append(s.Historias, id)
		if completada {
			s.Completadas = append(s.Completadas, id)
		}
		congelados[id] = nil
		if puntos.Valid {
			p := int(puntos.Int64)
			congelados[id] = &p
		}
	}
	if s.Estado == dominio.SprintCerrado {
		s.StoryPointsAlCerrar = congelados
	}
	return rows.Err()
}

func escanearSprint(f fila) (dominio.Sprint, error) {
	var s dominio.Sprint
	var estado string
	if err := f.Scan(&s.ID, &s.ProyectoID, &s.Numero, &s.FechaInicio, &s.FechaFin, &s.Goal, &estado); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s, err
		}
		return s, fmt.Errorf("leer Sprint: %w", err)
	}
	s.Estado = dominio.EstadoSprint(estado)
	return s, nil
}
