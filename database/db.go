package database

import (
	"database/sql"
)

type PhoneDb struct {
	Connection
	sqlDb *sql.DB
}

type PhoneNumber struct {
	ID     int
	Number string
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

func (pdb PhoneDb) InsertNumbers(numbers []string) error {
	const insert = `
	INSERT INTO phones (number)
	VALUES ($1)`

	for _, p := range numbers {
		_, err := pdb.sqlDb.Exec(insert, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pdb PhoneDb) UpdateNumbers(phone PhoneNumber) error {
	const update = `
	UPDATE phones
	SET number = $1
	WHERE id = $2`

	_, err := pdb.sqlDb.Exec(update, phone.Number, phone.ID)
	if err != nil {
		return err
	}

	return nil
}

func (pdb PhoneDb) GetNumbers() ([]PhoneNumber, error) {
	const query = `SELECT id, number FROM phones`

	rows, err := pdb.sqlDb.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	numbers := make([]PhoneNumber, 0)
	for rows.Next() {
		var pn PhoneNumber

		err = rows.Scan(&pn.ID, &pn.Number)
		if err != nil {
			return nil, err
		}

		numbers = append(numbers, pn)

	}

	return numbers, nil
}

func (pdb PhoneDb) Close() error {
	err := pdb.sqlDb.Close()
	return err
}
