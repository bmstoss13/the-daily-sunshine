package repository

import (
	"context"
	"fmt"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPostImageRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPostImageRepository(db *pgxpool.Pool) *PostgresPostImageRepository {
	return &PostgresPostImageRepository{
		db: db,
	}
}

func (r *PostgresPostImageRepository) GetPostImageByID(ctx context.Context, userID string, imageID string) (domain.PostImage, error) {
	pgImageID, err := StringToPgUUID(imageID)
	if err != nil {
		return domain.PostImage{}, fmt.Errorf("[post_image_repo.go] GetPostImageByID: failed to convert string to pg uuid: %w", err)
	}

	var fetchedPostImage domain.PostImage

	err = WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcPostImage, err := q.SelectPostImageByID(ctx, pgImageID)
		if err != nil {
			return fmt.Errorf("[post_image_repo.go] GetPostImageByID: An error occurred while retrieving post image from id: %w", err)
		}

		fetchedPostImage = ConvertPostImage(sqlcPostImage)
		return nil
	})

	return fetchedPostImage, err
}

func (r *PostgresPostImageRepository) GetCoverImageOfPost(ctx context.Context, userID string, postID string) (domain.PostImage, error) {
	pgImageID, err := StringToPgUUID(postID)
	if err != nil {
		return domain.PostImage{}, fmt.Errorf("[post_image_repo.go] GetCoverImageOfPost: failed to convert string to pg uuid: %w", err)
	}

	var fetchedImagePost domain.PostImage

	err = WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcImagePost, err := q.SelectCoverImage(ctx, pgImageID)
		if err != nil {
			return fmt.Errorf("[post_image_repo.go] GetCoverImageOfPost: An error occurred while retrieving cover image from post: %w", err)
		}

		fetchedImagePost = ConvertPostImage(sqlcImagePost)
		return nil
	})

	return fetchedImagePost, err
}

func (r *PostgresPostImageRepository) GetPostImagesByPost(ctx context.Context, userID string, postID string) ([]domain.PostImage, error) {
	pgPostID, err := StringToPgUUID(postID)
	if err != nil {
		return []domain.PostImage{}, fmt.Errorf("[post_image_repo.go] GetPostImagesByPost: failed to convert string to pd uuid: %w", err)
	}

	var postImageList []domain.PostImage
	err = WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcPostImageList, err := q.SelectPostImagesByPost(ctx, pgPostID)
		if err != nil {
			return fmt.Errorf("[post_image_repo.go] GetPostImagesByPost: failed to retrieve post images by post id: %w", err)
		}

		postImageList = make([]domain.PostImage, len(sqlcPostImageList))

		for i, pi := range sqlcPostImageList {
			postImageList[i] = ConvertPostImage(pi)
		}
		return nil
	})

	return postImageList, err
}

func (r *PostgresPostImageRepository) CreatePostImage(ctx context.Context, userID string, newPostImage domain.PostImage) (domain.PostImage, error) {
	var createdPostImage domain.PostImage
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgID, err := StringToPgUUID(newPostImage.ID)
		if err != nil {
			return fmt.Errorf("[post_image_repo.go] CreatePostImage: An error occurred while converting post image id %v to pg id: %w", newPostImage.ID, err)
		}

		postPgID, postErr := StringToPgUUID(newPostImage.PostID)
		if postErr != nil {
			return fmt.Errorf("[post_image_repo.go] CreatePostImage: An error occurred while converting post id %v to pg id: %w", newPostImage.PostID, err)
		}

		params := CreatePostImageParams{
			ID:               pgID,
			PostID:           postPgID,
			ImageUrl:         newPostImage.ImageURL,
			ImageDescription: &newPostImage.ImageDescription,
			AltText:          &newPostImage.AltText,
			IsCoverImage:     newPostImage.IsCoverImage,
		}

		sqlcPostImage, err := q.CreatePostImage(ctx, params)
		if err != nil {
			return fmt.Errorf("[post_image_repo.go] CreatePostImage: failed to insert post image: %w", err)
		}

		createdPostImage = ConvertPostImage(sqlcPostImage)
		return nil
	})

	return createdPostImage, err
}

func (r *PostgresPostImageRepository) UpdatePostImage(ctx context.Context, userID string, postImageWithUpdates domain.PostImage) (domain.PostImage, error) {
	var updatedPostImage domain.PostImage
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgID, err := StringToPgUUID(postImageWithUpdates.ID)
		if err != nil {
			return fmt.Errorf("[post_image_repo.go] UpdatePostImage: An error occurred while converting post image id %v to pg id: %w", postImageWithUpdates.ID, err)
		}

		postPgID, postErr := StringToPgUUID(postImageWithUpdates.PostID)
		if postErr != nil {
			return fmt.Errorf("[post_image_repo.go] UpdatePostImage: An error occurred while converting post id %v to pg id: %w", postImageWithUpdates.PostID, err)
		}

		params := UpdatePostImageParams{
			ID:               pgID,
			PostID:           postPgID,
			ImageUrl:         postImageWithUpdates.ImageURL,
			ImageDescription: &postImageWithUpdates.ImageDescription,
			AltText:          &postImageWithUpdates.AltText,
			IsCoverImage:     postImageWithUpdates.IsCoverImage,
		}

		sqlcPostImage, err := q.UpdatePostImage(ctx, params)
		if err != nil {
			return fmt.Errorf("[post_image_repo.go] UpdatePostImage: failed to update post image: %w", err)
		}

		updatedPostImage = ConvertPostImage(sqlcPostImage)
		return nil
	})

	return updatedPostImage, err
}

func (r *PostgresPostImageRepository) DeletePostImage(ctx context.Context, userID string, postID string, postImageID string) (domain.PostImage, error) {
	var deletedPostImage domain.PostImage
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgPostImageID, err := StringToPgUUID(postImageID)
		if err != nil {
			return fmt.Errorf("[post_image_repo.go] DeletePostImage: invalid post image id format: %w", err)
		}

		pgPostID, err := StringToPgUUID(postID)
		if err != nil {
			return fmt.Errorf("[post_image_repo.go] DeletePostImage: invalid post id format: %w", err)
		}

		sqlcPostImage, err := q.DeletePostImage(ctx, DeletePostImageParams{pgPostImageID, pgPostID})
		if err != nil {
			return fmt.Errorf("[post_image_repo.go] DeletePostImage: failed to delete post image: %w", err)
		}

		deletedPostImage = ConvertPostImage(sqlcPostImage)
		return nil
	})

	return deletedPostImage, err
}
