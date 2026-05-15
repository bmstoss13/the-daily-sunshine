package repository

import (
	"context"
	"fmt"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresSearchRepository struct {
	db *pgxpool.Pool
}

func NewPostgresSearchRepository(db *pgxpool.Pool) *PostgresSearchRepository {
	return &PostgresSearchRepository{
		db: db,
	}
}

func (r *PostgresSearchRepository) SearchProfiles(ctx context.Context, searchQuery string, limit int32, offset int32) ([]domain.Profile, error) {
	q := New(r.db)

	sqlcProfiles, err := q.SearchProfiles(ctx, SearchProfilesParams{searchQuery, limit, offset})
	if err != nil {
		return nil, fmt.Errorf("[search_repo.go] Failed to search profiles with search query %v: %w", searchQuery, err)
	}
	var profiles []domain.Profile
	for _, p := range sqlcProfiles {

		var profileImageUrl string
		if p.ProfileImageUrl != nil {
			profileImageUrl = *p.ProfileImageUrl
		} else {
			profileImageUrl = ""
		}

		profiles = append(profiles, domain.Profile{
			ID:              PgUUIDToString(p.ID),
			Username:        p.Username,
			FirstName:       p.FirstName,
			LastName:        p.LastName,
			SubscriberTier:  domain.SubscriberTier(p.SubscriberTier),
			ProfileImageUrl: profileImageUrl,
		})
	}

	return profiles, nil
}

func (r *PostgresSearchRepository) SearchPosts(ctx context.Context, searchQuery string, limit int32, offset int32) ([]domain.Post, error) {
	q := New(r.db)

	sqlcPosts, err := q.SearchPosts(ctx, SearchPostsParams{searchQuery, limit, offset})
	if err != nil {
		return nil, fmt.Errorf("[search_repo.go] Failed to search posts with search query %v: %w", searchQuery, err)
	}

	posts := make([]domain.Post, len(sqlcPosts))
	for i, p := range sqlcPosts {
		posts[i] = ConvertSearchPostRow(p)
	}

	return posts, nil
}
