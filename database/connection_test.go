package database

import (
	"database/sql"
	"errors"
	_ "errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

// Mock para DBConnector
type MockConnector struct {
	mock.Mock
}

func (m *MockConnector) Connect() (*sql.DB, error) {
	args := m.Called()
	return args.Get(0).(*sql.DB), args.Error(1)
}

func (m *MockConnector) Close() error {
	args := m.Called()
	return args.Error(0)
}

// GetConnectionString es una función mock para GetConnectionString
func (m *MockConnector) GetConnectionString() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestGetConnection(t *testing.T) {
	// Crear un mock para DBConnector
	mockConnector := new(MockConnector)

	// Configurar el comportamiento esperado del mock
	mockConnector.On("Connect").Return(&sql.DB{}, nil)

	// Usar el mock en tu función bajo prueba
	db, err := GetConnection(mockConnector)

	// Verificar que la función bajo prueba se comporte como se espera
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Asegurarse de que se haya llamado al método Connect del mock
	mockConnector.AssertExpectations(t)
}

func TestCloseConnection(t *testing.T) {
	// Crear un mock para DBConnector
	mockConnector := new(MockConnector)

	// Configurar el comportamiento esperado del mock
	mockConnector.On("Close").Return(nil)

	// Usar el mock en tu función bajo prueba
	err := CloseConnection(mockConnector)

	// Verificar que la función bajo prueba se comporte como se espera
	assert.NoError(t, err)

	// Asegurarse de que se haya llamado al método Close del mock
	mockConnector.AssertExpectations(t)
}

func TestNewConnectorWithConnectionString(t *testing.T) {
	var defaultConnectionString = "fakeConnectionString"
	// Caso de prueba 1: Configuración válida
	viper.Set("database.connectionString", defaultConnectionString)
	connection, err := GetConnectionString()
	assert.NoError(t, err)                               // Asegura que no haya error
	assert.Equal(t, defaultConnectionString, connection) // Asegura que la conexión sea la esperada

	tests := []struct {
		name    string
		want    *Connector
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "New Connector Instance was create successfully with connection string from config files",
			want:    &Connector{ConnectionString: defaultConnectionString},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConnectorWithConnectionString()
			if !tt.wantErr(t, err, fmt.Sprintf("NewConnectorWithConnectionString()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewConnectorWithConnectionString()")
		})
	}
}

func TestNewConnectorWithConnectionString_ErrorInGetConnectionString(t *testing.T) {
	var defaultConnectionString = ""
	// Caso de prueba 1: Configuración válida
	viper.Set("database.connectionString", "")
	connectionString, err := GetConnectionString()
	assert.Error(t, err)                                       // Asegura que no haya error
	assert.Equal(t, defaultConnectionString, connectionString) // Asegura que la conexión sea la esperada
	tests := []struct {
		name    string
		want    *Connector
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "New Connector Instance was not created successfully with connection string from config files empty ",
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConnectorWithConnectionString()
			if !tt.wantErr(t, err, fmt.Sprintf("NewConnectorWithConnectionString()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewConnectorWithConnectionString()")
		})
	}
}

func TestConnectionError_Error(t *testing.T) {
	// Crear una instancia de ConnectionError
	err := &ConnectionError{
		Reason: "Test error reason",
	}

	// Obtener el mensaje de error llamando a la función Error()
	errorMessage := err.Error()

	// Verificar que el mensaje de error sea el esperado
	expectedErrorMessage := "Connection error: Test error reason"
	assert.Equal(t, expectedErrorMessage, errorMessage)
}

func TestConfigureDBParameters(t *testing.T) {
	// Crear un objeto Connector
	connector := &Connector{}

	// Crear una instancia de sql.DB (puede ser en memoria para la prueba)
	newDB, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_lib_go?parseTime=true")
	if err != nil {
		t.Fatal(err)
	}
	defer func(newDB *sql.DB) {
		err := newDB.Close()
		if err != nil {

		}
	}(newDB)

	// Crear una configuración de parámetros
	config := DBParametersConfig{
		MaxIdleConns:    10,
		MaxOpenConns:    20,
		ConnMaxIdleTime: time.Minute * 5,
		ConnMaxLifetime: time.Minute * 10,
	}

	// Llamar a la función configureDBParameters
	connector.configureDBParameters(newDB, config)

	// Obtener los valores actuales de los parámetros
	actualMaxIdleConns := newDB.Stats().MaxIdleClosed
	actualMaxOpenConns := newDB.Stats().MaxOpenConnections
	actualConnMaxIdleTime := newDB.Stats().MaxIdleTimeClosed
	actualConnMaxLifetime := newDB.Stats().MaxLifetimeClosed

	// Verificar que los parámetros se hayan configurado correctamente
	assert.Equal(t, config.MaxIdleConns, actualMaxIdleConns)
	assert.Equal(t, config.MaxOpenConns, actualMaxOpenConns)
	assert.Equal(t, config.ConnMaxIdleTime, actualConnMaxIdleTime)
	assert.Equal(t, config.ConnMaxLifetime, actualConnMaxLifetime)
}

// Mock para simular la función LoadConfigFromJSON
type MockConfigLoader struct {
	mock.Mock
}

func (m *MockConfigLoader) LoadConfigFromJSON(fileName, objectName string) (DBParametersConfig, error) {
	args := m.Called(fileName, objectName)
	return args.Get(0).(DBParametersConfig), args.Error(1)
}

func TestConnect(t *testing.T) {
	// Configuración de prueba
	connectionString := "root:@tcp(localhost:3306)/db_lib_go?parseTime=true"
	config := DBParametersConfig{
		MaxIdleConns:    10,
		MaxOpenConns:    20,
		ConnMaxIdleTime: DefaultTimeDurationMinuteConfig,
		ConnMaxLifetime: DefaultTimeDurationMinuteConfig,
	}
	mockConfigLoader := new(MockConfigLoader)
	connector := &Connector{
		ConnectionString: connectionString,
	}

	// Caso 1: Configuración y conexión exitosas
	mockConfigLoader.On("LoadConfigFromJSON", DefaultConfigFileName, DefaultConfigObjectName).Return(config, nil)

	// Llamar a Connect
	db, err := connector.Connect()

	// Verificar resultados
	assert.NoError(t, err)
	assert.NotNil(t, db)
	//assert.Equal(t, config.MaxIdleConns, db.Stats().MaxIdleConns)
	//assert.Equal(t, config.MaxOpenConns, db.Stats().OpenConnections)
	//assert.Equal(t, config.ConnMaxIdleTime, db.Stats().MaxIdleTime)
	//assert.Equal(t, config.ConnMaxLifetime, db.Stats().MaxLifetime)

	// Caso 2: Error en la carga de configuración
	mockConfigLoader.On("LoadConfigFromJSON", DefaultConfigFileName, DefaultConfigObjectName).Return(DBParametersConfig{}, errors.New("dummy error"))

	// Llamar a Connect nuevamente
	db, err = connector.Connect()

	// Verificar resultados
	assert.Error(t, err)
	assert.Nil(t, db)
}
