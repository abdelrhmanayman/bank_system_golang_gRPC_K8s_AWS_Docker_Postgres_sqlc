package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const (
	dbDriver       = "postgres"
	dataSourceName = "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dataSourceName)

	if err != nil {
		log.Fatal("Testing can be executed, DB connection failed")
	}

	testQueries = New(conn)

	os.Exit(m.Run())

}
