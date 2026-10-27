package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioHistorias guarda las historias del backlog en la tabla historias.
type RepositorioHistorias struct {
	db *sql.DB
}

func NuevoRepositorioHistorias(db *sql.DB) *RepositorioHistorias {
	return &RepositorioHistorias{db: db}
}

// Guardar inserta la historia; el ID lo genera PostgreSQL (spec 04 RN-10).
func (r *RepositorioHistorias) Guardar(ctx context.Context, h dominio.Historia) (dominio.Historia, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO historias (proyecto_id, titulo, descripcion, prioridad, criterios, estado, story_points, horas_estimadas_centesimas)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id`,
		h.ProyectoID, h.Titulo, h.Descripcion, h.Prioridad, h.Criterios, string(h.Estado), h.StoryPoints, h.HorasEstimadas,
	).Scan(&h.ID)
	if err != nil {
		return dominio.Historia{}, fmt.Errorf("guardar historia: %w", err)
	}
	return h, nil
}

// Actualizar reemplaza los datos de la historia de ese proyecto, o devuelve ErrHistoriaNoExiste.
func (r *RepositorioHistorias) Actualizar(ctx context.Context, h dominio.Historia) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE historias
		 SET titulo = $1, descripcion = $2, prioridad = $3, criterios = $4, estado = $5, story_points = $6, horas_estimadas_centesimas = $7
		 WHERE id = $8 AND proyecto_id = $9`,
		h.Titulo, h.Descripcion, h.Prioridad, h.Criterios, string(h.Estado), h.StoryPoints, h.HorasEstimadas, h.ID, h.ProyectoID,
	)
	if err != nil {
		return fmt.Errorf("actualizar historia %d: %w", h.ID, err)
	}
	return unaFilaAfectada(res, fmt.Sprintf("actualizar historia %d", h.ID))
}

// Eliminar borra la historia de ese proyecto, o devuelve ErrHistoriaNoExiste.
func (r *RepositorioHistorias) Eliminar(ctx context.Context, proyectoID, id int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM historias WHERE id = $1 AND proyecto_id = $2`, id, proyectoID)
	if err != nil {
		return fmt.Errorf("eliminar historia %d: %w", id, err)
	}
	return unaFilaAfectada(res, fmt.Sprintf("eliminar historia %d", id))
}

// Listar devuelve las historias del proyecto en orden de creación. El orden del backlog
// (prioridad y desempate) lo aplica el dominio (spec 04 RN-8).
func (r *RepositorioHistorias) Listar(ctx context.Context, proyectoID int64) ([]dominio.Historia, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, proyecto_id, titulo, descripcion, prioridad, criterios, estado, story_points, horas_estimadas_centesimas
		 FROM historias
		 WHERE proyecto_id = $1
		 ORDER BY id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar historias: %w", err)
	}
	defer rows.Close()

	historias := []dominio.Historia{}
	for rows.Next() {
		h, err := leerHistoria(rows)
		if err != nil {
			return nil, err
		}
		historias = append(historias, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listar historias: %w", err)
	}
	return historias, nil
}

// Obtener devuelve la historia de ese proyecto, o ErrHistoriaNoExiste (spec 04 RN-10).
func (r *RepositorioHistorias) Obtener(ctx context.Context, proyectoID, id int64) (dominio.Historia, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, proyecto_id, titulo, descripcion, prioridad, criterios, estado, story_points, horas_estimadas_centesimas
		 FROM historias
		 WHERE id = $1 AND proyecto_id = $2`, id, proyectoID)
	h, err := leerHistoria(row)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Historia{}, dominio.ErrHistoriaNoExiste
	}
	return h, err
}

// fila es lo común entre *sql.Row y *sql.Rows para leer un registro (historias y Sprints).
type fila interface {
	Scan(dest ...any) error
}

// leerHistoria arma una historia desde una fila; story_points NULL queda como sin estimar (nil) y
// horas_estimadas_centesimas NULL, como sin horas estimadas (nil).
func leerHistoria(f fila) (dominio.Historia, error) {
	var h dominio.Historia
	var estado string
	var sp, horas sql.NullInt64
	err := f.Scan(&h.ID, &h.ProyectoID, &h.Titulo, &h.Descripcion, &h.Prioridad, &h.Criterios, &estado, &sp, &horas)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Historia{}, err
	}
	if err != nil {
		return dominio.Historia{}, fmt.Errorf("leer historia: %w", err)
	}
	h.Estado = dominio.EstadoHistoria(estado)
	if sp.Valid {
		v := int(sp.Int64)
		h.StoryPoints = &v
	}
	if horas.Valid {
		v := int(horas.Int64)
		h.HorasEstimadas = &v
	}
	return h, nil
}

// unaFilaAfectada traduce "no se tocó ninguna fila" en ErrHistoriaNoExiste.
func unaFilaAfectada(res sql.Result, operacion string) error {
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", operacion, err)
	}
	if n == 0 {
		return dominio.ErrHistoriaNoExiste
	}
	return nil
}
