package usecase

import (
	"context"
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"
	"isling-be/pkg/logger"
	"strings"

	"github.com/gosimple/slug"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type (
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

func (uc *RoomUC) GetManyRoomOfUser(c context.Context, accountID common_entity.AccountID) (*common_entity.Collection[*entity.Room], error) {
	filter := FindRoomFilter{
		OwnerID: &accountID,
	}

	return uc.roomRepo.FindMany(c, &filter)
}

func (uc *RoomUC) GetRoom(c context.Context, currentUserID common_entity.AccountID, slugName string) (*entity.Room, error) {
	room, err := uc.roomRepo.FindOneBySlug(c, slugName)
	if err != nil {
		return nil, err
	}

	if room.OwnerID != currentUserID {
		room.InviteCode = "******"
	}

	return room, nil
}

func (uc *RoomUC) UpdateRoom(c context.Context, currentUserID common_entity.AccountID, id int64, req *UpdateRoomRequest) (*entity.Room, error) {
	room, err := uc.roomRepo.FindOne(c, id)
	if err != nil {
		return nil, err
	}

	if room.OwnerID != currentUserID {
		return nil, entity.ErrMissingUpdatePerm
	}

	newRoom := UpdateRoomFromReq(room, req)

	if newRoom.Name != room.Name {
		slugPieces := strings.Split(room.Slug, "-")

		if len(slugPieces) == 0 {
			return nil, entity.ErrInvalidRoomSlug
		}

		slugID := slugPieces[len(slugPieces)-1]

		newRoom.Slug = slug.Make(newRoom.Name) + "-" + slugID
	}

	return uc.roomRepo.UpdateOne(c, newRoom)
}

func (uc *RoomUC) DeleteRoom(c context.Context, currentUserID common_entity.AccountID, id int64) error {
	room, err := uc.roomRepo.FindOne(c, id)
	if err != nil {
		return err
	}

	if room.OwnerID != currentUserID {
		return entity.ErrMissingDeletePerm
	}

	return uc.roomRepo.DeleteOne(c, id)
}
