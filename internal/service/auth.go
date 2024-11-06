package auth

import (
	"context"

	"github.com/Bitummit/booking_auth/internal/models"
)

type Service struct {
	DB Storage
}

type Storage interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	SetUserRole(ctx context.Context, user models.User) error
}

func New(storage Storage) *Service{
	return &Service{
		DB: storage,
	}
}

