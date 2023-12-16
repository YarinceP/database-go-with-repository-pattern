package main

import (
	"database-go-with-repository-pattern/database"
	"database-go-with-repository-pattern/repository"
	"log"

	"database/sql"
)

var Connection *sql.DB

func main() {

	Connection = database.GetConnection()
	commentRepository := repository.NewCommentRepository(Connection)
	c, err := repository.GetComment(commentRepository)
	if err != nil {
		log.Println("Error on retrieving a comment: ", err)
	}
	log.Printf("Comment %+v", c)
}
