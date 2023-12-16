package database

import (
	"errors"
	"github.com/spf13/viper"
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
