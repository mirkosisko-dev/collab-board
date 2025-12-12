package sqlc

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:password@localhost:5433/collab_board?sslmode=disable"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal("Couldn't connect to database:", err)
	}

	testQueries = New(pool)

	os.Exit(m.Run())
}
