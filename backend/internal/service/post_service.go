package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
)

type PostImageInput struct {
	ImageID        *string
	ImageBytes     []byte
	FileExtension  string
	AltText        *string
	Description    *string
	IsCoverImage   bool
	DeleteExisting bool
}

type PostVideoInput struct {
	YoutubeVideoID string
	VideoMetadata  string
}

type IPostRepository interface {
	GetPostByID(ctx context.Context, userID string, postID string) (domain.Post, error)
	GetPostBySlug(ctx context.Context, userID string, slug string) (domain.Post, error)
	GetListOfPosts(ctx context.Context, userID string, limit int32, offset int32) ([]domain.Post, error)
	GetPostsForToday(ctx context.Context, userID string, startOfDay time.Time, endOfDay time.Time) ([]domain.Post, error)
	CreatePost(ctx context.Context, userID string, newPost domain.Post) (domain.Post, error)
	UpdatePost(ctx context.Context, userID string, postWithUpdates domain.Post) (domain.Post, error)
	SoftDeletePost(ctx context.Context, userID string, postID string) (domain.Post, error)
	PermanentlyDeletePost(ctx context.Context, userID string, postID string) (domain.Post, error)
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
	CheckSlugExistsOtherPosts(ctx context.Context, slug string, postID string) (bool, error)
}

type IPostImageRepository interface {
	GetPostImageByID(ctx context.Context, userID string, imageID string) (domain.PostImage, error)
	GetCoverImageOfPost(ctx context.Context, userID string, postID string) (domain.PostImage, error)
	GetPostImagesByPost(ctx context.Context, userID string, postID string) ([]domain.PostImage, error)
	CreatePostImage(ctx context.Context, userID string, newPostImage domain.PostImage) (domain.PostImage, error)
	UpdatePostImage(ctx context.Context, userID string, postImageWithUpdates domain.PostImage) (domain.PostImage, error)
	DeletePostImage(ctx context.Context, userID string, postID string, postImageID string) (domain.PostImage, error)
}

type IPostVideoRepository interface {
	GetPostVideoByID(ctx context.Context, userID string, videoID string) (domain.PostVideo, error)
	GetPostVideoByPost(ctx context.Context, userID string, postID string) (domain.PostVideo, error)
	GetListOfPostVideos(ctx context.Context, userID string, limit int32, offset int32) ([]domain.PostVideo, error)
	CreatePostVideo(ctx context.Context, userID string, newVideo domain.PostVideo) (domain.PostVideo, error)
	UpdatePostVideo(ctx context.Context, userID string, videoToUpdate domain.PostVideo) (domain.PostVideo, error)
	DeletePostVideo(ctx context.Context, userID string, videoID string, postID string) (domain.PostVideo, error)
}

type IPostImageStorage interface {
	UploadPicture(ctx context.Context, fileBytes []byte, fileName string) (string, error)
	DeletePicture(ctx context.Context, fullImageURL string) error
}

type IPostCache interface {
	GetTopPostsOfDay(ctx context.Context, dayKey string) ([]domain.Post, bool, error)
	SetTopPostsOfDay(ctx context.Context, dayKey string, posts []domain.Post, ttl time.Duration) error
	InvalidateTopPostsOfDay(ctx context.Context, dayKey string) error
}

type PostService struct {
	repo         IPostRepository
	imageRepo    IPostImageRepository
	videoRepo    IPostVideoRepository
	imageStorage IPostImageStorage
	cache        IPostCache
}

const appDayLocationName = "America/New_York"

