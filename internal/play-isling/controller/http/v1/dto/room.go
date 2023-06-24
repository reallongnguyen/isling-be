package dto

import "isling-be/internal/play-isling/usecase"

type CreateRoomDTO struct {
	Name        string `json:"name" validate:"required,max=256"`
	Description string `json:"description" validate:"max=512"`
	Cover       string `json:"cover" validate:"required,max=256,http_url"`
}

func (dto *CreateRoomDTO) ToCreateRoomReq() *usecase.CreateRoomRequest {
	return &usecase.CreateRoomRequest{
		Name:        dto.Name,
		Description: dto.Description,
		Cover:       dto.Cover,
	}
}
