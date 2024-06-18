package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Forum-service/Forum-Service/config"
	"github.com/Forum-service/Forum-Service/genproto/tag"
	"github.com/Forum-service/Forum-Service/storage/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func newTestTag(t *testing.T) *postgres.TagDb {
	cfg := config.Load()

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		"localhost", // Replace if your DB is not on localhost
		5432,        // Replace if your DB port is different
		cfg.PostgresDatabase,
	)

	db, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return &postgres.TagDb{Db: db}
}

func createTestTag(t *testing.T, tDb *postgres.TagDb) *tag.Tag {
	testTag := &tag.CreateTagRequest{
		Name: "Test Tag",
	}

	createdTag, err := tDb.Create(context.Background(), testTag)
	if err != nil {
		t.Fatalf("Error creating tag: %v", err)
	}

	return createdTag.Tag
}

func TestCreateTag(t *testing.T) {
	tDb := newTestTag(t)
	testTag := &tag.CreateTagRequest{
		Name: "Test Tag",
	}

	createdTag, err := tDb.Create(context.Background(), testTag)
	if err != nil {
		t.Fatalf("Error creating tag: %v", err)
	}

	assert.NotEmpty(t, createdTag.Tag.Id)
	assert.Equal(t, testTag.Name, createdTag.Tag.Name)
	assert.NotEmpty(t, createdTag.Tag.CreatedAt)
	assert.NotEmpty(t, createdTag.Tag.UpdatedAt)
}

func TestGetTagById(t *testing.T) {
	tDb := newTestTag(t)
	createdTag := createTestTag(t, tDb)

	getTag, err := tDb.GetById(context.Background(), &tag.GetTagRequest{Id: createdTag.Id})
	if err != nil {
		t.Fatalf("Error getting tag by ID: %v", err)
	}

	assert.Equal(t, createdTag.Id, getTag.Tag.Id)
	assert.Equal(t, createdTag.Name, getTag.Tag.Name)
}

func TestUpdateTag(t *testing.T) {
	tDb := newTestTag(t)
	createdTag := createTestTag(t, tDb)

	updatedTagReq := &tag.UpdateTagRequest{
		Id:   createdTag.Id,
		Name: "Updated Tag Name",
	}

	updatedTag, err := tDb.Update(context.Background(), updatedTagReq)
	if err != nil {
		t.Fatalf("Error updating tag: %v", err)
	}

	assert.Equal(t, updatedTagReq.Id, updatedTag.Tag.Id)
	assert.Equal(t, updatedTagReq.Name, updatedTag.Tag.Name)
}

func TestDeleteTag(t *testing.T) {
	tDb := newTestTag(t)
	createdTag := createTestTag(t, tDb)

	_, err := tDb.Delete(context.Background(), &tag.DeleteTagRequest{Id: createdTag.Id})
	if err != nil {
		t.Fatalf("Error deleting tag: %v", err)
	}

	// Try to get the soft deleted tag
	_, err = tDb.GetById(context.Background(), &tag.GetTagRequest{Id: createdTag.Id})

	assert.Error(t, err, "Should be able to retrieve soft deleted post")
}

func TestGetAllTags(t *testing.T) {
	tDb := newTestTag(t)

	testTags := []*tag.CreateTagRequest{
		{Name: "Tag A"},
		{Name: "Tag B"},
		{Name: "Tag C"},
	}

	for _, tt := range testTags {
		_, err := tDb.Create(context.Background(), tt)
		if err != nil {
			t.Fatalf("Error creating test tag: %v", err)
		}
	}

	t.Run("ListTags without filters", func(t *testing.T) {
		resp, err := tDb.GetAllTags(context.Background(), &tag.GetAllTagsRequest{})
		if err != nil {
			t.Fatalf("Error listing all tags: %v", err)
		}
		assert.GreaterOrEqual(t, len(resp.Tags), len(testTags))
	})
}
