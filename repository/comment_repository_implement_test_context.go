package repository

import (
	"context"
	"database-go-with-repository-pattern/entity"
)

type TestContextInsert struct {
	ctx     context.Context
	comment entity.Comment
}

var contextInsert = TestContextInsert{
	ctx: context.Background(),
	comment: entity.Comment{
		Email:   "example@test.com",
		Comment: "Example",
	},
}

//Case maps

var MapTestContextInsert = map[string]TestContextInsert{
	"Inserting a Comment successfully": contextInsert,
}
