package domain

import "time"

type PostImage struct {
	ID               string    `json:"id"`
	PostID           string    `json:"post_id"`
	ImageURL         string    `json:"image_url"`
	ImageDescription string    `json:"image_description"`
	AltText          string    `json:"alt_text"`
	IsCoverImage     int32     `json:"is_cover_image"`
	CreatedAt        time.Time `json:"created_at"`
}
