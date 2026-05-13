package domain

import "time"

type Post struct {
	ID          string     `json:"id"`
	PublisherID string     `json:"publisher_id"`
	Title       string     `json:"title"`
	Subtitle    *string    `json:"subtitle"`
	Slug        string     `json:"slug"`
	Content     *string    `json:"post_content"`
	Status      PostStatus `json:"status"`
	NumRays     int64      `json:"num_rays"`
	NumComments int64      `json:"num_comments"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`

	Publisher  *Profile
	CoverImage *PostImage
	Video      *PostVideo
}

type PostStatus string

const (
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
	Archived  PostStatus = "archived"
)
