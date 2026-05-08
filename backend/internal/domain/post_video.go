package domain

type PostVideo struct {
	ID             string `json:"id"`
	PostID         string `json:"post_id"`
	YouTubeVideoID string `json:"youtube_video_id"`
	VideoMetadata  string `json:"video_metadata"`
	CreatedAt      string `json:"created_at"`
}
