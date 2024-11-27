package my_grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	my_jwt "github.com/Bitummit/booking_auth/internal/jwt"
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
		CheckUserRole(ctx context.Context, token string) (string, error)
		LoginUser(ctx context.Context, user *models.User) (string, error)
		RegistrateUser(ctx context.Context, user models.User) (string, error)
		CheckIsAdmin(ctx context.Context, token string) error
		UpdateUserRole(ctx context.Context, username, role string) error
		GetUserFromToken(ctx context.Context, token string) (*models.User, error)
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
		return nil, status.Error(codes.Internal, err.Error())
	}
	response := &auth.RegistrationResponse{
		Token: token,
	}
	return response, nil
}

func (a *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user := models.User{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
	token, err := a.Service.LoginUser(ctx, &user)
	if err != nil {
		if errors.Is(err, postgresql.ErrorUserNotExists) || errors.Is(err, authService.ErrorIncorrectPassword){
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &auth.LoginResponse{
		Token: token,
	}
	return response, nil
}

func (a *AuthServer) CheckRole(ctx context.Context, req *auth.CheckRoleRequest) (*auth.CheckRoleResponse, error) {
	token := req.GetToken()
	role, err := a.Service.CheckUserRole(ctx, token)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response := &auth.CheckRoleResponse{
		Role: role,
	}
	return response, nil
}

func (a *AuthServer) IsAdmin(ctx context.Context, req *auth.CheckTokenRequest) (*auth.EmptyResponse, error) {
	token := req.GetToken()
	if err := a.Service.CheckIsAdmin(ctx, token); err != nil {
		a.Log.Error("error while login:", logger.Err(err))
		if errors.Is(err, my_jwt.ErrorTokenDuration) || errors.Is(err, postgresql.ErrorNotFound){
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	}
	return &auth.EmptyResponse{}, nil
}

func (a *AuthServer) UpdateUserRole(ctx context.Context, req *auth.UpdateUserRoleRequest) (*auth.EmptyResponse, error) {
	username := req.GetUsername()
	newRole := req.GetRole()
	if err := a.Service.UpdateUserRole(ctx, username, newRole); err != nil {
		a.Log.Error("error while granting new role: ", logger.Err(err))
		
		return nil, status.Error(codes.Internal, fmt.Sprintf("%v", err))
	} 
	return &auth.EmptyResponse{}, nil
}

func (a *AuthServer) GetUser(ctx context.Context, req *auth.GetUserRequest) (*auth.GetUserResponse, error) {
	token := req.GetToken()
	if token == "" {
		a.Log.Info("anon user")
		return nil, nil
	}
	user, err := a.Service.GetUserFromToken(ctx, token)
	if err != nil {
		a.Log.Error("error getting user: ", logger.Err(err))

		return nil, status.Error(codes.NotFound, fmt.Sprintf("%v", err))
	}

	res := auth.GetUserResponse{
		Id: user.Id,
		Username: user.Username,
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
	}
	return &res, nil
}