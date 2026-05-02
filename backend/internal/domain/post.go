package domain

import "time"

type Post struct {
	ID               string    `json:"id"`
	PublisherID      string    `json:"publisher_id"`
	Title            string    `json:"title"`
	Subtitle         string    `json:"subtitle"`
	Slug             string    `json:"slug"`
	Content          string    `json:"post_content"`
	VideoId          string    `json:"video_id"`
	VideoMetadata    string    `json:"video_metadata"`
	ImageUrl         string    `json:"image_url"`
	ImageDescription string    `json:"image_description"`
	NumRays          int8      `json:"num_rays"`
	NumComments      int8      `json:"num_comments"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}
