package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thaian1234/simplebank/api"
	db "github.com/thaian1234/simplebank/db/sqlc"
	"github.com/thaian1234/simplebank/utils"
)

func main() {
	config, err := utils.NewEnvironment(".").GetConfig()
	if err != nil {
		log.Fatal("cannot load config file", err)
	}
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("cannot create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
