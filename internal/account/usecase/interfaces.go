// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase/request"
	common_entity "isling-be/internal/common/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	AccountUsecase interface {
		GetAccountByID(context.Context, common_entity.AccountID) (*entity.AccountWithoutPass, error)
		CreateAccount(context.Context, request.CreateAccountReq) (*entity.AccountWithoutPass, error)
	}

	AuthUsecase interface {
		GetTokenByPassword(context.Context, *request.GetTokenByPasswordRequest) (*request.GetTokenResponse, error)
		Logout(context.Context, common_entity.AccountID, string) error
	}

	AccountRepository interface {
		FindByUsername(context.Context, string) (*entity.Account, error)
		FindByID(context.Context, common_entity.AccountID) (*entity.Account, error)
		Store(context.Context, *entity.Account) (*entity.Account, error)
	}

	RefreshTokenRepository interface {
		Store(context.Context, *entity.RefreshTokens) (*entity.RefreshTokens, error)
		RevokeManyByAccountID(context.Context, common_entity.AccountID) (int64, error)
		RevokeOneByEncryptedToken(context.Context, string) (int64, error)
	}
)
