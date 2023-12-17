package database

import (
	_ "errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
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

func TestValidateNonZeroValues(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		errMsg string
	}{
		{
			name:   "No Zero Values",
			values: []int{1, 2, 3},
			errMsg: "",
		},
		{
			name:   "Single Zero Value",
			values: []int{1, 0, 3},
			errMsg: "parameter at index 1 cannot be zero",
		},
		{
			name:   "Multiple Zero Values",
			values: []int{0, 0, 0},
			errMsg: "parameter at index 0 cannot be zero",
		},
		{
			name:   "Empty Values",
			values: []int{},
			errMsg: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateNonZeroValues(test.values...)
			if test.errMsg == "" {
				assert.NoError(t, err, "Expected no error")
			} else {
				assert.EqualError(t, err, test.errMsg, "Error message does not match")
			}
		})
	}
}

/*
var viperReadInConfigFunc = viper.ReadInConfig

func TestLoadConfigFromJSON(t *testing.T) {
	tests := []struct {
		name           string
		fileContent    string
		objectName     string
		expectedConfig DBParametersConfig
		expectedError  string
		simulateError  bool // Nuevo campo para simular un error al leer la configuración
	}{
		{
			name: "Valid Configuration",
			fileContent: `
				{
					"database": {
						"config": {
							"max_idle_conns": 5,
							"max_open_conns": 10,
							"conn_max_idle_time": 3,
							"conn_max_lifetime": 12
						}
					}
				}
			`,
			objectName: "database",
			expectedConfig: DBParametersConfig{
				MaxIdleConns:    5,
				MaxOpenConns:    10,
				ConnMaxIdleTime: 3 * DefaultTimeDurationMinuteConfig,
				ConnMaxLifetime: 12 * DefaultTimeDurationMinuteConfig,
			},
			expectedError: "",
			simulateError: false,
		},
		{
			name:           "Invalid Configuration (Zero Values)",
			fileContent:    `{ "database": { "config": { "max_idle_conns": 0, "max_open_conns": 0, "conn_max_idle_time": 0, "conn_max_lifetime": 0 } } }`,
			objectName:     "database",
			expectedConfig: DBParametersConfig{},
			expectedError:  "Validation error: parameter at index 0 cannot be zero",
			simulateError:  false,
		},
		{
			name:           "Invalid Configuration (Missing Object)",
			fileContent:    `{}`,
			objectName:     "database",
			expectedConfig: DBParametersConfig{},
			expectedError:  "failed to read configuration from file",
			simulateError:  true,
		},
		{
			name:           "Error Reading Configuration",
			fileContent:    `{}`, // Contenido del archivo no importa, ya que simulará un error al leer la configuración
			objectName:     "database",
			expectedConfig: DBParametersConfig{},
			expectedError:  "failed to read configuration from file",
			simulateError:  true, // Indicar que se simulará un error al leer la configuración
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Crear un archivo temporal con el contenido proporcionado
			fileName := "temp_config.json"
			err := createTempFile(fileName, test.fileContent)
			assert.NoError(t, err, "Error creating temporary file")
			defer func(fileName string) {
				err := deleteTempFile(fileName)
				if err != nil {
					log.Println(err)
				}
			}(fileName)

			// Ejecutar la función LoadConfigFromJSON
			config, err := LoadConfigFromJSON(fileName, test.objectName)

			if err != nil {
				log.Printf("Test LoadConfigFromJSON Error: %v", err)
			}

			// Simular un error al leer la configuración si es necesario
			if test.simulateError {
				log.Println("Viper tendrá un error simulado al leer el archivo de configuración")
				viperReadInConfigFunc = func() error {
					return fmt.Errorf("failed to read configuration from file: simulated read error")
				}
				err = viperReadInConfigFunc()
			}

			// Verificar los resultados
			if test.expectedError != "" {
				assert.Error(t, err, "Expected an error")
				assert.Contains(t, err.Error(), test.expectedError, "Error message does not match")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, test.expectedConfig, config, "Config does not match expected values")
			}
		})
	}
}

// Función de utilidad para crear un archivo temporal con contenido
func createTempFile(fileName, content string) error {
	return ioutil.WriteFile(fileName, []byte(content), 0644)
}

// Función de utilidad para eliminar un archivo temporal
func deleteTempFile(fileName string) error {
	return os.Remove(fileName)
}

*/

