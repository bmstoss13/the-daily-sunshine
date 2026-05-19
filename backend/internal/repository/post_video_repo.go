package repository

import (
	"context"
	"fmt"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPostVideoRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPostVideoRepository(db *pgxpool.Pool) *PostgresPostVideoRepository {
	return &PostgresPostVideoRepository{
		db: db,
	}
}

func (r *PostgresPostVideoRepository) GetPostVideoByID(ctx context.Context, userID string, videoID string) (domain.PostVideo, error) {
	pgVideoID, err := StringToPgUUID(videoID)
	if err != nil {
		return domain.PostVideo{}, fmt.Errorf("[post_video_repo.go] GetPostVideoByID: failed to convert string to pg uuid: %w", err)
	}

	var fetchedPostVideo domain.PostVideo

	err = WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcPostVideo, err := q.SelectPostVideoByID(ctx, pgVideoID)
		if err != nil {
			return fmt.Errorf("[post_video_repo.go] GetPostVideoByID: An error occurred while retrieving post video from id: %w", err)
		}

		fetchedPostVideo = ConvertPostVideo(sqlcPostVideo)
		return nil
	})

	return fetchedPostVideo, err
}

func (r *PostgresPostVideoRepository) GetPostVideoByPost(ctx context.Context, userID string, postID string) (domain.PostVideo, error) {
	pgPostID, err := StringToPgUUID(postID)
	if err != nil {
		return domain.PostVideo{}, fmt.Errorf("[post_video_repo.go] GetPostVideoByPost: failed to convert string to pg uuid: %w", err)
	}

	var fetchedPostVideo domain.PostVideo
	err = WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcPostVideo, err := q.SelectPostVideoByPostID(ctx, pgPostID)
		if err != nil {
			return fmt.Errorf("[post_video_repo.go] GetPostByPost: An error occurred while retrieving post video from post: %w", err)
		}

		fetchedPostVideo = ConvertPostVideo(sqlcPostVideo)
		return nil
	})

	return fetchedPostVideo, err
}

func (r *PostgresPostVideoRepository) GetListOfPostVideos(ctx context.Context, userID string, limit int32, offset int32) ([]domain.PostVideo, error) {
	var postVideoList []domain.PostVideo
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcPostVideoList, err := q.ListVideoFeed(ctx, ListVideoFeedParams{limit, offset})
		if err != nil {
			return fmt.Errorf("[post_video_repo.go] GetListOfPostVideos: An error occurred while retrieving list of videos with limit %v and offset %v: %w", limit, offset, err)
		}

		postVideoList = make([]domain.PostVideo, len(sqlcPostVideoList))

		for i, pv := range sqlcPostVideoList {
			postVideoList[i] = ConvertListVideoFeedRow(pv)
		}
		return nil
	})

	return postVideoList, err
}

func (r *PostgresPostVideoRepository) CreatePostVideo(ctx context.Context, userID string, newVideo domain.PostVideo) (domain.PostVideo, error) {
	var createdPostVideo domain.PostVideo
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgPostID, postErr := StringToPgUUID(newVideo.PostID)
		if postErr != nil {
			return fmt.Errorf("[post_video_repo.go] CreatePostVideo: An error occurred while converting post id %v to pg id: %w", newVideo.PostID, postErr)
		}

		params := CreatePostVideoParams{
			PostID:         pgPostID,
			YoutubeVideoID: newVideo.YouTubeVideoID,
			VideoMetadata:  []byte(newVideo.VideoMetadata),
		}

		sqlcPostVideo, err := q.CreatePostVideo(ctx, params)
		if err != nil {
			return fmt.Errorf("[post_video_repo.go] CreatePostVideo: failed to insert post video: %w", err)
		}

		createdPostVideo = ConvertPostVideo(sqlcPostVideo)
		return nil
	})

	return createdPostVideo, err
}

func (r *PostgresPostVideoRepository) UpdatePostVideo(ctx context.Context, userID string, videoToUpdate domain.PostVideo) (domain.PostVideo, error) {
	var updatedPostVideo domain.PostVideo
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgVideoID, vidErr := StringToPgUUID(videoToUpdate.ID)
		if vidErr != nil {
			return fmt.Errorf("[post_video_repo.go] UpdatePostVideo: An error occurred while converting post video id %v to pg id: %w", videoToUpdate.ID, vidErr)
		}

		pgPostID, postErr := StringToPgUUID(videoToUpdate.PostID)
		if postErr != nil {
			return fmt.Errorf("[post_video_repo.go] UpdatePostVideo: An error occurred while converting post id %v to pg id: %w", videoToUpdate.PostID, vidErr)
		}

		pgUserID, userErr := StringToPgUUID(userID)
		if userErr != nil {
			return fmt.Errorf("[post_video_repo.go] UpdatePostVideo: An error occurred while converting user id %v to pg id: %w", userID, userErr)
		}

		params := UpdatePostVideoParams{
			ID:             pgVideoID,
			PostID:         pgPostID,
			YoutubeVideoID: videoToUpdate.YouTubeVideoID,
			VideoMetadata:  []byte(videoToUpdate.VideoMetadata),
			PublisherID:    pgUserID,
		}

		sqlcPostVideo, err := q.UpdatePostVideo(ctx, params)
		if err != nil {
			return fmt.Errorf("[post_video_repo.go] UpdatePostVideo: failed to update post video: %w", err)
		}

		updatedPostVideo = ConvertPostVideo(sqlcPostVideo)
		return nil
	})

	return updatedPostVideo, err
}

func (r *PostgresPostVideoRepository) DeletePostVideo(ctx context.Context, userID string, videoID string, postID string) (domain.PostVideo, error) {
	var deletedPostVideo domain.PostVideo
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgVideoID, vidErr := StringToPgUUID(videoID)
		if vidErr != nil {
			return fmt.Errorf("[post_video_repo.go] DeletePostVideo: An error occurred while converting post video id %v to pg id: %w", videoID, vidErr)
		}

		pgPostID, postErr := StringToPgUUID(postID)
		if postErr != nil {
			return fmt.Errorf("[post_video_repo.go] DeletePostVideo: An error occurred while converting post id %v to pg id: %w", postID, postErr)
		}

		pgUserID, userErr := StringToPgUUID(userID)
		if userErr != nil {
			return fmt.Errorf("[post_video_repo.go] DeletePostVideo: An error occurred while converting user id %v to pg id: %w", userID, userErr)
		}

		sqlcPostVideo, err := q.DeletePostVideo(ctx, DeletePostVideoParams{pgVideoID, pgPostID, pgUserID})
		if err != nil {
			return fmt.Errorf("[post_video_repo.go] DeletePostVideo: failed to delete post video: %w", err)
		}

		deletedPostVideo = ConvertPostVideo(sqlcPostVideo)
		return nil
	})

	return deletedPostVideo, err
}
