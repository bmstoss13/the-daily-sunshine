package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/bmstoss13/the-daily-sunshine/internal/helpers"
	"github.com/bmstoss13/the-daily-sunshine/internal/service"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	svc *service.PostService
}

type CreatePostRequest struct {
	Title    string  `form:"title" binding:"required"`
	Subtitle *string `form:"subtitle"`
	Slug     string  `form:"slug"`                      //maybe we should handle this on the server? No user wants to make a slug. Or the frontend can process.
	Content  *string `form:"content"`                   //Not required because can leave status as draft
	Status   string  `form:"status" binding:"required"` //leaving as string but frontend will be binding this as a type. Should still check.
}

type UpdatePostRequest struct {
	Title    string  `form:"title" binding:"required"`
	Subtitle *string `form:"subtitle"`
	Slug     string  `form:"slug"`
	Content  *string `form:"content"`
	Status   string  `form:"status" binding:"required"`
}

type ImageOpRequest struct {
	ImageID          *string `json:"image_id"`
	ImageDescription *string `json:"image_description"`
	AltText          *string `json:"alt_text"`
	IsCoverImage     bool    `json:"is_cover_image"`
	DeleteExisting   bool    `json:"delete_existing"`
	FileIndex        *int    `json:"file_index"`
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{
		svc: svc,
	}
}

func (p *PostHandler) GetPostByID(c *gin.Context) {
	targetPostID := c.Param("id")
	requestingUserID := c.GetString("userID")

	post, err := p.svc.FetchPostByID(c.Request.Context(), requestingUserID, targetPostID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch post",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

func (p *PostHandler) GetPostBySlug(c *gin.Context) {
	targetPostSlug := c.Param("slug")
	requestingUserID := c.GetString("userID")

	post, err := p.svc.FetchPostBySlug(c.Request.Context(), requestingUserID, targetPostSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch post by slug",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})
}

func (p *PostHandler) GetListOfPosts(c *gin.Context) {

	var limit int32 = 20
	var offset int32 = 0

	limitStr := c.Query("limit") //?limit=5 for example
	if limitStr != "" {
		parsedLimit, err := helpers.StringToBase10Int32(limitStr)
		if err != nil || parsedLimit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid limit parameter. Must be positive integer.",
			})
			return
		}
		limit = parsedLimit
	}

	offsetStr := c.Query("offset")
	if offsetStr != "" {
		parsedOffset, err := helpers.StringToBase10Int32(offsetStr)
		if err != nil || parsedOffset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid offset parameter. Cannot be negative.",
			})
			return
		}
		offset = parsedOffset
	}

	requestingUserID := c.GetString("userID")
	postList, err := p.svc.FetchListOfPosts(c.Request.Context(), requestingUserID, limit, offset)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch post list",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": postList,
	})
}

func (p *PostHandler) GetSlugAvailability(c *gin.Context) {
	slugToCheck := c.GetString("slug")
	isTaken, err := p.svc.IsSlugTaken(c.Request.Context(), slugToCheck)

	if err != nil {
		if slugToCheck == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Slug query is required",
			})
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to check slug availability",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"taken": isTaken,
	})
}

func (p *PostHandler) GetTopPostsOfDay(c *gin.Context) {
	requestingUserID := c.GetString("userID")
	topPosts, err := p.svc.FetchTopPostsOfTheDay(c.Request.Context(), requestingUserID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch top posts of the day",
			"details": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": topPosts,
	})
}

func (p *PostHandler) CreatePost(c *gin.Context) {
	requestingUserID := c.GetString("userID")

	var req CreatePostRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON payload",
			"details": err.Error(),
		})
		return
	}

	newPost := domain.Post{
		PublisherID: requestingUserID,
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Slug:        req.Slug,
		Content:     req.Content,
		Status:      domain.PostStatus(req.Status),
	}

	imageInputs, err := buildPostImageInputs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid image payload",
			"details": err.Error(),
		})
		return
	}

	var dummyVideo service.PostVideoInput //TODO: establish videos on posts using youtube data api
	createdPost, err := p.svc.CreateNewPost(c.Request.Context(), requestingUserID, newPost, imageInputs, dummyVideo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create new post",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": createdPost,
	})
}

func (p *PostHandler) UpdatePost(c *gin.Context) {
	targetPostID := c.Param("id")
	requestingUserID := c.GetString("userID")

	var req UpdatePostRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON payload",
			"details": err.Error(),
		})
		return
	}

	updatedPost := domain.Post{
		ID:          targetPostID,
		PublisherID: requestingUserID,
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Slug:        req.Slug,
		Content:     req.Content,
		Status:      domain.PostStatus(req.Status),
	}

	imageInputs, err := buildPostImageInputs(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid image payload",
			"details": err.Error(),
		})
		return
	}

	var dummyVideo service.PostVideoInput //TODO: establish videos on posts using youtube data api
	updatedPostResult, err := p.svc.UpdatePost(c.Request.Context(), requestingUserID, updatedPost, imageInputs, dummyVideo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create update",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": updatedPostResult,
	})
}

func (p *PostHandler) SoftDeletePost(c *gin.Context) {
	targetPostID := c.Param("id")
	requestingUserID := c.GetString("userID")

	softDeletedPost, err := p.svc.DeletePostSoft(c.Request.Context(), requestingUserID, targetPostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to soft delete post",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": softDeletedPost,
	})
}

func (p *PostHandler) DeletePost(c *gin.Context) {
	targetPostID := c.Param("id")
	requestingUserID := c.GetString("userID")

	deletedPost, err := p.svc.DeletePost(c.Request.Context(), requestingUserID, targetPostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete post",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": deletedPost,
	})
}

func buildPostImageInputs(c *gin.Context) ([]service.PostImageInput, error) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		return nil, fmt.Errorf("invalid multipart form: %w", err)
	}

	images := c.Request.MultipartForm.File["images"]
	imageOpsJSON := c.PostForm("image_ops")

	if imageOpsJSON == "" {
		return nil, nil
	}

	var imageOps []ImageOpRequest
	if err := json.Unmarshal([]byte(imageOpsJSON), &imageOps); err != nil {
		return nil, fmt.Errorf("invalid image metadata: %w", err)
	}

	imageInputs := make([]service.PostImageInput, 0, len(imageOps))

	for _, op := range imageOps {
		input := service.PostImageInput{
			ImageID:        op.ImageID,
			Description:    op.ImageDescription,
			AltText:        op.AltText,
			IsCoverImage:   op.IsCoverImage,
			DeleteExisting: op.DeleteExisting,
		}

		if op.FileIndex != nil {
			if *op.FileIndex < 0 || *op.FileIndex >= len(images) {
				return nil, fmt.Errorf("image file_index %d is out of range", *op.FileIndex)
			}

			fileHeader := images[*op.FileIndex]
			file, err := fileHeader.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to open uploaded image: %w", err)
			}

			fileBytes, err := io.ReadAll(file)
			file.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read uploaded image: %w", err)
			}

			input.ImageBytes = fileBytes
			input.FileExtension = filepath.Ext(fileHeader.Filename)
		}

		imageInputs = append(imageInputs, input)
	}

	return imageInputs, nil
}
