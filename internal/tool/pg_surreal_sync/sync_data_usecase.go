package pg_surreal_sync

import (
	"encoding/json"
	"fmt"
	"isling-be/pkg/facade"
	"isling-be/pkg/surreal"
	"strings"
)

type SyncDataUsecase struct {
	sr *surreal.Surreal
}

func NewSyncDataUsecase(sr *surreal.Surreal) *SyncDataUsecase {
	return &SyncDataUsecase{
		sr: sr,
	}
}

func (r *SyncDataUsecase) Handle(payload *Payload) error {
	facade.Log().Debug("pg_surreal_sync: receive msg: %w", *payload)

	switch payload.Table {
	case "profiles":
		profile := new(Profile)
		if err := json.Unmarshal([]byte(payload.Data), profile); err != nil {
			facade.Log().Error("pg_surreal_sync: parse profile mess: %w", err)

			return err
		}

		userID := fmt.Sprintf("users:%d", profile.AccountID)

		if payload.Type == "DELETE" {
			_, err := r.sr.Delete(userID)
			if err != nil {
				facade.Log().Error("SyncDataUsecase: delete user: %w", err)
			}

			break
		}

		user := SRUser{
			ID:          userID,
			FirstName:   profile.FirstName,
			LastName:    profile.LastName,
			Gender:      profile.Gender,
			DateOfBirth: profile.DateOfBirth,
			AvatarURL:   profile.AvatarURL,
		}

		_, err := r.sr.Update(user.ID, user)
		if err != nil {
			facade.Log().Error("SyncDataUsecase: update user: %w", err)
		}
	case "media_rooms":
		room := new(Room)
		if err := json.Unmarshal([]byte(payload.Data), room); err != nil {
			facade.Log().Error("pg_surreal_sync: parse media_rooms mess: %w", err)

			return err
		}

		slugPieces := strings.Split(room.Slug, "-")
		uid := slugPieces[len(slugPieces)-1]
		srRoomID := fmt.Sprintf("media_rooms:%s", uid)

		if payload.Type == "DELETE" {
			_, err := r.sr.Delete(srRoomID)
			if err != nil {
				facade.Log().Error("SyncDataUsecase: delete media_rooms: %w", err)
			}

			break
		}

		srRoom := SRRoom{
			ID:          srRoomID,
			OwnerID:     fmt.Sprintf("users:%d", room.OwnerID),
			Visibility:  room.Visibility,
			InviteCode:  room.InviteCode,
			Name:        room.Name,
			Slug:        room.Slug,
			Description: room.Description,
			Cover:       room.Cover,
		}

		_, err := r.sr.Update(srRoom.ID, srRoom)
		if err != nil {
			facade.Log().Error("SyncDataUsecase: update room: %w", err)
		}
	default:
	}

	return nil
}
