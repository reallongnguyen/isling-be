// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase/request"
	common_entity "isling-be/internal/common/entity"

	"github.com/jackc/pgx/v4"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	AccountUsecase interface {
		GetAccountByID(context.Context, common_entity.AccountID) (*entity.Account, error)
		CreateAccount(context.Context, request.CreateAccountReq) (*entity.Account, error)
		ChangePassword(context.Context, common_entity.AccountID, *request.ChangePasswordReq) error
	}

	AuthUsecase interface {
		GetTokenByPassword(context.Context, *request.GetTokenByPasswordRequest) (*request.GetTokenResponse, error)
		GetTokenByRefreshToken(context.Context, *request.GetTokenByRefreshTokenRequest) (*request.GetTokenResponse, error)
		Logout(context.Context, common_entity.AccountID, string) error
		SignUp(context.Context, request.CreateAccountReq) (*request.GetTokenResponse, error)
	}

	ProfileUsecase interface {
		GetProfile(context.Context, common_entity.AccountID) (*entity.Profile, error)
		UpsertProfile(context.Context, common_entity.AccountID, *request.CreateProfileReq) (*entity.Profile, error)
	}

	AccountRepository interface {
		FindByUsername(context.Context, string) (*entity.Account, error)
		FindByID(context.Context, common_entity.AccountID) (*entity.Account, error)
		Store(context.Context, *entity.Account) (*entity.Account, error)
		UpdateEncryptedPassword(context.Context, common_entity.AccountID, string) error
	}

	// TODO: hide postgres transaction logic in interface
	RefreshTokenRepository interface {
		Store(context.Context, pgx.Tx, *entity.RefreshTokens) (*entity.RefreshTokens, error)
		FindOneByEncryptedToken(context.Context, pgx.Tx, string) (*entity.RefreshTokens, error)
		RevokeManyByAccountID(context.Context, pgx.Tx, common_entity.AccountID) (int64, error)
		RevokeOneByEncryptedToken(context.Context, pgx.Tx, string) (int64, error)
		BeginTx(context.Context) (pgx.Tx, error)
	}

	ProfileRepository interface {
		FindOneProfileByID(context.Context, common_entity.AccountID) (*entity.Profile, error)
		UpsertProfile(context.Context, common_entity.AccountID, *request.CreateProfileReq) (*entity.Profile, error)
	}
)
