package repository

import (
	"context"
	"fmt"

	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/jackc/pgx/v5"
)

/*
* Fetch profile from the id
**/
func GetProfileFromID(ctx context.Context, userID string, profileID string) (domain.Profile, error) {

	pgProfileID, err := StringToPgUUID(profileID)
	if err != nil {
		return domain.Profile{}, err
	}

	var fetchedProfile domain.Profile

	err = WithRLS(ctx, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcProfile, err := q.GetProfileByID(ctx, pgProfileID)
		if err != nil {
			return fmt.Errorf("[profile_repo.go] An error occurred while retrieving profile from id: %w", err)
		}

		fetchedProfile = ConvertProfile(sqlcProfile)

		return nil
	})

	return fetchedProfile, err
}

func GetProfileFromUsername(ctx context.Context, userID string, username string) (domain.Profile, error) {
	var fetchedProfile domain.Profile

	err := WithRLS(ctx, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcProfile, err := q.GetProfileByUsername(ctx, username)
		if err != nil {
			return fmt.Errorf("[profile_repo.go] An error occurred while retrieving profile from username: %w", err)
		}

		fetchedProfile = ConvertProfile(sqlcProfile)

		return nil
	})

	return fetchedProfile, err
}
