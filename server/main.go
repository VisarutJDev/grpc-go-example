package main

import (
	"context"
	"log"
	"net"

	pb "github.com/VisarutJDev/grpc-go-example/go-proto-output"
	pbAuth "github.com/VisarutJDev/grpc-go-example/go-proto-output/auth"
	pbPost "github.com/VisarutJDev/grpc-go-example/go-proto-output/post"
	"github.com/VisarutJDev/grpc-go-example/server/auth"
	"github.com/VisarutJDev/grpc-go-example/server/post"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement
type server struct {
	pb.UnimplementedGreeterServer
	pbAuth.UnimplementedAuthServiceServer
	pbPost.UnimplementedPostServiceServer
	DBClient *mongo.Client
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	authSrv := &auth.AuthService{DbClient: client}
	postSrv := &post.PostService{DbClient: client}

	pb.RegisterGreeterServer(s, &server{})

	pbAuth.RegisterAuthServiceServer(s, authSrv)
	pbPost.RegisterPostServiceServer(s, postSrv)
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
