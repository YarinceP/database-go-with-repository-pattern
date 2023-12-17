package repository

import (
	"context"
	"database-go-with-repository-pattern/database"
	"database-go-with-repository-pattern/entity"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type commentRepositoryImplementation struct {
	DB *sql.DB
}

func NewCommentRepository(db *sql.DB) CommentRepository {
	// Validar que el parámetro db no sea nulo
	if db == nil {
		panic("Error creating CommentRepository: provided database connection is nil. Please ensure a valid database connection is provided.")
	}

	// Validar la conexión a la base de datos
	err := db.Ping()
	if err != nil {
		panic("Error creating CommentRepository: unable to connect to the database. Please check your database connection settings and ensure the database server is running.")
	}

	// Crear la instancia del repositorio
	return &commentRepositoryImplementation{DB: db}
}

func (repository *commentRepositoryImplementation) Insert(ctx context.Context, comment entity.Comment) error {
	query := CommentQueries.Insert

	// Ensure that the SQL query is not empty
	if err := database.EnsureNonEmptyQuery(query); err != nil {
		return fmt.Errorf("failed to execute SQL query: %v", err)
	}

	// Log the SQL query before executing
	log.Printf("Executing SQL query: %s", query)

	// Execute the SQL query
	result, err := repository.DB.ExecContext(ctx, query, comment.Email, comment.Comment)
	if err != nil {
		return fmt.Errorf("failed to execute SQL query: %v", err)
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get the number of affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("no rows were affected by the SQL query")
	}

	// Log the success of the operation
	log.Printf("Comment inserted successfully: %s", comment.Comment)

	return nil
}

func (repository *commentRepositoryImplementation) FindByID(ctx context.Context, id int32) (*entity.Comment, error) {
	script := "SELECT id, email, comment FROM comments WHERE id = ? LIMIT 1"
	result, err := repository.DB.QueryContext(ctx, script, id)
	comment := entity.Comment{}
	if err != nil {
		return &comment, err
	}

	defer func(result *sql.Rows) {
		err := result.Close()
		if err != nil {
			return
		}
	}(result)

	if result.Next() {
		err := result.Scan(&comment.Id, &comment.Email, &comment.Comment)
		if err != nil {
			return nil, err
		}
		return &comment, nil
	} else {
		return &comment, errors.New("ID " + strconv.Itoa(int(id)) + " Not found !!")
	}
}

func (repository *commentRepositoryImplementation) FindAll(ctx context.Context) ([]entity.Comment, error) {
	script := "SELECT id, email, comment FROM comments"
	result, err := repository.DB.QueryContext(ctx, script)
	if err != nil {
		return nil, err
	}

	defer func(result *sql.Rows) {
		err := result.Close()
		if err != nil {
			return
		}
	}(result)
	var comments []entity.Comment
	for result.Next() {
		comment := entity.Comment{}
		err := result.Scan(&comment.Id, &comment.Email, &comment.Comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
