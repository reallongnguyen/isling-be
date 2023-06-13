package usecase

import (
	"context"
	"errors"
	"fmt"
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

var _ AccountUsecase = (*AccountUC)(nil)

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
	return uc.repo.FindByID(ctx, accountID)
}

func (uc *AccountUC) ChangePassword(ctx context.Context, accountID common_entity.AccountID, changePassReq *request.ChangePasswordReq) error {
	account, err := uc.repo.FindByID(ctx, accountID)

	if err != nil && errors.Is(err, common_entity.ErrNoRows) {
		uc.log.Warn("account usecase: change password: not found account ID: %s", accountID)

		return common_entity.ErrAccountNotFound
	}

	if err != nil {
		uc.log.Info("account usecase: change password: find one account: %w", err)

		return fmt.Errorf("account usecase: change password: find an account %w", err)
	}

	if !common_uc.IsMatchHashAndPassword(account.EncryptedPassword, changePassReq.OldPassword) {
		uc.log.Warn("account usecase: change password: password not correct. Account ID: %s", accountID)

		return common_entity.ErrPasswordNotCorrect
	}

	newEncryptedPassword, err := common_uc.HashPassword(changePassReq.NewPassword)
	if err != nil {
		return fmt.Errorf("account usecase: change password: hash new password: %w", err)
	}

	if err := uc.repo.UpdateEncryptedPassword(ctx, accountID, newEncryptedPassword); err != nil {
		uc.log.Info("account usecase: change password: update encrypted password: %w", err)

		return err
	}

	return nil
}
