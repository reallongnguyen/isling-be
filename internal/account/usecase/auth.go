package usecase

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase/request"
	common_entity "isling-be/internal/common/entity"
	common_uc "isling-be/internal/common/usecase"
	"isling-be/pkg/logger"
	"time"

	"isling-be/config"

	"github.com/golang-jwt/jwt/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const (
	refreshTokenLength = 48
)

var cfg, _ = config.NewConfig()

type AuthUC struct {
	log              logger.Interface
	accountUC        AccountUsecase
	accountRepo      AccountRepository
	refreshTokenRepo RefreshTokenRepository
}

var _ AuthUsecase = (*AuthUC)(nil)

func NewAuthUsecase(log logger.Interface, accountUC AccountUsecase, accountRepo AccountRepository, refreshTokenRepo RefreshTokenRepository) AuthUsecase {
	return &AuthUC{
		log:              log,
		accountUC:        accountUC,
		accountRepo:      accountRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (authUC *AuthUC) GetTokenByPassword(c context.Context, credential *request.GetTokenByPasswordRequest) (*request.GetTokenResponse, error) {
	account, err := authUC.accountRepo.FindByUsername(c, credential.Email)

	if err != nil && errors.Is(err, common_entity.ErrNoRows) {
		authUC.log.Warn("auth usecase: get token by password: not found email: %s", credential.Email)

		return nil, common_entity.ErrAccountNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("auth usecase: get token by password: find an account %w", err)
	}

	if !common_uc.IsMatchHashAndPassword(account.EncryptedPassword, credential.Password) {
		authUC.log.Warn("auth usecase: get token by password: email password not match. Email: %s", credential.Email)

		return nil, common_entity.ErrEmailPasswordNotMatch
	}

	tokenRes, err := getTokenResponse(account)
	if err != nil {
		return nil, fmt.Errorf("auth usecase: get token by password: get token: %w", err)
	}

	refreshToken := entity.RefreshTokens{
		EncryptedToken: hashRefreshToken(tokenRes.RefreshToken),
		AccountID:      account.ID,
	}

	if _, err := authUC.refreshTokenRepo.Store(c, nil, &refreshToken); err != nil {
		return nil, fmt.Errorf("auth usecase: get token by password: store refresh token: %w", err)
	}

	return tokenRes, nil
}

func (authUC *AuthUC) GetTokenByRefreshToken(c context.Context, credential *request.GetTokenByRefreshTokenRequest) (*request.GetTokenResponse, error) {
	encryptedRefreshToken := hashRefreshToken(credential.RefreshToken)

	refreshToken, err := authUC.refreshTokenRepo.FindOneByEncryptedToken(c, nil, encryptedRefreshToken)
	if err != nil {
		return nil, err
	}

	if refreshToken.Revoked {
		return nil, common_entity.ErrRefreshTokenInvalid
	}

	account, err := authUC.accountRepo.FindByID(c, refreshToken.AccountID)
	if err != nil {
		return nil, err
	}

	tokenRes, err := getTokenResponse(account)
	if err != nil {
		return nil, fmt.Errorf("auth usecase: get token by refresh token: get token: %w", err)
	}

	// TODO: hide logic of transaction in use case
	tx, err := authUC.refreshTokenRepo.BeginTx(c)
	if err != nil {
		return nil, fmt.Errorf("auth usecase: get token by refresh token: create transaction: %w", err)
	}

	if _, err := authUC.refreshTokenRepo.RevokeOneByEncryptedToken(c, tx, encryptedRefreshToken); err != nil {
		tx.Rollback(c)

		return nil, fmt.Errorf("auth usecase: get token by refresh token: revoke token: %w", err)
	}

	newRefreshToken := entity.RefreshTokens{
		EncryptedToken: hashRefreshToken(tokenRes.RefreshToken),
		AccountID:      account.ID,
	}

	if _, err := authUC.refreshTokenRepo.Store(c, tx, &newRefreshToken); err != nil {
		tx.Rollback(c)

		return nil, fmt.Errorf("auth usecase: get token by refresh token: store refresh token: %w", err)
	}

	if err = tx.Commit(c); err != nil {
		return nil, fmt.Errorf("auth usecase: get token by refresh token: commit transaction: %w", err)
	}

	return tokenRes, nil
}

func (authUC *AuthUC) Logout(c context.Context, accountID common_entity.AccountID, refreshToken string) error {
	if refreshToken != "" {
		rowAffected, err := authUC.refreshTokenRepo.RevokeOneByEncryptedToken(c, nil, hashRefreshToken(refreshToken))

		if err != nil {
			return err
		}

		if rowAffected == 0 {
			return common_entity.ErrRefreshTokenNotFound
		}

		return nil
	}

	_, err := authUC.refreshTokenRepo.RevokeManyByAccountID(c, nil, accountID)

	return err
}

func (authUC *AuthUC) SignUp(ctx context.Context, createUserDto request.CreateAccountReq) (*request.GetTokenResponse, error) {
	acc, err := authUC.accountUC.CreateAccount(ctx, createUserDto)

	if err != nil {
		return nil, err
	}

	return getTokenResponse(acc)
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
		ExpiresIn:    cfg.JWT.AccessTokenEXP,
	}

	return &tokenRes, nil
}

func getAccessToken(account *entity.Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud":        []string{cfg.JWT.Audience},
		"exp":        time.Now().Add(time.Duration(cfg.JWT.AccessTokenEXP) * time.Second).Unix(),
		"iat":        time.Now().Unix(),
		"iss":        "sign in",
		"sub":        account.Email,
		"account_id": account.ID,
	})

	tokenString, err := token.SignedString([]byte(cfg.JWT.JWTSecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getRandomRefreshToken() (string, error) {
	return gonanoid.New(refreshTokenLength)
}

func hashRefreshToken(refreshToken string) string {
	bytes := sha512.Sum512([]byte(refreshToken))

	return base64.StdEncoding.EncodeToString(bytes[:])
}
