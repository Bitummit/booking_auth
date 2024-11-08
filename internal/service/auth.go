package authService

import (
	"context"
	"errors"
	"fmt"

	my_jwt "github.com/Bitummit/booking_auth/internal/jwt"
	"github.com/Bitummit/booking_auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Storage Storage
}

type Storage interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
	GetUser(ctx context.Context, username string) (*models.User, error)
	SetUserRole(ctx context.Context, user models.User) error
}

var ErrorHashingPassword = errors.New("error while hashing password")


func New(storage Storage) *Service{
	return &Service{
		Storage: storage,
	}
}

func (s *Service) RegistrateUser(ctx context.Context, user models.User) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generating password: %w", ErrorHashingPassword)
	}
	user.PasswordHashed = hashedPass

	id, err := s.Storage.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("registration user: %w", err)
	}
	user.Id = id

	token, err := my_jwt.NewToken(user)
	if err != nil {
		return "", fmt.Errorf("registration user: %w", err)
	}

	return token, nil
}
