package usecase

import (
	"context"
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"
)

type (
	RoomUsecase interface {
		CreateRoom(c context.Context, accountID common_entity.AccountID, req *CreateRoomRequest) (*entity.Room, error)
		GetManyRoomOfUser(c context.Context, accountID common_entity.AccountID) (*common_entity.Collection[*entity.Room], error)
		GetRoom(c context.Context, currentUserID common_entity.AccountID, slugName string) (*entity.Room, error)
		UpdateRoom(c context.Context, currentUserID common_entity.AccountID, id int64, req *UpdateRoomRequest) (*entity.Room, error)
		DeleteRoom(c context.Context, currentUserID common_entity.AccountID, id int64) error
	}

	RoomRepository interface {
		Create(c context.Context, room *entity.Room) (*entity.Room, error)
		FindMany(c context.Context, filter *FindRoomFilter) (*common_entity.Collection[*entity.Room], error)
		FindOne(c context.Context, id int64) (*entity.Room, error)
		FindOneBySlug(c context.Context, slug string) (*entity.Room, error)
		UpdateOne(c context.Context, room *entity.Room) (*entity.Room, error)
		DeleteOne(c context.Context, id int64) error
	}
)