func NewPostService(repo IPostRepository, imageRepo IPostImageRepository, videoRepo IPostVideoRepository, imageStorage IPostImageStorage, cache IPostCache) *PostService {
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

func (s *PostService) invalidateTopPostsOfDay(ctx context.Context) {
	dayKey, _, _, err := topPostsDayWindow(time.Now())
	if err != nil {
		log.Printf("[post_service.go] invalidateTopPostsOfDay: failed to compute top posts of day key: %v", err)
		return
	}

	if err := s.cache.InvalidateTopPostsOfDay(ctx, dayKey); err != nil {
		log.Printf("[post_service.go] invalidateTopPostsOfDay: failed to invalidate top posts cache: %v", err)
	}
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
		return domain.Post{}, fmt.Errorf("CreateNewPost: Slug already exists")
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

	s.invalidateTopPostsOfDay(ctx)
	return createdPost, nil
}

func (s *PostService) UpdatePost(ctx context.Context, userID string, postWithUpdates domain.Post, imageList []PostImageInput, video PostVideoInput) (domain.Post, error) {
	if postWithUpdates.Title == "" {
		return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: title is required.")
	}

	if postWithUpdates.Slug == "" {
		return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: Slug is required.")
	}

	slugExists, slugErr := s.repo.CheckSlugExistsOtherPosts(ctx, postWithUpdates.Slug, postWithUpdates.ID)
	if slugErr != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to check if slug %v already exists: %w", postWithUpdates.Slug, slugErr)
	}

	if slugExists {
		return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: Slug already exists")
	}

	if postWithUpdates.Status == domain.Published {
		if postWithUpdates.Content == nil {
			return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: Content is required.")
		}
	}

	coverCount := 0
	for _, image := range imageList {
		if image.IsCoverImage && !image.DeleteExisting {
			coverCount++
		}
	}
	if coverCount > 1 {
		return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: only one cover image is allowed")
	}

	updatedPost, err := s.repo.UpdatePost(ctx, userID, postWithUpdates)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to update in database: %w", err)
	}

	existingImages, err := s.imageRepo.GetPostImagesByPost(ctx, userID, updatedPost.ID)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to fetch existing post images: %w", err)
	}

	existingByID := make(map[string]domain.PostImage, len(existingImages))
	seenImageIDs := make(map[string]struct{}, len(imageList))

	for _, img := range existingImages {
		existingByID[img.ID] = img
	}

	for i, image := range imageList {
		switch {
		// fetch row, delete R2 object, delete DB row
		case image.DeleteExisting && image.ImageID != nil:
			existing, ok := existingByID[*image.ImageID]
			if !ok {
				return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: image %s does not belong to post %s", *image.ImageID, updatedPost.ID)
			}

			if err := s.imageStorage.DeletePicture(ctx, existing.ImageURL); err != nil {
				return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to delete post image from storage: %w", err)
			}

			if _, err := s.imageRepo.DeletePostImage(ctx, userID, updatedPost.ID, existing.ID); err != nil {
				return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to delete post image from db: %w", err)
			}

			seenImageIDs[existing.ID] = struct{}{}

		// replace existing image: upload new, update row, delete old object
		case image.ImageID != nil && len(image.ImageBytes) > 0:
			existing, ok := existingByID[*image.ImageID]
			if !ok {
				return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: image %s does not belong to post %s", *image.ImageID, updatedPost.ID)
			}

			fileName := fmt.Sprintf("posts/%s/%v", updatedPost.ID, existing.ID)
			publicURL, err := s.imageStorage.UploadPicture(ctx, image.ImageBytes, fileName)
			if err != nil {
				return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to upload replacement image: %w", err)
			}

			imageParams := domain.PostImage{
				PostID:           updatedPost.ID,
				ImageURL:         publicURL,
				ImageDescription: image.Description,
				AltText:          image.AltText,
				IsCoverImage:     image.IsCoverImage,
			}

			_, imgErr := s.imageRepo.CreatePostImage(ctx, userID, imageParams)
			if imgErr != nil {
				deleteImgErr := s.imageStorage.DeletePicture(ctx, publicURL)
				if deleteImgErr != nil {
					log.Printf("[post_service.go] UpdatePost: failed to delete post image with url %v: %v", publicURL, err)
				}
				return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to update post in database: %w", imgErr)
			}

			if err := s.imageStorage.DeletePicture(ctx, existing.ImageURL); err != nil {
				log.Printf("[post_service.go] UpdatePost: orphaned old image left in R2 for image %s: %v", existing.ID, err)
			}

			seenImageIDs[existing.ID] = struct{}{}

		// create new image row
		case image.ImageID == nil && len(image.ImageBytes) > 0:
			fileName := fmt.Sprintf("posts/%s/%v", updatedPost.ID, i)
			publicURL, err := s.imageStorage.UploadPicture(ctx, image.ImageBytes, fileName)
			if err != nil {
				return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to upload new image: %w", err)
			}

			imageParams := domain.PostImage{
				PostID:           updatedPost.ID,
				ImageURL:         publicURL,
				ImageDescription: image.Description,
				AltText:          image.AltText,
				IsCoverImage:     image.IsCoverImage,
			}

			createdImage, imgErr := s.imageRepo.UpdatePostImage(ctx, userID, imageParams)
			if imgErr != nil {
				deleteImgErr := s.imageStorage.DeletePicture(ctx, publicURL)
				if deleteImgErr != nil {
					log.Printf("[post_service.go] UpdatePost: failed to delete post image with url %v: %v", publicURL, err)
				}
				return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to insert image in database: %w", imgErr)
			}

			seenImageIDs[createdImage.ID] = struct{}{}

		// metadata-only update
		case image.ImageID != nil && len(image.ImageBytes) == 0:
			existing, ok := existingByID[*image.ImageID]
			if !ok {
				return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: image %s does not belong to post %s", *image.ImageID, updatedPost.ID)
			}

			imageParams := domain.PostImage{
				ID:               existing.ID,
				PostID:           updatedPost.ID,
				ImageURL:         existing.ImageURL,
				ImageDescription: image.Description,
				AltText:          image.AltText,
				IsCoverImage:     image.IsCoverImage,
			}

			if _, err := s.imageRepo.UpdatePostImage(ctx, userID, imageParams); err != nil {
				return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to update image metadata: %w", err)
			}

			seenImageIDs[existing.ID] = struct{}{}

		// invalid/no-op input
		default:
			return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: invalid image update payload")
		}
	}

	for _, existing := range existingImages {
		if _, ok := seenImageIDs[existing.ID]; ok {
			continue
		}

		if err := s.imageStorage.DeletePicture(ctx, existing.ImageURL); err != nil {
			return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to delete removed image from storage: %w", err)
		}

		if _, err := s.imageRepo.DeletePostImage(ctx, userID, updatedPost.ID, existing.ID); err != nil {
			return domain.Post{}, fmt.Errorf("[post_service.go] UpdatePost: failed to delete removed image from database: %w", err)
		}
	}

	s.invalidateTopPostsOfDay(ctx)
	return updatedPost, nil
}

