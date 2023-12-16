package main

import (
	"database-go-with-repository-pattern/repository"
	"log"
)

func SearchDefaultComment() {
	commentRepository := repository.NewCommentRepository(Connection)
	c, err := repository.GetComment(commentRepository, RequestedId)
	if err != nil {
		log.Println("Error on retrieving a comment: ", err)
	}
	log.Printf("Comment %+v", c)
}
