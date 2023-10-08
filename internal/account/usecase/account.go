package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase/request"
	common_entity "isling-be/internal/common/entity"
	common_uc "isling-be/internal/common/usecase"
	"isling-be/pkg/facade"
)

type AccountUC struct {
	repo AccountRepository
}

var _ AccountUsecase = (*AccountUC)(nil)

func NewAccountUC(repo AccountRepository) AccountUsecase {
	return &AccountUC{
		repo: repo,
	}
}

func (uc *AccountUC) CreateAccount(ctx context.Context, createAccountReq request.CreateAccountReq) (*entity.Account, error) {
	createAccountReq.Normalize()

	// check user exist
	usernameAvailable, err := uc.checkUsernameAvailable(ctx, createAccountReq.Email)
	if err != nil {
		facade.Log().Error("got error when checkUsernameAvailable" + err.Error())

		return nil, err
	}

	if !usernameAvailable {
		facade.Log().Info("username " + createAccountReq.Email + " already registered")

		return nil, common_entity.ErrEmailDuplicated
	}

	hashedPassword, err := common_uc.HashPassword(createAccountReq.Password)
	if err != nil {
		facade.Log().Info("hash password got error:" + err.Error())

		return nil, err
	}

	user := entity.NewAccount(0, createAccountReq.Email, hashedPassword, time.Now(), time.Now())

	account, err := uc.repo.Store(ctx, &user)
	if err != nil {
		return nil, err
	}

	go func() {
		accJSON, err := json.Marshal(account)
		if err != nil {
			return
		}

		facade.MsgBus().Publish("account.created", accJSON, nil)
	}()

	return account, nil
}

func (uc *AccountUC) checkUsernameAvailable(ctx context.Context, username string) (bool, error) {
	_, err := uc.repo.FindByUsername(ctx, username)

	if errors.Is(err, common_entity.ErrAccountNotFound) {
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
		facade.Log().Warn("account usecase: change password: not found account ID: %s", accountID)

		return common_entity.ErrAccountNotFound
	}

	if err != nil {
		facade.Log().Info("account usecase: change password: find one account: %w", err)

		return fmt.Errorf("account usecase: change password: find an account %w", err)
	}

	if !common_uc.IsMatchHashAndPassword(account.EncryptedPassword, changePassReq.OldPassword) {
		facade.Log().Warn("account usecase: change password: password not correct. Account ID: %s", accountID)

		return common_entity.ErrPasswordNotCorrect
	}

	newEncryptedPassword, err := common_uc.HashPassword(changePassReq.NewPassword)
	if err != nil {
		return fmt.Errorf("account usecase: change password: hash new password: %w", err)
	}

	if err := uc.repo.UpdateEncryptedPassword(ctx, accountID, newEncryptedPassword); err != nil {
		facade.Log().Info("account usecase: change password: update encrypted password: %w", err)

		return err
	}

	return nil
}
