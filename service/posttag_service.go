package service

import (
	"context"

	"github.com/Forum-service/Forum-Service/genproto/posttag"
	"github.com/Forum-service/Forum-Service/storage"
	"github.com/rs/zerolog/log"
)

// PostTagService implements the posttag.PostTagServiceServer interface.
type PostTagService struct {
	stg storage.StorageI
	posttag.UnimplementedPostTagServiceServer
}

// NewPostTagService creates a new PostTagService.
func NewPostTagService(stg storage.StorageI) *PostTagService {
	return &PostTagService{stg: stg}
}

// CreatePostTag creates a new post-tag association.
func (s *PostTagService) CreatePostTag(ctx context.Context, req *posttag.CreatePostTagRequest) (*posttag.CreatePostTagResponse, error) {
	log.Info().Msg("PostTagService: CreatePostTag called")

	resp, err := s.stg.PostTag().Create(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("PostTagService: Error creating post-tag association")
		return nil, err
	}
	return resp, nil
}

// DeletePostTag deletes a post-tag association.
func (s *PostTagService) DeletePostTag(ctx context.Context, req *posttag.DeletePostTagRequest) (*posttag.DeletePostTagResponse, error) {
	log.Info().Msg("PostTagService: DeletePostTag called")

	resp, err := s.stg.PostTag().Delete(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("PostTagService: Error deleting post-tag association")
		return nil, err
	}
	return resp, nil
}

// GetAllPostTags lists post-tag associations with filtering and pagination.
func (s *PostTagService) GetAllPostTags(ctx context.Context, req *posttag.GetAllPostTagsRequest) (*posttag.GetAllPostTagsResponse, error) {
	log.Info().Msg("PostTagService: GetAllPostTags called")

	resp, err := s.stg.PostTag().GetAllPostTags(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("PostTagService: Error getting all post-tag associations")
		return nil, err
	}
	return resp, nil
}

// GetPostsByTag retrieves posts associated with a specific tag ID.
func (s *PostTagService) GetPostsByTag(ctx context.Context, req *posttag.GetPostsByTagRequest) (*posttag.GetPostsByTagResponse, error) {
	log.Info().Msg("PostTagService: GetPostsByTag called")

	resp, err := s.stg.PostTag().GetPostsByTag(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("PostTagService: Error getting posts by tag ID")
		return nil, err
	}
	return resp, nil
}
