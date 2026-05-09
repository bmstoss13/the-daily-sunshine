package repository

import (
	"context"
	"fmt"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProfileRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProfileRepository(db *pgxpool.Pool) *PostgresProfileRepository {
	return &PostgresProfileRepository{
		db: db,
	}
}

/*
* Fetch profile from the id
**/
func (r *PostgresProfileRepository) GetProfileFromID(ctx context.Context, userID string, profileID string) (domain.Profile, error) {

	pgProfileID, err := StringToPgUUID(profileID)
	if err != nil {
		return domain.Profile{}, err
	}

	var fetchedProfile domain.Profile

	err = WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
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

// Function for getting profile from username
func (r *PostgresProfileRepository) GetProfileFromUsername(ctx context.Context, userID string, username string) (domain.Profile, error) {
	var fetchedProfile domain.Profile

	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
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

func (r *PostgresProfileRepository) CheckProfileUsernameExists(ctx context.Context, userID string, username string) (bool, error) {
	var doesExist bool
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		exists, checkErr := q.CheckUsernameExists(ctx, username)
		if checkErr != nil {
			return fmt.Errorf("[profile_repo.go] CheckProfileUsernameExists: An error occurred while checking if username, %v, is in use: %w", username, checkErr)
		}

		doesExist = exists
		return nil
	})

	return doesExist, err
}

func (r *PostgresProfileRepository) GetListOfProfiles(ctx context.Context, userID string, limit int32, offset int32) ([]domain.Profile, error) {
	var profileList []domain.Profile
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		sqlcProfileList, listErr := q.ListProfiles(ctx, ListProfilesParams{limit, offset})
		if listErr != nil {
			return fmt.Errorf("[profile_repo.go] GetListOfProfiles: An error occurred while retrieving list of profiles with limit %v and offset %v: %w", limit, offset, listErr)
		}

		profileList = make([]domain.Profile, len(sqlcProfileList))

		for i, p := range sqlcProfileList {
			profileList[i] = ConvertProfile(p)
		}
		return nil
	})

	return profileList, err
}

// takes in user profile with temp id
func (r *PostgresProfileRepository) CreateProfile(ctx context.Context, userID string, newProfile domain.Profile) (domain.Profile, error) {
	var createdProfile domain.Profile
	err := WithRLS(ctx, r.db, userID, func(tx pgx.Tx) error {
		q := New(tx)

		pgID, err := StringToPgUUID(newProfile.ID)
		if err != nil {
			return fmt.Errorf("[profile_repo.go] CreateProfile: invalid profile ID format: %w", err)
		}

		params := CreateProfileParams{
			ID:        pgID,
			FirstName: newProfile.FirstName,
			LastName:  newProfile.LastName,
			Username:  newProfile.Username,
			Role:      string(newProfile.Role),
		}

		if newProfile.ProfileImageUrl != "" {
			params.ProfileImageUrl = &newProfile.ProfileImageUrl
		}

		if newProfile.Bio != nil {
			params.ProfileBio = []byte(*newProfile.Bio)
		}

		sqlcProfile, err := q.CreateProfile(ctx, params)
		if err != nil {
			return fmt.Errorf("[profile_repo.go] CreateProfile: failed to insert profile: %w", err)
		}

		createdProfile = ConvertProfile(sqlcProfile)
		return nil
	})

	return createdProfile, err
}
