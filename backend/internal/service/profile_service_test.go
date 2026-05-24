package service

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
)

type mockProfileRepo struct {
	mockProfile                   domain.Profile
	mockProfileFromID             domain.Profile
	mockUpdatedProfile            domain.Profile
	mockSoftDeletedProfile        domain.Profile
	mockPermanentlyDeletedProfile domain.Profile
	mockProfileList               []domain.Profile
	mockUsernameExists            bool
	mockSubscriberTierProfile     domain.Profile
	mockSubscriberTier            domain.SubscriberTier
	mockRetrieveError             error
	mockGetByIDError              error
	mockCheckError                error
	mockCreateError               error
	mockUpdateError               error
	mockSoftDeleteError           error
	mockPermanentDeleteError      error
	mockSetSubscriberTierError    error
	mockGetSubscriberTierError    error

	getByIDCalled           bool
	getByUsernameCalled     bool
	checkUsernameCalled     bool
	getProfileListCalled    bool
	createProfileCalled     bool
	updateProfileCalled     bool
	softDeleteCalled        bool
	permanentDeleteCalled   bool
	setSubscriberTierCalled bool
	lastActorUserID         string
	lastUserID              string
	lastTargetProfileID     string
	lastUsername            string
	lastLimit               int32
	lastOffset              int32
	lastCreatedProfileArg   domain.Profile
	lastUpdatedProfileArg   domain.Profile
	lastSubscriberTier      domain.SubscriberTier
}

type mockImageStorage struct {
	mockURL         string
	mockUploadError error
	mockDeleteError error
	uploadCalled    bool
	deleteCalled    bool
	lastFileBytes   []byte
	lastFileName    string
	lastDeletedURL  string
}

func (m *mockImageStorage) UploadPicture(_ context.Context, fileBytes []byte, fileName string) (string, error) {
	m.uploadCalled = true
	m.lastFileBytes = append([]byte(nil), fileBytes...)
	m.lastFileName = fileName
	if m.mockUploadError != nil {
		return "", m.mockUploadError
	}
	return m.mockURL, nil
}

func (m *mockImageStorage) DeletePicture(_ context.Context, fullImageURL string) error {
	m.deleteCalled = true
	m.lastDeletedURL = fullImageURL
	if m.mockDeleteError != nil {
		return m.mockDeleteError
	}
	return m.mockDeleteError
}

func (m *mockProfileRepo) GetProfileFromID(_ context.Context, userID string, targetProfileID string) (domain.Profile, error) {
	m.getByIDCalled = true
	m.lastUserID = userID
	m.lastTargetProfileID = targetProfileID
	if m.mockGetByIDError != nil {
		return domain.Profile{}, m.mockGetByIDError
	}
	if m.mockProfileFromID != (domain.Profile{}) {
		return m.mockProfileFromID, nil
	}
	return m.mockProfile, m.mockRetrieveError
}

func (m *mockProfileRepo) GetProfileFromUsername(_ context.Context, userID string, username string) (domain.Profile, error) {
	m.getByUsernameCalled = true
	m.lastUserID = userID
	m.lastUsername = username
	return m.mockProfile, m.mockRetrieveError
}

func (m *mockProfileRepo) CheckProfileUsernameExists(_ context.Context, userID string, username string) (bool, error) {
	m.checkUsernameCalled = true
	m.lastUserID = userID
	m.lastUsername = username
	return m.mockUsernameExists, m.mockCheckError
}

func (m *mockProfileRepo) GetListOfProfiles(_ context.Context, userID string, limit int32, offset int32) ([]domain.Profile, error) {
	m.getProfileListCalled = true
	m.lastUserID = userID
	m.lastLimit = limit
	m.lastOffset = offset
	return m.mockProfileList, m.mockRetrieveError
}

func (m *mockProfileRepo) CreateProfile(_ context.Context, userID string, newProfile domain.Profile) (domain.Profile, error) {
	m.createProfileCalled = true
	m.lastUserID = userID
	m.lastCreatedProfileArg = newProfile
	if m.mockCreateError != nil {
		return domain.Profile{}, m.mockCreateError
	}
	return m.mockProfile, nil
}

func (m *mockProfileRepo) UpdateProfile(ctx context.Context, userID string, profileWithUpdates domain.Profile) (domain.Profile, error) {
	m.updateProfileCalled = true
	m.lastUserID = userID
	m.lastUpdatedProfileArg = profileWithUpdates
	if m.mockUpdateError != nil {
		return domain.Profile{}, m.mockUpdateError
	}
	if m.mockUpdatedProfile != (domain.Profile{}) {
		return m.mockUpdatedProfile, nil
	}
	return profileWithUpdates, nil
}

