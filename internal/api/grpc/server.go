package my_grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Bitummit/booking_auth/internal/models"
	authService "github.com/Bitummit/booking_auth/internal/service"
	"github.com/Bitummit/booking_auth/internal/storage/postgresql"
	"github.com/Bitummit/booking_auth/pkg/config"
	"github.com/Bitummit/booking_auth/pkg/logger"
	auth "github.com/Bitummit/booking_auth/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type (
	AuthServer struct {
		Cfg *config.Config
		Log *slog.Logger
		Service Service
		auth.UnimplementedAuthServer
	}

	Service interface {
		// CheckTokenUser(token string) error
		// CheckRoleUser(token string) (string, error)
		// LoginUser(cusername string, password string) (*string, error)
		RegistrateUser(ctx context.Context, user models.User) (string, error)
	}
)


func New(log *slog.Logger, cfg *config.Config, storage authService.Storage) *AuthServer {
	service := authService.New(storage)

	return &AuthServer{
		Cfg: cfg,
		Log: log,
		Service: service,
	}
}

func (a *AuthServer) Registration(ctx context.Context, req *auth.RegistrationRequest) (*auth.RegistrationResponse, error)  {
	// register client
	user := models.User{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email: req.GetEmail(),
		FirstName: req.GetFirstName(),
		LastName: req.GetLastName(),
		Role: "client",
	}
	token, err := a.Service.RegistrateUser(ctx, user)
	if err != nil {
		a.Log.Error("error while register user:", logger.Err(err))
		if errors.Is(err, postgresql.ErrorUserExists) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}
	response := &auth.RegistrationResponse{
		Token: token,
	}
	return response, nil
}
