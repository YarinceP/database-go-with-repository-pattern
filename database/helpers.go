package database

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

func GetConnectionString() (string, error) {
	// Get the connection string from viper
	connectionString := viper.GetString("database.connectionString")
	if connectionString == "" {
		// Raise an error indicating that the configuration is mandatory
		return "", errors.New("database connection string is not configured; make sure to provide it in the configuration file")
	}

	return connectionString, nil
}

func LoadConfigFromJSON(fileName, objectName string) (DBParametersConfig, error) {
	// Configurar Viper para leer el archivo JSON
	viper.SetConfigFile(fileName)

	// Leer la configuración
	if err := viper.ReadInConfig(); err != nil {
		return DBParametersConfig{}, fmt.Errorf("failed to read configuration from file %v", err)
	}

	// Obtener los valores
	maxIdleConns := viper.GetInt("database.config.max_idle_conns")
	maxOpenConns := viper.GetInt("database.config.max_open_conns")
	connMaxIdleTime := viper.GetInt("database.config.conn_max_idle_time")
	connMaxLifetime := viper.GetInt("database.config.conn_max_lifetime")

	// Validar los valores
	if err := validateNonZeroValues(maxIdleConns, maxOpenConns, connMaxIdleTime, connMaxLifetime); err != nil {
		return DBParametersConfig{}, fmt.Errorf("Validation error: %v\n", err)
	}

	config := NewDBParametersConfig(maxIdleConns, maxOpenConns, time.Duration(connMaxIdleTime)*DefaultTimeDurationMinuteConfig, time.Duration(connMaxLifetime)*DefaultTimeDurationMinuteConfig)

	return config, nil
}

// validateNonZeroValues verifica que los valores proporcionados no sean cero.
// Devuelve un error indicando qué parámetro no puede ser cero, si es el caso.
func validateNonZeroValues(values ...int) error {
	for i, value := range values {
		if value == 0 {
			return fmt.Errorf("parameter at index %d cannot be zero", i)
		}
	}
	return nil
}

// EnsureNonEmptyQuery garantiza que la consulta SQL proporcionada no esté vacía.
//
// Esta función toma una cadena de consulta SQL como entrada y verifica que no esté vacía.
// Si la cadena está vacía, se devuelve un error indicando que la consulta no puede estar vacía.
// En caso contrario, se devuelve nil, indicando que la consulta es válida.
//
// Parámetros:
//   - query (string): La consulta SQL que se va a validar.
//
// Devuelve:
//   - error: Un error si la consulta está vacía, nil si la consulta es válida.
//
// Ejemplo de uso:
//
//	err := EnsureNonEmptyQuery("SELECT * FROM table WHERE condition")
//	if err != nil {
//	    log.Printf("Error al validar la consulta SQL: %v", err)
//	    // Manejar el error, por ejemplo, devolver un error desde tu función
//	    return err
//	}
//
// Nota: Asegúrate de utilizar esta función antes de ejecutar la consulta SQL para evitar operaciones
// no deseadas o errores debido a consultas SQL vacías.
func EnsureNonEmptyQuery(query string) error {
	if query == "" {
		return errors.New("la consulta SQL no puede estar vacía")
	}
	return nil
}
