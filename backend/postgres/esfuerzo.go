package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioEsfuerzo guarda los registros de esfuerzo en la tabla registros_esfuerzo.
type RepositorioEsfuerzo struct {
	db *sql.DB
}

func NuevoRepositorioEsfuerzo(db *sql.DB) *RepositorioEsfuerzo {
	return &RepositorioEsfuerzo{db: db}
}

// Guardar inserta el registro; el ID lo genera PostgreSQL (spec 09 RN-8).
func (r *RepositorioEsfuerzo) Guardar(ctx context.Context, e dominio.RegistroEsfuerzo) (dominio.RegistroEsfuerzo, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO registros_esfuerzo (proyecto_id, historia_id, integrante_id, fecha, actividad, centesimas)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id`,
		e.ProyectoID, e.HistoriaID, e.IntegranteID, e.Fecha, e.Actividad, e.Centesimas,
	).Scan(&e.ID)
	if err != nil {
		return dominio.RegistroEsfuerzo{}, fmt.Errorf("guardar registro de esfuerzo: %w", err)
	}
	e.IntegranteNombre = "" // el nombre no se guarda con el registro
	return e, nil
}

// Listar devuelve los registros de la historia. El orden por fecha lo aplica el dominio (RN-7).
func (r *RepositorioEsfuerzo) Listar(ctx context.Context, proyectoID, historiaID int64) ([]dominio.RegistroEsfuerzo, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, proyecto_id, historia_id, integrante_id, fecha, actividad, centesimas
		 FROM registros_esfuerzo
		 WHERE proyecto_id = $1 AND historia_id = $2
		 ORDER BY id`, proyectoID, historiaID)
	if err != nil {
		return nil, fmt.Errorf("listar esfuerzo: %w", err)
	}
	defer rows.Close()

	registros := []dominio.RegistroEsfuerzo{}
	for rows.Next() {
		var e dominio.RegistroEsfuerzo
		if err := rows.Scan(&e.ID, &e.ProyectoID, &e.HistoriaID, &e.IntegranteID, &e.Fecha, &e.Actividad, &e.Centesimas); err != nil {
			return nil, fmt.Errorf("leer registro de esfuerzo: %w", err)
		}
		registros = append(registros, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listar esfuerzo: %w", err)
	}
	return registros, nil
}

// Eliminar borra el registro de esa historia, o devuelve ErrRegistroEsfuerzoNoExiste (RN-8).
func (r *RepositorioEsfuerzo) Eliminar(ctx context.Context, proyectoID, historiaID, id int64) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM registros_esfuerzo WHERE id = $1 AND proyecto_id = $2 AND historia_id = $3`,
		id, proyectoID, historiaID)
	if err != nil {
		return fmt.Errorf("eliminar registro de esfuerzo %d: %w", id, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("eliminar registro de esfuerzo %d: %w", id, err)
	}
	if n == 0 {
		return dominio.ErrRegistroEsfuerzoNoExiste
	}
	return nil
}
