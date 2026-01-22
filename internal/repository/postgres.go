package repository

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/qaZar1/test/autogen"
)

//go:generate mockgen -source=postgres.go -package mocks -destination ../../autogen/mocks/postgres.go
type IRepository interface {
	GetWallet(walletID string) (*autogen.Wallet, error)
	UpsertWallet(wallet autogen.WalletUpdate) error
	Close() error
}

type Config struct {
	Hostname string
	Port     uint64
	Database string
	User     string
	Password string
}

type postgres struct {
	db *sqlx.DB
}

func NewRepository(cfg Config) IRepository {
	return &postgres{db: newPostgres(cfg)}
}

func newPostgres(cfg Config) *sqlx.DB {
	const driver = "pgx"

	pattern := fmt.Sprintf(
		"host=%s port=%d database=%s user=%s password=%s sslmode=disable",
		cfg.Hostname, cfg.Port, cfg.Database, cfg.User, cfg.Password,
	)

	config, err := pgx.ParseConfig(pattern)
	if err != nil {
		panic(err)
	}

	connection := stdlib.RegisterConnConfig(config)
	db, err := sqlx.Open(driver, connection)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)

	return db
}

func (pg *postgres) Close() error {
	return pg.db.Close()
}
