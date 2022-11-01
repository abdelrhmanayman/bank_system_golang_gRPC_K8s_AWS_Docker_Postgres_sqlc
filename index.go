package main

import (
	"banksystem/api"
	db "banksystem/db/sqlc"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbDriver       = "postgres"
	dataSourceName = "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable"
)

func main() {
	dbConn, err := sql.Open(dbDriver, dataSourceName)

	if err != nil {

		log.Fatal("Can not open database connection", err)
	}

	store := db.NewStore(dbConn)
	server := api.SetupRoutes(store)

	err = server.StartServer(":8080")

	if err != nil {
		log.Fatal("Can't start the server")
	}

}
