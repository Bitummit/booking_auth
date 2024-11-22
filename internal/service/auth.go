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
	GetUser(ctx context.Context, user *models.User) (*models.User, error)
	SetUserRole(ctx context.Context, username, role string) error
}

var ErrorHashingPassword = errors.New("error while hashing password")
var ErrorIncorrectPassword = errors.New("invalid password")
var ErrorNotAdmin = errors.New("only admin allowed")


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

func (s *Service) LoginUser(ctx context.Context, user *models.User) (string, error) {
	var token string

	user, err := s.Storage.GetUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("login user: %w", err)
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHashed, []byte(user.Password)); if err != nil {
		return "", fmt.Errorf("login user: %w", ErrorIncorrectPassword)
	}

	token, err = my_jwt.NewToken(*user)
	if err != nil {
		return "", fmt.Errorf("registration user: %w", err)
	}

	return token, nil
}

func (s *Service) CheckUserRole(ctx context.Context, token string) (string, error) {

	user, err := s.checkUser(ctx, token)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return user.Role, nil
}

func (s *Service) CheckIsAdmin(ctx context.Context, token string) error {
	user, err := s.checkUser(ctx, token)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if user.Role != "admin" {
		return ErrorNotAdmin
	}

	return nil
}

func (s *Service) checkUser(ctx context.Context, token string) (*models.User, error) {
	user, err := my_jwt.ParseToken(token)
	if err != nil {
		return nil, fmt.Errorf("check user token: %w", err)
	}

	_, err = s.Storage.GetUser(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("check user token: %w", err)
	}
	return &user, nil
}

func (s *Service) UpdateUserRole(ctx context.Context, username, role string) error {
	if err := s.Storage.SetUserRole(ctx, username, role); err != nil {
		return fmt.Errorf("updating user role: %w", err)
	}
	return nil
}
