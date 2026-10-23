package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// codigoViolacionUnica es el código de PostgreSQL para "se violó una restricción UNIQUE".
const codigoViolacionUnica = "23505"

// RepositorioIntegrantes guarda los integrantes en la tabla integrantes.
type RepositorioIntegrantes struct {
	db *sql.DB
}

func NuevoRepositorioIntegrantes(db *sql.DB) *RepositorioIntegrantes {
	return &RepositorioIntegrantes{db: db}
}

// Guardar inserta el integrante; el ID lo genera PostgreSQL (spec 02 RN-8). Si el email ya está
// en el proyecto, la restricción integrantes_email_unico lo rechaza (RN-6).
func (r *RepositorioIntegrantes) Guardar(ctx context.Context, i dominio.Integrante) (dominio.Integrante, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO integrantes (proyecto_id, nombre, email, rol)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		i.ProyectoID, i.Nombre, i.Email, string(i.Rol),
	).Scan(&i.ID)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == codigoViolacionUnica && pgErr.ConstraintName == "integrantes_email_unico" {
		return dominio.Integrante{}, dominio.ErrIntegranteYaRegistrado
	}
	if err != nil {
		return dominio.Integrante{}, fmt.Errorf("guardar integrante: %w", err)
	}
	return i, nil
}

// Listar devuelve los integrantes del proyecto en orden de alta (RN-7): el id es creciente.
func (r *RepositorioIntegrantes) Listar(ctx context.Context, proyectoID int64) ([]dominio.Integrante, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, proyecto_id, nombre, email, rol
		 FROM integrantes
		 WHERE proyecto_id = $1
		 ORDER BY id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar integrantes: %w", err)
	}
	defer rows.Close()

	integrantes := []dominio.Integrante{}
	for rows.Next() {
		i, err := leerIntegrante(rows)
		if err != nil {
			return nil, err
		}
		integrantes = append(integrantes, i)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listar integrantes: %w", err)
	}
	return integrantes, nil
}

// Obtener devuelve el integrante de ese proyecto, o ErrIntegranteNoExiste (spec 02 RN-8).
func (r *RepositorioIntegrantes) Obtener(ctx context.Context, proyectoID, id int64) (dominio.Integrante, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, proyecto_id, nombre, email, rol
		 FROM integrantes
		 WHERE id = $1 AND proyecto_id = $2`, id, proyectoID)
	i, err := leerIntegrante(row)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Integrante{}, dominio.ErrIntegranteNoExiste
	}
	return i, err
}

// leerIntegrante arma un integrante desde una fila.
func leerIntegrante(f fila) (dominio.Integrante, error) {
	var i dominio.Integrante
	var rol string
	err := f.Scan(&i.ID, &i.ProyectoID, &i.Nombre, &i.Email, &rol)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Integrante{}, err
	}
	if err != nil {
		return dominio.Integrante{}, fmt.Errorf("leer integrante: %w", err)
	}
	i.Rol = dominio.RolIntegrante(rol)
	return i, nil
}
