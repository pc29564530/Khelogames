package main

import (
	"database/sql"
	"khelogames/api"
	db "khelogames/db/sqlc"
	"khelogames/logger"
	"khelogames/util"

	_ "github.com/lib/pq"
)

func main() {
	newLogger := logger.NewLogger()
	config, err := util.LoadConfig(".")

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		newLogger.Error("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	log := logger.NewLogger()
	server, err := api.NewServer(config, store, log)
	if err != nil {
		newLogger.Error("Server does not created", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		newLogger.Error("cannot start server", err)
	}
}
