package request

import "strings"

type GetTokenByPasswordRequest struct {
	Email    string
	Password string
}

func (req *GetTokenByPasswordRequest) Normalize() *GetTokenByPasswordRequest {
	req.Email = strings.ToLower(req.Email)

	return req
}

type GetTokenResponse struct {
	RefreshToken string
	AccessToken  string
	TokenType    string
	ExpiresIn    int
}

type GetTokenByRefreshTokenRequest struct {
	RefreshToken string
}
