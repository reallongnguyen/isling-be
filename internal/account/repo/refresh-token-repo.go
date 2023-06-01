package repo

import (
	"context"
	"errors"
	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase"
	common_entity "isling-be/internal/common/entity"
	"isling-be/pkg/postgres"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
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

func (repo *RefreshTokenRepo) BeginTx(c context.Context) (pgx.Tx, error) {
	return repo.Pool.Begin(c)
}

func (repo *RefreshTokenRepo) Store(c context.Context, tx pgx.Tx, refreshToken *entity.RefreshTokens) (*entity.RefreshTokens, error) {
	var querier pgxtype.Querier = repo.Pool

	if tx != nil {
		querier = tx
	}

	sql := `
		INSERT INTO refresh_tokens (account_id, encrypted_token)
		VALUES ($1, $2)
		RETURNING id
	`

	refreshTokenID := 0

	err := querier.QueryRow(c, sql, refreshToken.AccountID, refreshToken.EncryptedToken).Scan(&refreshTokenID)
	if err != nil {
		return nil, err
	}

	refreshToken.ID = refreshTokenID

	return refreshToken, nil
}

func (repo *RefreshTokenRepo) FindOneByEncryptedToken(c context.Context, tx pgx.Tx, encryptedRefreshToken string) (*entity.RefreshTokens, error) {
	var querier pgxtype.Querier = repo.Pool

	if tx != nil {
		querier = tx
	}

	sql := `
		SELECT
			id,
			account_id,
			encrypted_token,
			revoked,
			created_at
		FROM refresh_tokens
		WHERE encrypted_token = $1
		LIMIT 1
	`

	row := querier.QueryRow(c, sql, encryptedRefreshToken)

	return rowToRefreshToken(row)
}

func (repo *RefreshTokenRepo) RevokeManyByAccountID(c context.Context, tx pgx.Tx, accountID common_entity.AccountID) (int64, error) {
	var querier pgxtype.Querier = repo.Pool

	if tx != nil {
		querier = tx
	}

	sql := `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE account_id = $1 AND revoked IS FALSE
	`

	tag, err := querier.Exec(c, sql, accountID)

	return tag.RowsAffected(), err
}

func (repo *RefreshTokenRepo) RevokeOneByEncryptedToken(c context.Context, tx pgx.Tx, encryptedRefreshToken string) (int64, error) {
	var querier pgxtype.Querier = repo.Pool

	if tx != nil {
		querier = tx
	}

	sql := `
		UPDATE refresh_tokens
		SET revoked = TRUE
		WHERE encrypted_token = $1
	`

	tag, err := querier.Exec(c, sql, encryptedRefreshToken)

	return tag.RowsAffected(), err
}

func rowToRefreshToken(row pgx.Row) (*entity.RefreshTokens, error) {
	refreshToken := entity.RefreshTokens{}

	err := row.Scan(
		&refreshToken.ID,
		&refreshToken.AccountID,
		&refreshToken.EncryptedToken,
		&refreshToken.Revoked,
		&refreshToken.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, common_entity.ErrRefreshTokenNotFound
	}

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}
