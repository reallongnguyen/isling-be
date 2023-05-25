package dto

import "isling-be/internal/account/usecase/request"

type GetTokenByPasswordRequestDTO struct {
	Email    string `json:"email" validate:"required,email" example:"isling@isling.me"`
	Password string `json:"password" validate:"required" example:"wakaranai"`
}

type GetTokenResponseDTO struct {
	RefreshToken string `json:"refreshToken" example:"a-refresh-token"`
	AccessToken  string `json:"accessToken" example:"an-access-token"`
	TokenType    string `json:"tokenType" example:"bearer"`
	ExpiresIn    int    `json:"expiresIn" example:"3600"`
}

func (dto *GetTokenByPasswordRequestDTO) ToRequest() *request.GetTokenByPasswordRequest {
	return &request.GetTokenByPasswordRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func FromGetTokenRequestToDTO(req *request.GetTokenResponse) *GetTokenResponseDTO {
	return &GetTokenResponseDTO{
		RefreshToken: req.RefreshToken,
		AccessToken:  req.AccessToken,
		TokenType:    req.TokenType,
		ExpiresIn:    req.ExpiresIn,
	}
}
