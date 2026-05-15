package domain

import "time"

type Profile struct {
	ID              string         `json:"id"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Username        string         `json:"username"`
	ProfileImageUrl string         `json:"profile_image_url"`
	Bio             *string        `json:"profile_bio"`
	NumRays         int            `json:"num_rays_received"`
	SubscriberTier  SubscriberTier `json:"subscriber_tier"`
	AppRole         AppRole        `json:"app_role"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       *time.Time     `json:"deleted_at"`
}

type SubscriberTier string

const (
	Member     SubscriberTier = "Member"
	Subscriber SubscriberTier = "Subscriber"
)

type AppRole string

const (
	Admin AppRole = "Admin"
	User  AppRole = "User"
)
