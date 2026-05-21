package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
)

type PostImageInput struct {
	ImageBytes    []byte
	FileExtension string
	AltText       *string
	Description   *string
	IsCoverImage  bool
}

type PostVideoInput struct {
	YoutubeVideoID string
	VideoMetadata  string
}

type PostRepository interface {
	GetPostByID(ctx context.Context, userID string, postID string) (domain.Post, error)
	GetPostBySlug(ctx context.Context, userID string, slug string) (domain.Post, error)
	GetListOfPosts(ctx context.Context, userID string, limit int32, offset int32) ([]domain.Post, error)
	GetPostsForToday(ctx context.Context, userID string, startOfDay time.Time, endOfDay time.Time) ([]domain.Post, error)
	CreatePost(ctx context.Context, userID string, newPost domain.Post) (domain.Post, error)
	UpdatePost(ctx context.Context, userID string, postWithUpdates domain.Post) (domain.Post, error)
	SoftDeletePost(ctx context.Context, userID string, postID string) (domain.Post, error)
	PermanentlyDeletePost(ctx context.Context, userID string, postID string) (domain.Post, error)
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
}

type PostImageRepository interface {
	GetPostImageByID(ctx context.Context, userID string, imageID string) (domain.PostImage, error)
	GetCoverImageOfPost(ctx context.Context, userID string, postID string) (domain.PostImage, error)
	GetPostImagesByPost(ctx context.Context, userID string, postID string) ([]domain.PostImage, error)
	CreatePostImage(ctx context.Context, userID string, newPostImage domain.PostImage) (domain.PostImage, error)
	UpdatePostImage(ctx context.Context, userID string, postImageWithUpdates domain.PostImage) (domain.PostImage, error)
	DeletePostImage(ctx context.Context, userID string, postID string, postImageID string) (domain.PostImage, error)
}

type PostVideoRepository interface {
	GetPostVideoByID(ctx context.Context, userID string, videoID string) (domain.PostVideo, error)
	GetPostVideoByPost(ctx context.Context, userID string, postID string) (domain.PostVideo, error)
	GetListOfPostVideos(ctx context.Context, userID string, limit int32, offset int32) ([]domain.PostVideo, error)
	CreatePostVideo(ctx context.Context, userID string, newVideo domain.PostVideo) (domain.PostVideo, error)
	UpdatePostVideo(ctx context.Context, userID string, videoToUpdate domain.PostVideo) (domain.PostVideo, error)
	DeletePostVideo(ctx context.Context, userID string, videoID string, postID string) (domain.PostVideo, error)
}

type PostImageStorage interface {
	UploadPicture(ctx context.Context, fileBytes []byte, fileName string) (string, error)
	DeletePicture(ctx context.Context, fullImageURL string) error
}

type PostCache interface {
	GetTopPostsOfDay(ctx context.Context, dayKey string) ([]domain.Post, bool, error)
	SetTopPostsOfDay(ctx context.Context, dayKey string, posts []domain.Post, ttl time.Duration) error
	InvalidateTopPostsOfDay(ctx context.Context, dayKey string) error
}

type PostService struct {
	repo         PostRepository
	imageRepo    PostImageRepository
	videoRepo    PostVideoRepository
	imageStorage PostImageStorage
	cache        PostCache
}

const appDayLocationName = "America/New_York"

func NewPostService(repo PostRepository, imageRepo PostImageRepository, videoRepo PostVideoRepository, imageStorage PostImageStorage, cache PostCache) *PostService {
	return &PostService{
		repo:         repo,
		imageRepo:    imageRepo,
		videoRepo:    videoRepo,
		imageStorage: imageStorage,
		cache:        cache,
	}
}

func (s *PostService) FetchPostByID(ctx context.Context, userID string, postID string) (domain.Post, error) {
	if postID == "" {
		return domain.Post{}, fmt.Errorf("[post_service.go] FetchPostByID: post id is required.")
	}

	post, err := s.repo.GetPostByID(ctx, userID, postID)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] FetchPostByID: failed to fetch post %v for user %v: %w", postID, userID, err)
	}

	return post, nil
}

func (s *PostService) FetchPostBySlug(ctx context.Context, userID string, slug string) (domain.Post, error) {
	if slug == "" {
		return domain.Post{}, fmt.Errorf("[post_service.go] FetchPostBySlug: slug is required.")
	}

	post, err := s.repo.GetPostBySlug(ctx, userID, slug)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] FetchPostBySlug: failed to fetch post with slug %v for user %v: %w", slug, userID, err)
	}

	return post, nil
}

func (s *PostService) FetchListOfPosts(ctx context.Context, userID string, limit int32, offset int32) ([]domain.Post, error) {
	if limit > 100 {
		limit = 100
	}
	postList, err := s.repo.GetListOfPosts(ctx, userID, limit, offset)
	if err != nil {
		return []domain.Post{}, fmt.Errorf("[post_service.go] FetchListOfPosts: failed to get list of posts with limit %v and offset %v for user %v: %w", limit, offset, userID, err)
	}

	return postList, nil
}

func topPostsDayWindow(now time.Time) (string, time.Time, time.Time, error) {
	loc, err := time.LoadLocation(appDayLocationName)
	if err != nil {
		return "", time.Time{}, time.Time{}, fmt.Errorf("[post_service.go] topPostsDayWindow: failed to load location %s: %w", appDayLocationName, err)
	}

	localNow := now.In(loc)
	startOfDayLocal := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), 0, 0, 0, 0, loc)
	endOfDayLocal := startOfDayLocal.Add(24 * time.Hour)
	dayKey := "top-posts:" + startOfDayLocal.Format("2006-01-02")

	return dayKey, startOfDayLocal.UTC(), endOfDayLocal.UTC(), nil
}

