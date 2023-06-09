package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"khelogames/api"
	db "khelogames/db/sqlc"
	"khelogames/util"
	"log"
)

func main() {

	config, err := util.LoadConfig(".")

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Server does not created", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
