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

	profile, err := repository.GetProfileFromUsername(ctx, requestingUserID, targetUsername)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] failed to get profile from username %s with id %s: %w", targetUsername, requestingUserID, err)
	}

	return profile, nil
}
