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
