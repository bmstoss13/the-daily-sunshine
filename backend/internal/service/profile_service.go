package service

import (
	"context"
	"fmt"

	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/repository"
)

func FetchProfileByID(ctx context.Context, requestingUserID string, targetProfileID string) (domain.Profile, error) {
	if targetProfileID == "" {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] target profile ID is required.")
	}

	// Add business logic in the future (i.e. blocked users)

	profile, err := repository.GetProfileFromID(ctx, requestingUserID, targetProfileID)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] failed to get profile from id %s with id %s: %w", targetProfileID, requestingUserID, err)
	}

	return profile, nil
}

func FetchProfileByUsername(ctx context.Context, requestingUserID string, targetUsername string) (domain.Profile, error) {
	if targetUsername == "" {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] target username is required.")
	}

	// Add business logic in the future (i.e. blocked users)

	profile, err := repository.GetProfileFromUsername(ctx, requestingUserID, targetUsername)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] failed to get profile from username %s with id %s: %w", targetUsername, requestingUserID, err)
	}

	return profile, nil
}

func IsUsernameTaken(ctx context.Context, requestingUserID string, usernameToCheck string) (bool, error) {
	if usernameToCheck == "" {
		return true, fmt.Errorf("[profile_service.go] IsUsernameTaken: username is required")
	}

	doesExist, err := repository.CheckProfileUsernameExists(ctx, requestingUserID, usernameToCheck)
	if err != nil {
		return true, fmt.Errorf("[profile_service.go] IsUsernameTaken: failed to get check if username %s exists from user with id %s: %w", usernameToCheck, requestingUserID, err)
	}

	return doesExist, nil
}

func FetchListOfProfiles(ctx context.Context, requestingUserID string, limit int32, offset int32) ([]domain.Profile, error) {
	//guardrail to prevent pulling over 100 profiles at a time
	if limit > 100 {
		limit = 100
	}
	profileList, err := repository.GetListOfProfiles(ctx, requestingUserID, limit, offset)
	if err != nil {
		return []domain.Profile{}, fmt.Errorf("[profile_service.go] FetchListOfProfiles: failed to get list of profiles with limit %v and offset %v with id %s: %w", limit, offset, requestingUserID, err)
	}

	return profileList, nil
}
