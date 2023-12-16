package repository

import (
	"context"
	"database-go-with-repository-pattern/entity"
	"errors"
	"log"
	"time"
)

// GetComment recupera un comentario utilizando el CommentRepository proporcionado.
// La función utiliza un contexto con un tiempo de espera de 5 segundos para la operación.
// Devuelve un puntero a entity. Comment si la operación es exitosa.
// Si commentRepository es nulo, se devuelve un error indicando que la operación no puede continuar sin un repositorio válido.
// Si se produce un error al recuperar el comentario, se registra con un mensaje detallado que incluye
// información sobre el contexto de la operación, el ID del comentario solicitado y la naturaleza específica del error.
// La función garantiza la cancelación adecuada del contexto incluso si se produce un error.
func GetComment(commentRepository CommentRepository, commentID int32) (*entity.Comment, error) {
	if commentRepository == nil {
		return nil, errors.New("GetComment: commentRepository no puede ser nulo. No se puede continuar sin un repositorio válido. ")
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	comment, err := commentRepository.FindByID(ctx, commentID)
	if err != nil {
		log.Printf("GetComment: Error al recuperar comentario con ID %v. Contexto: %v. Error: %v", commentID, ctx, err)
		return nil, err
	}

	return comment, nil
}
