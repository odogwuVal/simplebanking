package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/odogwuVal/simplebanking/api"
	db "github.com/odogwuVal/simplebanking/db/sqlc"
	"github.com/odogwuVal/simplebanking/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("could not connect to config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.Address)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
