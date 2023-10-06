package pg_surreal_sync

import "time"

type Payload struct {
	ID    int64  `json:"id"`
	Table string `json:"table"`
	Type  string `json:"type"`
	Data  string `json:"data"`
}

type Profile struct {
	AccountID   int64     `json:"account_id"`
	FirstName   string    `json:"first_name"`
	LastName    *string   `json:"last_name,omitempty"`
	Gender      *string   `json:"gender,omitempty"`
	DateOfBirth *string   `json:"date_of_birth,omitempty"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Room struct {
	ID            int64     `json:"id"`
	OwnerID       int64     `json:"owner_id"`
	Visibility    string    `json:"visibility"`
	InviteCode    string    `json:"invite_code"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Description   string    `json:"description"`
	Cover         string    `json:"cover"`
	AudienceCount int       `json:"audience_count"`
	Audiences     []int64   `json:"audiences"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type SRUser struct {
	ID          string  `json:"id,omitempty"`
	FirstName   string  `json:"firstName"`
	LastName    *string `json:"lastName,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	DateOfBirth *string `json:"dateOfBirth,omitempty"`
	AvatarURL   *string `json:"avatarUrl,omitempty"`
}

type SRRoom struct {
	ID          string `json:"id"`
	OwnerID     string `json:"ownerID"`
	Visibility  string `json:"visibility"`
	InviteCode  string `json:"inviteCode"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Cover       string `json:"cover"`
}