func (s *PostService) DeletePostSoft(ctx context.Context, userID string, postID string) (domain.Post, error) {
	softDeletedPost, err := s.repo.SoftDeletePost(ctx, userID, postID)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] DeletePostSoft: failed to soft delete post %v: %w", postID, err)
	}

	s.invalidateTopPostsOfDay(ctx)
	return softDeletedPost, nil
}

func (s *PostService) DeletePost(ctx context.Context, userID string, postID string) (domain.Post, error) {
	postImages, imgErr := s.imageRepo.GetPostImagesByPost(ctx, userID, postID)
	if imgErr != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] DeletePost: failed to fetch post images: %w", imgErr)
	}

	for _, image := range postImages {
		if err := s.imageStorage.DeletePicture(ctx, image.ImageURL); err != nil {
			return domain.Post{}, fmt.Errorf("[post_service.go] DeletePost: failed to delete image from storage: %w", err)
		}
	}

	for _, image := range postImages {
		if _, err := s.imageRepo.DeletePostImage(ctx, userID, postID, image.ID); err != nil {
			return domain.Post{}, fmt.Errorf("[post_service.go] DeletePost: failed to delete image from db: %w", err)
		}
	}

	deletedPost, err := s.repo.PermanentlyDeletePost(ctx, userID, postID)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_service.go] DeletePost: failed to permanently delete post: %w", err)
	}

	s.invalidateTopPostsOfDay(ctx)
	return deletedPost, nil
}
