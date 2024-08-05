package main

import (
	"context"
	"fmt"
	"log"

	"github.com/VisarutJDev/grpc-go-example/go-proto-output/auth"
	"github.com/VisarutJDev/grpc-go-example/go-proto-output/post"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	authClient := auth.NewAuthServiceClient(conn)
	postClient := post.NewPostServiceClient(conn)

	// Login and get token
	loginResp, err := authClient.Login(context.Background(), &auth.LoginRequest{Username: "user1", Password: "password1"})
	if err != nil {
		log.Fatalf("could not login: %v", err)
	}
	fmt.Printf("Login response: %s\n", loginResp.Token)

	// Create a new post with the token
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", loginResp.Token)
	createPostResp, err := postClient.CreatePost(ctx, &post.CreatePostRequest{Content: "This is my first post!"})
	if err != nil {
		log.Fatalf("could not create post: %v", err)
	}
	fmt.Printf("Create post response: %s\n", createPostResp.Message)

	// Get the post
	getPostResp, err := postClient.GetPost(ctx, &post.GetPostRequest{Id: createPostResp.Post.Id})
	if err != nil {
		log.Fatalf("could not get post: %v", err)
	}
	log.Printf("get Post: %s", getPostResp)

	// Get the posts
	getPostsResp, err := postClient.GetPosts(ctx, &post.GetPostsRequest{})
	if err != nil {
		log.Fatalf("could not get post: %v", err)
	}
	log.Printf("get many Posts: %s", getPostsResp)

	// update the post
	updateResp, err := postClient.UpdatePost(ctx, &post.UpdatePostRequest{Id: createPostResp.Post.Id, Content: "Update content"})
	if err != nil {
		log.Fatalf("could not get post: %v", err)
	}
	log.Printf("Update post: %s", updateResp)

	// delete the post
	deleteResp, err := postClient.DeletePost(ctx, &post.DeletePostRequest{Id: createPostResp.Post.Id})
	if err != nil {
		log.Fatalf("could not get post: %v", err)
	}
	log.Printf("Delete post: %s", deleteResp)
}
