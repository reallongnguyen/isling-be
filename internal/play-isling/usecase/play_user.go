package usecase

import (
	"context"
	common_entity "isling-be/internal/common/entity"

	"golang.org/x/exp/slices"
)

type PlayUser struct {
	repo PlayUserRepository
}

var _ PlayUserUsecase = (*PlayUser)(nil)

func NewPlayUserUC(repo PlayUserRepository) PlayUserUsecase {
	return &PlayUser{
		repo: repo,
	}
}

func (uc *PlayUser) InsertRecentlyJoinedRoom(c context.Context, accountID common_entity.AccountID, roomID int64) error {
	playUser, err := uc.repo.GetOne(c, accountID)
	if err != nil {
		return err
	}

	recentlyJoinedRooms := playUser.RecentlyJoinedRooms
	recentlyJoinedRooms = slices.DeleteFunc(recentlyJoinedRooms, func(a int64) bool {
		return a == roomID
	})
	recentlyJoinedRooms = append([]int64{roomID}, recentlyJoinedRooms...)

	playUser.RecentlyJoinedRooms = recentlyJoinedRooms

	if len(recentlyJoinedRooms) > 8 {
		playUser.RecentlyJoinedRooms = recentlyJoinedRooms[:8]
	}

	return uc.repo.Update(c, accountID, playUser)
}
