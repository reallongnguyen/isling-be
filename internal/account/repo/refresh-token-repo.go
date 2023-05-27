package repo

import (
	"context"
	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase"
	"isling-be/pkg/postgres"
)

type RefreshTokenRepo struct {
	*postgres.Postgres
}

var _ usecase.RefreshTokenRepository = (*RefreshTokenRepo)(nil)

func NewRefreshTokenRepo(pg *postgres.Postgres) *RefreshTokenRepo {
	return &RefreshTokenRepo{
		Postgres: pg,
	}
}

func (repo *RefreshTokenRepo) Store(c context.Context, refreshToken *entity.RefreshTokens) (*entity.RefreshTokens, error) {
	sql := `
		INSERT INTO refresh_tokens (account_id, encrypted_token)
		VALUES ($1, $2)
		RETURNING id
	`

	refreshTokenID := 0

	err := repo.Pool.QueryRow(c, sql, refreshToken.AccountID, refreshToken.EncryptedToken).Scan(&refreshTokenID)
	if err != nil {
		return nil, err
	}

	refreshToken.ID = refreshTokenID

	return refreshToken, nil
}
