package repository

import (
	"context"
	"database-go-with-repository-pattern/entity"
	"database-go-with-repository-pattern/mocks"
)

func MockInsert() *mocks.CommentRepository {
	mock := mocks.CommentRepository{}
	ctx := context.Background()

	mock.On("Insert", ctx, entity.Comment{
		Email:   "example@test.com",
		Comment: "Example comment",
	}).Return(&entity.Comment{
		Id:      1,
		Email:   "example@test.com",
		Comment: "Example comment",
	}, nil)

	return &mock
}
