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

		fetchedProfile = domain.Profile{
			ID:        PgUUIDToString(sqlcProfile.ID),
			FirstName: sqlcProfile.FirstName,
			LastName:  sqlcProfile.LastName,
			Username:  sqlcProfile.Username,
			NumRays:   int(sqlcProfile.NumRaysReceived),
			Role:      domain.MembershipRole(sqlcProfile.Role),

			CreatedAt: sqlcProfile.CreatedAt.Time,
			UpdatedAt: sqlcProfile.UpdatedAt.Time,
		}

		if sqlcProfile.ProfileImageUrl != nil {
			fetchedProfile.ProfileImage = *sqlcProfile.ProfileImageUrl
		} else {
			fetchedProfile.ProfileImage = "" // Or a default URL string
		}

		if sqlcProfile.ProfileBio != nil {
			bioStr := string(sqlcProfile.ProfileBio)
			fetchedProfile.Bio = &bioStr
		}

		if sqlcProfile.DeletedAt.Valid {
			fetchedProfile.DeletedAt = &sqlcProfile.DeletedAt.Time
		}

		return nil
	})

	return fetchedProfile, err
}
