package usecase

import (
	"context"
	"errors"
	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase/request"
	common_entity "isling-be/internal/common/entity"
	common_uc "isling-be/internal/common/usecase"
	"isling-be/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	accessTokenExpiresIn = 3600
	jwtSecretKey         = "isling.me"
	audience             = "isling.me"
	refreshTokenLength   = 48
)

type AuthUC struct {
	log         logger.Interface
	accountRepo AccountRepository
}

func NewAuthUsecase(log logger.Interface, accountRepo AccountRepository) AuthUsecase {
	return &AuthUC{
		log:         log,
		accountRepo: accountRepo,
	}
}

func (authUC *AuthUC) GetTokenByPassword(c context.Context, credential *request.GetTokenByPasswordRequest) (*request.GetTokenResponse, error) {
	account, err := authUC.accountRepo.FindByUsername(c, credential.Email)

	if err != nil && errors.Is(err, common_entity.ErrNoRows) {
		authUC.log.Warn("auth usecase: sign in error: not found email: %s", credential.Email)

		return nil, common_entity.ErrAccountNotFound
	}

	if err != nil {
		authUC.log.Error("auth usecase: sign in error: unexpected error when find account: %s", err.Error())

		return nil, err
	}

	if !common_uc.IsMatchHashAndPassword(account.EncryptedPassword, credential.Password) {
		authUC.log.Warn("auth usecase: sign in error: email password not match. Email: %s", credential.Email)

		return nil, common_entity.ErrEmailPasswordNotMatch
	}

	tokenRes, err := getTokenResponse(account)
	if err != nil {
		authUC.log.Error("auth usecase: sign in error: unexpected error when get token response: %s", err.Error())

		return nil, err
	}

	return tokenRes, nil
}

func getTokenResponse(account *entity.Account) (*request.GetTokenResponse, error) {
	accessToken, err := getAccessToken(account)

	if err != nil {
		return nil, err
	}

	refreshToken, err := getRandomRefreshToken()

	if err != nil {
		return nil, err
	}

	tokenRes := request.GetTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "bearer",
		ExpiresIn:    accessTokenExpiresIn,
	}

	return &tokenRes, nil
}

func getAccessToken(account *entity.Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":        []string{audience},
		"exp":        time.Now().Add(accessTokenExpiresIn * time.Second).Unix(),
		"iat":        time.Now().Unix(),
		"iss":        "sign in",
		"sub":        account.Email,
		"account_id": account.ID,
	})

	tokenString, err := token.SignedString([]byte(jwtSecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getRandomRefreshToken() (string, error) {
	return gonanoid.New(refreshTokenLength)
}
