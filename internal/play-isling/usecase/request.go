package usecase

import (
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"
	"time"
)

type CreateRoomRequest struct {
	Name        string
	Description string
	Cover       string
}

type UpdateRoomRequest struct {
	Name        string
	Description string
	Cover       string
}

func UpdateRoomFromReq(room *entity.Room, req *UpdateRoomRequest) *entity.Room {
	newRoom := *room

	newRoom.Name = req.Name
	newRoom.Description = req.Description
	newRoom.Cover = req.Cover

	return &newRoom
}

type HomePageResponse struct {
	Collections []*entity.RoomCollection `json:"collections"`
}

type CreateActionRequest struct {
	AccountID common_entity.AccountID
	Type      string
	ObjectID  *string
	Timestamp time.Time
}
