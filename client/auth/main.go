package main

import (
	"context"
	"fmt"
	"log"

	"github.com/VisarutJDev/grpc-go-example/go-proto-output/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	authClient := auth.NewAuthServiceClient(conn)

	// register and get token
	resgisResp, err := authClient.Register(context.Background(), &auth.RegisterRequest{Username: "user1", Password: "password1"})
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}
	fmt.Printf("Register response: %s\n", resgisResp.Token)

	// Login and get token
	loginResp, err := authClient.Login(context.Background(), &auth.LoginRequest{Username: "user1", Password: "password1"})
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}
	fmt.Printf("Login response: %s\n", loginResp.Token)
}
