package repository

import (
	"context"
	"database-go-with-repository-pattern/entity"
	"log"
	"time"
)

func GetComment(commentRepository CommentRepository) (*entity.Comment, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	comment, err := commentRepository.FindByID(ctx, 1)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return comment, nil
}
