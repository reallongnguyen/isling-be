package entity

import (
	common_entity "isling-be/internal/common/entity"
	"time"
)

type VisibilityType string

const (
	VisibilityPublic VisibilityType = "public"
	VisibilityMember VisibilityType = "member"
)

type RoomOwner struct {
	ID        common_entity.AccountID `json:"id"`
	FirstName *string                 `json:"firstName,omitempty"`
	LastName  *string                 `json:"lastName,omitempty"`
	AvatarURL *string                 `json:"avatarUrl,omitempty"`
}

type Room struct {
	ID            int64                     `json:"id"`
	OwnerID       common_entity.AccountID   `json:"ownerId,omitempty"`
	Owner         *RoomOwner                `json:"owner,omitempty"`
	Visibility    VisibilityType            `json:"visibility"`
	InviteCode    string                    `json:"inviteCode"`
	Name          string                    `json:"name"`
	Slug          string                    `json:"slug"`
	Description   string                    `json:"description"`
	Cover         string                    `json:"cover"`
	AudienceCount int                       `json:"audienceCount"`
	Audiences     []common_entity.AccountID `json:"audiences"`
	CreatedAt     time.Time                 `json:"createdAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
}

func (room *Room) ToRoomPublic() *RoomPublic {
	return &RoomPublic{
		ID:            room.ID,
		OwnerID:       room.OwnerID,
		Owner:         room.Owner,
		Visibility:    room.Visibility,
		Name:          room.Name,
		Slug:          room.Slug,
		Description:   room.Description,
		Cover:         room.Cover,
		AudienceCount: room.AudienceCount,
		Audiences:     room.Audiences,
		CreatedAt:     room.CreatedAt,
		UpdatedAt:     room.UpdatedAt,
	}
}

type RoomPublic struct {
	ID            int64                     `json:"id"`
	OwnerID       common_entity.AccountID   `json:"ownerId"`
	Owner         *RoomOwner                `json:"owner,omitempty"`
	Visibility    VisibilityType            `json:"visibility"`
	Name          string                    `json:"name"`
	Slug          string                    `json:"slug"`
	Description   string                    `json:"description"`
	Cover         string                    `json:"cover"`
	AudienceCount int                       `json:"audienceCount"`
	Audiences     []common_entity.AccountID `json:"audiences"`
	CreatedAt     time.Time                 `json:"createdAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
}
