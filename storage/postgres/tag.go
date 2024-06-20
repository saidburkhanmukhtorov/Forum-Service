package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Forum-service/Forum-Service/genproto/tag" // Your tag proto package
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// ErrTagNotFound is returned when a tag is not found.
var ErrTagNotFound = errors.New("tag not found")

// TagDb provides database operations for tags.
type TagDb struct {
	Db *pgx.Conn
}

// NewTag creates a new instance of TagDb.
func NewTag(db *pgx.Conn) *TagDb {
	return &TagDb{Db: db}
}

// Create creates a new tag in the database.
func (tDb *TagDb) Create(ctx context.Context, req *tag.CreateTagRequest) (*tag.CreateTagResponse, error) {
	tagID := uuid.New().String()
	query := `
		INSERT INTO 
			tags (
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
		dbTag     tag.Tag
		createdAt time.Time
		updatedAt time.Time
	)

	err := tDb.Db.QueryRow(ctx, query, tagID, req.Name).Scan(
		&dbTag.Id,
		&dbTag.Name,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		log.Error().Err(err).Msg("Error creating tag")
		return nil, err
	}

	dbTag.CreatedAt = createdAt.Format(time.RFC3339)
	dbTag.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &tag.CreateTagResponse{Tag: &dbTag}, nil
}

// GetById gets a tag by its ID.
func (tDb *TagDb) GetById(ctx context.Context, req *tag.GetTagRequest) (*tag.GetTagResponse, error) {
	var (
		dbTag     tag.Tag
		createdAt time.Time
		updatedAt time.Time
	)

	query := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM 
			tags 
		WHERE 
			id = $1
		AND 
			deleted_at = 0
	`
	err := tDb.Db.QueryRow(ctx, query, req.Id).Scan(
		&dbTag.Id,
		&dbTag.Name,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Tag not found")
			return nil, ErrTagNotFound
		}
		log.Error().Err(err).Msg("Error getting tag by ID")
		return nil, err
	}

	dbTag.CreatedAt = createdAt.Format(time.RFC3339)
	dbTag.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &tag.GetTagResponse{Tag: &dbTag}, nil
}

// Update updates an existing tag in the database.
func (tDb *TagDb) Update(ctx context.Context, req *tag.UpdateTagRequest) (*tag.UpdateTagResponse, error) {
	var args []interface{}
	count := 1
	query := `
		UPDATE 
			tags 
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

	_, err := tDb.Db.Exec(ctx, query, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Tag not found")
			return nil, ErrTagNotFound
		}
		log.Error().Err(err).Msg("Error updating tag")
		return nil, err
	}

	// After successful update, get the updated tag from the database
	updatedTag, err := tDb.GetById(ctx, &tag.GetTagRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &tag.UpdateTagResponse{Tag: updatedTag.Tag}, nil
}

// Delete soft deletes a tag by setting its deleted_at field to the current Unix timestamp.
func (tDb *TagDb) Delete(ctx context.Context, req *tag.DeleteTagRequest) (*tag.DeleteTagResponse, error) {
	query := `
		UPDATE 
			tags 
		SET 
			deleted_at = $1 
		WHERE 
			id = $2
	`
	_, err := tDb.Db.Exec(ctx, query, time.Now().Unix(), req.Id)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Error().Err(err).Msg("Tag not found")
			return nil, ErrTagNotFound
		}
		log.Error().Err(err).Msg("Error soft deleting tag")
		return nil, err
	}
	return &tag.DeleteTagResponse{Message: "Tag soft deleted successfully"}, nil
}

// GetAllTags retrieves a list of non-deleted tags with optional pagination.
func (tDb *TagDb) GetAllTags(ctx context.Context, req *tag.GetAllTagsRequest) (*tag.GetAllTagsResponse, error) {
	var args []interface{}
	query := `
		SELECT
			id,
			name,
			created_at,
			updated_at
		FROM 
			tags
		WHERE 
			deleted_at = 0
	`
	if req.Name != "" {
		query += " AND name ILIKE $1 "
		args = append(args, "%"+req.Name+"%")
	}

	// Apply pagination
	if req.Limit <= 0 {
		req.Limit = 10 // Default limit
	}
	if req.Page <= 0 {
		req.Page = 1 // Default page
	}
	offset := (req.Page - 1) * req.Limit
	query += fmt.Sprintf(" OFFSET %d LIMIT %d", offset, req.Limit)

	rows, err := tDb.Db.Query(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Error listing tags")
		return nil, err
	}
	defer rows.Close()

	var tags []*tag.Tag
	for rows.Next() {
		var (
			createdAt time.Time
			updatedAt time.Time
		)
		dbTag := &tag.Tag{}
		err := rows.Scan(
			&dbTag.Id,
			&dbTag.Name,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning tag row")
			return nil, err
		}
		dbTag.CreatedAt = createdAt.Format(time.RFC3339)
		dbTag.UpdatedAt = updatedAt.Format(time.RFC3339)

		tags = append(tags, dbTag)
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over tag rows")
		return nil, err
	}

	return &tag.GetAllTagsResponse{Tags: tags}, nil
}

// GetFamousTags retrieves a list of famous tags with optional pagination and sorting.
func (tDb *TagDb) GetFamousTags(ctx context.Context, req *tag.GetFamousTagsReq) (*tag.GetFamousTagsRes, error) {
	var args []interface{}
	query := `
		SELECT name, count(name)
		FROM tags
		WHERE deleted_at = 0
	`
	if req.Name != "" {
		query += " AND name ILIKE $1 "
		args = append(args, "%"+req.Name+"%")
	}

	// Apply sorting
	if req.Desc {
		query += " GROUP BY name ORDER BY count(name) DESC"
	} else {
		query += " GROUP BY name ORDER BY count(name) ASC"
	}
	// Apply pagination
	if req.Limit <= 0 {
		req.Limit = 10 // Default limit
	}
	if req.Page <= 0 {
		req.Page = 1 // Default page
	}
	offset := (req.Page - 1) * req.Limit
	query += fmt.Sprintf(" OFFSET %d LIMIT %d", offset, req.Limit)

	rows, err := tDb.Db.Query(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("Error getting famous tags")
		return nil, err
	}
	defer rows.Close()

	var famousTags []*tag.FamousTag
	for rows.Next() {
		var (
			name  string
			count int32
		)
		err := rows.Scan(
			&name,
			&count,
		)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning famous tag row")
			return nil, err
		}

		famousTags = append(famousTags, &tag.FamousTag{Name: name, Count: count})
	}

	if err = rows.Err(); err != nil {
		log.Error().Err(err).Msg("Error iterating over famous tag rows")
		return nil, err
	}

	return &tag.GetFamousTagsRes{Tags: famousTags}, nil
}
