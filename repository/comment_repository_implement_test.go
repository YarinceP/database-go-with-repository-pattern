package repository

import (
	"context"
	"database-go-with-repository-pattern/entity"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestInsertComment(t *testing.T) {
	// Crear un DB mock y un objeto de simulación
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear el mock de la base de datos: %v", err)
	}
	defer db.Close()

	// Crear un repositorio con la base de datos simulada
	repository := &commentRepositoryImplementation{
		DB: db,
	}

	// Comentario de prueba
	testComment := entity.Comment{
		Email:   "test@example.com",
		Comment: "This is a test comment",
	}

	// Caso de prueba exitoso
	mock.ExpectExec("INSERT INTO comments").WithArgs(testComment.Email, testComment.Comment).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repository.Insert(context.Background(), testComment)
	assert.NoError(t, err, "No se esperaba un error para un caso de prueba exitoso")

	// Caso de prueba con consulta SQL vacía
	repositoryEmptyQuery := &commentRepositoryImplementation{
		DB: db,
	}

	errEmptyQuery := repositoryEmptyQuery.Insert(context.Background(), testComment)
	assert.Error(t, errEmptyQuery, "Se esperaba un error para una consulta SQL vacía")
	assert.Contains(t, errEmptyQuery.Error(), "failed to execute SQL query", "Mensaje de error inesperado")

	// Caso de prueba con cero filas afectadas
	repositoryNoRowsAffected := &commentRepositoryImplementation{
		DB: db,
	}

	mock.ExpectExec("INSERT INTO comments").WithArgs(testComment.Email, testComment.Comment).WillReturnResult(sqlmock.NewResult(0, 0))

	errNoRowsAffected := repositoryNoRowsAffected.Insert(context.Background(), testComment)
	assert.Error(t, errNoRowsAffected, "Se esperaba un error para cero filas afectadas")
	assert.Contains(t, errNoRowsAffected.Error(), "no rows were affected", "Mensaje de error inesperado")
}

// MockQueryValidator es un mock para EnsureNonEmptyQuery
type MockQueryValidator struct {
	mock.Mock
}

// EnsureNonEmptyQuery es la implementación del mock para EnsureNonEmptyQuery
func (m *MockQueryValidator) EnsureNonEmptyQuery(query string) error {
	args := m.Called(query)
	return args.Error(0)
}

func Test_commentRepositoryImplementation_Insert(t *testing.T) {
	mockValidator := new(MockQueryValidator)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error al crear el mock de la base de datos: %v", err)
	}

	ctx := context.Background()
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx     context.Context
		comment entity.Comment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Comment inserted successfully",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: ctx,
				comment: entity.Comment{
					Email:   "test@example.com",
					Comment: "Example comment",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Error inserting comment",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: ctx,
				comment: entity.Comment{
					Email:   "test@example.com",
					Comment: "Example comment",
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "Error in EnsureNonEmptyQuery",
			fields: fields{
				DB: db,
			},
			args: args{
				ctx: ctx,
				comment: entity.Comment{
					Email:   "test@example.com",
					Comment: "Example comment",
				},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &commentRepositoryImplementation{
				DB: tt.fields.DB,
			}

			if tt.name == "Comment inserted successfully" {
				// Configurar expectativas para la ejecución de la consulta SQL
				mock.ExpectExec("INSERT INTO comments").
					WithArgs(tt.args.comment.Email, tt.args.comment.Comment).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}
			if tt.name == "Error inserting comment" {
				mock.ExpectExec("INSERT INTO comments").
					WithArgs(tt.args.comment.Email, tt.args.comment.Comment).
					WillReturnError(errors.New("Error inserting comment"))
			}
			if tt.name == "Error in EnsureNonEmptyQuery" {
				mockValidator.On("EnsureNonEmptyQuery", "").Return(errors.New("failed to execute SQL query: la consulta SQL no puede estar vacía"))
			}

			// Llamar a la función bajo prueba
			tt.wantErr(t, repository.Insert(tt.args.ctx, tt.args.comment), fmt.Sprintf("Insert(%v, %v)", tt.args.ctx, tt.args.comment))

			// Verificar que no haya más expectativas no cumplidas
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("Expectativas no cumplidas: %s", err)
			}
		})
	}
}
