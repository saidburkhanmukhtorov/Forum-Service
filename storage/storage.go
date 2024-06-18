package storage

import (
	"context"

	"github.com/Forum-service/Forum-Service/genproto/category"
	"github.com/Forum-service/Forum-Service/genproto/comment"
	"github.com/Forum-service/Forum-Service/genproto/post"
	"github.com/Forum-service/Forum-Service/genproto/posttag"
	"github.com/Forum-service/Forum-Service/genproto/tag"
)

// StorageI defines the interface for interacting with the forum service storage.
type StorageI interface {
	Category() CategoryRepo
	Tag() TagRepo
	Post() PostRepo
	Comment() CommentRepo
	PostTag() PostTagRepo
}

// CategoryRepo defines methods for managing categories.
type CategoryRepo interface {
	Create(ctx context.Context, req *category.CreateCategoryRequest) (*category.CreateCategoryResponse, error)
	GetById(ctx context.Context, req *category.GetCategoryRequest) (*category.GetCategoryResponse, error)
	Update(ctx context.Context, req *category.UpdateCategoryRequest) (*category.UpdateCategoryResponse, error)
	Delete(ctx context.Context, req *category.DeleteCategoryRequest) (*category.DeleteCategoryResponse, error)
	GetAllCategories(ctx context.Context, req *category.GetAllCategoriesRequest) (*category.GetAllCategoriesResponse, error)
}

// TagRepo defines methods for managing tags.
type TagRepo interface {
	Create(ctx context.Context, req *tag.CreateTagRequest) (*tag.CreateTagResponse, error)
	GetById(ctx context.Context, req *tag.GetTagRequest) (*tag.GetTagResponse, error)
	Update(ctx context.Context, req *tag.UpdateTagRequest) (*tag.UpdateTagResponse, error)
	Delete(ctx context.Context, req *tag.DeleteTagRequest) (*tag.DeleteTagResponse, error)
	GetAllTags(ctx context.Context, req *tag.GetAllTagsRequest) (*tag.GetAllTagsResponse, error)
}

// PostRepo defines methods for managing posts.
type PostRepo interface {
	Create(ctx context.Context, req *post.CreatePostRequest) (*post.CreatePostResponse, error)
	GetById(ctx context.Context, req *post.GetPostRequest) (*post.GetPostResponse, error)
	Update(ctx context.Context, req *post.UpdatePostRequest) (*post.UpdatePostResponse, error)
	Delete(ctx context.Context, req *post.DeletePostRequest) (*post.DeletePostResponse, error)
	GetAllPosts(ctx context.Context, req *post.GetAllPostsRequest) (*post.GetAllPostsResponse, error)
}

// CommentRepo defines methods for managing comments.
type CommentRepo interface {
	Create(ctx context.Context, req *comment.CreateCommentRequest) (*comment.CreateCommentResponse, error)
	GetById(ctx context.Context, req *comment.GetCommentRequest) (*comment.GetCommentResponse, error)
	Update(ctx context.Context, req *comment.UpdateCommentRequest) (*comment.UpdateCommentResponse, error)
	Delete(ctx context.Context, req *comment.DeleteCommentRequest) (*comment.DeleteCommentResponse, error)
	GetAllComments(ctx context.Context, req *comment.GetAllCommentsRequest) (*comment.GetAllCommentsResponse, error)
}

// PostTagRepo defines methods for managing post-tag associations.
type PostTagRepo interface {
	Create(ctx context.Context, req *posttag.CreatePostTagRequest) (*posttag.CreatePostTagResponse, error)
	Delete(ctx context.Context, req *posttag.DeletePostTagRequest) (*posttag.DeletePostTagResponse, error)
	GetAllPostTags(ctx context.Context, req *posttag.GetAllPostTagsRequest) (*posttag.GetAllPostTagsResponse, error)
	GetPostsByTag(ctx context.Context, req *posttag.GetPostsByTagRequest) (*posttag.GetPostsByTagResponse, error)
}
