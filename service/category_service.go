package service

import (
	"context"

	"github.com/Forum-service/Forum-Service/genproto/category"
	"github.com/Forum-service/Forum-Service/storage"
	"github.com/rs/zerolog/log"
)

// CategoryService implements the category.CategoryServiceServer interface.
type CategoryService struct {
	stg storage.StorageI
	category.UnimplementedCategoryServiceServer
}

// NewCategoryService creates a new CategoryService.
func NewCategoryService(stg storage.StorageI) *CategoryService {
	return &CategoryService{stg: stg}
}

// CreateCategory creates a new category.
func (s *CategoryService) CreateCategory(ctx context.Context, req *category.CreateCategoryRequest) (*category.CreateCategoryResponse, error) {
	log.Info().Msg("CategoryService: CreateCategory called")

	resp, err := s.stg.Category().Create(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("CategoryService: Error creating category")
		return nil, err
	}
	return resp, nil
}

// GetCategoryById gets a category by its ID.
func (s *CategoryService) GetCategory(ctx context.Context, req *category.GetCategoryRequest) (*category.GetCategoryResponse, error) {
	log.Info().Msg("CategoryService: GetCategoryById called")

	resp, err := s.stg.Category().GetById(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("CategoryService: Error getting category by ID")
		return nil, err
	}
	return resp, nil
}

// UpdateCategory updates a category.
func (s *CategoryService) UpdateCategory(ctx context.Context, req *category.UpdateCategoryRequest) (*category.UpdateCategoryResponse, error) {
	log.Info().Msg("CategoryService: UpdateCategory called")

	resp, err := s.stg.Category().Update(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("CategoryService: Error updating category")
		return nil, err
	}
	return resp, nil
}

// DeleteCategory deletes a category.
func (s *CategoryService) DeleteCategory(ctx context.Context, req *category.DeleteCategoryRequest) (*category.DeleteCategoryResponse, error) {
	log.Info().Msg("CategoryService: DeleteCategory called")

	resp, err := s.stg.Category().Delete(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("CategoryService: Error deleting category")
		return nil, err
	}
	return resp, nil
}

// GetAllCategories lists categories with pagination.
func (s *CategoryService) GetAllCategories(ctx context.Context, req *category.GetAllCategoriesRequest) (*category.GetAllCategoriesResponse, error) {
	log.Info().Msg("CategoryService: GetAllCategories called")

	resp, err := s.stg.Category().GetAllCategories(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("CategoryService: Error getting all categories")
		return nil, err
	}
	return resp, nil
}
