package service

import (
	"context"

	"github.com/Forum-service/Forum-Service/genproto/post"
	"github.com/Forum-service/Forum-Service/storage"
	"github.com/rs/zerolog/log"
)

// PostService implements the post.PostServiceServer interface.
type PostService struct {
	stg storage.StorageI
	post.UnimplementedPostServiceServer
}

// NewPostService creates a new PostService.
func NewPostService(stg storage.StorageI) *PostService {
	return &PostService{stg: stg}
}

// CreatePost creates a new post.
func (s *PostService) CreatePost(ctx context.Context, req *post.CreatePostRequest) (*post.CreatePostResponse, error) {
	log.Info().Msg("PostService: CreatePost called")

	resp, err := s.stg.Post().Create(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("PostService: Error creating post")
		return nil, err
	}
	return resp, nil
}

// GetPostById gets a post by its ID.
func (s *PostService) GetPost(ctx context.Context, req *post.GetPostRequest) (*post.GetPostResponse, error) {
	log.Info().Msg("PostService: GetPostById called")

	resp, err := s.stg.Post().GetById(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("PostService: Error getting post by ID")
		return nil, err
	}
	return resp, nil
}

// UpdatePost updates a post.
func (s *PostService) UpdatePost(ctx context.Context, req *post.UpdatePostRequest) (*post.UpdatePostResponse, error) {
	log.Info().Msg("PostService: UpdatePost called")

	resp, err := s.stg.Post().Update(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("PostService: Error updating post")
		return nil, err
	}
	return resp, nil
}

// DeletePost deletes a post.
func (s *PostService) DeletePost(ctx context.Context, req *post.DeletePostRequest) (*post.DeletePostResponse, error) {
	log.Info().Msg("PostService: DeletePost called")

	resp, err := s.stg.Post().Delete(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("PostService: Error deleting post")
		return nil, err
	}
	return resp, nil
}

// GetAllPosts lists posts with filtering and pagination.
func (s *PostService) GetAllPosts(ctx context.Context, req *post.GetAllPostsRequest) (*post.GetAllPostsResponse, error) {
	log.Info().Msg("PostService: GetAllPosts called")

	resp, err := s.stg.Post().GetAllPosts(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("PostService: Error getting all posts")
		return nil, err
	}
	return resp, nil
}
