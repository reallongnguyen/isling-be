package dto

import (
	"github.com/btcs-longnp/isling-be/internal/account/usecase/request"
)

type CreateAccountDto struct {
	Username string `json:"username" validate:"required,email" example:"isling@isling.me"`
	Password string `json:"password" validate:"required,min=8,max=255" example:"wakaranai"`
}

func (dto *CreateAccountDto) ToCreateAccountRequest() request.CreateAccountReq {
	return request.CreateAccountReq{
		Username: dto.Username,
		Password: dto.Password,
	}
}
