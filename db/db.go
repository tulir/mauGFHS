// mauGFHS - A server that can serve as a backend for many kinds of services that only require file hosting.
// Copyright (C) 2017 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

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
