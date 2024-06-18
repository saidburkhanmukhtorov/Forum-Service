package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	// Update with your actual package path
	"github.com/Forum-service/Forum-Service/genproto/category"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// ErrCategoryNotFound is returned when a category is not found.
var ErrCategoryNotFound = errors.New("category not found")

// CategoryDb provides database operations for categories.
type CategoryDb struct {
	Db *pgx.Conn
}

// NewCategory creates a new instance of CategoryDb.
func NewCategory(db *pgx.Conn) *CategoryDb {
	return &CategoryDb{Db: db}
}

// Create creates a new category in the database.
func (cDb *CategoryDb) Create(ctx context.Context, req *category.CreateCategoryRequest) (*category.CreateCategoryResponse, error) {
	categoryID := uuid.New().String()
	query := `
		INSERT INTO 
			categories (
				id,
				name
			) 
		VALUES (
				$1, 
				$2
			)
		RETURNING 
			id,
			name,
			created_at,
			updated_at
	`
	var (
		dbCategory category.Category
		createdAt  time.Time
		updatedAt  time.Time
	)

	err := cDb.Db.QueryRow(ctx, query,
		categoryID,
		req.Name,
	).Scan(
		&dbCategory.Id,
		&dbCategory.Name,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		log.Error().Err(err).Msg("Error creating category")
		return nil, err
	}

	dbCategory.CreatedAt = createdAt.Format(time.RFC3339)
	dbCategory.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &category.CreateCategoryResponse{Category: &dbCategory}, nil
}

// GetById gets a category by its ID.
func (cDb *CategoryDb) GetById(ctx context.Context, req *category.GetCategoryRequest) (*category.GetCategoryResponse, error) {
	var (
		dbCategory category.Category
		createdAt  time.Time
		updatedAt  time.Time
	)

	query := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM 
			categories 
		WHERE 
			id = $1
		AND
			deleted_at = 0
	`
	err := cDb.Db.QueryRow(ctx, query, req.Id).Scan(
		&dbCategory.Id,
		&dbCategory.Name,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Category not found")
			return nil, ErrCategoryNotFound
		}
		log.Error().Err(err).Msg("Error getting category by ID")
		return nil, err

	}
	dbCategory.CreatedAt = createdAt.Format(time.RFC3339)
	dbCategory.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &category.GetCategoryResponse{Category: &dbCategory}, nil
}

// Update updates an existing category in the database.
func (cDb *CategoryDb) Update(ctx context.Context, req *category.UpdateCategoryRequest) (*category.UpdateCategoryResponse, error) {
	var args []interface{}
	count := 1
	query := `
		UPDATE 
			categories 
		SET `
	filter := ``

	if len(req.Name) > 0 {
		filter += fmt.Sprintf(" name = $%d, ", count)
		args = append(args, req.Name)
		count++
	}

	if filter == "" {
		log.Error().Msg("No fields provided for update.")
		return nil, errors.New("no fields provided for update")
	}

	filter += fmt.Sprintf(`
			updated_at = NOW()
		WHERE
			id = $%d
		AND 
			deleted_at = 0
			`, count)

	args = append(args, req.Id)
	query += filter

	_, err := cDb.Db.Exec(ctx, query, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Category not found")
			return nil, ErrCategoryNotFound
		}
		log.Error().Err(err).Msg("Error updating category")
		return nil, err
	}

	// After successful update, get the updated category from the database
	updatedCategory, err := cDb.GetById(ctx, &category.GetCategoryRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &category.UpdateCategoryResponse{Category: updatedCategory.Category}, nil
}

// Delete soft deletes a category by setting its deleted_at field to the current time.
func (cDb *CategoryDb) Delete(ctx context.Context, req *category.DeleteCategoryRequest) (*category.DeleteCategoryResponse, error) {
	query := `
		UPDATE 
			categories 
		SET 
			deleted_at = $1
		WHERE 
			id = $2
	`
	_, err := cDb.Db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Category not found")
			return nil, ErrCategoryNotFound
		}
		log.Error().Err(err).Msg("Error soft deleting category")
		return nil, err
	}
	return &category.DeleteCategoryResponse{Message: "Category soft deleted successfully"}, nil
}

// GetAllCategories retrieves a list of non-deleted categories with optional pagination.
func (cDb *CategoryDb) GetAllCategories(ctx context.Context, req *category.GetAllCategoriesRequest) (*category.GetAllCategoriesResponse, error) {
	var args []interface{}
	query := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM 
			categories
		WHERE 
			deleted_at = 0
	`

	// Apply pagination
	if req.Limit <= 0 {
		req.Limit = 10 // Default limit
	}
	if req.Page <= 0 {
		req.Page = 1 // Default page
	}
	offset := (req.Page - 1) * req.Limit
	query += fmt.Sprintf(" OFFSET %d LIMIT %d", offset, req.Limit)

	rows, err := cDb.Db.Query(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Error listing categories")
		return nil, err
	}
	defer rows.Close()

	var categories []*category.Category
	for rows.Next() {
		var (
			createdAt time.Time
			updatedAt time.Time
		)
		dbCategory := &category.Category{}
		err := rows.Scan(
			&dbCategory.Id,
			&dbCategory.Name,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning category row")
			return nil, err
		}
		dbCategory.CreatedAt = createdAt.Format(time.RFC3339)
		dbCategory.UpdatedAt = updatedAt.Format(time.RFC3339)

		categories = append(categories, dbCategory)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over category rows")
		return nil, err
	}

	return &category.GetAllCategoriesResponse{Categories: categories}, nil
}
