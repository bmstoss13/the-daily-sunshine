package domain

import "time"

type Comment struct {
	ID              string     `json:"id"`
	PostID          string     `json:"post_id"`
	PublisherID     string     `json:"publisher_id"`
	ParentCommentID *string    `json:"parent_comment_id"`
	CommentBody     string     `json:"comment_body"`
	CommentBodyText string     `json:"comment_body_text"` //Maybe use if comments aren't rich text
	NumRays         int64      `json:"num_rays"`
	NumReplies      int64      `json:"num_replies"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}
