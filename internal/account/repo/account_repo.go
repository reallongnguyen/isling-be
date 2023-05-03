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
		INSERT INTO accounts (username, password)
		VALUES ($1, $2)
		RETURNING id
	`

	var id common_entity.AccountId
	// Use repo.pool.QueryRow to scan returning data
	err := repo.Pool.QueryRow(ctx, sql, user.Username, user.Password).Scan(&id)

	if err != nil {
		return nil, err
	}

	user.Id = id

	return user, nil
}

func (repo *AccountRepo) FindByUsername(ctx context.Context, username string) (*entity.Account, error) {
	sql := `
		SELECT
		  id,
			username,
			password,
			created_at,
			updated_at
		FROM accounts
		WHERE username = $1
		LIMIT 1
	`

	return rowToAccount(repo.Pool.QueryRow(ctx, sql, username))
}

func (repo *AccountRepo) FindByID(ctx context.Context, accountId common_entity.AccountId) (*entity.Account, error) {
	sql := `
		SELECT
			id,
			username,
			password,
			created_at,
			updated_at
		FROM accounts
		WHERE id = $1
		LIMIT 1
	`

	return rowToAccount(repo.Pool.QueryRow(ctx, sql, accountId))
}

func rowToAccount(row pgx.Row) (*entity.Account, error) {
	user := entity.Account{}
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
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
