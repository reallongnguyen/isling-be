package usecase

import (
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"
)

type (
	CreateRoomRequest struct {
		Name        string
		Description string
		Cover       string
	}

	FindRoomFilter struct {
		OwnerID *common_entity.AccountID
	}

	UpdateRoomRequest struct {
		Name        string
		Description string
		Cover       string
	}
)

func UpdateRoomFromReq(room *entity.Room, req *UpdateRoomRequest) *entity.Room {
	newRoom := *room

	newRoom.Name = req.Name
	newRoom.Description = req.Description
	newRoom.Cover = req.Cover

	return &newRoom
}
