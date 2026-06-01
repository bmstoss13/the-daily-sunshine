package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPostRepository struct {
	db *pgxpool.Pool
}

func NewPostgresPostRepository(db *pgxpool.Pool) *PostgresPostRepository {
	return &PostgresPostRepository{
		db: db,
	}
}

func (r *PostgresPostRepository) GetPostByID(ctx context.Context, userID string, postID string) (domain.Post, error) {
	pgPostID, err := StringToPgUUID(postID)
	if err != nil {
		return domain.Post{}, fmt.Errorf("[post_repo.go] GetPostByID: failed to convert string to pg uuid: %w", err)
	}

	var fetchedPost domain.Post

	err = WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcPost, err := q.SelectPostByID(ctx, pgPostID)
		if err != nil {
			return fmt.Errorf("[post_repo.go] An error occurred while retrieving post from id: %w", err)
		}

		fetchedPost = ConvertRichPost(sqlcPost)
		return nil
	})

	return fetchedPost, err
}

func (r *PostgresPostRepository) GetPostBySlug(ctx context.Context, userID string, slug string) (domain.Post, error) {
	var fetchedPost domain.Post

	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcPost, err := q.SelectPostBySlug(ctx, slug)
		if err != nil {
			return fmt.Errorf("[post_repo.go] An error occurred while retrieving post from id: %w", err)
		}

		fetchedPost = ConvertRichPostBySlug(sqlcPost)
		return nil
	})

	return fetchedPost, err
}

func (r *PostgresPostRepository) GetListOfPosts(ctx context.Context, userID string, limit int32, offset int32) ([]domain.Post, error) {
	var postList []domain.Post
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcPostList, listErr := q.ListPosts(ctx, ListPostsParams{limit, offset})
		if listErr != nil {
			return fmt.Errorf("[post_repo.go] GetListOfPosts: An error occurred while retrievign list of posts with limit %v and offset %v: %w", limit, offset, listErr)
		}

		postList = make([]domain.Post, len(sqlcPostList))

		for i, p := range sqlcPostList {
			postList[i] = ConvertRichPostForList(p)
		}
		return nil
	})

	return postList, err
}

func (r *PostgresPostRepository) GetPostsForToday(ctx context.Context, userID string, startOfDay time.Time, endOfDay time.Time) ([]domain.Post, error) {
	var postList []domain.Post
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgStartOfDay := TimeToPgTimestamp(startOfDay)
		if !pgStartOfDay.Valid {
			return fmt.Errorf("[post_repo.go] GetPostsForToday: failed to convert start of day to pgtype.timestamptz. Start of day date: %v", pgStartOfDay)
		}

		pgEndOfDay := TimeToPgTimestamp(endOfDay)
		if !pgStartOfDay.Valid {
			return fmt.Errorf("[post_repo.go] GetPostsForToday: failed to convert end of day to pgtype.timestamptz. End of day date: %v", pgEndOfDay)
		}

		sqlcPostList, listErr := q.SelectPostsOfTheDay(ctx)
		if listErr != nil {
			return fmt.Errorf("[post_repo.go] GetPostsForToday: An error occurred while retrieving list of post for today: %w", listErr)
		}

		postList = make([]domain.Post, len(sqlcPostList))

		for i, p := range sqlcPostList {
			postList[i] = ConvertRichPostForDay(p)
		}
		return nil
	})

	return postList, err
}

func (r *PostgresPostRepository) CreatePost(ctx context.Context, userID string, newPost domain.Post) (domain.Post, error) {
	var createdPost domain.Post
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		publisherPgID, pubErr := StringToPgUUID(userID)
		if pubErr != nil {
			return fmt.Errorf("[post_repo.go] CreatePost: An error occurred while converting user id %v to pg id: %w", userID, pubErr)
		}

		// safely handle nullable content
		var contentBytes []byte
		if newPost.Content != nil {
			contentBytes = []byte(*newPost.Content)
		}

		params := CreatePostParams{
			PublisherID: publisherPgID,
			Title:       newPost.Title,
			Subtitle:    newPost.Subtitle,
			Slug:        newPost.Slug,
			Status:      PostStatus(newPost.Status),
			PostContent: contentBytes,
			NumRays:     int32(newPost.NumRays),
			NumComments: int32(newPost.NumComments),
		}

		sqlcPost, err := q.CreatePost(ctx, params)
		if err != nil {
			return fmt.Errorf("[post_repo.go] CreatePost: failed to insert post: %w", err)
		}

		createdPost = ConvertPost(sqlcPost)
		return nil
	})

	return createdPost, err
}

