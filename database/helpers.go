package database

func GetConnectionString() string {
	return "root:@tcp(localhost:3306)/db_lib_go?parseTime=true" //Utilizar package config json para obtener
}
