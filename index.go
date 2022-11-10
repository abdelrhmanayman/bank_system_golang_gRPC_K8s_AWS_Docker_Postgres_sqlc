package main

import (
	"banksystem/api"
	db "banksystem/db/sqlc"
	"banksystem/util"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	appConfig, err := util.LoadConfig()

	if err != nil {
		log.Fatal("Can not load application configs", err)
	}

	dbConn, err := sql.Open(appConfig.DbDriver, appConfig.DbSourceName)

	if err != nil {
		log.Fatal("Can not open database connection", err)
	}

	store := db.NewStore(dbConn)
	server, err := api.SetupRoutes(appConfig, store)

	if err != nil {
		log.Fatal("Can not setup routes", err)
	}

	err = server.StartServer(":" + appConfig.Port)

	if err != nil {
		log.Fatal("Can't start the server")
	}

}
