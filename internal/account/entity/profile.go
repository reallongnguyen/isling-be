package entity

import (
	"isling-be/internal/common/entity"
	"time"
)

type GenderIdentity string

const (
	GenderMale    GenderIdentity = "male"
	GenderFemale  GenderIdentity = "female"
	GenderOther   GenderIdentity = "other"
	GenderUnknown GenderIdentity = "unknown"
)

type Profile struct {
	AccountID   entity.AccountID `json:"accountId"`
	Email       string           `json:"email"`
	FirstName   *string          `json:"firstName,omitempty"`
	LastName    *string          `json:"lastName,omitempty"`
	Gender      *GenderIdentity  `json:"gender,omitempty"`
	DateOfBirth *time.Time       `json:"dateOfBirth,omitempty"`
	AvatarURL   *string          `json:"avatarUrl"`
	CreatedAt   *time.Time       `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time       `json:"updatedAt,omitempty"`
}
