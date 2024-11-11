package main

import (
	"context"
	"log"
	"net"

	"fcache"

	pb "github.com/srinivas-boga/fcache/proto"

	"google.golang.org/grpc"
)

type server struct {
	cache *fcache.Cache
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {

	value, err := s.cache.Get(req.Key)
	if err != nil {
		return &pb.GetResponse{Value: ""}, nil
	}
	return &pb.GetResponse{Value: value}, nil
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {

	s.cache.Set(req.Key, req.Value)
	return &pb.SetResponse{Success: true}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCacheServiceServer(s, &server{
		cache: fcache.NewCache(),
	})
	log.Println("Server is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
