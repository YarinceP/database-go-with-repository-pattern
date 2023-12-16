package users

import (
	"database-go-with-repository-pattern/database"
	"database/sql"
	_ "database/sql"
	_ "errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	// Configura el mock de la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			return
		}
	}(db)

	// Establece la conexión de la base de datos en el paquete 'database'
	database.DB = db

	// Define las filas esperadas y el resultado
	expectedUserID := 1
	expectedUserName := "John Doe"
	mock.ExpectQuery("SELECT \\* FROM users WHERE id = \\?").
		WithArgs(expectedUserID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(expectedUserID, expectedUserName))

	// Llama a la función que deseas probar
	resultUser, err := NewUserService(db).GetUserByID(expectedUserID)
	if err != nil {
		log.Println("Error: ", err)
	}

	// Verifica el resultado y los errores esperados
	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.Equal(t, expectedUserID, resultUser.ID)
	assert.Equal(t, expectedUserName, resultUser.Name)

	// Asegúrate de que no haya más interacciones con la base de datos
	assert.NoError(t, mock.ExpectationsWereMet())
}
