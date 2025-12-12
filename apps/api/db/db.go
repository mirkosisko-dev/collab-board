package pool

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mirkosisko-dev/api/config"
	sqlc "github.com/mirkosisko-dev/api/db/sqlc"
)

type Database struct {
	Pool  *pgxpool.Pool
	Query *sqlc.Queries
}

func NewPostgreSQLStorage() (*Database, error) {
	pool, err := pgxpool.New(context.Background(), config.Envs.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	query := sqlc.New(pool)

	db := Database{
		Pool:  pool,
		Query: query,
	}

	fmt.Println("Successfully connected!")

	return &db, nil
}
