package run

import (
	"context"
	"net"
	"os/signal"
	"sync"

	"syscall"

	my_grpc "github.com/Bitummit/booking_auth/internal/api/grpc"
	"github.com/Bitummit/booking_auth/internal/interceptors"
	"github.com/Bitummit/booking_auth/internal/storage/postgresql"
	"github.com/Bitummit/booking_auth/pkg/config"
	"github.com/Bitummit/booking_auth/pkg/logger"
	auth "github.com/Bitummit/booking_auth/pkg/proto"
	"google.golang.org/grpc"
)


func Run() {
	ctx, stop  := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	wg := &sync.WaitGroup{}

	cfg := config.InitConfig()
	log := logger.NewLogger()
	log.Info("Initializing config success")

	log.Info("Connecting database ...")
	storage, err := postgresql.New(ctx)
	if err != nil {
		log.Error("Error connecting to DB", logger.Err(err))
		return
	}
	log.Info("Connecting database success")

	wg.Add(1)
	log.Info("Starting server ...")
	server := my_grpc.New(log, cfg, storage)
	go startServer(ctx, wg, server) 

	<-ctx.Done()
	wg.Wait()
	storage.DB.Close()
	log.Info("Database stopped")
}

func startServer(ctx context.Context, wg *sync.WaitGroup, server *my_grpc.AuthServer) {	
	listener, err := net.Listen("tcp", server.Cfg.GrpcAddress)
	if err != nil {
		server.Log.Error("failed to listen", logger.Err(err))
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.UnaryLogRequest(server.Log),
		),
	}
	grpcServer := grpc.NewServer(opts...)
	auth.RegisterAuthServer(grpcServer, server)

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			server.Log.Error("error starting server", logger.Err(err))
		}
	}()
	
	<-ctx.Done()
	defer wg.Done()
	grpcServer.GracefulStop()
	server.Log.Info("Server stopped")
}
