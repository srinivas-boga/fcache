package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/srinivas-boga/fcache/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCacheServiceServer
	cache map[string]string
	mu    sync.RWMutex
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.cache[req.Key]
	if !ok {
		return &pb.GetResponse{Value: ""}, nil
	}
	return &pb.GetResponse{Value: value}, nil
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[req.Key] = req.Value
	return &pb.SetResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCacheServiceServer(s, &server{
		cache: make(map[string]string),
	})
	log.Println("Server is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
