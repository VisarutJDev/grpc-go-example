package post

import (
	"context"
	"errors"

	pb "github.com/VisarutJDev/grpc-go-example/go-proto-output/post"
	"github.com/VisarutJDev/grpc-go-example/server/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostService struct {
	pb.UnimplementedPostServiceServer
	DbClient *mongo.Client
}

func (s *PostService) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	username, err := auth.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	postsCollection := s.DbClient.Database("socialdb").Collection("posts")

	result, err := postsCollection.InsertOne(ctx, bson.M{
		"content": req.Content,
		"author":  username,
	})
	if err != nil {
		return nil, err
	}
	var post Post
	err = postsCollection.FindOne(ctx, bson.M{"_id": result.InsertedID.(primitive.ObjectID)}).Decode(&post)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return &pb.CreatePostResponse{
		Message: "Post created successfully",
		Post: &pb.Post{
			Id:      post.ID.Hex(),
			Content: post.Content,
			Author:  post.Author,
		}}, nil
}

func (s *PostService) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	postsCollection := s.DbClient.Database("socialdb").Collection("posts")

	var post Post
	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, errors.New("error while trying to parse object id bson")
	}

	err = postsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&post)
	if err != nil {
		return nil, errors.New("post not found")
	}

	return &pb.GetPostResponse{Post: &pb.Post{
		Id:      post.ID.Hex(),
		Content: post.Content,
		Author:  post.Author,
	}}, nil
}

func (s *PostService) GetPosts(ctx context.Context, req *pb.GetPostsRequest) (*pb.GetPostsResponse, error) {
	postsCollection := s.DbClient.Database("socialdb").Collection("posts")

	cursor, err := postsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*pb.Post
	for cursor.Next(ctx) {
		var post Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, &pb.Post{
			Id:      post.ID.Hex(),
			Content: post.Content,
			Author:  post.Author,
		})
	}

	return &pb.GetPostsResponse{Posts: posts}, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	username, err := auth.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	postsCollection := s.DbClient.Database("socialdb").Collection("posts")
	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, errors.New("error while trying to parse object id bson")
	}
	filter := bson.M{"_id": objID, "author": username}

	update := bson.M{
		"$set": bson.M{"content": req.Content},
	}

	result, err := postsCollection.UpdateOne(ctx, filter, update)
	if err != nil || result.MatchedCount == 0 {
		return nil, errors.New("not authorized to update this post or post not found")
	}
	var post Post
	err = postsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&post)
	if err != nil {
		return nil, errors.New("post not found")
	}
	return &pb.UpdatePostResponse{
		Message: "Post updated successfully",
		Post: &pb.Post{
			Id:      post.ID.Hex(),
			Content: post.Content,
			Author:  post.Author,
		},
	}, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	username, err := auth.Authenticate(ctx)
	if err != nil {
		return nil, err
	}

	postsCollection := s.DbClient.Database("socialdb").Collection("posts")
	objID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, errors.New("error while trying to parse object id bson")
	}
	filter := bson.M{"_id": objID, "author": username}

	result, err := postsCollection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount == 0 {
		return nil, errors.New("not authorized to delete this post or post not found")
	}

	return &pb.DeletePostResponse{Message: "Post deleted successfully"}, nil
}
