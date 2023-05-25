package request

type GetTokenByPasswordRequest struct {
	Email    string
	Password string
}

type GetTokenResponse struct {
	RefreshToken string
	AccessToken  string
	TokenType    string
	ExpiresIn    int
}
