package domain

import "time"

type Profile struct {
	ID           string         `json:"id"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Username     string         `json:"username"`
	ProfileImage string         `json:"profile_image_url"`
	Bio          string         `json:"profile_bio"`
	NumRays      string         `json:"num_rays_received"`
	Role         MembershipRole `json:"role"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    time.Time      `json:"deleted_at"`
}

type MembershipRole string

const (
	Member     MembershipRole = "Member"
	Subscriber MembershipRole = "Subscriber"
	Admin      MembershipRole = "Admin"
)
