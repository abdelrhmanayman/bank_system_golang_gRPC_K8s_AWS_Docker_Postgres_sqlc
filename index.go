package main

import (
	"banksystem/api"
	db "banksystem/db/sqlc"
	"banksystem/util"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbDriver       = "postgres"
	dataSourceName = "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable"
)

func main() {
	appConfig, err := util.LoadConfig()

	dbConn, err := sql.Open(appConfig.DbDriver, appConfig.DbSourceName)

	if err != nil {
		log.Fatal("Can not open database connection", err)
	}

	store := db.NewStore(dbConn)
	server := api.SetupRoutes(store)

	err = server.StartServer(":" + appConfig.Port)

	if err != nil {
		log.Fatal("Can't start the server")
	}

}
