package usecase

import (
	"context"
	"errors"
	"time"

	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase/request"
	common_entity "isling-be/internal/common/entity"
	common_uc "isling-be/internal/common/usecase"
	"isling-be/pkg/logger"
)

type AccountUC struct {
	repo AccountRepository
	log  logger.Interface
}

func NewAccountUC(repo AccountRepository, log logger.Interface) AccountUsecase {
	return &AccountUC{
		repo: repo,
		log:  log,
	}
}

func (uc *AccountUC) CreateAccount(ctx context.Context, createUserDto request.CreateAccountReq) (*entity.Account, error) {
	// check user exist
	usernameAvailable, err := uc.checkUsernameAvailable(ctx, createUserDto.Email)
	if err != nil {
		uc.log.Error("got error when checkUsernameAvailable" + err.Error())

		return nil, err
	}

	if !usernameAvailable {
		uc.log.Info("username " + createUserDto.Email + " already registered")

		return nil, common_entity.ErrDuplicated
	}

	hashedPassword, err := common_uc.HashPassword(createUserDto.Password)
	if err != nil {
		uc.log.Info("hash password got error:" + err.Error())

		return nil, err
	}

	user := entity.NewAccount(0, createUserDto.Email, hashedPassword, time.Now(), time.Now())

	account, err := uc.repo.Store(ctx, &user)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (uc *AccountUC) checkUsernameAvailable(ctx context.Context, username string) (bool, error) {
	_, err := uc.repo.FindByUsername(ctx, username)

	if errors.Is(err, common_entity.ErrNoRows) {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	return false, nil
}

func (uc *AccountUC) GetAccountByID(ctx context.Context, accountID common_entity.AccountID) (*entity.Account, error) {
	account, err := uc.repo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return account, nil
}
