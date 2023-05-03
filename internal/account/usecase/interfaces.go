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
		GetAccountByID(context.Context, common_entity.AccountId) (*entity.AccountWithoutPass, error)
		CreateAccount(context.Context, request.CreateAccountReq) (*entity.AccountWithoutPass, error)
	}

	AccountRepository interface {
		FindByUsername(context.Context, string) (*entity.Account, error)
		FindByID(context.Context, common_entity.AccountId) (*entity.Account, error)
		Store(context.Context, *entity.Account) (*entity.Account, error)
	}
)
