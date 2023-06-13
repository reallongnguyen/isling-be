package dto

import (
	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase/request"
)

type (
	CreateProfileReqDTO struct {
		FirstName   string                `json:"firstName" validate:"required,max=64,alphaunicodespace"`
		LastName    string                `json:"lastName" validate:"required,max=64,alphaunicodespace"`
		Gender      entity.GenderIdentity `json:"gender" validate:"required,oneof=male female other unknown"`
		DateOfBirth string                `json:"dateOfBirth" validate:"required,datetime=2006-01-02,beforenow"`
	}
)

func (dto *CreateProfileReqDTO) ToRequest() *request.CreateProfileReq {
	return &request.CreateProfileReq{
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Gender:      dto.Gender,
		DateOfBirth: dto.DateOfBirth,
	}
}
