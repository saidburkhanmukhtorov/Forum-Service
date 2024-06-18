package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Forum-service/Forum-Service/genproto/post" // Your post proto package
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// ErrPostNotFound is returned when a post is not found.
var ErrPostNotFound = errors.New("post not found")

// PostDb provides database operations for posts.
type PostDb struct {
	Db *pgx.Conn
}

// NewPost creates a new instance of PostDb.
func NewPost(db *pgx.Conn) *PostDb {
	return &PostDb{Db: db}
}

// Create creates a new post in the database.
func (pDb *PostDb) Create(ctx context.Context, req *post.CreatePostRequest) (*post.CreatePostResponse, error) {
	postID := uuid.New().String()
	query := `
		INSERT INTO 
			posts (
				id,
				user_id,
				title,
				body,
				category_id 
			) 
		VALUES (
				$1, 
				$2, 
				$3,
				$4,
				$5
			)
		RETURNING 
			id,
			user_id,
			title,
			body,
			category_id,
			created_at,
			updated_at
	`
	var (
		dbPost    post.Post
		createdAt time.Time
		updatedAt time.Time
	)

	err := pDb.Db.QueryRow(ctx, query, postID, req.UserId, req.Title, req.Body, req.CategoryId).Scan(
		&dbPost.Id,
		&dbPost.UserId,
		&dbPost.Title,
		&dbPost.Body,
		&dbPost.CategoryId,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		log.Error().Err(err).Msg("Error creating post")
		return nil, err
	}

	dbPost.CreatedAt = createdAt.Format(time.RFC3339)
	dbPost.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &post.CreatePostResponse{Post: &dbPost}, nil
}

// GetById gets a post by its ID.
func (pDb *PostDb) GetById(ctx context.Context, req *post.GetPostRequest) (*post.GetPostResponse, error) {
	var (
		dbPost    post.Post
		createdAt time.Time
		updatedAt time.Time
	)

	query := `
		SELECT
			id,
			user_id,
			title,
			body,
			category_id,
			created_at,
			updated_at
		FROM 
			posts 
		WHERE 
			id = $1 
		AND 
			deleted_at = 0
	`
	err := pDb.Db.QueryRow(ctx, query, req.Id).Scan(
		&dbPost.Id,
		&dbPost.UserId,
		&dbPost.Title,
		&dbPost.Body,
		&dbPost.CategoryId,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Post not found")
			return nil, ErrPostNotFound
		}
		log.Error().Err(err).Msg("Error getting post by ID")
		return nil, err
	}

	dbPost.CreatedAt = createdAt.Format(time.RFC3339)
	dbPost.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &post.GetPostResponse{Post: &dbPost}, nil
}

// Update updates an existing post in the database.
func (pDb *PostDb) Update(ctx context.Context, req *post.UpdatePostRequest) (*post.UpdatePostResponse, error) {
	var args []interface{}
	count := 1
	query := `
		UPDATE 
			posts 
		SET `
	filter := ``

	if len(req.Title) > 0 {
		filter += fmt.Sprintf(" title = $%d, ", count)
		args = append(args, req.Title)
		count++
	}

	if len(req.Body) > 0 {
		filter += fmt.Sprintf(" body = $%d, ", count)
		args = append(args, req.Body)
		count++
	}

	if len(req.CategoryId) > 0 {
		filter += fmt.Sprintf(" category_id = $%d, ", count)
		args = append(args, req.CategoryId)
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

	_, err := pDb.Db.Exec(ctx, query, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Post not found")
			return nil, ErrPostNotFound
		}
		log.Error().Err(err).Msg("Error updating post")
		return nil, err
	}

	// After successful update, get the updated post from the database
	updatedPost, err := pDb.GetById(ctx, &post.GetPostRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &post.UpdatePostResponse{Post: updatedPost.Post}, nil
}

// Delete soft deletes a post by setting its deleted_at field to the current Unix timestamp.
func (pDb *PostDb) Delete(ctx context.Context, req *post.DeletePostRequest) (*post.DeletePostResponse, error) {
	query := `
		UPDATE 
			posts 
		SET 
			deleted_at = $1
		WHERE 
			id = $2
	`
	_, err := pDb.Db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Post not found")
			return nil, ErrPostNotFound
		}
		log.Error().Err(err).Msg("Error soft deleting post")
		return nil, err
	}
	return &post.DeletePostResponse{Message: "Post soft deleted successfully"}, nil
}

// GetAllPosts retrieves a list of non-deleted posts with optional filtering and pagination.
func (pDb *PostDb) GetAllPosts(ctx context.Context, req *post.GetAllPostsRequest) (*post.GetAllPostsResponse, error) {
	var (
		args  []interface{}
		count int = 1
	)
	query := `
		SELECT
			id,
			user_id,
			title,
			body,
			category_id,
			created_at,
			updated_at
		FROM 
			posts
		WHERE deleted_at = 0
	`
	filter := ""

	if req.UserId != "" {
		filter += fmt.Sprintf(" AND user_id = $%d", count)
		args = append(args, req.UserId)
		count++
	}

	if req.Title != "" {
		filter += fmt.Sprintf(" AND title ILIKE $%d", count)
		args = append(args, "%"+req.Title+"%")
		count++
	}

	if req.Body != "" {
		filter += fmt.Sprintf(" AND body ILIKE $%d", count)
		args = append(args, "%"+req.Body+"%")
		count++
	}

	if req.CategoryId != "" {
		filter += fmt.Sprintf(" AND category_id = $%d", count)
		args = append(args, req.CategoryId)
		count++
	}

	query += filter

	// Apply pagination
	if req.Limit <= 0 {
		req.Limit = 10 // Default limit
	}
	if req.Page <= 0 {
		req.Page = 1 // Default page
	}
	offset := (req.Page - 1) * req.Limit
	query += fmt.Sprintf(" OFFSET %d LIMIT %d", offset, req.Limit)

	rows, err := pDb.Db.Query(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Error listing posts")
		return nil, err
	}
	defer rows.Close()

	var posts []*post.Post
	for rows.Next() {
		var (
			createdAt time.Time
			updatedAt time.Time
		)
		dbPost := &post.Post{}
		err := rows.Scan(
			&dbPost.Id,
			&dbPost.UserId,
			&dbPost.Title,
			&dbPost.Body,
			&dbPost.CategoryId,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning post row")
			return nil, err
		}
		dbPost.CreatedAt = createdAt.Format(time.RFC3339)
		dbPost.UpdatedAt = updatedAt.Format(time.RFC3339)

		posts = append(posts, dbPost)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over post rows")
		return nil, err
	}

	return &post.GetAllPostsResponse{Posts: posts}, nil
}
