package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver       = "postgres"
	dataSourceName = "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable"
)

var testQueries *Queries
var DBTestConnection *sql.DB

func TestMain(m *testing.M) {
	var err error

	DBTestConnection, err = sql.Open(dbDriver, dataSourceName)

	if err != nil {
		log.Fatal("Testing can be executed, DB connection failed")
	}

	testQueries = New(DBTestConnection)

	os.Exit(m.Run())

}
