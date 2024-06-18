package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Forum-service/Forum-Service/genproto/comment" // Your comment proto package
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// ErrCommentNotFound is returned when a comment is not found.
var ErrCommentNotFound = errors.New("comment not found")

// CommentDb provides database operations for comments.
type CommentDb struct {
	Db *pgx.Conn
}

// NewComment creates a new instance of CommentDb.
func NewComment(db *pgx.Conn) *CommentDb {
	return &CommentDb{Db: db}
}

// Create creates a new comment in the database.
func (cDb *CommentDb) Create(ctx context.Context, req *comment.CreateCommentRequest) (*comment.CreateCommentResponse, error) {
	commentID := uuid.New().String()
	query := `
		INSERT INTO 
			comments (
				id,
				post_id,
				user_id,
				body
			) 
		VALUES (
				$1, 
				$2, 
				$3,
				$4
			)
		RETURNING 
			id,
			post_id,
			user_id,
			body,
			created_at,
			updated_at
	`
	var (
		dbComment comment.Comment
		createdAt time.Time
		updatedAt time.Time
	)

	err := cDb.Db.QueryRow(ctx, query, commentID, req.PostId, req.UserId, req.Body).Scan(
		&dbComment.Id,
		&dbComment.PostId,
		&dbComment.UserId,
		&dbComment.Body,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		log.Error().Err(err).Msg("Error creating comment")
		return nil, err
	}

	dbComment.CreatedAt = createdAt.Format(time.RFC3339)
	dbComment.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &comment.CreateCommentResponse{Comment: &dbComment}, nil
}

// GetById gets a comment by its ID.
func (cDb *CommentDb) GetById(ctx context.Context, req *comment.GetCommentRequest) (*comment.GetCommentResponse, error) {
	var (
		dbComment comment.Comment
		createdAt time.Time
		updatedAt time.Time
	)

	query := `
		SELECT
			id,
			post_id,
			user_id,
			body,
			created_at,
			updated_at
		FROM 
			comments 
		WHERE 
			id = $1 
		AND 
			deleted_at = 0
	`
	err := cDb.Db.QueryRow(ctx, query, req.Id).Scan(
		&dbComment.Id,
		&dbComment.PostId,
		&dbComment.UserId,
		&dbComment.Body,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Comment not found")
			return nil, ErrCommentNotFound
		}
		log.Error().Err(err).Msg("Error getting comment by ID")
		return nil, err
	}

	dbComment.CreatedAt = createdAt.Format(time.RFC3339)
	dbComment.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &comment.GetCommentResponse{Comment: &dbComment}, nil
}

// Update updates an existing comment in the database.
func (cDb *CommentDb) Update(ctx context.Context, req *comment.UpdateCommentRequest) (*comment.UpdateCommentResponse, error) {
	var args []interface{}
	count := 1
	query := `
		UPDATE 
			comments 
		SET `
	filter := ``

	if len(req.Body) > 0 {
		filter += fmt.Sprintf(" body = $%d, ", count)
		args = append(args, req.Body)
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
			log.Error().Err(err).Msg("Comment not found")
			return nil, ErrCommentNotFound
		}
		log.Error().Err(err).Msg("Error updating comment")
		return nil, err
	}

	// After successful update, get the updated comment from the database
	updatedComment, err := cDb.GetById(ctx, &comment.GetCommentRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &comment.UpdateCommentResponse{Comment: updatedComment.Comment}, nil
}

// Delete soft deletes a comment by setting its deleted_at field to the current Unix timestamp.
func (cDb *CommentDb) Delete(ctx context.Context, req *comment.DeleteCommentRequest) (*comment.DeleteCommentResponse, error) {
	query := `
		UPDATE 
			comments 
		SET 
			deleted_at = $1
		WHERE 
			id = $2
	`
	_, err := cDb.Db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Comment not found")
			return nil, ErrCommentNotFound
		}
		log.Error().Err(err).Msg("Error soft deleting comment")
		return nil, err
	}
	return &comment.DeleteCommentResponse{Message: "Comment soft deleted successfully"}, nil
}

// GetAllComments retrieves a list of non-deleted comments with optional filtering and pagination.
func (cDb *CommentDb) GetAllComments(ctx context.Context, req *comment.GetAllCommentsRequest) (*comment.GetAllCommentsResponse, error) {
	var (
		args  []interface{}
		count int = 1
	)
	query := `
		SELECT
			id,
			post_id,
			user_id,
			body,
			created_at,
			updated_at
		FROM 
			comments
		WHERE deleted_at = 0
	`
	filter := ""

	if req.PostId != "" {
		filter += fmt.Sprintf(" AND post_id = $%d", count)
		args = append(args, req.PostId)
		count++
	}

	if req.UserId != "" {
		filter += fmt.Sprintf(" AND user_id = $%d", count)
		args = append(args, req.UserId)
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

	rows, err := cDb.Db.Query(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Error listing comments")
		return nil, err
	}
	defer rows.Close()

	var comments []*comment.Comment
	for rows.Next() {
		var (
			createdAt time.Time
			updatedAt time.Time
		)
		dbComment := &comment.Comment{}
		err := rows.Scan(
			&dbComment.Id,
			&dbComment.PostId,
			&dbComment.UserId,
			&dbComment.Body,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning comment row")
			return nil, err
		}
		dbComment.CreatedAt = createdAt.Format(time.RFC3339)
		dbComment.UpdatedAt = updatedAt.Format(time.RFC3339)

		comments = append(comments, dbComment)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over comment rows")
		return nil, err
	}

	return &comment.GetAllCommentsResponse{Comments: comments}, nil
}
