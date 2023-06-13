package dto

import (
	"isling-be/internal/account/usecase/request"
)

type CreateAccountDto struct {
	Email    string `json:"email" validate:"required,email" example:"isling@isling.me"`
	Password string `json:"password" validate:"required,min=8,max=255,printascii" example:"wakaranai"`
}

func (dto *CreateAccountDto) ToCreateAccountRequest() request.CreateAccountReq {
	return request.CreateAccountReq{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

type ChangePasswordDto struct {
	OldPassword string `json:"oldPassword" validate:"required" example:"himitsu"`
	NewPassword string `json:"newPassword" validate:"required,min=8,max=255,printascii" example:"wakaranai"`
}

func (dto *ChangePasswordDto) ToChangePasswordRequest() *request.ChangePasswordReq {
	return &request.ChangePasswordReq{
		OldPassword: dto.OldPassword,
		NewPassword: dto.NewPassword,
	}
}
