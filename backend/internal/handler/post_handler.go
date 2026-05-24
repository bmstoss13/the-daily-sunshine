package handler

import "github.com/bmstoss13/the-daily-sunshine/internal/service"

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}