func TestLoadConfigFromJSON(t *testing.T) {
	// Crear un archivo temporal para la prueba
	tmpfile, err := os.Create("config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Logf("Error removing file: %v", err)
		}
	}(tmpfile.Name()) // Eliminar el archivo temporal al finalizar la prueba

	// Escribir datos de configuración JSON en el archivo temporal
	configData := `{
		"database": {
			"config": {
				"max_idle_conns": 10,
				"max_open_conns": 20,
				"conn_max_idle_time": 30,
				"conn_max_lifetime": 40
			}
		}
	}`
	if _, err := tmpfile.Write([]byte(configData)); err != nil {
		t.Fatal(err)
	}
	err = tmpfile.Close()
	if err != nil {
		t.Logf("Error closing file: %v", err)
	}

	// Llamar a la función LoadConfigFromJSON con el archivo temporal
	config, err := LoadConfigFromJSON(tmpfile.Name(), "database.config")
	assert.NoError(t, err)

	// Verificar los valores devueltos por la función
	assert.Equal(t, 10, config.MaxIdleConns)
	assert.Equal(t, 20, config.MaxOpenConns)
	assert.Equal(t, 30, int(config.ConnMaxIdleTime/time.Minute))
	assert.Equal(t, 40, int(config.ConnMaxLifetime/time.Minute))
}

func TestLoadConfigFromJSON_InvalidFile(t *testing.T) {
	// Llamar a la función LoadConfigFromJSON con un archivo inexistente
	config, err := LoadConfigFromJSON("nonexistent.json", "database.config")
	assert.Error(t, err)
	assert.Equal(t, DBParametersConfig{}, config)
}

func TestLoadConfigFromJSON_InvalidValues(t *testing.T) {
	// Crear un archivo temporal con valores inválidos
	tmpfile, err := os.CreateTemp("", "invalid_config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Eliminar el archivo temporal al finalizar la prueba

	invalidConfigData := `{
		"database": {
			"config": {
				"max_idle_conns": 0, // Valor inválido
				"max_open_conns": 20,
				"conn_max_idle_time": 30,
				"conn_max_lifetime": 40
			}
		}
	}`
	if _, err := tmpfile.Write([]byte(invalidConfigData)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Llamar a la función LoadConfigFromJSON con valores inválidos
	config, err := LoadConfigFromJSON(tmpfile.Name(), "database.config")
	assert.Error(t, err)
	assert.Equal(t, DBParametersConfig{}, config)
}

func TestLoadConfigFromJSON_InvalidValuesValidation(t *testing.T) {
	// Crear un archivo temporal con valores inválidos
	tmpfile, err := os.Create("config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Logf("Error removing file: %v", err)
		}
	}(tmpfile.Name()) // Eliminar el archivo temporal al finalizar la prueba

	// Escribir datos de configuración JSON en el archivo temporal con valores inválidos
	configData := `{
		"database": {
			"config": {
				"max_idle_conns": 0,
				"max_open_conns": 20,
				"conn_max_idle_time": 30,
				"conn_max_lifetime": 40
			}
		}
	}`
	if _, err := tmpfile.Write([]byte(configData)); err != nil {
		t.Fatal(err)
	}
	err = tmpfile.Close()
	if err != nil {
		t.Logf("Error closing file: %v", err)
	}

	// Llamar a la función LoadConfigFromJSON con valores inválidos
	config, err := LoadConfigFromJSON(tmpfile.Name(), "database.config")

	// Verificar que se haya devuelto un error con el mensaje esperado
	expectedErrMsg := "Validation error: parameter at index 0 cannot be zero"
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErrMsg)
	assert.Equal(t, DBParametersConfig{}, config)
}

func TestEnsureNonEmptyQuery(t *testing.T) {
	quryExample := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
		{
			name:    "Error",
			args:    args{query: ""},
			wantErr: assert.Error,
		},
		{
			name:    "ok",
			args:    args{query: quryExample},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, EnsureNonEmptyQuery(tt.args.query), fmt.Sprintf("EnsureNonEmptyQuery(%v)", tt.args.query))
		})
	}
}
