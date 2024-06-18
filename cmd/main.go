package main

import (
	"fmt"
	"net"

	"github.com/Forum-service/Forum-Service/config"
	"github.com/Forum-service/Forum-Service/genproto/category"
	"github.com/Forum-service/Forum-Service/genproto/comment"
	"github.com/Forum-service/Forum-Service/genproto/post"
	"github.com/Forum-service/Forum-Service/genproto/posttag"
	"github.com/Forum-service/Forum-Service/genproto/tag"
	"github.com/Forum-service/Forum-Service/service"
	"github.com/Forum-service/Forum-Service/storage/postgres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	pgStorage, err := postgres.NewStorage(&cfg)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to postgres: %v", err))
	}
	defer pgStorage.Close()

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		panic(fmt.Sprintf("Error listening on port %s: %v", "8082", err))
	}

	s := grpc.NewServer()

	// Register your gRPC services here
	category.RegisterCategoryServiceServer(s, service.NewCategoryService(pgStorage))
	tag.RegisterTagServiceServer(s, service.NewTagService(pgStorage))
	post.RegisterPostServiceServer(s, service.NewPostService(pgStorage))
	comment.RegisterCommentServiceServer(s, service.NewCommentService(pgStorage))
	posttag.RegisterPostTagServiceServer(s, service.NewPostTagService(pgStorage))

	reflection.Register(s) // Enable reflection for debugging

	fmt.Println("gRPC Server listening on", "8082")
	if err := s.Serve(lis); err != nil {
		panic(fmt.Sprintf("Failed to start gRPC server: %v", err))
	}
}
