package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Forum-service/Forum-Service/config"
	"github.com/Forum-service/Forum-Service/genproto/comment"
	"github.com/Forum-service/Forum-Service/storage/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func newTestComment(t *testing.T) *postgres.CommentDb {
	cfg := config.Load()

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		"localhost",
		5432,
		cfg.PostgresDatabase,
	)

	db, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return &postgres.CommentDb{Db: db}
}

func createTestComment(t *testing.T, cDb *postgres.CommentDb, postID string) *comment.Comment {

	testComment := &comment.CreateCommentRequest{
		PostId: postID,
		UserId: uuid.New().String(),
		Body:   "This is a test comment body.",
	}

	createdComment, err := cDb.Create(context.Background(), testComment)
	if err != nil {
		t.Fatalf("Error creating comment: %v", err)
	}

	return createdComment.Comment
}

func TestCreateComment(t *testing.T) {
	cDb := newTestComment(t)
	testPostID := "0257605c-b5a7-4480-8571-52c101da352b" // Your specified post ID
	createdComment := createTestComment(t, cDb, testPostID)

	assert.NotEmpty(t, createdComment.Id)
	assert.Equal(t, testPostID, createdComment.PostId)
	assert.NotEmpty(t, createdComment.UserId)
	assert.Equal(t, "This is a test comment body.", createdComment.Body)
	assert.NotEmpty(t, createdComment.CreatedAt)
	assert.NotEmpty(t, createdComment.UpdatedAt)
}

func TestGetCommentById(t *testing.T) {
	cDb := newTestComment(t)
	testPostID := "0257605c-b5a7-4480-8571-52c101da352b"
	createdComment := createTestComment(t, cDb, testPostID)

	getComment, err := cDb.GetById(context.Background(), &comment.GetCommentRequest{Id: createdComment.Id})
	if err != nil {
		t.Fatalf("Error getting comment by ID: %v", err)
	}

	assert.Equal(t, createdComment.Id, getComment.Comment.Id)
	assert.Equal(t, createdComment.PostId, getComment.Comment.PostId)
	assert.Equal(t, createdComment.UserId, getComment.Comment.UserId)
	assert.Equal(t, createdComment.Body, getComment.Comment.Body)
}

func TestUpdateComment(t *testing.T) {
	cDb := newTestComment(t)
	testPostID := "bd986d7f-a81e-4a68-8a01-b80c0b1178a4" // Updated post ID
	createdComment := createTestComment(t, cDb, testPostID)

	updatedCommentReq := &comment.UpdateCommentRequest{
		Id:   createdComment.Id,
		Body: "This is the updated comment body.",
	}

	updatedComment, err := cDb.Update(context.Background(), updatedCommentReq)
	if err != nil {
		t.Fatalf("Error updating comment: %v", err)
	}

	assert.Equal(t, updatedCommentReq.Id, updatedComment.Comment.Id)
	assert.Equal(t, updatedCommentReq.Body, updatedComment.Comment.Body)
}

func TestDeleteComment(t *testing.T) {
	cDb := newTestComment(t)
	testPostID := "0257605c-b5a7-4480-8571-52c101da352b"
	createdComment := createTestComment(t, cDb, testPostID)

	_, err := cDb.Delete(context.Background(), &comment.DeleteCommentRequest{Id: createdComment.Id})
	if err != nil {
		t.Fatalf("Error deleting comment: %v", err)
	}

	// Try to get the soft deleted comment
	_, err = cDb.GetById(context.Background(), &comment.GetCommentRequest{Id: createdComment.Id})

	assert.Error(t, err, "Should be able to retrieve soft deleted comment")
}

func TestGetAllComments(t *testing.T) {
	cDb := newTestComment(t)
	testPostID := "0257605c-b5a7-4480-8571-52c101da352b" // Post ID for testing

	// Create some test comments for the specified post
	testComments := []*comment.CreateCommentRequest{
		{
			PostId: testPostID,
			UserId: uuid.New().String(),
			Body:   "Comment A",
		},
		{
			PostId: testPostID,
			UserId: uuid.New().String(),
			Body:   "Comment B",
		},
	}

	for _, tc := range testComments {
		_, err := cDb.Create(context.Background(), tc)
		if err != nil {
			t.Fatalf("Error creating test comment: %v", err)
		}
	}

	t.Run("GetAllComments without filters", func(t *testing.T) {
		resp, err := cDb.GetAllComments(context.Background(), &comment.GetAllCommentsRequest{})
		if err != nil {
			t.Fatalf("Error listing all comments: %v", err)
		}
		assert.GreaterOrEqual(t, len(resp.Comments), len(testComments))
	})

	t.Run("Filter by Post ID", func(t *testing.T) {
		resp, err := cDb.GetAllComments(context.Background(), &comment.GetAllCommentsRequest{PostId: testPostID})
		if err != nil {
			t.Fatalf("Error listing comments by post ID: %v", err)
		}

		// Assert that all retrieved comments belong to the specified post
		for _, c := range resp.Comments {
			assert.Equal(t, testPostID, c.PostId)
		}
	})

}
