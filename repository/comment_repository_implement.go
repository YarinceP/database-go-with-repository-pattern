package repository

import (
	"context"
	"database-go-with-repository-pattern/entity"
	"database/sql"
	"errors"
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

func (repository commentRepositoryImplementation) Insert(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	script := "INSERT INTO comments(email, comment) VALUES(?,?)"
	result, err := repository.DB.ExecContext(ctx, script, comment.Email, comment.Comment)
	if err != nil {
		return comment, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return comment, err
	}
	comment.Id = int32(lastId)
	return comment, nil
}

func (repository commentRepositoryImplementation) FindByID(ctx context.Context, id int32) (*entity.Comment, error) {
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

func (repository commentRepositoryImplementation) FindAll(ctx context.Context) ([]entity.Comment, error) {
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
