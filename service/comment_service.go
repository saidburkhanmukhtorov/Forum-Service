package service

import (
	"context"

	"github.com/Forum-service/Forum-Service/genproto/comment"
	"github.com/Forum-service/Forum-Service/storage"
	"github.com/rs/zerolog/log"
)

// CommentService implements the comment.CommentServiceServer interface.
type CommentService struct {
	stg storage.StorageI
	comment.UnimplementedCommentServiceServer
}

// NewCommentService creates a new CommentService.
func NewCommentService(stg storage.StorageI) *CommentService {
	return &CommentService{stg: stg}
}

// CreateComment creates a new comment.
func (s *CommentService) CreateComment(ctx context.Context, req *comment.CreateCommentRequest) (*comment.CreateCommentResponse, error) {
	log.Info().Msg("CommentService: CreateComment called")

	resp, err := s.stg.Comment().Create(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("CommentService: Error creating comment")
		return nil, err
	}
	return resp, nil
}

// GetCommentById gets a comment by its ID.
func (s *CommentService) GetComment(ctx context.Context, req *comment.GetCommentRequest) (*comment.GetCommentResponse, error) {
	log.Info().Msg("CommentService: GetCommentById called")

	resp, err := s.stg.Comment().GetById(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("CommentService: Error getting comment by ID")
		return nil, err
	}
	return resp, nil
}

// UpdateComment updates a comment.
func (s *CommentService) UpdateComment(ctx context.Context, req *comment.UpdateCommentRequest) (*comment.UpdateCommentResponse, error) {
	log.Info().Msg("CommentService: UpdateComment called")

	resp, err := s.stg.Comment().Update(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("CommentService: Error updating comment")
		return nil, err
	}
	return resp, nil
}

// DeleteComment deletes a comment.
func (s *CommentService) DeleteComment(ctx context.Context, req *comment.DeleteCommentRequest) (*comment.DeleteCommentResponse, error) {
	log.Info().Msg("CommentService: DeleteComment called")

	resp, err := s.stg.Comment().Delete(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("CommentService: Error deleting comment")
		return nil, err
	}
	return resp, nil
}

// GetAllComments lists comments with filtering and pagination.
func (s *CommentService) GetAllComments(ctx context.Context, req *comment.GetAllCommentsRequest) (*comment.GetAllCommentsResponse, error) {
	log.Info().Msg("CommentService: GetAllComments called")

	resp, err := s.stg.Comment().GetAllComments(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("CommentService: Error getting all comments")
		return nil, err
	}
	return resp, nil
}
