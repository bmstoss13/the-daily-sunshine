package service

import (
	"context"
	"testing"
	"time"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
)

type mockPostRepo struct {
	mockPost            domain.Post
	mockPostFromID      domain.Post
	mockPostFromSlug    domain.Post
	mockPostList        []domain.Post
	mockTopPosts        []domain.Post
	mockCreatePost      domain.Post
	mockUpdatedPost     domain.Post
	mockSoftDeletedPost domain.Post
	mockDeletedPost     domain.Post
	mockSlugExists      bool

	mockRetrieveError error
	mockErrFromID     error
	mockErrFromSlug   error
	// mockErrFromPostList          error
	mockErrFromTopPosts       error
	mockErrFromCreatePost     error
	mockErrFromUpdatePost     error
	mockErrFromSoftDeletePost error
	mockErrFromDeletePost     error
	mockErrFromSlugExists     error

	hasMockPostFromID    bool
	hasMockPostFromSlug  bool
	hasMockPostList      bool
	hasMockTopPosts      bool
	isMockCreatePost     bool
	isMockUpdatePost     bool
	isMockSoftDeletePost bool
	isMockDeletePost     bool
	isMockSlugExists     bool

	lastUserID       string
	lastPostID       string
	lastSlug         string
	lastPostToCreate domain.Post
	lastUpdatedPost  domain.Post
	lastLimit        int32
	lastOffset       int32
	lastStartOfDay   time.Time
	lastEndOfDay     time.Time
}

type mockPostImageRepo struct {
	mockPostImageByID       domain.PostImage
	mockCoverImageOfPost    domain.PostImage
	mockPostImageListByPost []domain.PostImage
	mockCreatePostImage     domain.PostImage
	mockUpdatedPostImage    domain.PostImage
	mockDeletedPostImage    domain.PostImage

	mockRetrieveError error
	mockErrFromID     error
	mockErrFromPost   error
	mockErrFromCreate error
	mockErrFromUpdate error
	mockErrFromDelete error

	hasMockImageByID    bool
	hasListImagesByPost bool
	isCreatePostImage   bool
	isUpdatePostImage   bool
	isDeletePostImage   bool

	lastUserID       string
	lastPostID       string
	lastImageID      string
	lastUpdatedImage string
}

type mockPostVideoRepo struct {
	mockPostVideo       domain.PostVideo
	mockPostVideoByID   domain.PostVideo
	mockPostVideoByPost domain.PostVideo
	mockCreatePostVideo domain.PostVideo
	mockUpdatePostVideo domain.PostVideo
	mockDeletePostVideo domain.PostVideo

	mockRetrieveError error
	mockErrFromID     error
	mockErrFromPost   error
	mockErrFromCreate error
	mockErrFromUpdate error
	mockErrFromDelete error

	hasMockVideoByID      bool
	hasMockVideoBySlug    bool
	isMockCreatePostVideo bool
	isMockUpdatePostVideo bool
	isMockDeletePostVideo bool
}

type mockPostStorage struct {
	mockUrl string

	mockErrFromDelete error

	lastUploadedBytes    []byte
	lastUploadedFileName string
}

type mockPostCache struct {
	mockTopPosts []domain.Post

	mockErrFromGet        error
	mockErrFromSet        error
	mockErrFromInvalidate error

	mockFound bool

	lastInvalidatedKey string
}

// POST REPO
func (m *mockPostRepo) GetPostByID(_ context.Context, userID string, postID string) (domain.Post, error) {
	m.hasMockPostFromID = true
	m.lastUserID = userID
	m.lastPostID = postID
	if m.mockErrFromID != nil {
		return domain.Post{}, m.mockErrFromID
	}
	if m.hasMockPostFromID {
		return m.mockPostFromID, nil
	}
	return m.mockPost, m.mockRetrieveError
}

func (m *mockPostRepo) GetPostBySlug(_ context.Context, userID string, slug string) (domain.Post, error) {
	m.hasMockPostFromSlug = true
	m.lastUserID = userID
	m.lastSlug = slug
	if m.mockErrFromSlug != nil {
		return domain.Post{}, m.mockErrFromSlug
	}
	if m.hasMockPostFromSlug {
		return m.mockPostFromSlug, nil
	}
	return m.mockPost, m.mockRetrieveError
}

func (m *mockPostRepo) GetListOfPosts(_ context.Context, userID string, limit int32, offset int32) ([]domain.Post, error) {
	m.hasMockPostList = true
	m.lastUserID = userID
	m.lastLimit = limit
	m.lastOffset = offset
	return m.mockPostList, m.mockRetrieveError
}

func (m *mockPostRepo) GetPostsForToday(_ context.Context, userID string, startOfDay time.Time, endOfDay time.Time) ([]domain.Post, error) {
	m.hasMockTopPosts = true
	m.lastUserID = userID
	m.lastStartOfDay = startOfDay
	m.lastEndOfDay = endOfDay
	return m.mockTopPosts, m.mockRetrieveError
}

func (m *mockPostRepo) CheckSlugExists(_ context.Context, slug string) (bool, error) {
	m.isMockSlugExists = true
	return m.mockSlugExists, m.mockErrFromSlugExists
}

func (m *mockPostRepo) CheckSlugExistsOtherPosts(_ context.Context, slug string, postID string) (bool, error) {
	m.isMockSlugExists = true
	m.lastPostID = postID
	return m.mockSlugExists, m.mockErrFromSlugExists
}

func (m *mockPostRepo) CreatePost(_ context.Context, userID string, newPost domain.Post) (domain.Post, error) {
	m.isMockCreatePost = true
	m.lastUserID = userID
	m.lastPostToCreate = newPost
	return m.mockCreatePost, m.mockErrFromCreatePost
}

func (m *mockPostRepo) UpdatePost(_ context.Context, userID string, postWithUpdates domain.Post) (domain.Post, error) {
	m.isMockUpdatePost = true
	m.lastUserID = userID
	m.lastUpdatedPost = postWithUpdates
	return m.mockUpdatedPost, m.mockErrFromUpdatePost
}

func (m *mockPostRepo) SoftDeletePost(_ context.Context, userID string, postID string) (domain.Post, error) {
	m.isMockSoftDeletePost = true
	m.lastUserID = userID
	m.lastPostID = postID
	return m.mockSoftDeletedPost, m.mockErrFromSoftDeletePost
}

func (m *mockPostRepo) PermanentlyDeletePost(_ context.Context, userID string, postID string) (domain.Post, error) {
	m.isMockDeletePost = true
	m.lastUserID = userID
	m.lastPostID = postID
	return m.mockDeletedPost, m.mockErrFromDeletePost
}

// IMAGE REPO
func (m *mockPostImageRepo) GetPostImageByID(ctx context.Context, userID string, imageID string) (domain.PostImage, error) {
	m.hasMockImageByID = true
	m.lastUserID = userID
	m.lastImageID = imageID
	return m.mockPostImageByID, m.mockErrFromID
}

func TestPostService_FetchPostByID(t *testing.T) {
	// t.Run("returns validation error when target post ID is empty", func(t *testing.T) {
	// 	repo := &mockPostRepo{}
	// 	imageRepo := &mockPostImageRepo{}
	// 	videoRepo := &mockPostVideoRepo{}
	// 	imageStorage := &mockPostStorage{}
	// 	postCache := &mockPostCache{}

	// 	service := NewPostService(repo, imageRepo, videoRepo, imageStorage, postCache)
	// })
}
