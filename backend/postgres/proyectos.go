// Package postgres implementa los repositorios del dominio sobre PostgreSQL,
// con database/sql y el driver pgx. Las tablas están en esquema/*.sql.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // registra el driver "pgx" para database/sql

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// Abrir conecta con la base indicada en url (DATABASE_URL) y verifica que responda.
func Abrir(url string) (*sql.DB, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("abrir la base: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("conectar con la base: %w", err)
	}
	return db, nil
}

// RepositorioProyectos guarda los proyectos en la tabla proyectos.
type RepositorioProyectos struct {
	db *sql.DB
}

func NuevoRepositorioProyectos(db *sql.DB) *RepositorioProyectos {
	return &RepositorioProyectos{db: db}
}

// Guardar inserta el proyecto; el ID lo genera PostgreSQL (spec 01 RN-6).
func (r *RepositorioProyectos) Guardar(ctx context.Context, p dominio.Proyecto) (dominio.Proyecto, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO proyectos (nombre, descripcion, fecha_inicio, fecha_fin)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		p.Nombre, p.Descripcion, p.FechaInicio, p.FechaFin,
	).Scan(&p.ID)
	if err != nil {
		return dominio.Proyecto{}, fmt.Errorf("guardar proyecto: %w", err)
	}
	return p, nil
}

// Actualizar reemplaza los datos del proyecto, o devuelve ErrProyectoNoExiste.
func (r *RepositorioProyectos) Actualizar(ctx context.Context, p dominio.Proyecto) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE proyectos
		 SET nombre = $1, descripcion = $2, fecha_inicio = $3, fecha_fin = $4
		 WHERE id = $5`,
		p.Nombre, p.Descripcion, p.FechaInicio, p.FechaFin, p.ID,
	)
	if err != nil {
		return fmt.Errorf("actualizar proyecto %d: %w", p.ID, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("actualizar proyecto %d: %w", p.ID, err)
	}
	if n == 0 {
		return dominio.ErrProyectoNoExiste
	}
	return nil
}

// Listar devuelve todos los proyectos en orden de creación (el ID es creciente).
func (r *RepositorioProyectos) Listar(ctx context.Context) ([]dominio.Proyecto, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, nombre, descripcion, fecha_inicio, fecha_fin
		 FROM proyectos
		 ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("listar proyectos: %w", err)
	}
	defer rows.Close()

	proyectos := []dominio.Proyecto{}
	for rows.Next() {
		var p dominio.Proyecto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.FechaInicio, &p.FechaFin); err != nil {
			return nil, fmt.Errorf("leer proyecto: %w", err)
		}
		proyectos = append(proyectos, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listar proyectos: %w", err)
	}
	return proyectos, nil
}

// Obtener devuelve el proyecto con ese ID, o ErrProyectoNoExiste.
func (r *RepositorioProyectos) Obtener(ctx context.Context, id int64) (dominio.Proyecto, error) {
	var p dominio.Proyecto
	err := r.db.QueryRowContext(ctx,
		`SELECT id, nombre, descripcion, fecha_inicio, fecha_fin
		 FROM proyectos
		 WHERE id = $1`, id,
	).Scan(&p.ID, &p.Nombre, &p.Descripcion, &p.FechaInicio, &p.FechaFin)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Proyecto{}, dominio.ErrProyectoNoExiste
	}
	if err != nil {
		return dominio.Proyecto{}, fmt.Errorf("obtener proyecto %d: %w", id, err)
	}
	return p, nil
}
