package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Forum-service/Forum-Service/config"
	"github.com/Forum-service/Forum-Service/genproto/post"
	"github.com/Forum-service/Forum-Service/storage/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func newTestPost(t *testing.T) *postgres.PostDb {
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
	return &postgres.PostDb{Db: db}
}

func createTestPost(t *testing.T, pDb *postgres.PostDb) *post.Post {
	testCategoryID := "a5a171a3-e5dd-491e-94a1-94eefc9fb320"
	testUserID := uuid.New().String()

	testPost := &post.CreatePostRequest{
		UserId:     testUserID,
		Title:      "Test Post",
		Body:       "This is a test post body.",
		CategoryId: testCategoryID,
	}

	createdPost, err := pDb.Create(context.Background(), testPost)
	if err != nil {
		t.Fatalf("Error creating post: %v", err)
	}

	return createdPost.Post
}

func TestCreatePost(t *testing.T) {
	pDb := newTestPost(t)
	createdPost := createTestPost(t, pDb)

	assert.NotEmpty(t, createdPost.Id)
	assert.NotEmpty(t, createdPost.UserId)
	assert.Equal(t, "Test Post", createdPost.Title)
	assert.Equal(t, "This is a test post body.", createdPost.Body)
	assert.Equal(t, "a5a171a3-e5dd-491e-94a1-94eefc9fb320", createdPost.CategoryId)
	assert.NotEmpty(t, createdPost.CreatedAt)
	assert.NotEmpty(t, createdPost.UpdatedAt)
}

func TestGetPostById(t *testing.T) {
	pDb := newTestPost(t)
	createdPost := createTestPost(t, pDb)

	getPost, err := pDb.GetById(context.Background(), &post.GetPostRequest{Id: createdPost.Id})
	if err != nil {
		t.Fatalf("Error getting post by ID: %v", err)
	}

	assert.Equal(t, createdPost.Id, getPost.Post.Id)
	assert.Equal(t, createdPost.UserId, getPost.Post.UserId)
	// ... (other assertions for Title, Body, CategoryId, etc.)
}

func TestUpdatePost(t *testing.T) {
	pDb := newTestPost(t)
	createdPost := createTestPost(t, pDb)

	updatedPostReq := &post.UpdatePostRequest{
		Id:         createdPost.Id,
		Title:      "Updated Post Title",
		Body:       "This is the updated post body.",
		CategoryId: "e8650976-1985-48fa-bbbf-9685180f8ff8",
	}

	updatedPost, err := pDb.Update(context.Background(), updatedPostReq)
	if err != nil {
		t.Fatalf("Error updating post: %v", err)
	}

	assert.Equal(t, updatedPostReq.Id, updatedPost.Post.Id)
	assert.Equal(t, updatedPostReq.Title, updatedPost.Post.Title)
	assert.Equal(t, updatedPostReq.Body, updatedPost.Post.Body)
	assert.Equal(t, updatedPostReq.CategoryId, updatedPost.Post.CategoryId)
}

func TestDeletePost(t *testing.T) {
	pDb := newTestPost(t)
	createdPost := createTestPost(t, pDb)

	_, err := pDb.Delete(context.Background(), &post.DeletePostRequest{Id: createdPost.Id})
	if err != nil {
		t.Fatalf("Error deleting post: %v", err)
	}

	// Try to get the soft deleted post
	_, err = pDb.GetById(context.Background(), &post.GetPostRequest{Id: createdPost.Id})

	assert.Error(t, err, "Should be able to retrieve soft deleted post")
}

func TestGetAllPosts(t *testing.T) {
	pDb := newTestPost(t)
	testCategoryID := "a5a171a3-e5dd-491e-94a1-94eefc9fb320"

	// Create some test posts with the same category ID
	testPosts := []*post.CreatePostRequest{
		{
			UserId:     uuid.New().String(),
			Title:      "Post A",
			Body:       "Body of Post A",
			CategoryId: testCategoryID,
		},
		{
			UserId:     uuid.New().String(),
			Title:      "Post B",
			Body:       "Body of Post B",
			CategoryId: testCategoryID,
		},
		// ... add more test posts if you want
	}
	for _, tp := range testPosts {
		_, err := pDb.Create(context.Background(), tp)
		if err != nil {
			t.Fatalf("Error creating test post: %v", err)
		}
	}

	t.Run("GetAllPosts without filters", func(t *testing.T) {
		resp, err := pDb.GetAllPosts(context.Background(), &post.GetAllPostsRequest{})
		if err != nil {
			t.Fatalf("Error listing all posts: %v", err)
		}

		// Assert that you get at least the posts you created
		assert.GreaterOrEqual(t, len(resp.Posts), len(testPosts))
	})

	t.Run("Filter by Category ID", func(t *testing.T) {
		resp, err := pDb.GetAllPosts(context.Background(), &post.GetAllPostsRequest{CategoryId: testCategoryID})
		if err != nil {
			t.Fatalf("Error listing posts by category ID: %v", err)
		}

		// Assert that all retrieved posts have the correct category ID
		for _, p := range resp.Posts {
			assert.Equal(t, testCategoryID, p.CategoryId)
		}
	})

	// ... (add more test cases for other filters: UserId, Title, Pagination)
}
