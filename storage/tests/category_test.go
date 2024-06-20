package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Forum-service/Forum-Service/config"
	"github.com/Forum-service/Forum-Service/genproto/category"
	"github.com/Forum-service/Forum-Service/storage/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func newTestCategory(t *testing.T) *postgres.CategoryDb {
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
	return &postgres.CategoryDb{Db: db}
}

func createTestCategory(t *testing.T, cDb *postgres.CategoryDb) *category.Category {
	testCategory := &category.CreateCategoryRequest{
		Name: "Test Category",
	}

	createdCategory, err := cDb.Create(context.Background(), testCategory)
	if err != nil {
		t.Fatalf("Error creating category: %v", err)
	}

	return createdCategory.Category
}

func TestCreateCategory(t *testing.T) {
	cDb := newTestCategory(t)
	testCategory := &category.CreateCategoryRequest{
		Name: "Test Category",
	}

	createdCategory, err := cDb.Create(context.Background(), testCategory)
	if err != nil {
		t.Fatalf("Error creating category: %v", err)
	}

	assert.NotEmpty(t, createdCategory.Category.Id)
	assert.Equal(t, testCategory.Name, createdCategory.Category.Name)
	assert.NotEmpty(t, createdCategory.Category.CreatedAt)
	assert.NotEmpty(t, createdCategory.Category.UpdatedAt)
}

func TestGetCategoryById(t *testing.T) {
	cDb := newTestCategory(t)
	createdCategory := createTestCategory(t, cDb)

	getCategory, err := cDb.GetById(context.Background(), &category.GetCategoryRequest{Id: createdCategory.Id})
	if err != nil {
		t.Fatalf("Error getting category by ID: %v", err)
	}

	assert.Equal(t, createdCategory.Id, getCategory.Category.Id)
	assert.Equal(t, createdCategory.Name, getCategory.Category.Name)
}

func TestUpdateCategory(t *testing.T) {
	cDb := newTestCategory(t)
	createdCategory := createTestCategory(t, cDb)

	updatedCategoryReq := &category.UpdateCategoryRequest{
		Id:   createdCategory.Id,
		Name: "Updated Category Name",
	}
	updatedCategory, err := cDb.Update(context.Background(), updatedCategoryReq)
	if err != nil {
		t.Fatalf("Error updating category: %v", err)
	}

	assert.Equal(t, updatedCategoryReq.Id, updatedCategory.Category.Id)
	assert.Equal(t, updatedCategoryReq.Name, updatedCategory.Category.Name)

}

func TestDeleteCategory(t *testing.T) {
	cDb := newTestCategory(t)
	createdCategory := createTestCategory(t, cDb)

	_, err := cDb.Delete(context.Background(), &category.DeleteCategoryRequest{Id: createdCategory.Id})
	if err != nil {
		t.Fatalf("Error deleting category: %v", err)
	}

	// Try to get the soft deleted category
	_, err = cDb.GetById(context.Background(), &category.GetCategoryRequest{Id: createdCategory.Id})

	assert.Error(t, err, "Should be able to retrieve soft deleted post")
}

func TestGetAllCategories(t *testing.T) {
	cDb := newTestCategory(t)

	// Create some test categories
	testCategories := []*category.CreateCategoryRequest{
		{Name: "Category A"},
		{Name: "Category B"},
		{Name: "Category C"},
	}

	for _, tc := range testCategories {
		_, err := cDb.Create(context.Background(), tc)
		if err != nil {
			t.Fatalf("Error creating test category: %v", err)
		}
	}

	t.Run("ListCategories without filters", func(t *testing.T) {
		resp, err := cDb.GetAllCategories(context.Background(), &category.GetAllCategoriesRequest{})
		if err != nil {
			t.Fatalf("Error listing all categories: %v", err)
		}
		assert.GreaterOrEqual(t, len(resp.Categories), len(testCategories))
	})
}
