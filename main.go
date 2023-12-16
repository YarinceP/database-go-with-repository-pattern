package main

import (
	"database-go-with-repository-pattern/database"
	"database-go-with-repository-pattern/users"
	"log"
)

func main() {
	database.ConnectDB()
	defer database.DisconnectDB()

	users.NewUserService(database.GetConnectionInstance())
	user, err := users.NewUserService(database.GetConnectionInstance()).GetUserByID(1)
	if err != nil {
		return
	}

	log.Printf("User: %+v", user)
}
