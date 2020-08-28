package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var uname, pass, schema = "username", "password", "phonedb"
	connection := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", uname, pass, schema)
	db, err := sql.Open("mysql", connection)
	check(err)
	conerr := db.Ping()
	check(conerr)
	defer db.Close()
	// DB is ready

	var phnumber string
	rows, err := db.Query("select * from phone_numbers")
	check(err)
	defer rows.Close()

	// Display all data
	var phnumbers = []string{}
	for rows.Next() {
		err := rows.Scan(&phnumber)
		check(err)
		phnumbers = append(phnumbers, phnumber)
	}

	// Clean the data
	var cleaned = map[string]string{}
	for _, number := range phnumbers {
		cleaned[number] = clean(number)
		fmt.Printf("%14s - %s\n", number, cleaned[number])
	}

	// Update the data
	var phno string
	var delstatement = `
	DELETE FROM phone_numbers
	WHERE phone_number = ?;`
	var updatestatement = `
	UPDATE phone_numbers
	SET phone_number = ?
	WHERE phone_number = ?;`
	for k, v := range cleaned {
		// Check if exists, skip
		err := db.QueryRow("select * from phone_numbers where phone_number = ?", v).Scan(&phno)
		if err == nil {
			if k == v {
				continue
			}
			fmt.Printf("Deleting %14s\n", k)
			_, exerr := db.Exec(delstatement, k)
			check(exerr)
		} else if err == sql.ErrNoRows {
			fmt.Printf("Updating %14s\n", k)
			_, exerr := db.Exec(updatestatement, v, k)
			check(exerr)
		} else {
			log.Fatal(err)
		}
	}
}

func clean(phno string) string {
	re := regexp.MustCompile("[\\(\\-\\)\\s]")
	matched := re.ReplaceAllString(phno, "")
	return matched
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
