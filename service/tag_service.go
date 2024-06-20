package service

import (
	"context"

	"github.com/Forum-service/Forum-Service/genproto/tag"
	"github.com/Forum-service/Forum-Service/storage"
	"github.com/rs/zerolog/log"
)

// TagService implements the tag.TagServiceServer interface.
type TagService struct {
	stg storage.StorageI
	tag.UnimplementedTagServiceServer
}

// NewTagService creates a new TagService.
func NewTagService(stg storage.StorageI) *TagService {
	return &TagService{stg: stg}
}

// CreateTag creates a new tag.
func (s *TagService) CreateTag(ctx context.Context, req *tag.CreateTagRequest) (*tag.CreateTagResponse, error) {
	log.Info().Msg("TagService: CreateTag called")

	resp, err := s.stg.Tag().Create(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("TagService: Error creating tag")
		return nil, err
	}
	return resp, nil
}

// GetTagById gets a tag by its ID.
func (s *TagService) GetTag(ctx context.Context, req *tag.GetTagRequest) (*tag.GetTagResponse, error) {
	log.Info().Msg("TagService: GetTagById called")
	resp, err := s.stg.Tag().GetById(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("TagService: Error getting tag by ID")
		return nil, err
	}
	return resp, nil
}

// UpdateTag updates a tag.
func (s *TagService) UpdateTag(ctx context.Context, req *tag.UpdateTagRequest) (*tag.UpdateTagResponse, error) {
	log.Info().Msg("TagService: UpdateTag called")

	resp, err := s.stg.Tag().Update(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("TagService: Error updating tag")
		return nil, err
	}
	return resp, nil
}

// DeleteTag deletes a tag.
func (s *TagService) DeleteTag(ctx context.Context, req *tag.DeleteTagRequest) (*tag.DeleteTagResponse, error) {
	log.Info().Msg("TagService: DeleteTag called")

	resp, err := s.stg.Tag().Delete(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("TagService: Error deleting tag")
		return nil, err
	}
	return resp, nil
}

// GetAllTags lists tags with pagination.
func (s *TagService) GetAllTags(ctx context.Context, req *tag.GetAllTagsRequest) (*tag.GetAllTagsResponse, error) {
	log.Info().Msg("TagService: GetAllTags called")

	resp, err := s.stg.Tag().GetAllTags(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("TagService: Error getting all tags")
		return nil, err
	}
	return resp, nil
}

// GetFamousTags gets famous tags with optional pagination and sorting.
func (s *TagService) GetFamousTags(ctx context.Context, req *tag.GetFamousTagsReq) (*tag.GetFamousTagsRes, error) {
	log.Info().Msg("TagService: GetFamousTags called")

	resp, err := s.stg.Tag().GetFamousTags(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("TagService: Error getting famous tags")
		return nil, err
	}
	return resp, nil
}