func (s *PostService) FetchTopPostsOfTheDay(ctx context.Context, userID string) ([]domain.Post, error) {
	dayKey, startOfDay, endOfDay, err := topPostsDayWindow(time.Now())
	if err != nil {
		return []domain.Post{}, err
	}

	redisPostList, found, err := s.cache.GetTopPostsOfDay(ctx, dayKey)
	if err != nil {
		log.Printf("[post_service.go] FetchTopPostsOfTheDay: failed to fetch top posts of the day from redis: %v", err)
	}

	if found {
		return redisPostList, nil
	}

	dbPostList, dbErr := s.repo.GetPostsForToday(ctx, userID, startOfDay, endOfDay)
	if dbErr != nil {
		return []domain.Post{}, fmt.Errorf("[post_service.go] FetchTopPostsOfTheDay: failed to fetch top posts of the day from db: %w", dbErr)
	}

	for i := range dbPostList {
		postImages, imgErr := s.imageRepo.GetPostImagesByPost(ctx, userID, dbPostList[i].ID)
		if imgErr != nil {
			fmt.Printf("failed to fetch images for post %s: %s", dbPostList[i].ID, imgErr)
		}
		dbPostList[i].Images = append(dbPostList[i].Images, postImages...)

		postVideo, vidErr := s.videoRepo.GetPostVideoByPost(ctx, userID, dbPostList[i].ID)
		if vidErr != nil {
			fmt.Printf("failed to fetch video for post %v: %v", dbPostList[i].ID, vidErr)
		}
		if postVideo != (domain.PostVideo{}) {
			dbPostList[i].Video = &postVideo
		}
	}

	// setting ttl to 1 hour. This does not need to be super active and can honestly be a greater value if need be
	redisErr := s.cache.SetTopPostsOfDay(ctx, dayKey, dbPostList, time.Hour*1)
	if redisErr != nil {
		fmt.Printf("[post_service.go] FetchTopPostsOfTheDay: failed to set top posts in redis: %v", redisErr)
	}

	return dbPostList, nil
}

func (s *PostService) IsSlugTaken(ctx context.Context, slug string) (bool, error) {
	if slug == "" {
		return true, fmt.Errorf("[post_service.go] IsSlugTaken: slug is required")
	}

	doesExist, err := s.repo.CheckSlugExists(ctx, slug)
	if err != nil {
		return true, fmt.Errorf("[post_service.go] IsSlug: failed to check if slug %s exists: %w", slug, err)
	}

	return doesExist, nil
}

func (s *PostService) CreateNewPost(ctx context.Context, userID string, newPost domain.Post, imageList []PostImageInput, video PostVideoInput) (domain.Post, error) {
	if newPost.Title == "" {
		return domain.Post{}, fmt.Errorf("[post_service.go] CreateNewPost: title is required.")
	}

	if newPost.Slug == "" {
		return domain.Post{}, fmt.Errorf("[post_service.go] CreateNewPost: slug is required.")
	}

	slugExists, err := s.IsSlugTaken(ctx, newPost.Slug)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] CreateNewPost: failed to check if slug %v already exists: %w", newPost.Slug, err)
	}

	if slugExists {
		return domain.Post{}, fmt.Errorf("Slug already exists")
	}

	if newPost.Status == domain.Published {
		if newPost.Content == nil {
			return domain.Post{}, fmt.Errorf("[post_service.go] CreateNewPost: content is required.")
		}
	}

	createdPost, err := s.repo.CreatePost(ctx, userID, newPost)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] CreateNewPost: failed to insert into database: %w", err)
	}

	for i, image := range imageList {
		if len(image.ImageBytes) > 0 {
			fileName := fmt.Sprintf("posts/%s/%v", createdPost.ID, i)
			publicURL, err := s.imageStorage.UploadPicture(ctx, image.ImageBytes, fileName)
			if err != nil {
				return domain.Post{}, fmt.Errorf("[post_service.go] CreateNewPost: failed to upload post picture: %w", err)
			}

			imageParams := domain.PostImage{
				PostID:           createdPost.ID,
				ImageURL:         publicURL,
				ImageDescription: image.Description,
				AltText:          image.AltText,
				IsCoverImage:     image.IsCoverImage,
			}

			newPostImage, imgErr := s.imageRepo.CreatePostImage(ctx, userID, imageParams)
			if imgErr != nil {
				deletePost, deletePostErr := s.repo.PermanentlyDeletePost(ctx, userID, createdPost.ID)
				if deletePostErr != nil {
					log.Printf("[post_service.go] CreateNewPost: failed to delete created post %v: %v", deletePost.ID, deletePostErr)
				}

				deleteImgErr := s.imageStorage.DeletePicture(ctx, publicURL)
				if deleteImgErr != nil {
					log.Printf("[post_service.go] CreateNewPost: failed to delete post image with url %v: %v", newPostImage.ImageURL, err)
				}

				return domain.Post{}, fmt.Errorf("[post_service.go] CreateNewPost: failed to create post in database: %w", imgErr)
			}
		}
	}

	return createdPost, nil
}
