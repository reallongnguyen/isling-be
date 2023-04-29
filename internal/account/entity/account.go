package entity

import (
	"time"

	common_entity "github.com/btcs-longnp/isling-be/internal/common/entity"
)

type Account struct {
	Id        common_entity.AccountId `json:"id" example:"1"`
	Username  string                  `json:"username" example:"isling@isling.me"`
	Password  string                  `json:"password" example:"wakaranai"`
	CreatedAt time.Time               `json:"createdAt" example:"2022-12-12T12:12:12"`
	UpdatedAt time.Time               `json:"updatedAt" example:"2022-12-12T12:12:12"`
}

func NewAccount(
	id common_entity.AccountId,
	username,
	password string,
	createdAt, updatedAt time.Time,
) Account {
	return Account{
		Id:        id,
		Username:  username,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type CreateAccountDto struct {
	Username string `json:"username" validate:"required,email" example:"isling@isling.me"`
	Password string `json:"password" validate:"required,min=8,max=255" example:"wakaranai"`
}

type AccountWithoutPass struct {
	Id        common_entity.AccountId `json:"id" example:"1"`
	Username  string                  `json:"username" example:"isling@isling.me"`
	CreatedAt time.Time               `json:"createdAt" example:"2022-12-12T12:12:12"`
	UpdatedAt time.Time               `json:"updatedAt" example:"2022-12-12T12:12:12"`
}

func (a *Account) ToAccountWithoutPass() *AccountWithoutPass {
	return &AccountWithoutPass{
		Id:        a.Id,
		Username:  a.Username,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
