package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Forum-service/Forum-Service/config"
	"github.com/Forum-service/Forum-Service/storage"
	"github.com/jackc/pgx/v5"
)

// Storage struct holds the database connection and interfaces for each table.
type Storage struct {
	db *pgx.Conn

	categoryRepo storage.CategoryRepo
	tagRepo      storage.TagRepo
	postRepo     storage.PostRepo
	commentRepo  storage.CommentRepo
	postTagRepo  storage.PostTagRepo
}

// NewStorage establishes a connection to the Postgres database and returns a Storage struct.
func NewStorage(cfg *config.Config) (*Storage, error) {
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	)

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		slog.Error("Unable to connect to database", err)
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		slog.Error("Failed to ping database", err)
		return nil, err
	}

	slog.Info("Connected to PostgreSQL database")

	return &Storage{
		db:           conn,
		categoryRepo: NewCategory(conn),
		tagRepo:      NewTag(conn),
		postRepo:     NewPost(conn),
		commentRepo:  NewComment(conn),
		postTagRepo:  NewPostTag(conn),
	}, nil
}

// Close closes the database connection.
func (s *Storage) Close() {
	if err := s.db.Close(context.Background()); err != nil {
		slog.Error("Error closing database connection", err)
	} else {
		slog.Info("Database connection closed successfully")
	}
}

// Category returns the CategoryRepo.
func (s *Storage) Category() storage.CategoryRepo {
	return s.categoryRepo
}

// Tag returns the TagRepo.
func (s *Storage) Tag() storage.TagRepo {
	return s.tagRepo
}

// Post returns the PostRepo.
func (s *Storage) Post() storage.PostRepo {
	return s.postRepo
}

// Comment returns the CommentRepo.
func (s *Storage) Comment() storage.CommentRepo {
	return s.commentRepo
}

// PostTag returns the PostTagRepo.
func (s *Storage) PostTag() storage.PostTagRepo {
	return s.postTagRepo
}
