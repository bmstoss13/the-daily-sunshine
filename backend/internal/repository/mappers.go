package repository

import (
	"fmt"

	"github.com/bmstoss13/the-daily-sunshine/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

// StringToPgUUID converts a standard string to a pgtype.UUID.
func StringToPgUUID(id string) (pgtype.UUID, error) {
	var pgUUID pgtype.UUID
	err := pgUUID.Scan(id)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return pgUUID, nil
}

// PgUUIDToString safely converts a pgtype.UUID back to a string.
// If the UUID is null or invalid, it returns an empty string.
func PgUUIDToString(pgUUID pgtype.UUID) string {
	if !pgUUID.Valid {
		return ""
	}

	// pgtype.UUID stores the UUID as a [16]byte array.
	// format it into the standard 8-4-4-4-12 string representation.
	b := pgUUID.Bytes
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func ConvertProfile(sqlcProfile Profile) domain.Profile {
	convertedProfile := domain.Profile{
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
		convertedProfile.ProfileImageUrl = *sqlcProfile.ProfileImageUrl
	} else {
		convertedProfile.ProfileImageUrl = "" // Or a default URL string
	}

	if sqlcProfile.ProfileBio != nil {
		bioStr := string(sqlcProfile.ProfileBio)
		convertedProfile.Bio = &bioStr
	}

	if sqlcProfile.DeletedAt.Valid {
		convertedProfile.DeletedAt = &sqlcProfile.DeletedAt.Time
	}

	return convertedProfile
}

func ConvertPost(sqlcPost Post) domain.Post {

	convertedPost := domain.Post{
		ID:          PgUUIDToString(sqlcPost.ID),
		PublisherID: PgUUIDToString(sqlcPost.PublisherID),
		Title:       sqlcPost.Title,
		Subtitle:    sqlcPost.Subtitle,
		Slug:        sqlcPost.Slug,
		Status:      domain.PostStatus(sqlcPost.Status),
		NumRays:     int64(sqlcPost.NumRays),
		NumComments: int64(sqlcPost.NumComments),
		CreatedAt:   sqlcPost.CreatedAt.Time,
		UpdatedAt:   sqlcPost.UpdatedAt.Time,
	}

	if sqlcPost.PostContent != nil {
		postContent := string(sqlcPost.PostContent)
		convertedPost.Content = &postContent
	}

	if sqlcPost.DeletedAt.Valid {
		convertedPost.DeletedAt = &sqlcPost.DeletedAt.Time
	}

	return convertedPost
}

func ConvertRichPost(sqlcRow SelectPostByIDRow) domain.Post {
	convertedPost := domain.Post{
		ID:          PgUUIDToString(sqlcRow.ID),
		PublisherID: PgUUIDToString(sqlcRow.PublisherID),
		Title:       sqlcRow.Title,
		Subtitle:    sqlcRow.Subtitle,
		Slug:        sqlcRow.Slug,
		Status:      domain.PostStatus(sqlcRow.Status),
		NumRays:     int64(sqlcRow.NumRays),
		NumComments: int64(sqlcRow.NumComments),
		CreatedAt:   sqlcRow.CreatedAt.Time,
		UpdatedAt:   sqlcRow.UpdatedAt.Time,
	}

	if len(sqlcRow.PostContent) > 0 {
		postContent := string(sqlcRow.PostContent)
		convertedPost.Content = &postContent
	}

	if sqlcRow.DeletedAt.Valid {
		convertedPost.DeletedAt = &sqlcRow.DeletedAt.Time
	}

	if sqlcRow.PublisherUsername != nil {
		convertedPost.Publisher = &domain.Profile{
			ID:              PgUUIDToString(sqlcRow.PublisherID),
			FirstName:       *sqlcRow.PublisherFirstName,
			LastName:        *sqlcRow.PublisherLastName,
			Username:        *sqlcRow.PublisherUsername,
			Role:            domain.MembershipRole(*sqlcRow.PublisherRole),
			ProfileImageUrl: "",
		}
		if sqlcRow.PublisherAvatar != nil {
			convertedPost.Publisher.ProfileImageUrl = *sqlcRow.PublisherAvatar
		}
	}

	if sqlcRow.ImageID.Valid {
		convertedPost.CoverImage = &domain.PostImage{
			ID:       PgUUIDToString(sqlcRow.ImageID),
			PostID:   PgUUIDToString(sqlcRow.ID),
			ImageURL: "",
			AltText:  "",
		}
		if sqlcRow.CoverImageUrl != nil {
			convertedPost.CoverImage.ImageURL = *sqlcRow.CoverImageUrl
		}
		if sqlcRow.CoverImageAlt != nil {
			convertedPost.CoverImage.AltText = *sqlcRow.CoverImageAlt
		}
	}

	if sqlcRow.VideoID.Valid {
		convertedPost.Video = &domain.PostVideo{
			ID:             PgUUIDToString(sqlcRow.VideoID),
			PostID:         PgUUIDToString(sqlcRow.ID),
			YouTubeVideoID: "",
		}

		if sqlcRow.YoutubeVideoID != nil {
			convertedPost.Video.YouTubeVideoID = *sqlcRow.YoutubeVideoID
		}
		if len(sqlcRow.VideoMetadata) > 0 {
			metaStr := string(sqlcRow.VideoMetadata)
			convertedPost.Video.VideoMetadata = metaStr
		}
	}
	return convertedPost
}

func ConvertRichPostBySlug(sqlcRow SelectPostBySlugRow) domain.Post {
	convertedPost := domain.Post{
		ID:          PgUUIDToString(sqlcRow.ID),
		PublisherID: PgUUIDToString(sqlcRow.PublisherID),
		Title:       sqlcRow.Title,
		Subtitle:    sqlcRow.Subtitle,
		Slug:        sqlcRow.Slug,
		Status:      domain.PostStatus(sqlcRow.Status),
		NumRays:     int64(sqlcRow.NumRays),
		NumComments: int64(sqlcRow.NumComments),
		CreatedAt:   sqlcRow.CreatedAt.Time,
		UpdatedAt:   sqlcRow.UpdatedAt.Time,
	}

	if len(sqlcRow.PostContent) > 0 {
		postContent := string(sqlcRow.PostContent)
		convertedPost.Content = &postContent
	}

	if sqlcRow.DeletedAt.Valid {
		convertedPost.DeletedAt = &sqlcRow.DeletedAt.Time
	}

	if sqlcRow.PublisherUsername != nil {
		convertedPost.Publisher = &domain.Profile{
			ID:              PgUUIDToString(sqlcRow.PublisherID),
			FirstName:       *sqlcRow.PublisherFirstName,
			LastName:        *sqlcRow.PublisherLastName,
			Username:        *sqlcRow.PublisherUsername,
			Role:            domain.MembershipRole(*sqlcRow.PublisherRole),
			ProfileImageUrl: "",
		}
		if sqlcRow.PublisherAvatar != nil {
			convertedPost.Publisher.ProfileImageUrl = *sqlcRow.PublisherAvatar
		}
	}

	if sqlcRow.ImageID.Valid {
		convertedPost.CoverImage = &domain.PostImage{
			ID:       PgUUIDToString(sqlcRow.ImageID),
			PostID:   PgUUIDToString(sqlcRow.ID),
			ImageURL: "",
			AltText:  "",
		}
		if sqlcRow.CoverImageUrl != nil {
			convertedPost.CoverImage.ImageURL = *sqlcRow.CoverImageUrl
		}
		if sqlcRow.CoverImageAlt != nil {
			convertedPost.CoverImage.AltText = *sqlcRow.CoverImageAlt
		}
	}

	if sqlcRow.VideoID.Valid {
		convertedPost.Video = &domain.PostVideo{
			ID:             PgUUIDToString(sqlcRow.VideoID),
			PostID:         PgUUIDToString(sqlcRow.ID),
			YouTubeVideoID: "",
		}

		if sqlcRow.YoutubeVideoID != nil {
			convertedPost.Video.YouTubeVideoID = *sqlcRow.YoutubeVideoID
		}
		if len(sqlcRow.VideoMetadata) > 0 {
			metaStr := string(sqlcRow.VideoMetadata)
			convertedPost.Video.VideoMetadata = metaStr
		}
	}
	return convertedPost
}

func ConvertRichPostForList(sqlcRow ListPostsRow) domain.Post {
	convertedPost := domain.Post{
		ID:          PgUUIDToString(sqlcRow.ID),
		PublisherID: PgUUIDToString(sqlcRow.PublisherID),
		Title:       sqlcRow.Title,
		Subtitle:    sqlcRow.Subtitle,
		Slug:        sqlcRow.Slug,
		Status:      domain.PostStatus(sqlcRow.Status),
		NumRays:     int64(sqlcRow.NumRays),
		NumComments: int64(sqlcRow.NumComments),
		CreatedAt:   sqlcRow.CreatedAt.Time,
	}

	if len(sqlcRow.PostContent) > 0 {
		postContent := string(sqlcRow.PostContent)
		convertedPost.Content = &postContent
	}

	if sqlcRow.PublisherUsername != nil {
		convertedPost.Publisher = &domain.Profile{
			Username:        *sqlcRow.PublisherUsername,
			ProfileImageUrl: "",
		}
		if sqlcRow.PublisherAvatar != nil {
			convertedPost.Publisher.ProfileImageUrl = *sqlcRow.PublisherAvatar
		}
	}

	if sqlcRow.CoverImageUrl != nil {
		convertedPost.CoverImage = &domain.PostImage{
			PostID:   PgUUIDToString(sqlcRow.ID), // Link it back to the parent post
			ImageURL: *sqlcRow.CoverImageUrl,
			AltText:  "", // Default empty
		}
		if sqlcRow.CoverImageAlt != nil {
			convertedPost.CoverImage.AltText = *sqlcRow.CoverImageAlt
		}
	}

	return convertedPost
}
