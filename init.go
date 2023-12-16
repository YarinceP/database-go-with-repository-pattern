package main

import "github.com/spf13/viper"

func init() {
	// Configurar viper para leer configuraciones desde diferentes fuentes
	viper.SetConfigFile("config/config.json") // Nombre del archivo de configuración (puedes cambiarlo según tu preferencia)
	viper.SetConfigType("json")
	viper.AutomaticEnv() // Permitir la lectura de variables de entorno con prefijo DB_LIB_GO_

	// Lee la configuración del archivo si está presente
	if err := viper.ReadInConfig(); err != nil {
		// Puedes manejar el error si el archivo de configuración no está presente
	}

}