func (m *mockProfileRepo) SoftDeleteProfile(ctx context.Context, userID string) (domain.Profile, error) {
	m.softDeleteCalled = true
	m.lastUserID = userID
	if m.mockSoftDeleteError != nil {
		return domain.Profile{}, m.mockSoftDeleteError
	}
	if m.mockSoftDeletedProfile != (domain.Profile{}) {
		return m.mockSoftDeletedProfile, nil
	}
	return m.mockProfile, nil
}

func (m *mockProfileRepo) PermanentlyDeleteProfile(ctx context.Context, userID string) (domain.Profile, error) {
	m.permanentDeleteCalled = true
	m.lastUserID = userID
	if m.mockPermanentDeleteError != nil {
		return domain.Profile{}, m.mockPermanentDeleteError
	}
	if m.mockPermanentlyDeletedProfile != (domain.Profile{}) {
		return m.mockPermanentlyDeletedProfile, nil
	}
	return m.mockProfile, nil
}

func (m *mockProfileRepo) SetProfileSubscriberTier(ctx context.Context, adminID string, userID string, tier domain.SubscriberTier) (domain.Profile, error) {
	m.setSubscriberTierCalled = true
	m.lastActorUserID = adminID
	m.lastUserID = userID
	m.lastSubscriberTier = tier
	if m.mockSetSubscriberTierError != nil {
		return domain.Profile{}, m.mockSetSubscriberTierError
	}
	if m.mockSubscriberTierProfile != (domain.Profile{}) {
		return m.mockSubscriberTierProfile, nil
	}
	return m.mockProfile, nil
}

func (m *mockProfileRepo) GetProfileSubscriberTier(_ context.Context, userID string) (domain.SubscriberTier, error) {
	m.getByIDCalled = true
	m.lastUserID = userID
	if m.mockGetByIDError != nil {
		var zeroTier domain.SubscriberTier
		return zeroTier, m.mockGetSubscriberTierError
	}

	var zeroTier domain.SubscriberTier
	if m.mockSubscriberTier != (zeroTier) {
		return m.mockSubscriberTier, nil
	}
	return m.mockSubscriberTier, m.mockRetrieveError
}

