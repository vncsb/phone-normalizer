package database

import "fmt"

type Connection struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
}

func (conn Connection) makeConnectionString() string {
	connString := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s sslmode=disable",
		conn.Host, conn.Port, conn.User, conn.Password)

	if conn.hasDatabaseName() {
		connString += fmt.Sprintf(" dbname=%s", conn.DBname)
	}

	return connString
}

func (conn Connection) hasDatabaseName() bool {
	if conn.DBname != "" {
		return true
	}

	return false
}
