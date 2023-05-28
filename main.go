package main

import (
	"database/sql"
	"khelogames/api"
	db "khelogames/db/sqlc"
	"log"

	_ "github.com/lib/pq"
)

// github.com/lib/pq we cannot connect to the database
const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:Bharat@12@localhost:5432/khelogames?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("Server does not created", err)
	}
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
