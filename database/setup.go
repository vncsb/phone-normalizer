package database

import (
	"database/sql"
)

type PhoneDb struct {
	Connection
	sqlDb *sql.DB
}

const defaultDatabaseName = "exercise_db"

func OpenConnection(conn Connection) (PhoneDb, error) {
	db, err := sql.Open("postgres", conn.makeConnectionString())
	if err != nil {
		return PhoneDb{}, err
	}

	phoneDb := PhoneDb{
		Connection: conn,
		sqlDb:      db,
	}

	return phoneDb.setupDb()
}

func (pdb PhoneDb) setupDb() (PhoneDb, error) {
	if !pdb.hasDatabaseName() {
		pdb.DBname = defaultDatabaseName

		err := pdb.resetDB()
		if err != nil {
			return PhoneDb{}, err
		}

		err = pdb.changeDatabase()
		if err != nil {
			return PhoneDb{}, err
		}
	}

	err := pdb.createTable()
	if err != nil {
		return PhoneDb{}, err
	}

	return pdb, nil
}

func (pdb *PhoneDb) changeDatabase() error {
	pdb.Close()

	db, err := sql.Open("postgres", pdb.makeConnectionString())
	if err != nil {
		return err
	}

	pdb.sqlDb = db
	return nil
}

func (pdb PhoneDb) resetDB() error {
	drop := `DROP DATABASE IF EXISTS ` + pdb.DBname

	_, err := pdb.sqlDb.Exec(drop)
	if err != nil {
		return err
	}

	return pdb.createDB()
}

func (pdb PhoneDb) createDB() error {
	create := "CREATE DATABASE " + pdb.DBname

	_, err := pdb.sqlDb.Exec(create)
	if err != nil {
		return err
	}

	return nil
}

func (pdb PhoneDb) createTable() error {
	const query = `
		CREATE TABLE IF NOT EXISTS phones (
			id SERIAL PRIMARY KEY,
			number TEXT NOT NULL
		)`

	_, err := pdb.sqlDb.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
