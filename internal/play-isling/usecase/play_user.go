package usecase

import (
	"context"
	common_entity "isling-be/internal/common/entity"
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
	return uc.repo.InsertRecentlyJoinedRoom(c, accountID, roomID)
}
