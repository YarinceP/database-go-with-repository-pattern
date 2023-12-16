package main

import (
	"database-go-with-repository-pattern/database"
	"database/sql"
	"log"
)

var Connection *sql.DB
var err error

// RequestedId Variable just for the example
var RequestedId int32 = 1

func main() {
	Connection, err = database.GetConnection()
	if err != nil {
		log.Println("Error acquiring the database connection:", err)
	}

	defer func() {
		err := database.CloseConnection()
		if err != nil {
			log.Println("Error closing the database connection:", err)
		}
	}()

	SearchDefaultComment()
}
