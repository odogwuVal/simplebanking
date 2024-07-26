package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/odogwuVal/simplebanking/api"
	db "github.com/odogwuVal/simplebanking/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:root@127.0.0.1:5432/simple_bank?sslmode=disable"
	address  = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(*store)

	err = server.Start(address)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}
}
