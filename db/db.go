package db

import (
	"database/sql"
	"fmt"

	"maunium.net/go/mauGFHS/db/config"

	// Import MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Open opens a database connection with the given details.
func Open(config dbconfig.DBConfig) error {
	var err error
	db, err = sql.Open("mysql", config.GetDSN())
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
	createTable("namespaces", namespacesSchema)
	createTable("files", filesSchema)
	createTable("filepermissions", filePermissionsSchema)
	createTable("nspermissions", nsPermissionsSchema)
}
