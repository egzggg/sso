package main

import (
	"context"
	"log"
	"net"

	pb "github.com/egzggg/protos_sso/gen/go/sso"
	"google.golang.org/grpc"
)

type authServer struct {
	pb.UnimplementedAuthServer
}

func (s *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{UserId: 123}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":44044")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServer(s, &authServer{})

	log.Println("gRPC server running on :44044")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
