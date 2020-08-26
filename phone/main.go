package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "lakedragon:lakedragon@tcp(127.0.0.1:3306)/company")
	if err != nil {
		log.Fatal(err)
	}
	conerr := db.Ping()
	if conerr != nil {
		log.Fatal(conerr)
	}
	defer db.Close()
}
