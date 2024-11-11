package main

import (
	"context"
	"log"
	"net"

	"fcache"

	pb "fcache/proto/cacheService"

	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedCacheServiceServer
	cache *fcache.Cache
}

func (s *Server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {

	// convert string to byte array
	key := []byte(req.Key)
	value, err := s.cache.Get(key)
	if err != nil {
		return &pb.GetResponse{Value: ""}, nil
	}
	return &pb.GetResponse{Value: string(value)}, nil
}

func (s *Server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {

	// convert string to byte array
	key := []byte(req.Key)
	value := []byte(req.Value)

	s.cache.Set(key, value)
	return &pb.SetResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCacheServiceServer(s, &Server{
		cache: fcache.NewCache(),
	})

	log.Println("Server is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
