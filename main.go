package main

import (
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/vncsb/phone-normalizer/database"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "root"
	dbname   = ""
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
	conn := database.Connection{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBname:   dbname,
	}

	db, err := database.OpenConnection(conn)
	must(err)
	defer db.Close()

	must(db.InsertNumbers(phones))

	numbers, err := db.GetNumbers()
	must(err)

	fmt.Println("Before")
	for _, pn := range numbers {
		fmt.Println(pn)
		pn.Number = normalize(pn.Number)
		must(db.UpdateNumbers(pn))
	}

	numbers, err = db.GetNumbers()
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

func must(err error) {
	if err != nil {
		panic(err)
	}
}
