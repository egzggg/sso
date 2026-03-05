package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/egzggg/protos_sso/gen/go/sso"
)

func main() {
	cc, err := grpc.Dial(
		"127.0.0.1:44044",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	defer cc.Close()

	client := pb.NewAuthClient(cc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Register(ctx, &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}

	log.Println("Register succeeded, userID:", resp.GetUserId())
}
