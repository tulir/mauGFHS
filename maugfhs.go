package main

import (
	"maunium.net/go/mauGFHS/db"
)

func main() {
	err := db.Open("root", "password", "localhost", 3306, "maugfhs")
	if err != nil {
		panic(err)
	}

	db.CreateTables()
	err = db.Close()
	if err != nil {
		panic(err)
	}
}
