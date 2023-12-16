// Package database proporciona funciones para gestionar la conexión y cierre de la base de datos.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DBInterface is an interface that represents the methods used from sql.DB
type DBInterface interface {
	Ping() error
	Close() error
}

// GetDB returns a DBInterface
func GetDB() DBInterface {
	// Use the actual SQL database in your application code
	return db
}

// La variable db es una variable en paquete para almacenar la conexión activa a la base de datos.
var db *sql.DB

// GetConnection devuelve una instancia de sql.DB que representa la conexión a la base de datos.
// Si ya existe una conexión, se devuelve la existente; de lo contrario, se crea una nueva conexión.
// La función maneja automáticamente la configuración de la conexión y establece los parámetros relevantes.
// Si se produce algún error al abrir la conexión, se devuelve un error descriptivo.
func GetConnection() (*sql.DB, error) {
	// Si ya hay una conexión, devolverla
	if db != nil {
		return db, nil
	}

	// Obtener la cadena de conexión
	connectionString, err := GetConnectionString()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Intentar abrir una nueva conexión a la base de datos
	newDB, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("GetConnection: failed to open database connection. Connection string: %s. Error: %w", connectionString, err)
	}

	// Configurar parámetros de la conexión
	newDB.SetMaxIdleConns(10)
	newDB.SetMaxOpenConns(100)
	newDB.SetConnMaxIdleTime(5 * time.Minute)
	newDB.SetConnMaxLifetime(15 * time.Minute)

	// Almacenar la conexión para futuros usos
	db = newDB
	return db, nil
}

// CloseConnection cierra la conexión activa a la base de datos.
// Si la conexión ya está cerrada, la función retorna sin errores.
// La función también verifica si la conexión está cerrada antes de intentar cerrarla para evitar errores adicionales.
// Si se produce algún error durante la verificación o el cierre de la conexión, se devuelve un error descriptivo.
func CloseConnection() error {
	if db != nil {
		// Verificar si la conexión ya está cerrada
		if err := db.Ping(); err != nil {
			if errors.Is(err, sql.ErrConnDone) {
				return nil // La conexión ya está cerrada
			}
			return fmt.Errorf("CloseConnection: error checking connection before closing: %w", err)
		}

		// Cerrar la conexión
		if err := db.Close(); err != nil {
			// Log o manejar el error según sea necesario
			return fmt.Errorf("CloseConnection: error closing database connection: %w", err)
		}
	}

	return nil
}
