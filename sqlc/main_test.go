package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/crackz/simple-bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatal("Couldn't load config")
	}

	testDb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't Connect To DB : ", err)
	}

	testQueries = New(testDb)
	os.Exit(m.Run())
}
