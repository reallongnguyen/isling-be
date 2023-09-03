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

type Room struct {
	ID            int64                     `json:"id"`
	OwnerID       common_entity.AccountID   `json:"ownerId"`
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
