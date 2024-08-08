package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thaian1234/simplebank/utils"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	var config utils.Config
	env := utils.NewEnviroment("../..")
	config, err = env.GetConfig()

	if err != nil {
		log.Fatal("cannot load config file", err)
	}

	ctx := context.Background()
	testDB, err = pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}

	testQueries = New(testDB)

	exitCode := m.Run()

	testDB.Close()

	os.Exit(exitCode)
}
