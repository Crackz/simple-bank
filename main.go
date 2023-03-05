package main

import (
	"database/sql"
	"log"

	"github.com/crackz/simple-bank/api"
	db "github.com/crackz/simple-bank/db/sqlc"
	"github.com/crackz/simple-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Couldn't load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't Connect To DB : ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Couldn't Start Server : ", err)
	}
}
