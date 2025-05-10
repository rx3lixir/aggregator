package main

import (
	"context"
	"net"
	"os"

	"github.com/charmbracelet/log"
	"github.com/rx3lixir/agg-api/agg-grpc/pb"
	"github.com/rx3lixir/agg-api/agg-grpc/server"
	"github.com/rx3lixir/agg-api/config"
	"github.com/rx3lixir/agg-api/internal/db"

	"github.com/ianschenck/envflag"
	"google.golang.org/grpc"
)

func main() {

	var (
		svcAddr = envflag.String("SVC_ADDR", "0.0.0.0:9091", "address where grpc service is listening")
	)

	// Базовый контекст приложения
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		os.Exit(1)
	}

	pool, err := db.CreatePostgresPool(ctx, cfg)
	if err != nil {
		os.Exit(1)
	}
	store := db.NewPosgresStore(pool)

	srv := server.NewServer(store)

	grpcServer := grpc.NewServer()
	pb.RegisterAggregatorServer(grpcServer, srv)

	listener, err := net.Listen("tcp", *svcAddr)
	if err != nil {
		log.Error("listener failed", "error", err)
	}

	log.Info("server listening", "address", *svcAddr)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Error("failed to serve", "error", err)
		os.Exit(1)
	}
}
