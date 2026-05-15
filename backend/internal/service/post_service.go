package service

import (
	"context"
	"fmt"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
)

type PostRepository interface {
	GetPostByID(ctx context.Context, userID string, postID string) (domain.Post, error)
	GetPostBySlug(ctx context.Context, userID string, slug string) (domain.Post, error)
	GetListOfPosts(ctx context.Context, userID string, limit int32, offset int32) ([]domain.Post, error)
	CreatePost(ctx context.Context, userID string, newPost domain.Post) (domain.Post, error)
	UpdatePost(ctx context.Context, userID string, postWithUpdates domain.Post) (domain.Post, error)
	SoftDeletePost(ctx context.Context, userID string, postID string) (domain.Post, error)
	PermanentlyDeletePost(ctx context.Context, userID string, postID string) (domain.Post, error)
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
}

type PostImageStorage interface {
	UploadPicture(ctx context.Context, fileBytes []byte, fileName string) (string, error)
	DeletePicture(ctx context.Context, fullImageURL string) error
}

type PostService struct {
	repo         PostRepository
	imageStorage PostImageStorage
}

func NewPostService(repo PostRepository, imageStorage PostImageStorage) *PostService {
	return &PostService{
		repo:         repo,
		imageStorage: imageStorage,
	}
}

func (s *PostService) FetchPostByID(ctx context.Context, userID string, postID string) (domain.Post, error) {
	if postID == "" {
		return domain.Post{}, fmt.Errorf("[post_service.go] FetchPostByID: post id is required.")
	}

	post, err := s.repo.GetPostByID(ctx, userID, postID)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] FetchPostByID: failed to fetch post %v for user %v: %w", postID, userID, err)
	}

	return post, nil
}

func (s *PostService) FetchPostBySlug(ctx context.Context, userID string, slug string) (domain.Post, error) {
	if slug == "" {
		return domain.Post{}, fmt.Errorf("[post_service.go] FetchPostBySlug: slug is required.")
	}

	post, err := s.repo.GetPostBySlug(ctx, userID, slug)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] FetchPostBySlug: failed to fetch post with slug %v for user %v: %w", slug, userID, err)
	}

	return post, nil
}

func (s *PostService) FetchListOfPosts(ctx context.Context, userID string, limit int32, offset int32) ([]domain.Post, error) {
	if limit > 100 {
		limit = 100
	}
	postList, err := s.repo.GetListOfPosts(ctx, userID, limit, offset)
	if err != nil {
		return []domain.Post{}, fmt.Errorf("[post_service.go] FetchListOfPosts: failed to get list of posts with limit %v and offset %v for user %v: %w", limit, offset, userID, err)
	}

	return postList, nil
}

func (s *PostService) IsSlugTaken(ctx context.Context, slug string) (bool, error) {
	if slug == "" {
		return true, fmt.Errorf("[post_service.go] IsSlugTaken: slug is required")
	}

	doesExist, err := s.repo.CheckSlugExists(ctx, slug)
	if err != nil {
		return true, fmt.Errorf("[post_service.go] IsSlug: failed to check if slug %s exists: %w", slug, err)
	}

	return doesExist, nil
}

func (s *PostService) CreateNewPost(ctx context.Context, userID string, newPost domain.Post, imageBytes []byte, fileExtension string, videoID string) (domain.Post, error) {
	if newPost.Content == nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] CreateNewPost: content is required.")
	}

	slugExists, err := s.IsSlugTaken(ctx, newPost.Slug)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] CreateNewPost: failed to check if slug %v already exists: %w", newPost.Slug, err)
	}

	if slugExists {
		return domain.Post{}, fmt.Errorf("Slug exists")
	}

	// probably need to account for an array of images, new small struct for imageBytes []byte
	// Need to also handle page_image, which should account for many
	// for _, images := range postImages
	// if len(imageBytes) > 0 {

	// }
	return domain.Post{}, nil
}
