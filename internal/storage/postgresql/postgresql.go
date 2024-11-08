package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	// "github.com/jackc/pgx/v5"
	"github.com/Bitummit/booking_auth/internal/models"
	"github.com/jackc/pgx/v5"
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

func (s *Storage) CreateUser(ctx context.Context, user models.User) (int64, error) {
	var id int64
	args := pgx.NamedArgs{
		"username": user.Username,
	}
	resp, err := s.DB.Exec(ctx, InsertUserStmt, args)
	if err != nil {
		return 0, fmt.Errorf("checking user: unknown error %w", err)
	}
	if resp.RowsAffected() != 0 {
		return 0, fmt.Errorf("inserting user: %w", ErrorUserExists)
	}
	args = pgx.NamedArgs{
		"username": user.Username,
		"password": user.PasswordHashed,
		"email": user.Email,
		"firstName": user.FirstName,
		"lastName": user.LastName,
	}
	err = s.DB.QueryRow(ctx, InsertUserStmt, args).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("inserting user: unknown error %w", err)
	}

	return id, nil
}