package service

import (
	"context"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
)

type mockPostRepo struct {
	mockPost            domain.Post
	mockPostFromID      domain.Post
	mockUpdatedPost     domain.Post
	mockSoftDeletedPost domain.Post

	mockErrFromID     error
	mockRetrieveError error

	hasMockPostFromID bool
	lastUserID        string
	lastPostID        string
}

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
