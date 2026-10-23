package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/martinabatista/seguimiento-medicion-bcg/backend/dominio"
)

// RepositorioVotaciones guarda las votaciones de Planning Poker en las tablas votaciones, rondas y
// votos (spec 08).
type RepositorioVotaciones struct {
	db *sql.DB
}

func NuevoRepositorioVotaciones(db *sql.DB) *RepositorioVotaciones {
	return &RepositorioVotaciones{db: db}
}

// Guardar inserta la votación con sus rondas en una transacción; el ID lo genera PostgreSQL.
func (r *RepositorioVotaciones) Guardar(ctx context.Context, v dominio.Votacion) (dominio.Votacion, error) {
	err := enTransaccion(ctx, r.db, func(tx *sql.Tx) error {
		err := tx.QueryRowContext(ctx,
			`INSERT INTO votaciones (proyecto_id, historia_id, estado, estimacion_acordada)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			v.ProyectoID, v.HistoriaID, string(v.Estado), v.Acordada,
		).Scan(&v.ID)
		if err != nil {
			return err
		}
		return insertarRondas(ctx, tx, v)
	})
	if err != nil {
		return dominio.Votacion{}, fmt.Errorf("guardar votación: %w", err)
	}
	return v, nil
}

// Actualizar reemplaza el estado, la acordada y las rondas con sus votos, o devuelve
// ErrVotacionNoExiste. Todo en una transacción: o cambia todo, o nada.
func (r *RepositorioVotaciones) Actualizar(ctx context.Context, v dominio.Votacion) error {
	err := enTransaccion(ctx, r.db, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx,
			`UPDATE votaciones SET estado = $1, estimacion_acordada = $2
			 WHERE id = $3 AND proyecto_id = $4`,
			string(v.Estado), v.Acordada, v.ID, v.ProyectoID)
		if err != nil {
			return err
		}
		if n, err := res.RowsAffected(); err != nil {
			return err
		} else if n == 0 {
			return dominio.ErrVotacionNoExiste
		}
		// Las rondas se reescriben enteras (los votos se borran en cascada).
		if _, err := tx.ExecContext(ctx, `DELETE FROM rondas WHERE votacion_id = $1`, v.ID); err != nil {
			return err
		}
		return insertarRondas(ctx, tx, v)
	})
	if errors.Is(err, dominio.ErrVotacionNoExiste) {
		return err
	}
	if err != nil {
		return fmt.Errorf("actualizar votación %d: %w", v.ID, err)
	}
	return nil
}

// Listar devuelve las votaciones del proyecto con sus rondas, en orden de creación.
func (r *RepositorioVotaciones) Listar(ctx context.Context, proyectoID int64) ([]dominio.Votacion, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, proyecto_id, historia_id, estado, estimacion_acordada
		 FROM votaciones
		 WHERE proyecto_id = $1
		 ORDER BY id`, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("listar votaciones: %w", err)
	}
	votaciones := []dominio.Votacion{}
	for rows.Next() {
		v, err := leerVotacion(rows)
		if err != nil {
			rows.Close()
			return nil, err
		}
		votaciones = append(votaciones, v)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listar votaciones: %w", err)
	}
	for k := range votaciones {
		if votaciones[k].Rondas, err = r.rondas(ctx, votaciones[k].ID); err != nil {
			return nil, err
		}
	}
	return votaciones, nil
}

// Obtener devuelve la votación del proyecto con sus rondas, o ErrVotacionNoExiste (spec 08 RN-11).
func (r *RepositorioVotaciones) Obtener(ctx context.Context, proyectoID, id int64) (dominio.Votacion, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, proyecto_id, historia_id, estado, estimacion_acordada
		 FROM votaciones
		 WHERE id = $1 AND proyecto_id = $2`, id, proyectoID)
	v, err := leerVotacion(row)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Votacion{}, dominio.ErrVotacionNoExiste
	}
	if err != nil {
		return dominio.Votacion{}, err
	}
	if v.Rondas, err = r.rondas(ctx, v.ID); err != nil {
		return dominio.Votacion{}, err
	}
	return v, nil
}

// rondas lee las rondas de una votación en orden, cada una con sus votos.
func (r *RepositorioVotaciones) rondas(ctx context.Context, votacionID int64) ([]dominio.Ronda, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT r.numero, r.revelada, v.integrante_id, v.valor
		 FROM rondas r
		 LEFT JOIN votos v ON v.votacion_id = r.votacion_id AND v.ronda = r.numero
		 WHERE r.votacion_id = $1
		 ORDER BY r.numero, v.integrante_id`, votacionID)
	if err != nil {
		return nil, fmt.Errorf("leer rondas de la votación %d: %w", votacionID, err)
	}
	defer rows.Close()

	var rondas []dominio.Ronda
	for rows.Next() {
		var numero int
		var revelada bool
		var integranteID, valor sql.NullInt64
		if err := rows.Scan(&numero, &revelada, &integranteID, &valor); err != nil {
			return nil, fmt.Errorf("leer ronda: %w", err)
		}
		if len(rondas) == 0 || rondas[len(rondas)-1].Numero != numero {
			rondas = append(rondas, dominio.Ronda{Numero: numero, Revelada: revelada})
		}
		if integranteID.Valid { // una ronda sin votos viene con integrante NULL (LEFT JOIN)
			actual := &rondas[len(rondas)-1]
			actual.Votos = append(actual.Votos, dominio.Voto{IntegranteID: integranteID.Int64, Valor: int(valor.Int64)})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("leer rondas de la votación %d: %w", votacionID, err)
	}
	return rondas, nil
}

// insertarRondas escribe las rondas de la votación y sus votos.
func insertarRondas(ctx context.Context, tx *sql.Tx, v dominio.Votacion) error {
	for _, ronda := range v.Rondas {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO rondas (votacion_id, numero, revelada) VALUES ($1, $2, $3)`,
			v.ID, ronda.Numero, ronda.Revelada); err != nil {
			return err
		}
		for _, voto := range ronda.Votos {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO votos (votacion_id, ronda, integrante_id, valor) VALUES ($1, $2, $3, $4)`,
				v.ID, ronda.Numero, voto.IntegranteID, voto.Valor); err != nil {
				return err
			}
		}
	}
	return nil
}

// leerVotacion arma una votación (sin rondas) desde una fila; la acordada NULL queda en nil.
func leerVotacion(f fila) (dominio.Votacion, error) {
	var v dominio.Votacion
	var estado string
	var acordada sql.NullInt64
	err := f.Scan(&v.ID, &v.ProyectoID, &v.HistoriaID, &estado, &acordada)
	if errors.Is(err, sql.ErrNoRows) {
		return dominio.Votacion{}, err
	}
	if err != nil {
		return dominio.Votacion{}, fmt.Errorf("leer votación: %w", err)
	}
	v.Estado = dominio.EstadoVotacion(estado)
	if acordada.Valid {
		a := int(acordada.Int64)
		v.Acordada = &a
	}
	return v, nil
}

// enTransaccion ejecuta hacer dentro de una transacción: si devuelve error, se deshace todo.
func enTransaccion(ctx context.Context, db *sql.DB, hacer func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := hacer(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
