package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// DB is a global variable to hold db connection
var DB *sql.DB

// ConnectDB opens a connection to the database
func ConnectDB() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_lib_go")
	if err != nil {
		panic(err.Error())
	}

	DB = db
}

// DisconnectDB cierra la conexión a la base de datos
func DisconnectDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			panic(err.Error())
		}
	}
}

func GetConnectionInstance() *sql.DB {
	return DB
}
