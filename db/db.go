package db

import (
	"database/sql"
	"fmt"

	// Import MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Open opens a database connection with the given details.
func Open(username, password, address string, port int, database string) error {
	var err error
	fmt.Printf("%s:%s@%s:%d/%s\n", username, password, address, port, database)
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, address, port, database))
	if err != nil {
		return err
	}
	return nil
}

func createTable(name, schema string) {
	_, err := db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( %s ) ENGINE=InnoDB DEFAULT CHARSET=utf8", name, schema))
	if err != nil {
		panic(err)
	}
}

// Close closes the database connection.
func Close() error {
	if db != nil {
		err := db.Close()
		if err != nil {
			return err
		}
		db = nil
		return nil
	}
	return nil
}

// CreateTables creates all the tables that are needed.
func CreateTables() {
	createTable("users", usersSchema)
	createTable("authtokens", authTokensSchema)
	createTable("files", filesSchema)
	createTable("permissions", permissionsSchema)
}
