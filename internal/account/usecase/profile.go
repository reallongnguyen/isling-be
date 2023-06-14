package usecase

import (
	"context"

	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase/request"
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

func (uc *ProfileUC) UpsertProfile(ctx context.Context, accountID common_entity.AccountID, createProfileReq *request.CreateProfileReq) (*entity.Profile, error) {
	createProfileReq.Normalize()

	return uc.repo.UpsertProfile(ctx, accountID, createProfileReq)
}
