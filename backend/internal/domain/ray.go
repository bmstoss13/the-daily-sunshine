package domain

import "time"

// Likes on posts/comments
type Ray struct {
	ID          string    `json:"id"`
	PublisherID string    `json:"publisher_id"`
	PostID      *string   `json:"post_id"`
	CommentID   *string   `json:"comment_id"`
	CreatedAt   time.Time `json:"created_at"`
}
