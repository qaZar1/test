package main

import (
	"os"
	"strconv"

	postgres "github.com/qaZar1/test/microservices/wallet/internal/postgres"
)

func main() {
	const (
		address = "ADDRESS"
		dbHost  = "DATABASE_HOST"
		dbPort  = "DATABASE_PORT"
		dbName  = "DATABASE_NAME"
		dbUser  = "DATABASE_USER"
		dbPass  = "DATABASE_PASSWORD"
		base    = 10
		size    = 64
	)

	port, err := strconv.ParseUint(os.Getenv(dbPort), base, size)
	if err != nil {
		panic(err)
	}

	db := postgres.NewPostgres(postgres.Config{
		Hostname: os.Getenv(dbHost),
		Database: os.Getenv(dbName),
		User:     os.Getenv(dbUser),
		Password: os.Getenv(dbPass),
		Port:     port,
	})

	println(db)
}
