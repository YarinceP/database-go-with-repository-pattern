// Package database proporciona funciones para gestionar la conexión y cierre de la base de datos.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const DefaultConfigFileName = "config/config.json"
const DefaultConfigObjectName = "config"
const DefaultTimeDurationMinuteConfig = time.Minute

// DBConnector es una interfaz para las operaciones de conexión a la base de datos.
type DBConnector interface {
	Connect() (*sql.DB, error)
	Close() error
}

type ConnectionError struct {
	Reason string
}

func (e *ConnectionError) Error() string {
	return fmt.Sprintf("Connection error: %s", e.Reason)
}

// GetConnection Modifica la función GetConnection para aceptar un parámetro DBConnector
func GetConnection(connector DBConnector) (*sql.DB, error) {
	return connector.Connect()
}

// CloseConnection Modifica la función CloseConnection para aceptar un parámetro DBConnector
func CloseConnection(connector DBConnector) error {
	return connector.Close()
}

// Connector implementa la interfaz DBConnector para MySQL.
type Connector struct {
	db               *sql.DB
	ConnectionString string
}

func NewConnectorWithConnectionString() (*Connector, error) {
	connectionString, err := GetConnectionString()
	if err != nil {
		return nil, &ConnectionError{Reason: fmt.Sprintf("NewConnectorWithConnectionString: error getting connection string: %v", err)}
	}
	return &Connector{ConnectionString: connectionString}, nil
}

func (c *Connector) Connect() (*sql.DB, error) {
	// Lógica de conexión aquí
	if c.db != nil {
		return c.db, nil
	}

	// Cargar configuración desde JSON
	config, err := LoadConfigFromJSON(DefaultConfigFileName, DefaultConfigObjectName)
	if err != nil {
		return nil, &ConnectionError{
			Reason: fmt.Sprintf("Connector.Connect: error loading configuration from JSON: %v", err),
		}
	}

	// Intentar abrir una nueva conexión a la base de datos
	newDB, err := sql.Open("mysql", c.ConnectionString)
	if err != nil {
		return nil, &ConnectionError{Reason: fmt.Sprintf("Connector.Connect: failed to open database connection. Connection string: %s. Error: %v", c.ConnectionString, err)}
	}

	c.configureDBParameters(newDB, config)

	// Almacenar la conexión para futuros usos
	c.db = newDB
	return c.db, nil
}

type DBParametersConfig struct {
	MaxIdleConns    int           `json:"max_idle_conns"`
	MaxOpenConns    int           `json:"max_open_conns"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
}

func NewDBParametersConfig(maxIdleConns int, maxOpenConns int, connMaxIdleTime time.Duration, connMaxLifetime time.Duration) DBParametersConfig {
	return DBParametersConfig{
		MaxIdleConns:    maxIdleConns,
		MaxOpenConns:    maxOpenConns,
		ConnMaxIdleTime: connMaxIdleTime,
		ConnMaxLifetime: connMaxLifetime}
}

func (c *Connector) configureDBParameters(newDB *sql.DB, config DBParametersConfig) {
	// Configurar parámetros de la conexión
	newDB.SetMaxIdleConns(config.MaxIdleConns)
	newDB.SetMaxOpenConns(config.MaxOpenConns)
	newDB.SetConnMaxIdleTime(config.ConnMaxIdleTime * DefaultTimeDurationMinuteConfig)
	newDB.SetConnMaxLifetime(config.ConnMaxLifetime * DefaultTimeDurationMinuteConfig)
}

func (c *Connector) Close() error {
	if c.db != nil {
		// Verificar si la conexión ya está cerrada
		if err := c.db.Ping(); err != nil {
			if errors.Is(err, sql.ErrConnDone) {
				return nil // La conexión ya está cerrada
			}
			return &ConnectionError{Reason: fmt.Sprintf("Connector.Close: error checking connection before closing: %v", err)}
		}

		// Cerrar la conexión
		if err := c.db.Close(); err != nil {
			// Log o manejar el error según sea necesario
			return &ConnectionError{Reason: fmt.Sprintf("Connector.Close: error closing database connection: %v", err)}
		}
	}

	return nil
}
