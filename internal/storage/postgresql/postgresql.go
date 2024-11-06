package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	// "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type Storage struct {
	DB *pgxpool.Pool
}

var ErrorNotFound = errors.New("not found")
var ErrorUserExists = errors.New("user exists")


func New(ctx context.Context) (*Storage, error) {
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
	defer cancel()

	dbPath := os.Getenv("DB_URL")
	dbConn, err := pgxpool.New(ctx, dbPath)
	if err != nil {
		return nil, fmt.Errorf("creating pool: %w", err)
	}

	if err := dbConn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &Storage{DB: dbConn}, nil
}

