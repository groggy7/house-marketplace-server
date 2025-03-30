package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func CreateDBConnection() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		panic("DB_URL environment variable is not set")
	}

	pgxConfig, err := pgx.ParseConfig(dbUrl)
	if err != nil {
		panic(err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), pgxConfig)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
}

func main() {
	CreateDBConnection()
}