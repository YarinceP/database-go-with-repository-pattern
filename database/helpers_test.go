package database

import (
	_ "errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetConnectionString(t *testing.T) {
	// Caso de prueba 1: Configuración válida
	viper.Set("database.connectionString", "valid_connection_string")
	connection, err := GetConnectionString()
	assert.NoError(t, err)                                 // Asegura que no haya error
	assert.Equal(t, "valid_connection_string", connection) // Asegura que la conexión sea la esperada

	// Caso de prueba 2: Configuración no válida
	viper.Set("database.connectionString", "") // Configura una cadena de conexión vacía
	_, err = GetConnectionString()
	assert.Error(t, err)                                                                                                         // Asegura que se haya producido un error
	assert.EqualError(t, err, "database connection string is not configured; make sure to provide it in the configuration file") // Asegura que el error sea el esperado
}
