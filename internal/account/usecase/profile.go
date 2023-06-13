package usecase

import (
	"context"

	"isling-be/internal/account/entity"
	common_entity "isling-be/internal/common/entity"
	"isling-be/pkg/logger"
)

type ProfileUC struct {
	repo ProfileRepository
	log  logger.Interface
}

var _ AccountUsecase = (*AccountUC)(nil)

func NewProfileUC(repo ProfileRepository, log logger.Interface) ProfileUsecase {
	return &ProfileUC{
		repo: repo,
		log:  log,
	}
}

func (uc *ProfileUC) GetProfile(ctx context.Context, accountID common_entity.AccountID) (*entity.Profile, error) {
	return uc.repo.FindOneProfileByID(ctx, accountID)
}
