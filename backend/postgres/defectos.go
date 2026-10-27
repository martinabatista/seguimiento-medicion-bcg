package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioDefectos guarda los defectos en la tabla defectos. El Sprint de resolución es NULL
// mientras el defecto no está Resuelto; en el dominio, 0.
type RepositorioDefectos struct {
	db *sql.DB
}

func NuevoRepositorioDefectos(db *sql.DB) *RepositorioDefectos {
	return &RepositorioDefectos{db: db}
}

const columnasDefecto = `id, proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id`

// Guardar inserta el defecto; el ID lo genera PostgreSQL (spec 11 RN-12).
func (r *RepositorioDefectos) Guardar(ctx context.Context, d dominio.Defecto) (dominio.Defecto, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO defectos (proyecto_id, historia_id, descripcion, severidad, estado, sprint_deteccion_id, sprint_resolucion_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		d.ProyectoID, d.HistoriaID, d.Descripcion, string(d.Severidad), string(d.Estado), d.SprintDeteccionID, sprintONulo(d.SprintResolucionID),
	).Scan(&d.ID)
	if err != nil {
		return dominio.Defecto{}, fmt.Errorf("guardar defecto: %w", err)
	}
	// Lo que se completa al consultar no se guarda con el defecto.
	d.HistoriaTitulo, d.SprintDeteccionNumero, d.SprintResolucionNumero = "", 0, 0
	return d, nil
}

// Actualizar guarda el estado y el Sprint de resolución: lo demás no cambia (spec 11 Alcance).
func (r *RepositorioDefectos) Actualizar(ctx context.Context, d dominio.Defecto) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE defectos SET estado = $1, sprint_resolucion_id = $2 WHERE id = $3 AND proyecto_id = $4`,
		string(d.Estado), sprintONulo(d.SprintResolucionID), d.ID, d.ProyectoID)
	if err != nil {
		return fmt.Errorf("actualizar defecto %d: %w", d.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("actualizar defecto %d: %w", d.ID, err)
	}
	if n == 0 {
		return dominio.ErrDefectoNoExiste
	}
	return nil
}

// Listar devuelve los defectos del proyecto en el orden en que se registraron (el ID es creciente).
func (r *RepositorioDefectos) Listar(ctx context.Context, proyectoID int64) ([]dominio.Defecto, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+columnasDefecto+` FROM defectos WHERE proyecto_id = $1 ORDER BY id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar defectos: %w", err)
	}
	defer rows.Close()

	defectos := []dominio.Defecto{}
	for rows.Next() {
		d, err := leerDefecto(rows)
		if err != nil {
			return nil, err
		}
		defectos = append(defectos, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listar defectos: %w", err)
	}
	return defectos, nil
}

// Obtener devuelve el defecto del proyecto, o ErrDefectoNoExiste (RN-12).
func (r *RepositorioDefectos) Obtener(ctx context.Context, proyectoID, id int64) (dominio.Defecto, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT `+columnasDefecto+` FROM defectos WHERE id = $1 AND proyecto_id = $2`, id, proyectoID)
	d, err := leerDefecto(row)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Defecto{}, dominio.ErrDefectoNoExiste
	}
	return d, err
}

// leerDefecto lee una fila con columnasDefecto, de un QueryRow o de un Query.
func leerDefecto(fila interface{ Scan(...any) error }) (dominio.Defecto, error) {
	var d dominio.Defecto
	var severidad, estado string
	var resolucion sql.NullInt64
	if err := fila.Scan(&d.ID, &d.ProyectoID, &d.HistoriaID, &d.Descripcion, &severidad, &estado, &d.SprintDeteccionID, &resolucion); err != nil {
		return dominio.Defecto{}, fmt.Errorf("leer defecto: %w", err) // %w: Obtener reconoce sql.ErrNoRows
	}
	d.Severidad = dominio.Severidad(severidad)
	d.Estado = dominio.EstadoDefecto(estado)
	d.SprintResolucionID = resolucion.Int64 // 0 si es NULL
	return d, nil
}

// sprintONulo guarda "sin Sprint de resolución" (0 en el dominio) como NULL.
func sprintONulo(id int64) sql.NullInt64 {
	return sql.NullInt64{Int64: id, Valid: id != 0}
}
