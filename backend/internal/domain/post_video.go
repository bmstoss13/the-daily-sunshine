package domain

import "time"

type PostVideo struct {
	ID             string    `json:"id"`
	PostID         string    `json:"post_id"`
	YouTubeVideoID string    `json:"youtube_video_id"`
	VideoMetadata  string    `json:"video_metadata"`
	CreatedAt      time.Time `json:"created_at"`

	Post *Post
}
