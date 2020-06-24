package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type PhoneNumber struct {
	ID     int
	Number string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "exercise_db"
)

var phones = []string{"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892"}

func main() {
	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s sslmode=disable",
		host, port, user, password)

	db, err := sql.Open("postgres", connString)
	must(err)

	defer db.Close()

	err = setupDB(db)
	must(err)

	numbers, err := getNumbers(db)
	must(err)

	fmt.Println("Before")
	for _, pn := range numbers {
		fmt.Println(pn)
		pn.Number = normalize(pn.Number)
		err := updateNumbers(pn, db)
		if err != nil {
			panic(err)
		}
	}

	numbers, err = getNumbers(db)
	must(err)

	fmt.Println("After")
	for _, pn := range numbers {
		fmt.Println(pn)
	}
}

func normalize(number string) string {
	var buf strings.Builder
	for _, ch := range number {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}

func setupDB(db *sql.DB) error {
	err := createDB(dbname, db)
	must(err)

	err = createTable(db)
	must(err)

	err = resetDB(dbname, db)
	must(err)

	err = insertNumbers(phones, db)
	must(err)

	return nil
}

func createDB(name string, db *sql.DB) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}

	return nil
}

func createTable(db *sql.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS $1.phones (
			id SERIAL PRIMARY KEY,
			number TEXT NOT NULL
		)`

	_, err := db.Exec(query, dbname)
	if err != nil {
		return err
	}

	return nil
}

func insertNumbers(numbers []string, db *sql.DB) error {
	const insert = `
	INSERT INTO $1.phones (number)
	VALUES ($2)`

	for _, p := range numbers {
		_, err := db.Exec(insert, dbname, p)
		if err != nil {
			return err
		}
	}

	return nil
}

func updateNumbers(phone PhoneNumber, db *sql.DB) error {
	const update = `
	UPDATE $1.phones
	SET number = $2
	WHERE id = $3`

	_, err := db.Exec(update, dbname, phone.Number, phone.ID)
	if err != nil {
		return err
	}

	return nil
}

func resetDB(name string, db *sql.DB) error {
	var drop = `DROP DATABASE IF EXISTS ` + name

	_, err := db.Exec(drop)
	if err != nil {
		return err
	}

	return nil
}

func getNumbers(db *sql.DB) ([]PhoneNumber, error) {
	const query = `SELECT id, number FROM $1.phones`

	rows, err := db.Query(query, dbname)
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
