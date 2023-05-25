package dto

import (
	"isling-be/internal/account/usecase/request"
)

type CreateAccountDto struct {
	Email    string `json:"email" validate:"required,email" example:"isling@isling.me"`
	Password string `json:"password" validate:"required,min=8,max=255" example:"wakaranai"`
}

func (dto *CreateAccountDto) ToCreateAccountRequest() request.CreateAccountReq {
	return request.CreateAccountReq{
		Email:    dto.Email,
		Password: dto.Password,
	}
}
