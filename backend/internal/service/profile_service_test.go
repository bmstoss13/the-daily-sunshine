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
	mockProfile        domain.Profile
	mockProfileList    []domain.Profile
	mockUsernameExists bool
	mockRetrieveError  error
	mockCheckError     error
	mockCreateError    error

	getByIDCalled         bool
	getByUsernameCalled   bool
	checkUsernameCalled   bool
	getProfileListCalled  bool
	createProfileCalled   bool
	lastUserID            string
	lastTargetProfileID   string
	lastUsername          string
	lastLimit             int32
	lastOffset            int32
	lastCreatedProfileArg domain.Profile
}

func (m *mockProfileRepo) GetProfileFromID(_ context.Context, userID string, targetProfileID string) (domain.Profile, error) {
	m.getByIDCalled = true
	m.lastUserID = userID
	m.lastTargetProfileID = targetProfileID
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

func TestProfileService_FetchProfileByID(t *testing.T) {
	t.Run("returns validation error when target profile ID is empty", func(t *testing.T) {
		repo := &mockProfileRepo{}
		service := NewProfileService(repo)

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
		service := NewProfileService(repo)

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
		service := NewProfileService(repo)

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
		service := NewProfileService(repo)

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
		service := NewProfileService(repo)

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
		service := NewProfileService(repo)

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
		service := NewProfileService(repo)

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
		service := NewProfileService(repo)

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
		service := NewProfileService(repo)

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
		service := NewProfileService(repo)

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
		service := NewProfileService(repo)

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
		ID:        "profile-9",
		FirstName: "Daily",
		LastName:  "Sunshine",
		Username:  "sunshine",
		Role:      domain.Member,
	}

	t.Run("returns validation error when username is empty", func(t *testing.T) {
		repo := &mockProfileRepo{}
		service := NewProfileService(repo)

		created, err := service.CreateUserProfile(context.Background(), "requester-1", domain.Profile{})
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
		service := NewProfileService(repo)

		_, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile)
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
		service := NewProfileService(repo)

		created, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile)
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
		service := NewProfileService(repo)

		_, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile)
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
	})

	t.Run("creates profile when username is available", func(t *testing.T) {
		expected := baseProfile
		repo := &mockProfileRepo{
			mockUsernameExists: false,
			mockProfile:        expected,
		}
		service := NewProfileService(repo)

		created, err := service.CreateUserProfile(context.Background(), "requester-1", baseProfile)
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
	})
}
