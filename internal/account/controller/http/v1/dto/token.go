package dto

import "isling-be/internal/account/usecase/request"

const (
	GrantTypePassword     = "password"
	GrantTypeRefreshToken = "refresh_token"
)

type GetTokenByPasswordRequestDTO struct {
	Email    string `validate:"required,email" example:"isling@isling.me"`
	Password string `validate:"required" example:"wakaranai"`
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

type GetTokenByRefreshTokenRequestDTO struct {
	RefreshToken string `validate:"required" example:"himitsu"`
}

func (dto *GetTokenByRefreshTokenRequestDTO) ToRequest() *request.GetTokenByRefreshTokenRequest {
	return &request.GetTokenByRefreshTokenRequest{
		RefreshToken: dto.RefreshToken,
	}
}
