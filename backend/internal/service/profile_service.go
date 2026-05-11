package service

import (
	"context"
	"fmt"
	"log"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
)

type ProfileRepository interface {
	GetProfileFromID(ctx context.Context, userID string, targetProfileID string) (domain.Profile, error)
	GetProfileFromUsername(ctx context.Context, userID string, username string) (domain.Profile, error)
	CheckProfileUsernameExists(ctx context.Context, userID string, username string) (bool, error)
	GetListOfProfiles(ctx context.Context, userID string, limit int32, offset int32) ([]domain.Profile, error)
	CreateProfile(ctx context.Context, userID string, newProfile domain.Profile) (domain.Profile, error)
	UpdateProfile(ctx context.Context, userID string, profileWithUpdates domain.Profile) (domain.Profile, error)
	SoftDeleteProfile(ctx context.Context, userID string) (domain.Profile, error)
	PermanentlyDeleteProfile(ctx context.Context, userID string) (domain.Profile, error)
}

type ImageStorage interface {
	UploadProfilePicture(ctx context.Context, fileBytes []byte, fileName string) (string, error)
	DeleteProfilePicture(ctx context.Context, fullImageURL string) error
}

type ProfileService struct {
	repo         ProfileRepository
	imageStorage ImageStorage
}

func NewProfileService(repo ProfileRepository, imageStorage ImageStorage) *ProfileService {
	return &ProfileService{
		repo:         repo,
		imageStorage: imageStorage,
	}
}

func (s *ProfileService) FetchProfileByID(ctx context.Context, requestingUserID string, targetProfileID string) (domain.Profile, error) {
	if targetProfileID == "" {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] target profile ID is required.")
	}

	// Add business logic in the future (i.e. blocked users)

	profile, err := s.repo.GetProfileFromID(ctx, requestingUserID, targetProfileID)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] failed to get profile from id %s with id %s: %w", targetProfileID, requestingUserID, err)
	}

	return profile, nil
}

func (s *ProfileService) FetchProfileByUsername(ctx context.Context, requestingUserID string, targetUsername string) (domain.Profile, error) {
	if targetUsername == "" {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] target username is required.")
	}

	// Add business logic in the future (i.e. blocked users)

	profile, err := s.repo.GetProfileFromUsername(ctx, requestingUserID, targetUsername)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] failed to get profile from username %s with id %s: %w", targetUsername, requestingUserID, err)
	}

	return profile, nil
}

func (s *ProfileService) IsUsernameTaken(ctx context.Context, requestingUserID string, usernameToCheck string) (bool, error) {
	if usernameToCheck == "" {
		return true, fmt.Errorf("[profile_service.go] IsUsernameTaken: username is required")
	}

	doesExist, err := s.repo.CheckProfileUsernameExists(ctx, requestingUserID, usernameToCheck)
	if err != nil {
		return true, fmt.Errorf("[profile_service.go] IsUsernameTaken: failed to check if username %s exists from user with id %s: %w", usernameToCheck, requestingUserID, err)
	}

	return doesExist, nil
}

func (s *ProfileService) FetchListOfProfiles(ctx context.Context, requestingUserID string, limit int32, offset int32) ([]domain.Profile, error) {
	//guardrail to prevent pulling over 100 profiles at a time
	if limit > 100 {
		limit = 100
	}
	profileList, err := s.repo.GetListOfProfiles(ctx, requestingUserID, limit, offset)
	if err != nil {
		return []domain.Profile{}, fmt.Errorf("[profile_service.go] FetchListOfProfiles: failed to get list of profiles with limit %v and offset %v with id %s: %w", limit, offset, requestingUserID, err)
	}

	return profileList, nil
}