func (r *PostgresPostRepository) UpdatePost(ctx context.Context, userID string, postWithUpdates domain.Post) (domain.Post, error) {
	var updatedPost domain.Post
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgID, err := StringToPgUUID(postWithUpdates.ID)
		if err != nil {
			return fmt.Errorf("[post_repo.go] UpdatePost: invalid post ID format: %w", err)
		}

		publisherPgID, err := StringToPgUUID(userID)
		if err != nil {
			return fmt.Errorf("[post_repo.go] UpdatePost: invalid user ID format: %w", err)
		}

		// safely handle nullable content
		var contentBytes []byte
		if postWithUpdates.Content != nil {
			contentBytes = []byte(*postWithUpdates.Content)
		}

		params := UpdatePostParams{
			ID:          pgID,
			PublisherID: publisherPgID,
			Title:       postWithUpdates.Title,
			Subtitle:    postWithUpdates.Subtitle,
			Slug:        postWithUpdates.Slug,
			Status:      PostStatus(postWithUpdates.Status),
			PostContent: contentBytes,
			NumRays:     int32(postWithUpdates.NumRays),
			NumComments: int32(postWithUpdates.NumComments),
		}

		sqlcPost, err := q.UpdatePost(ctx, params)
		if err != nil {
			return fmt.Errorf("[post_repo.go] UpdatePost: failed to update post: %w", err)
		}

		updatedPost = ConvertPost(sqlcPost)
		return nil
	})

	return updatedPost, err
}

func (r *PostgresPostRepository) SoftDeletePost(ctx context.Context, userID string, postID string) (domain.Post, error) {
	var deletedPost domain.Post
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgID, err := StringToPgUUID(userID)
		if err != nil {
			return fmt.Errorf("[post_repo.go] SoftDeletePost: invalid user id format: %w", err)
		}

		pgPostID, err := StringToPgUUID(postID)
		if err != nil {
			return fmt.Errorf("[post_repo.go] SoftDeletePost: invalid post id format: %w", err)
		}

		sqlcPost, err := q.SoftDeletePost(ctx, SoftDeletePostParams{pgPostID, pgID})
		if err != nil {
			return fmt.Errorf("[post_repo.go] SoftDeletePost: failed to delete post: %w", err)
		}

		deletedPost = ConvertPost(sqlcPost)
		return nil
	})

	return deletedPost, err
}

func (r *PostgresPostRepository) PermanentlyDeletePost(ctx context.Context, userID string, postID string) (domain.Post, error) {
	var deletedPost domain.Post
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgID, err := StringToPgUUID(userID)
		if err != nil {
			return fmt.Errorf("[post_repo.go] PermanentlyDeletePost: invalid user id format: %w", err)
		}

		pgPostID, err := StringToPgUUID(postID)
		if err != nil {
			return fmt.Errorf("[post_repo.go] PermanentlyDeletePost: invalid post id format: %w", err)
		}

		sqlcPost, err := q.PermanentlyDeletePost(ctx, PermanentlyDeletePostParams{pgPostID, pgID})
		if err != nil {
			return fmt.Errorf("[post_repo.go] PermanentlyDeletePost: failed to delete post: %w", err)
		}

		deletedPost = ConvertPost(sqlcPost)
		return nil
	})

	return deletedPost, err
}

func (r *PostgresPostRepository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	q := New(r.db)

	exists, err := q.CheckSlugExists(ctx, slug)
	if err != nil {
		return false, fmt.Errorf("[post_repo.go] CheckIfSlugExists: failed to check slug: %w", err)
	}

	return exists, nil
}

func (r *PostgresPostRepository) CheckSlugExistsOtherPosts(ctx context.Context, slug string, postID string) (bool, error) {
	q := New(r.db)

	pgID, idErr := StringToPgUUID(postID)
	if idErr != nil {
		return false, fmt.Errorf("[post_repo.go] CheckSlugExistsOtherPosts: failed to convert post id: %w", idErr)
	}

	exists, err := q.CheckSlugExistsOtherPosts(ctx, CheckSlugExistsOtherPostsParams{slug, pgID})
	if err != nil {
		return false, fmt.Errorf("[post_repo.go] CheckSlugExistsOtherPosts: failed to check slug: %w", err)
	}

	return exists, nil
}
