package repository

import (
	"context"
	db "database-go-with-repository-pattern/database"
	"database-go-with-repository-pattern/entity"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestInsertComment(t *testing.T) {
	connection, err := db.GetConnection()
	if err != nil {
		return
	}
	commentRepository := NewCommentRepository(connection)

	ctx := context.Background()
	comment := entity.Comment{
		Email:   "repository@test.com",
		Comment: "Test repository",
	}
	result, err := commentRepository.Insert(ctx, comment)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func TestFindByIdComment(t *testing.T) {
	connection, err := db.GetConnection()
	if err != nil {
		return
	}
	commentRepository := NewCommentRepository(connection)

	ctx := context.Background()

	result, err := commentRepository.FindByID(ctx, 2)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}

func TestFindAllComment(t *testing.T) {
	connection, err := db.GetConnection()
	if err != nil {
		return
	}
	commentRepository := NewCommentRepository(connection)

	ctx := context.Background()

	result, err := commentRepository.FindAll(ctx)
	if err != nil {
		panic(err)
	}

	for _, comment := range result {
		fmt.Println(comment)
	}
}
