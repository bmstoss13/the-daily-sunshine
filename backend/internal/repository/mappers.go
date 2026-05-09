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
		convertedProfile.ProfileImage = *sqlcProfile.ProfileImageUrl
	} else {
		convertedProfile.ProfileImage = "" // Or a default URL string
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
