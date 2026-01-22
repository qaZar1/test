package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Address  string
	Hostname string
	Port     uint64
	Database string
	User     string
	Password string
}

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

func New() *Config {
	if err := godotenv.Load("./config.env"); err != nil {
		panic(err)
	}

	port, err := strconv.ParseUint(os.Getenv(dbPort), base, size)
	if err != nil {
		panic(err)
	}

	address := ":" + os.Getenv(address)

	return &Config{
		Address:  address,
		Hostname: os.Getenv(dbHost),
		Port:     port,
		Database: os.Getenv(dbName),
		User:     os.Getenv(dbUser),
		Password: os.Getenv(dbPass),
	}
}
