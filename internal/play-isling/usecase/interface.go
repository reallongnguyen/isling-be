package usecase

import (
	"context"
	account_entity "isling-be/internal/account/entity"
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"

	"github.com/zhenghaoz/gorse/client"
)

type (
	RoomUsecase interface {
		CreateRoom(c context.Context, accountID common_entity.AccountID, req *CreateRoomRequest) (*entity.Room, error)
		GetManyRoomOfUser(c context.Context, accountID common_entity.AccountID) (*common_entity.Collection[*entity.Room], error)
		GetRoom(c context.Context, currentUserID common_entity.AccountID, slugName string) (*entity.Room, error)
		UpdateRoom(c context.Context, currentUserID common_entity.AccountID, id int64, req *UpdateRoomRequest) (*entity.Room, error)
		DeleteRoom(c context.Context, currentUserID common_entity.AccountID, id int64) error
	}

	FindRoomFilter struct {
		OwnerID *common_entity.AccountID
		IDIn    *[]int64
	}

	Order struct {
		Field     string
		Direction string
	}

	RoomRepository interface {
		Create(c context.Context, room *entity.Room) (*entity.Room, error)
		FindMany(c context.Context, filter *FindRoomFilter, order *Order) (*common_entity.Collection[*entity.Room], error)
		FindOne(c context.Context, id int64) (*entity.Room, error)
		FindOneBySlug(c context.Context, slug string) (*entity.Room, error)
		UpdateOne(c context.Context, room *entity.Room) (*entity.Room, error)
		DeleteOne(c context.Context, id int64) error
	}

	HomeUsecase interface {
		Show(c context.Context, accountID common_entity.AccountID) (*HomePageResponse, error)
		ShowGuest(c context.Context) (*HomePageResponse, error)
	}

	PlayUserRepository interface {
		GetOne(c context.Context, accountID common_entity.AccountID) (*entity.PlayUser, error)
		Create(c context.Context, accountID common_entity.AccountID) (*entity.PlayUser, error)
		Update(c context.Context, accountID common_entity.AccountID, playUser *entity.PlayUser) error
		InsertRecentlyJoinedRoom(c context.Context, accountID common_entity.AccountID, roomID int64) error
	}

	PlayUserUsecase interface {
		InsertRecentlyJoinedRoom(c context.Context, accountID common_entity.AccountID, roomID int64) error
	}

	RecommendationUsecase interface {
		InsertUser(c context.Context, account *account_entity.Account) error
		InsertRoom(c context.Context, room *entity.Room) error
		HideItem(c context.Context, itemID string) error
		InsertFeedback(c context.Context, feedback client.Feedback) error
	}

	SearchRepository interface {
		SearchRoom(c context.Context, accountID common_entity.AccountID, req *SearchRequest) ([]entity.RoomSearchResult, error)
		GetTotalRoomMatches(c context.Context, accountID common_entity.AccountID, req *SearchRequest) (int, error)
	}

	SearchUsecase interface {
		Search(c context.Context, accountID common_entity.AccountID, req *SearchRequest) (*SearchResponse, error)
	}
)
