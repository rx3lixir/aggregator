package server

import (
	"context"

	"github.com/rx3lixir/agg-api/agg-grpc/pb"
	"github.com/rx3lixir/agg-api/internal/db"
)

type Server struct {
	storer db.Storage
	pb.UnimplementedAggregatorServer
}

func NewServer(storer db.Storage) *Server {
	return &Server{
		storer: storer,
	}
}

func (s *Server) CreateEvent(ctx context.Context, req *pb.EventReq) (*pb.EventRes, error) {
	return &pb.EventRes{}, nil
}
