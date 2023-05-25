package repo

import (
	"context"
	"errors"

	pgx "github.com/jackc/pgx/v4"
	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase"
	common_entity "isling-be/internal/common/entity"
	"isling-be/pkg/postgres"
)

type AccountRepo struct {
	*postgres.Postgres
}

func NewAccountRepo(pg *postgres.Postgres) usecase.AccountRepository {
	return &AccountRepo{Postgres: pg}
}

func (repo *AccountRepo) Store(ctx context.Context, user *entity.Account) (*entity.Account, error) {
	sql := `
		INSERT INTO accounts (email, encrypted_password)
		VALUES ($1, $2)
		RETURNING id
	`

	var id common_entity.AccountID
	// Use repo.pool.QueryRow to scan returning data
	err := repo.Pool.QueryRow(ctx, sql, user.Email, user.EncryptedPassword).Scan(&id)
	if err != nil {
		return nil, err
	}

	user.ID = id

	return user, nil
}

func (repo *AccountRepo) FindByUsername(ctx context.Context, email string) (*entity.Account, error) {
	sql := `
		SELECT
		  id,
			email,
			encrypted_password,
			created_at,
			updated_at
		FROM accounts
		WHERE email = $1
		LIMIT 1
	`

	return rowToAccount(repo.Pool.QueryRow(ctx, sql, email))
}

func (repo *AccountRepo) FindByID(ctx context.Context, accountID common_entity.AccountID) (*entity.Account, error) {
	sql := `
		SELECT
			id,
			email,
			encrypted_password,
			created_at,
			updated_at
		FROM accounts
		WHERE id = $1
		LIMIT 1
	`

	return rowToAccount(repo.Pool.QueryRow(ctx, sql, accountID))
}

func rowToAccount(row pgx.Row) (*entity.Account, error) {
	user := entity.Account{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.EncryptedPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, common_entity.ErrNoRows
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
