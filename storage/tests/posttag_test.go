package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Forum-service/Forum-Service/config"
	"github.com/Forum-service/Forum-Service/genproto/post"
	"github.com/Forum-service/Forum-Service/genproto/posttag"
	"github.com/Forum-service/Forum-Service/storage/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

// ... other test setup functions ...

func newTestPostTag(t *testing.T) *postgres.PostTagDb {
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
	return &postgres.PostTagDb{Db: db}
}

func TestCreatePostTag(t *testing.T) {
	ptDb := newTestPostTag(t)
	testPostID := "0257605c-b5a7-4480-8571-52c101da352b"
	testTagID := "2e62d3bd-8fab-43e2-a809-62f48e9824c1"

	createReq := &posttag.CreatePostTagRequest{
		PostId: testPostID,
		TagId:  testTagID,
	}

	createdPostTag, err := ptDb.Create(context.Background(), createReq)
	if err != nil {
		t.Fatalf("Error creating post_tag association: %v", err)
	}

	assert.Equal(t, createReq.PostId, createdPostTag.PostTag.PostId)
	assert.Equal(t, createReq.TagId, createdPostTag.PostTag.TagId)
	assert.NotEmpty(t, createdPostTag.PostTag.CreatedAt)
}

func TestDeletePostTag(t *testing.T) {
	ptDb := newTestPostTag(t)
	testPostID := "741e6fed-1ce6-475d-befe-8be1c34513fd"
	testTagID := "d5658a02-07f8-45bc-893f-da4e22879a16"

	// Create the association first
	_, err := ptDb.Create(context.Background(), &posttag.CreatePostTagRequest{
		PostId: testPostID,
		TagId:  testTagID,
	})
	if err != nil {
		t.Fatalf("Error creating post_tag association for delete test: %v", err)
	}

	// Then delete it
	deleteReq := &posttag.DeletePostTagRequest{
		PostId: testPostID,
		TagId:  testTagID,
	}

	_, err = ptDb.Delete(context.Background(), deleteReq)
	assert.NoError(t, err, "Error deleting post_tag association")

	res, err := ptDb.GetAllPostTags(context.Background(), &posttag.GetAllPostTagsRequest{
		PostId: testPostID,
		TagId:  testTagID,
	})
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	assert.Equal(t, 0, len(res.PostTags))

}

func TestGetAllPostTags(t *testing.T) {
	ptDb := newTestPostTag(t)
	testPostID := "0257605c-b5a7-4480-8571-52c101da352b"
	testTagID1 := "2e62d3bd-8fab-43e2-a809-62f48e9824c1"
	testTagID2 := "21974b88-d64a-4dc0-b1e9-9bee9c4a9c3b"

	// Create some test post_tag associations
	_, err := ptDb.Create(context.Background(), &posttag.CreatePostTagRequest{
		PostId: testPostID,
		TagId:  testTagID1,
	})
	if err != nil {
		t.Fatalf("Error creating test post_tag association: %v", err)
	}
	_, err = ptDb.Create(context.Background(), &posttag.CreatePostTagRequest{
		PostId: testPostID,
		TagId:  testTagID2,
	})
	if err != nil {
		t.Fatalf("Error creating test post_tag association: %v", err)
	}

	t.Run("GetAllPostTags without filters", func(t *testing.T) {
		resp, err := ptDb.GetAllPostTags(context.Background(), &posttag.GetAllPostTagsRequest{})
		if err != nil {
			t.Fatalf("Error getting all post_tags: %v", err)
		}
		assert.GreaterOrEqual(t, len(resp.PostTags), 2)
	})

	t.Run("Filter by Post ID", func(t *testing.T) {
		resp, err := ptDb.GetAllPostTags(context.Background(), &posttag.GetAllPostTagsRequest{PostId: testPostID})
		if err != nil {
			t.Fatalf("Error getting post_tags by Post ID: %v", err)
		}
		assert.GreaterOrEqual(t, len(resp.PostTags), 2) // Should get 2 associations for the test post ID
		for _, pt := range resp.PostTags {
			assert.Equal(t, testPostID, pt.PostId)
		}
	})

	t.Run("Filter by Tag ID", func(t *testing.T) {
		resp, err := ptDb.GetAllPostTags(context.Background(), &posttag.GetAllPostTagsRequest{TagId: testTagID1})
		if err != nil {
			t.Fatalf("Error getting post_tags by Tag ID: %v", err)
		}
		assert.GreaterOrEqual(t, len(resp.PostTags), 1)
		assert.Equal(t, testTagID1, resp.PostTags[0].TagId)
	})

}

func TestGetPostsByTag(t *testing.T) {
	ptDb := newTestPostTag(t)
	pDB := newTestPost(t)
	testTagID := "2e62d3bd-8fab-43e2-a809-62f48e9824c1"

	// Create some test posts and associate them with the tag
	testPosts := []*post.CreatePostRequest{
		{
			UserId:     uuid.New().String(),
			Title:      "Post A",
			Body:       "Body of Post A",
			CategoryId: "e8650976-1985-48fa-bbbf-9685180f8ff8", // Replace with valid category ID
		},
		{
			UserId:     uuid.New().String(),
			Title:      "Post B",
			Body:       "Body of Post B",
			CategoryId: "a5a171a3-e5dd-491e-94a1-94eefc9fb320", // Replace with valid category ID
		},
	}

	for _, postReq := range testPosts {
		createdPost, err := pDB.Create(context.Background(), postReq) // Use pDb (PostDb) to create posts
		if err != nil {
			t.Fatalf("Error creating test post: %v", err)
		}

		_, err = ptDb.Create(context.Background(), &posttag.CreatePostTagRequest{
			PostId: createdPost.Post.Id,
			TagId:  testTagID,
		})
		if err != nil {
			t.Fatalf("Error creating test post_tag association: %v", err)
		}
	}

	getReq := &posttag.GetPostsByTagRequest{
		TagId: testTagID,
	}

	resp, err := ptDb.GetPostsByTag(context.Background(), getReq)
	assert.NoError(t, err, "Error getting posts by tag ID")

	assert.LessOrEqual(t, len(testPosts), len(resp.Posts), "Should return the correct number of posts")
}
