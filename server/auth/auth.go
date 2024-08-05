package auth

import (
	"context"
	"errors"
	"time"

	pb "github.com/VisarutJDev/grpc-go-example/go-proto-output/auth"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/metadata"
)

var jwtKey = []byte("my_secret_key")

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	DbClient *mongo.Client
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	usersCollection := s.DbClient.Database("socialdb").Collection("users")

	var existingUser bson.M
	err := usersCollection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	_, err = usersCollection.InsertOne(ctx, bson.M{
		"username": req.Username,
		"password": req.Password,
	})
	if err != nil {
		return nil, err
	}

	token, err := generateJWT(req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{Token: token}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	usersCollection := s.DbClient.Database("socialdb").Collection("users")

	var user bson.M
	err := usersCollection.FindOne(ctx, bson.M{"username": req.Username, "password": req.Password}).Decode(&user)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := generateJWT(req.Username)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: token}, nil
}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func Authenticate(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}
	token := md["authorization"]
	if len(token) == 0 {
		return "", errors.New("missing token")
	}
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token[0], claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		return "", errors.New("invalid token")
	}
	return claims.Username, nil
}
