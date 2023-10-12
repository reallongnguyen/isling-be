package entity

import (
	common_entity "isling-be/internal/common/entity"
	"strconv"
	"strings"
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

type RoomOwnerSurreal struct {
	RoomOwner
	ID       string `json:"id"`
	FullName string `json:"fullName,omitempty"`
}

type RoomSurreal struct {
	Room
	ID         string            `json:"id"`
	OriginalID int64             `json:"originalID"`
	OwnerID    string            `json:"ownerID,omitempty"`
	Owner      *RoomOwnerSurreal `json:"owner,omitempty"`
}

func (r *RoomSurreal) ToRoom() *Room {
	ownerID, _ := strconv.Atoi(strings.Split(r.OwnerID, ":")[1])

	var roomOwner *RoomOwner

	if r.Owner != nil {
		roomOwner = &RoomOwner{
			ID:        common_entity.AccountID(ownerID),
			FirstName: r.Owner.FirstName,
			LastName:  r.Owner.LastName,
			AvatarURL: r.Owner.AvatarURL,
		}
	}

	return &Room{
		ID:            r.OriginalID,
		OwnerID:       common_entity.AccountID(ownerID),
		Owner:         roomOwner,
		Visibility:    r.Visibility,
		InviteCode:    r.InviteCode,
		Name:          r.Name,
		Slug:          r.Slug,
		Description:   r.Description,
		Cover:         r.Cover,
		AudienceCount: r.AudienceCount,
		Audiences:     r.Audiences,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}
