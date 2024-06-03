package model

import (
	"database/sql"
	"log"
	"os"
	"usermanagement/utils"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectDb() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	username := utils.DotEnvVariable(utils.USERNAME)
	password := utils.DotEnvVariable(utils.PASSWORD)
	host := utils.DotEnvVariable(utils.HOST)
	port := utils.DotEnvVariable(utils.PORT)
	dbName := utils.DotEnvVariable(utils.DATABASE_NAME)

	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName
	client, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database \n", err)
		os.Exit(1)
	}
	log.Println("Connect")
	db = client
}
