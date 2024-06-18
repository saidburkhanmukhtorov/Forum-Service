package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Forum-service/Forum-Service/genproto/post"
	"github.com/Forum-service/Forum-Service/genproto/posttag"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// ErrPostTagNotFound is returned when a post_tag record is not found.
var ErrPostTagNotFound = errors.New("post_tag record not found")

// PostTagDb provides database operations for post_tags.
type PostTagDb struct {
	Db *pgx.Conn
}

// NewPostTag creates a new instance of PostTagDb.
func NewPostTag(db *pgx.Conn) *PostTagDb {
	return &PostTagDb{Db: db}
}

// Create creates a new post_tag association in the database.
func (ptDb *PostTagDb) Create(ctx context.Context, req *posttag.CreatePostTagRequest) (*posttag.CreatePostTagResponse, error) {
	query := `
		INSERT INTO 
			post_tags (
				post_id,
				tag_id 
			) 
		VALUES (
				$1, 
				$2 
			)
		RETURNING 
			post_id,
			tag_id,
			created_at
	`
	var (
		dbPostTag posttag.PostTag
		createdAt time.Time
	)

	err := ptDb.Db.QueryRow(ctx, query, req.PostId, req.TagId).Scan(
		&dbPostTag.PostId,
		&dbPostTag.TagId,
		&createdAt,
	)
	if err != nil {
		log.Error().Err(err).Msg("Error creating post_tag association")
		return nil, err
	}

	dbPostTag.CreatedAt = createdAt.Format(time.RFC3339)

	return &posttag.CreatePostTagResponse{PostTag: &dbPostTag}, nil
}

// Delete removes a post_tag association from the database.
func (ptDb *PostTagDb) Delete(ctx context.Context, req *posttag.DeletePostTagRequest) (*posttag.DeletePostTagResponse, error) {
	query := `
		DELETE FROM 
			post_tags 
		WHERE 
			post_id = $1 
		AND 
			tag_id = $2
	`
	_, err := ptDb.Db.Exec(ctx, query, req.PostId, req.TagId)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Post_tag association not found")
			return nil, ErrPostTagNotFound
		}
		log.Error().Err(err).Msg("Error deleting post_tag association")
		return nil, err
	}
	return &posttag.DeletePostTagResponse{Message: "Post_tag association deleted successfully"}, nil
}

// GetAllPostTags retrieves all post_tag associations with optional filtering and pagination.
func (ptDb *PostTagDb) GetAllPostTags(ctx context.Context, req *posttag.GetAllPostTagsRequest) (*posttag.GetAllPostTagsResponse, error) {
	var (
		args  []interface{}
		count int = 1
	)
	query := `
		SELECT
			post_id,
			tag_id,
			created_at
		FROM 
			post_tags
	`
	filter := ""

	if req.PostId != "" {
		filter += fmt.Sprintf(" WHERE post_id = $%d", count)
		args = append(args, req.PostId)
		count++
	}

	if req.TagId != "" {
		if len(args) > 0 {
			filter += " AND "
		} else {
			filter += " WHERE "
		}
		filter += fmt.Sprintf("tag_id = $%d", count)
		args = append(args, req.TagId)
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

	rows, err := ptDb.Db.Query(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Error listing post_tags")
		return nil, err
	}
	defer rows.Close()

	var postTags []*posttag.PostTag
	for rows.Next() {
		var (
			createdAt time.Time
		)
		dbPostTag := &posttag.PostTag{}
		err := rows.Scan(
			&dbPostTag.PostId,
			&dbPostTag.TagId,
			&createdAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning post_tag row")
			return nil, err
		}
		dbPostTag.CreatedAt = createdAt.Format(time.RFC3339)

		postTags = append(postTags, dbPostTag)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over post_tag rows")
		return nil, err
	}

	return &posttag.GetAllPostTagsResponse{PostTags: postTags}, nil
}

// GetPostsByTag retrieves posts associated with a specific tag ID.
func (ptDb *PostTagDb) GetPostsByTag(ctx context.Context, req *posttag.GetPostsByTagRequest) (*posttag.GetPostsByTagResponse, error) {
	var args []interface{}
	query := `
        SELECT 
            p.id,
            p.user_id,
            p.title,
            p.body,
            p.category_id,
            p.created_at,
            p.updated_at
        FROM 
            posts p
        INNER JOIN post_tags pt ON p.id = pt.post_id
        WHERE 
            pt.tag_id = $1
        AND p.deleted_at = 0
    `
	args = append(args, req.TagId)

	// Apply pagination
	if req.Limit <= 0 {
		req.Limit = 10 // Default limit
	}
	if req.Page <= 0 {
		req.Page = 1 // Default page
	}
	offset := (req.Page - 1) * req.Limit
	query += fmt.Sprintf(" OFFSET %d LIMIT %d", offset, req.Limit)

	rows, err := ptDb.Db.Query(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Error getting posts by tag ID")
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
		log.Error().Err(err).Msg("Error iterating over posts")
		return nil, err
	}

	return &posttag.GetPostsByTagResponse{Posts: posts}, nil
}
