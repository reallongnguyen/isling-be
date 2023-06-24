package usecase

import (
	"context"
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"
	"isling-be/pkg/logger"

	"github.com/gosimple/slug"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type (
	CreateRoomRequest struct {
		Name        string
		Description string
		Cover       string
	}

	RoomUsecase interface {
		CreateRoom(c context.Context, accountID common_entity.AccountID, req *CreateRoomRequest) (*entity.Room, error)
	}

	RoomRepository interface {
		Create(c context.Context, room *entity.Room) (*entity.Room, error)
	}

	RoomUC struct {
		log      logger.Interface
		roomRepo RoomRepository
	}
)

const (
	base56Characters = "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz"
	slugIDLen        = 8
	inviteCodeLen    = 8
)

var _ RoomUsecase = (*RoomUC)(nil)

func NewRoomUsecase(log logger.Interface, roomRepo RoomRepository) *RoomUC {
	return &RoomUC{
		log:      log,
		roomRepo: roomRepo,
	}
}

func (uc *RoomUC) CreateRoom(c context.Context, accountID common_entity.AccountID, req *CreateRoomRequest) (*entity.Room, error) {
	slugID, err := gonanoid.Generate(base56Characters, slugIDLen)
	if err != nil {
		return nil, err
	}

	slugName := slug.Make(req.Name) + "-" + slugID

	inviteCode, err := gonanoid.Generate(base56Characters, inviteCodeLen)
	if err != nil {
		return nil, err
	}

	room := &entity.Room{
		OwnerID:     accountID,
		Visibility:  entity.VisibilityPublic,
		InviteCode:  inviteCode,
		Name:        req.Name,
		Slug:        slugName,
		Description: req.Description,
		Cover:       req.Cover,
	}

	newRoom, err := uc.roomRepo.Create(c, room)
	if err != nil {
		return nil, err
	}

	return newRoom, nil
}