func (s *ProfileService) CreateUserProfile(ctx context.Context, requestingUserID string, newProfile domain.Profile, imageBytes []byte, fileExtension string) (domain.Profile, error) {
	if newProfile.Username == "" {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] CreateUserProfile: username is required")
	}

	usernameExists, err := s.repo.CheckProfileUsernameExists(ctx, requestingUserID, newProfile.Username)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] CreateUserProfile: failed to check if username %s exists from user with id %s: %w", newProfile.Username, requestingUserID, err)
	}

	if usernameExists {
		return domain.Profile{}, fmt.Errorf("the username '%s' is already taken", newProfile.Username)
	}

	if len(imageBytes) > 0 {
		// unique, organized file name: "profiles/user-uuid/avatar.jpg"
		fileName := fmt.Sprintf("profiles/%s/avatar%s", requestingUserID, fileExtension)

		publicURL, err := s.imageStorage.UploadProfilePicture(ctx, imageBytes, fileName)
		if err != nil {
			// If the image fails to upload, abort profile creation
			return domain.Profile{}, fmt.Errorf("[profile_service.go] CreateUserProfile: failed to upload profile picture: %w", err)
		}

		newProfile.ProfileImageUrl = publicURL
	}

	createdProfile, err := s.repo.CreateProfile(ctx, requestingUserID, newProfile)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] CreateUserProfile: failed to insert into database: %w", err)
	}

	return createdProfile, nil
}

func (s *ProfileService) UpdateUserProfile(ctx context.Context, requestingUserID string, profileWithUpdates domain.Profile, imageBytes []byte, fileExtension string) (domain.Profile, error) {
	if profileWithUpdates.Username == "" {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] UpdateUserProfile: username is required")
	}

	currentProfile, err := s.repo.GetProfileFromID(ctx, requestingUserID, requestingUserID)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] UpdateUserProfile: failed to fetch current profile: %w", err)
	}

	if currentProfile.Username != profileWithUpdates.Username {
		usernameExists, err := s.repo.CheckProfileUsernameExists(ctx, requestingUserID, profileWithUpdates.Username)
		if err != nil {
			return domain.Profile{}, fmt.Errorf("[profile_service.go] UpdateUserProfile: failed to check if username %s exists from user with id %s: %w", profileWithUpdates.Username, requestingUserID, err)
		}
		if usernameExists {
			return domain.Profile{}, fmt.Errorf("the username '%s' is already taken", profileWithUpdates.Username)
		}
	}

	if len(imageBytes) > 0 {
		// unique, organized file name: "profiles/user-uuid/avatar.jpg"
		fileName := fmt.Sprintf("profiles/%s/avatar%s", requestingUserID, fileExtension)

		publicURL, err := s.imageStorage.UploadProfilePicture(ctx, imageBytes, fileName)
		if err != nil {
			// If the image fails to upload, abort profile creation
			return domain.Profile{}, fmt.Errorf("[profile_service.go] UpdateUserProfile: failed to upload profile picture: %w", err)
		}

		profileWithUpdates.ProfileImageUrl = publicURL
	} else {
		profileWithUpdates.ProfileImageUrl = currentProfile.ProfileImageUrl
	}

	updatedProfile, err := s.repo.UpdateProfile(ctx, requestingUserID, profileWithUpdates)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] UpdateUserProfile: failed to update profile in database: %w", err)
	}

	return updatedProfile, nil
}

func (s *ProfileService) SoftDeleteUserProfile(ctx context.Context, requestingUserID string) (domain.Profile, error) {
	softDeletedProfile, err := s.repo.SoftDeleteProfile(ctx, requestingUserID)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] SoftDeleteUserProfile: failed to soft delete profile: %w", err)
	}

	return softDeletedProfile, nil
}

func (s *ProfileService) PermanentlyDeleteUserProfile(ctx context.Context, requestingUserID string) (domain.Profile, error) {
	profileToDelete, err := s.repo.GetProfileFromID(ctx, requestingUserID, requestingUserID)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] PermanentlyDeleteUserProfile: failed to fetch profile prior to deletion: %w", err)
	}
	deletedProfile, err := s.repo.PermanentlyDeleteProfile(ctx, requestingUserID)
	if err != nil {
		return domain.Profile{}, fmt.Errorf("[profile_service.go] PermanentlyDeleteUserProfile: failed to permanently delete profile: %w", err)
	}

	if profileToDelete.ProfileImageUrl != "" {
		imageErr := s.imageStorage.DeleteProfilePicture(ctx, profileToDelete.ProfileImageUrl)
		if imageErr != nil {
			log.Printf("WARNING: Orphaned profile image left in Cloudflare R2 for user %s. Error: %v\n", requestingUserID, imageErr)
		}
	}

	return deletedProfile, nil
}