func TestProfileService_FetchProfileByID(t *testing.T) {
	t.Run("returns validation error when target profile ID is empty", func(t *testing.T) {
		repo := &mockProfileRepo{}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profile, err := service.FetchProfileByID(context.Background(), "requester-1", "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "target profile ID is required") {
			t.Fatalf("expected validation error, got %v", err)
		}
		if profile != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", profile)
		}
		if repo.getByIDCalled {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("returns wrapped repository error", func(t *testing.T) {
		repoErr := errors.New("database timeout")
		repo := &mockProfileRepo{mockRetrieveError: repoErr}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		_, err := service.FetchProfileByID(context.Background(), "requester-1", "profile-9")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to get profile from id profile-9 with id requester-1") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
		if !repo.getByIDCalled {
			t.Fatal("expected repository to be called")
		}
	})

	t.Run("returns profile when repository succeeds", func(t *testing.T) {
		expected := domain.Profile{ID: "profile-9", Username: "sunshine"}
		repo := &mockProfileRepo{mockProfile: expected}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profile, err := service.FetchProfileByID(context.Background(), "requester-1", "profile-9")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(profile, expected) {
			t.Fatalf("expected %#v, got %#v", expected, profile)
		}
		if repo.lastUserID != "requester-1" || repo.lastTargetProfileID != "profile-9" {
			t.Fatalf("expected repo to receive requester-1/profile-9, got %s/%s", repo.lastUserID, repo.lastTargetProfileID)
		}
	})
}

func TestProfileService_FetchProfileByUsername(t *testing.T) {
	t.Run("returns validation error when username is empty", func(t *testing.T) {
		repo := &mockProfileRepo{}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profile, err := service.FetchProfileByUsername(context.Background(), "requester-1", "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "target username is required") {
			t.Fatalf("expected validation error, got %v", err)
		}
		if profile != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", profile)
		}
		if repo.getByUsernameCalled {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("returns wrapped repository error", func(t *testing.T) {
		repoErr := errors.New("lookup failed")
		repo := &mockProfileRepo{mockRetrieveError: repoErr}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		_, err := service.FetchProfileByUsername(context.Background(), "requester-1", "sunshine")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to get profile from username sunshine with id requester-1") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
	})

	t.Run("returns profile when repository succeeds", func(t *testing.T) {
		expected := domain.Profile{ID: "profile-9", Username: "sunshine"}
		repo := &mockProfileRepo{mockProfile: expected}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profile, err := service.FetchProfileByUsername(context.Background(), "requester-1", "sunshine")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(profile, expected) {
			t.Fatalf("expected %#v, got %#v", expected, profile)
		}
		if repo.lastUserID != "requester-1" || repo.lastUsername != "sunshine" {
			t.Fatalf("expected repo to receive requester-1/sunshine, got %s/%s", repo.lastUserID, repo.lastUsername)
		}
	})
}

func TestProfileService_IsUsernameTaken(t *testing.T) {
	t.Run("returns validation error when username is empty", func(t *testing.T) {
		repo := &mockProfileRepo{}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		taken, err := service.IsUsernameTaken(context.Background(), "requester-1", "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !taken {
			t.Fatal("expected taken to default to true on validation error")
		}
		if !strings.Contains(err.Error(), "username is required") {
			t.Fatalf("expected validation error, got %v", err)
		}
		if repo.checkUsernameCalled {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("returns wrapped repository error", func(t *testing.T) {
		repoErr := errors.New("check failed")
		repo := &mockProfileRepo{mockCheckError: repoErr}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		taken, err := service.IsUsernameTaken(context.Background(), "requester-1", "sunshine")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !taken {
			t.Fatal("expected taken to default to true on repository error")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to check if username sunshine exists from user with id requester-1") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
	})

	t.Run("returns repository result when username check succeeds", func(t *testing.T) {
		repo := &mockProfileRepo{mockUsernameExists: false}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		taken, err := service.IsUsernameTaken(context.Background(), "requester-1", "sunshine")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if taken {
			t.Fatal("expected username to be available")
		}
		if repo.lastUserID != "requester-1" || repo.lastUsername != "sunshine" {
			t.Fatalf("expected repo to receive requester-1/sunshine, got %s/%s", repo.lastUserID, repo.lastUsername)
		}
	})
}

func TestProfileService_FetchListOfProfiles(t *testing.T) {
	t.Run("returns wrapped repository error", func(t *testing.T) {
		repoErr := errors.New("query failed")
		repo := &mockProfileRepo{mockRetrieveError: repoErr}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profiles, err := service.FetchListOfProfiles(context.Background(), "requester-1", 25, 10)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if len(profiles) != 0 {
			t.Fatalf("expected empty slice, got %#v", profiles)
		}
		if !strings.Contains(err.Error(), "failed to get list of profiles with limit 25 and offset 10 with id requester-1") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
	})

	t.Run("caps limit at one hundred before calling repository", func(t *testing.T) {
		expected := []domain.Profile{{ID: "1", Username: "sunny"}, {ID: "2", Username: "brighter"}}
		repo := &mockProfileRepo{mockProfileList: expected}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profiles, err := service.FetchListOfProfiles(context.Background(), "requester-1", 250, 5)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(profiles, expected) {
			t.Fatalf("expected %#v, got %#v", expected, profiles)
		}
		if repo.lastLimit != 100 || repo.lastOffset != 5 {
			t.Fatalf("expected repo to receive limit 100 and offset 5, got %d and %d", repo.lastLimit, repo.lastOffset)
		}
	})
}

func TestProfileService_CreateUserProfile(t *testing.T) {
	baseProfile := domain.Profile{
		ID:             "profile-9",
		FirstName:      "Daily",
		LastName:       "Sunshine",
		Username:       "sunshine",
		SubscriberTier: domain.Member,
	}

	t.Run("returns validation error when username is empty", func(t *testing.T) {
		repo := &mockProfileRepo{}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		created, err := service.CreateUserProfile(context.Background(), "requester-1", domain.Profile{}, nil, "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "username is required") {
			t.Fatalf("expected validation error, got %v", err)
		}
		if created != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", created)
		}
		if repo.checkUsernameCalled || repo.createProfileCalled {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("returns wrapped error when username lookup fails", func(t *testing.T) {
		repoErr := errors.New("lookup failed")
		repo := &mockProfileRepo{mockCheckError: repoErr}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		_, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile, nil, "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to check if username sunshine exists from user with id requester-1") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
		if repo.createProfileCalled {
			t.Fatal("expected create not to be called when username lookup fails")
		}
	})

	t.Run("returns error when username is already taken", func(t *testing.T) {
		repo := &mockProfileRepo{mockUsernameExists: true}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		created, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile, nil, "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "the username 'sunshine' is already taken") {
			t.Fatalf("expected duplicate username error, got %v", err)
		}
		if created != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", created)
		}
		if repo.createProfileCalled {
			t.Fatal("expected create not to be called when username is taken")
		}
	})

	t.Run("returns wrapped error when create fails", func(t *testing.T) {
		repoErr := errors.New("insert failed")
		repo := &mockProfileRepo{
			mockUsernameExists: false,
			mockCreateError:    repoErr,
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		_, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile, nil, "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to insert into database") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
		if !repo.createProfileCalled {
			t.Fatal("expected create to be called")
		}
		if !reflect.DeepEqual(repo.lastCreatedProfileArg, baseProfile) {
			t.Fatalf("expected create to receive %#v, got %#v", baseProfile, repo.lastCreatedProfileArg)
		}
		if imageStorage.uploadCalled {
			t.Fatal("expected image storage not to be called without image bytes")
		}
	})

	t.Run("creates profile when username is available without image upload", func(t *testing.T) {
		expected := baseProfile
		repo := &mockProfileRepo{
			mockUsernameExists: false,
			mockProfile:        expected,
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		created, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile, nil, "")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(created, expected) {
			t.Fatalf("expected %#v, got %#v", expected, created)
		}
		if !repo.checkUsernameCalled || !repo.createProfileCalled {
			t.Fatal("expected username check and create to both be called")
		}
		if repo.lastUserID != "requester-1" {
			t.Fatalf("expected repo to receive requester-1, got %s", repo.lastUserID)
		}
		if imageStorage.uploadCalled {
			t.Fatal("expected image storage not to be called without image bytes")
		}
	})

	t.Run("returns wrapped error when image upload fails", func(t *testing.T) {
		repo := &mockProfileRepo{mockUsernameExists: false}
		imageStorage := &mockImageStorage{mockUploadError: errors.New("r2 unavailable")}
		service := NewProfileService(repo, imageStorage)
		imageBytes := []byte{0x01, 0x02, 0x03}

		created, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile, imageBytes, ".jpg")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "failed to upload profile picture") {
			t.Fatalf("expected upload error message, got %v", err)
		}
		if created != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", created)
		}
		if !imageStorage.uploadCalled {
			t.Fatal("expected image storage to be called")
		}
		if imageStorage.lastFileName != "profiles/requester-1/avatar.jpg" {
			t.Fatalf("expected upload file name profiles/requester-1/avatar.jpg, got %s", imageStorage.lastFileName)
		}
		if !reflect.DeepEqual(imageStorage.lastFileBytes, imageBytes) {
			t.Fatalf("expected uploaded bytes %#v, got %#v", imageBytes, imageStorage.lastFileBytes)
		}
		if repo.createProfileCalled {
			t.Fatal("expected repository create not to be called when upload fails")
		}
	})

	t.Run("uploads image and saves returned URL before creating profile", func(t *testing.T) {
		expected := baseProfile
		expected.ProfileImageUrl = "https://cdn.example.com/profiles/requester-1/avatar.jpg"
		repo := &mockProfileRepo{
			mockUsernameExists: false,
			mockProfile:        expected,
		}
		imageStorage := &mockImageStorage{mockURL: expected.ProfileImageUrl}
		service := NewProfileService(repo, imageStorage)
		imageBytes := []byte{0xFF, 0xD8, 0xFF}

		created, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile, imageBytes, ".jpg")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(created, expected) {
			t.Fatalf("expected %#v, got %#v", expected, created)
		}
		if !imageStorage.uploadCalled {
			t.Fatal("expected image storage to be called")
		}
		if imageStorage.lastFileName != "profiles/requester-1/avatar.jpg" {
			t.Fatalf("expected upload file name profiles/requester-1/avatar.jpg, got %s", imageStorage.lastFileName)
		}
		if !reflect.DeepEqual(imageStorage.lastFileBytes, imageBytes) {
			t.Fatalf("expected uploaded bytes %#v, got %#v", imageBytes, imageStorage.lastFileBytes)
		}
		if repo.lastCreatedProfileArg.ProfileImageUrl != expected.ProfileImageUrl {
			t.Fatalf("expected created profile image URL %s, got %s", expected.ProfileImageUrl, repo.lastCreatedProfileArg.ProfileImageUrl)
		}
	})

	t.Run("deletes uploaded image when create fails after upload", func(t *testing.T) {
		repoErr := errors.New("insert failed")
		repo := &mockProfileRepo{
			mockUsernameExists: false,
			mockCreateError:    repoErr,
		}
		imageURL := "https://cdn.example.com/profiles/requester-1/avatar.jpg"
		imageStorage := &mockImageStorage{mockURL: imageURL}
		service := NewProfileService(repo, imageStorage)
		imageBytes := []byte{0xAA, 0xBB, 0xCC}

		created, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile, imageBytes, ".jpg")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if created != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", created)
		}
		if !imageStorage.uploadCalled {
			t.Fatal("expected upload to be attempted")
		}
		if !imageStorage.deleteCalled {
			t.Fatal("expected uploaded image to be cleaned up")
		}
		if imageStorage.lastDeletedURL != imageURL {
			t.Fatalf("expected cleanup URL %s, got %s", imageURL, imageStorage.lastDeletedURL)
		}
	})
}

func TestProfileService_UpdateUserProfile(t *testing.T) {
	baseProfile := domain.Profile{
		ID:              "requester-1",
		FirstName:       "Daily",
		LastName:        "Sunshine",
		Username:        "sunshine",
		ProfileImageUrl: "https://cdn.example.com/profiles/requester-1/avatar-old.jpg",
		SubscriberTier:  domain.Member,
	}

	t.Run("returns validation error when username is empty", func(t *testing.T) {
		repo := &mockProfileRepo{}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		updated, err := service.UpdateUserProfile(context.Background(), "requester-1", domain.Profile{}, nil, "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "UpdateUserProfile: username is required") {
			t.Fatalf("expected validation error, got %v", err)
		}
		if updated != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", updated)
		}
		if repo.getByIDCalled || repo.updateProfileCalled {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("returns wrapped error when fetching current profile fails", func(t *testing.T) {
		repoErr := errors.New("fetch current profile failed")
		repo := &mockProfileRepo{mockGetByIDError: repoErr}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		_, err := service.UpdateUserProfile(context.Background(), "requester-1", baseProfile, nil, "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to fetch current profile") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
		if repo.updateProfileCalled {
			t.Fatal("expected update not to be called when current profile lookup fails")
		}
	})

	t.Run("returns wrapped error when checking new username fails", func(t *testing.T) {
		repoErr := errors.New("username lookup failed")
		repo := &mockProfileRepo{
			mockProfileFromID: baseProfile,
			mockCheckError:    repoErr,
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profileWithUpdates := baseProfile
		profileWithUpdates.Username = "freshsunshine"

		_, err := service.UpdateUserProfile(context.Background(), "requester-1", profileWithUpdates, nil, "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to check if username freshsunshine exists from user with id requester-1") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
		if repo.updateProfileCalled {
			t.Fatal("expected update not to be called when username check fails")
		}
	})

	t.Run("returns error when new username is already taken", func(t *testing.T) {
		repo := &mockProfileRepo{
			mockProfileFromID:  baseProfile,
			mockUsernameExists: true,
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profileWithUpdates := baseProfile
		profileWithUpdates.Username = "freshsunshine"

		updated, err := service.UpdateUserProfile(context.Background(), "requester-1", profileWithUpdates, nil, "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "the username 'freshsunshine' is already taken") {
			t.Fatalf("expected duplicate username error, got %v", err)
		}
		if updated != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", updated)
		}
		if repo.updateProfileCalled {
			t.Fatal("expected update not to be called when username is taken")
		}
	})

	t.Run("preserves current image when no new image is uploaded", func(t *testing.T) {
		repo := &mockProfileRepo{
			mockProfileFromID: baseProfile,
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profileWithUpdates := baseProfile
		profileWithUpdates.FirstName = "Brighter"
		profileWithUpdates.ProfileImageUrl = ""

		updated, err := service.UpdateUserProfile(context.Background(), "requester-1", profileWithUpdates, nil, "")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if updated.ProfileImageUrl != baseProfile.ProfileImageUrl {
			t.Fatalf("expected image URL %s, got %s", baseProfile.ProfileImageUrl, updated.ProfileImageUrl)
		}
		if repo.lastUpdatedProfileArg.ProfileImageUrl != baseProfile.ProfileImageUrl {
			t.Fatalf("expected update arg image URL %s, got %s", baseProfile.ProfileImageUrl, repo.lastUpdatedProfileArg.ProfileImageUrl)
		}
		if repo.checkUsernameCalled {
			t.Fatal("expected username existence check to be skipped when username is unchanged")
		}
		if imageStorage.uploadCalled {
			t.Fatal("expected image storage not to be called without image bytes")
		}
	})

	t.Run("returns wrapped error when new profile image upload fails", func(t *testing.T) {
		repo := &mockProfileRepo{
			mockProfileFromID: baseProfile,
		}
		imageStorage := &mockImageStorage{mockUploadError: errors.New("r2 upload failed")}
		service := NewProfileService(repo, imageStorage)
		imageBytes := []byte{0x10, 0x20, 0x30}

		updated, err := service.UpdateUserProfile(context.Background(), "requester-1", baseProfile, imageBytes, ".png")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "UpdateUserProfile: failed to upload profile picture") {
			t.Fatalf("expected upload error message, got %v", err)
		}
		if updated != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", updated)
		}
		if !imageStorage.uploadCalled {
			t.Fatal("expected image storage to be called")
		}
		if imageStorage.lastFileName != "profiles/requester-1/avatar.png" {
			t.Fatalf("expected upload file name profiles/requester-1/avatar.png, got %s", imageStorage.lastFileName)
		}
		if repo.updateProfileCalled {
			t.Fatal("expected update not to be called when image upload fails")
		}
	})

	t.Run("uploads new image and updates profile", func(t *testing.T) {
		profileWithUpdates := baseProfile
		profileWithUpdates.Username = "freshsunshine"

		expected := profileWithUpdates
		expected.ProfileImageUrl = "https://cdn.example.com/profiles/requester-1/avatar.jpg"

		repo := &mockProfileRepo{
			mockProfileFromID:  baseProfile,
			mockUpdatedProfile: expected,
		}
		imageStorage := &mockImageStorage{mockURL: expected.ProfileImageUrl}
		service := NewProfileService(repo, imageStorage)
		imageBytes := []byte{0xFF, 0xD8, 0xFF}

		updated, err := service.UpdateUserProfile(context.Background(), "requester-1", profileWithUpdates, imageBytes, ".jpg")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(updated, expected) {
			t.Fatalf("expected %#v, got %#v", expected, updated)
		}
		if !repo.checkUsernameCalled {
			t.Fatal("expected username existence check for changed username")
		}
		if !repo.updateProfileCalled {
			t.Fatal("expected update to be called")
		}
		if repo.lastUpdatedProfileArg.ProfileImageUrl != expected.ProfileImageUrl {
			t.Fatalf("expected update arg image URL %s, got %s", expected.ProfileImageUrl, repo.lastUpdatedProfileArg.ProfileImageUrl)
		}
		if !imageStorage.uploadCalled {
			t.Fatal("expected image storage to be called")
		}
	})

	t.Run("returns wrapped error when update repository call fails", func(t *testing.T) {
		repoErr := errors.New("update failed")
		repo := &mockProfileRepo{
			mockProfileFromID: baseProfile,
			mockUpdateError:   repoErr,
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		_, err := service.UpdateUserProfile(context.Background(), "requester-1", baseProfile, nil, "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to update profile in database") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
	})

	t.Run("deletes newly uploaded image when update fails after upload", func(t *testing.T) {
		repoErr := errors.New("update failed")
		repo := &mockProfileRepo{
			mockProfileFromID: baseProfile,
			mockUpdateError:   repoErr,
		}
		imageURL := "https://cdn.example.com/profiles/requester-1/avatar-new.jpg"
		imageStorage := &mockImageStorage{mockURL: imageURL}
		service := NewProfileService(repo, imageStorage)
		imageBytes := []byte{0x01, 0x02, 0x03}

		updated, err := service.UpdateUserProfile(context.Background(), "requester-1", baseProfile, imageBytes, ".jpg")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if updated != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", updated)
		}
		if !imageStorage.uploadCalled {
			t.Fatal("expected upload to be attempted")
		}
		if !imageStorage.deleteCalled {
			t.Fatal("expected uploaded image to be cleaned up")
		}
		if imageStorage.lastDeletedURL != imageURL {
			t.Fatalf("expected cleanup URL %s, got %s", imageURL, imageStorage.lastDeletedURL)
		}
	})
}

func TestProfileService_SoftDeleteUserProfile(t *testing.T) {
	t.Run("returns wrapped error when soft delete fails", func(t *testing.T) {
		repoErr := errors.New("soft delete failed")
		repo := &mockProfileRepo{mockSoftDeleteError: repoErr}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		deleted, err := service.SoftDeleteUserProfile(context.Background(), "requester-1")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "SoftDeleteUserProfile: failed to soft delete profile") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
		if deleted != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", deleted)
		}
	})

	t.Run("returns soft deleted profile when repository succeeds", func(t *testing.T) {
		expected := domain.Profile{ID: "requester-1", Username: "sunshine"}
		repo := &mockProfileRepo{mockSoftDeletedProfile: expected}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		deleted, err := service.SoftDeleteUserProfile(context.Background(), "requester-1")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(deleted, expected) {
			t.Fatalf("expected %#v, got %#v", expected, deleted)
		}
		if !repo.softDeleteCalled {
			t.Fatal("expected soft delete to be called")
		}
		if imageStorage.uploadCalled || imageStorage.deleteCalled {
			t.Fatal("expected image storage not to be used during soft delete")
		}
	})
}

func TestProfileService_PermanentlyDeleteUserProfile(t *testing.T) {
	t.Run("returns wrapped error when fetching profile before deletion fails", func(t *testing.T) {
		repoErr := errors.New("fetch failed")
		repo := &mockProfileRepo{mockGetByIDError: repoErr}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		deleted, err := service.PermanentlyDeleteUserProfile(context.Background(), "requester-1")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to fetch profile prior to deletion") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
		if deleted != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", deleted)
		}
		if repo.permanentDeleteCalled {
			t.Fatal("expected permanent delete not to be called when fetch fails")
		}
	})

	t.Run("returns wrapped error when repository delete fails", func(t *testing.T) {
		repoErr := errors.New("delete failed")
		repo := &mockProfileRepo{
			mockProfileFromID:        domain.Profile{ID: "requester-1", ProfileImageUrl: "https://cdn.example.com/profile.jpg"},
			mockPermanentDeleteError: repoErr,
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		deleted, err := service.PermanentlyDeleteUserProfile(context.Background(), "requester-1")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to permanently delete profile") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
		if deleted != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", deleted)
		}
		if imageStorage.deleteCalled {
			t.Fatal("expected image delete not to run when repo delete fails")
		}
	})

	t.Run("deletes image when profile has a stored image", func(t *testing.T) {
		profileToDelete := domain.Profile{
			ID:              "requester-1",
			Username:        "sunshine",
			ProfileImageUrl: "https://cdn.example.com/profiles/requester-1/avatar.jpg",
		}
		expected := domain.Profile{ID: "requester-1", Username: "sunshine"}
		repo := &mockProfileRepo{
			mockProfileFromID:             profileToDelete,
			mockPermanentlyDeletedProfile: expected,
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		deleted, err := service.PermanentlyDeleteUserProfile(context.Background(), "requester-1")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(deleted, expected) {
			t.Fatalf("expected %#v, got %#v", expected, deleted)
		}
		if !repo.permanentDeleteCalled {
			t.Fatal("expected permanent delete to be called")
		}
		if !imageStorage.deleteCalled {
			t.Fatal("expected image delete to be called")
		}
		if imageStorage.lastDeletedURL != profileToDelete.ProfileImageUrl {
			t.Fatalf("expected deleted image URL %s, got %s", profileToDelete.ProfileImageUrl, imageStorage.lastDeletedURL)
		}
	})

	t.Run("does not attempt image deletion when profile has no image", func(t *testing.T) {
		profileToDelete := domain.Profile{ID: "requester-1", Username: "sunshine"}
		expected := domain.Profile{ID: "requester-1", Username: "sunshine"}
		repo := &mockProfileRepo{
			mockProfileFromID:             profileToDelete,
			mockPermanentlyDeletedProfile: expected,
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		deleted, err := service.PermanentlyDeleteUserProfile(context.Background(), "requester-1")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(deleted, expected) {
			t.Fatalf("expected %#v, got %#v", expected, deleted)
		}
		if imageStorage.deleteCalled {
			t.Fatal("expected image delete not to be called")
		}
	})

	t.Run("still succeeds when image cleanup fails", func(t *testing.T) {
		profileToDelete := domain.Profile{
			ID:              "requester-1",
			Username:        "sunshine",
			ProfileImageUrl: "https://cdn.example.com/profiles/requester-1/avatar.jpg",
		}
		expected := domain.Profile{ID: "requester-1", Username: "sunshine"}
		repo := &mockProfileRepo{
			mockProfileFromID:             profileToDelete,
			mockPermanentlyDeletedProfile: expected,
		}
		imageStorage := &mockImageStorage{mockDeleteError: errors.New("r2 delete failed")}
		service := NewProfileService(repo, imageStorage)

		deleted, err := service.PermanentlyDeleteUserProfile(context.Background(), "requester-1")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(deleted, expected) {
			t.Fatalf("expected %#v, got %#v", expected, deleted)
		}
		if !imageStorage.deleteCalled {
			t.Fatal("expected image delete to be attempted")
		}
	})
}

func TestProfileService_IsAdmin(t *testing.T) {
	t.Run("returns validation error when user ID is empty", func(t *testing.T) {
		repo := &mockProfileRepo{}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		isAdmin, err := service.IsAdmin(context.Background(), "")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if isAdmin {
			t.Fatal("expected false when validation fails")
		}
		if !strings.Contains(err.Error(), "user ID is required") {
			t.Fatalf("expected validation error, got %v", err)
		}
		if repo.getByIDCalled {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("returns wrapped error when loading profile fails", func(t *testing.T) {
		repoErr := errors.New("lookup failed")
		repo := &mockProfileRepo{mockGetByIDError: repoErr}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		isAdmin, err := service.IsAdmin(context.Background(), "requester-1")
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if isAdmin {
			t.Fatal("expected false on repository failure")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to load profile for user requester-1") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
	})

	t.Run("returns true when app role is admin", func(t *testing.T) {
		repo := &mockProfileRepo{
			mockProfileFromID: domain.Profile{ID: "requester-1", AppRole: domain.Admin},
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		isAdmin, err := service.IsAdmin(context.Background(), "requester-1")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !isAdmin {
			t.Fatal("expected true for admin role")
		}
		if repo.lastUserID != "requester-1" || repo.lastTargetProfileID != "requester-1" {
			t.Fatalf("expected repo to receive requester-1/requester-1, got %s/%s", repo.lastUserID, repo.lastTargetProfileID)
		}
	})

	t.Run("returns false when app role is user", func(t *testing.T) {
		repo := &mockProfileRepo{
			mockProfileFromID: domain.Profile{ID: "requester-1", AppRole: domain.User},
		}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		isAdmin, err := service.IsAdmin(context.Background(), "requester-1")
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if isAdmin {
			t.Fatal("expected false for non-admin role")
		}
	})
}

func TestProfileService_SetSubscriberTier(t *testing.T) {
	t.Run("returns validation error when target user ID is empty", func(t *testing.T) {
		repo := &mockProfileRepo{}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profile, err := service.SetSubscriberTier(context.Background(), "admin-1", "", domain.Subscriber)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if profile != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", profile)
		}
		if !strings.Contains(err.Error(), "target profile ID is required") {
			t.Fatalf("expected validation error, got %v", err)
		}
		if repo.setSubscriberTierCalled {
			t.Fatal("expected repository not to be called")
		}
	})

	t.Run("returns wrapped repository error", func(t *testing.T) {
		repoErr := errors.New("update tier failed")
		repo := &mockProfileRepo{mockSetSubscriberTierError: repoErr}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profile, err := service.SetSubscriberTier(context.Background(), "admin-1", "user-123", domain.Subscriber)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if profile != (domain.Profile{}) {
			t.Fatalf("expected empty profile, got %#v", profile)
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected wrapped repo error, got %v", err)
		}
		if !strings.Contains(err.Error(), "failed to set subscriber tier for user with id user-123") {
			t.Fatalf("expected contextual error message, got %v", err)
		}
		if repo.lastActorUserID != "admin-1" {
			t.Fatalf("expected actor user ID admin-1, got %s", repo.lastActorUserID)
		}
		if repo.lastUserID != "user-123" {
			t.Fatalf("expected target user ID user-123, got %s", repo.lastUserID)
		}
	})

	t.Run("updates subscriber tier when repository succeeds", func(t *testing.T) {
		expected := domain.Profile{
			ID:             "user-123",
			SubscriberTier: domain.Subscriber,
		}
		repo := &mockProfileRepo{mockSubscriberTierProfile: expected}
		imageStorage := &mockImageStorage{}
		service := NewProfileService(repo, imageStorage)

		profile, err := service.SetSubscriberTier(context.Background(), "admin-1", "user-123", domain.Subscriber)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if !reflect.DeepEqual(profile, expected) {
			t.Fatalf("expected %#v, got %#v", expected, profile)
		}
		if !repo.setSubscriberTierCalled {
			t.Fatal("expected repository to be called")
		}
		if repo.lastActorUserID != "admin-1" {
			t.Fatalf("expected actor user ID admin-1, got %s", repo.lastActorUserID)
		}
		if repo.lastUserID != "user-123" {
			t.Fatalf("expected target user ID user-123, got %s", repo.lastUserID)
		}
		if repo.lastSubscriberTier != domain.Subscriber {
			t.Fatalf("expected tier %q, got %q", domain.Subscriber, repo.lastSubscriberTier)
		}
	})
}
