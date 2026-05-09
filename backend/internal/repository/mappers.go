package repository

import (
	"fmt"

	"github.com/github.com/bmstoss13/the-daily-sunshine/internal/domain"
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

// func ConvertProfileToSQL(profile domain.Profile) (Profile, error) {
// 	profileId, err := StringToPgUUID(profile.ID)
// 	if err != nil {
// 		return Profile{}, fmt.Errorf("[mappers.go] ConvertProfileToSQL: an error occurred converting profile id to pg uuid: %w", err)
// 	}
// 	convertedProfile := Profile{
// 		ID:              profileId,
// 		FirstName:       profile.FirstName,
// 		LastName:        profile.LastName,
// 		Username:        profile.Username,
// 		NumRaysReceived: int32(profile.NumRays),
// 		Role:            string(profile.Role),

// 		CreatedAt: pgtype.Timestamptz{profile.CreatedAt, pgtype.Infinity, true},
// 		UpdatedAt: pgtype.Timestamptz{profile.UpdatedAt, pgtype.Infinity, true},
// 	}

// 	if profile.ProfileImageUrl != "" {
// 		convertedProfile.ProfileImageUrl = &profile.ProfileImageUrl
// 	} else {
// 		convertedProfile.ProfileImageUrl = nil // Or a default URL string
// 	}

// 	if profile.Bio != nil {
// 		bioStr := *profile.Bio
// 		convertedProfile.ProfileBio = []byte(bioStr)
// 	}

// 	if profile.DeletedAt != nil {
// 		convertedProfile.DeletedAt = pgtype.Timestamptz{*profile.DeletedAt, pgtype.Infinity, true}
// 	}

// 	return convertedProfile, nil
// }
