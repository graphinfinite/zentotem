package migrations

import "github.com/jmoiron/sqlx"

func Up(dbPool *sqlx.DB) error {
	q := `
		CREATE TABLE IF NOT EXISTS "user" (
		id SERIAL PRIMARY KEY,
		createdAt timestamp DEFAULT current_timestamp NOT NULL,
		name VARCHAR(64) UNIQUE NOT NULL,  
		age integer NOT NULL,
		`
	_, errDb := dbPool.Exec(q)
	if errDb != nil {
		return errDb
	}
	return nil
}
