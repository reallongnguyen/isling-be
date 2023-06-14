package repo

import (
	"context"
	"errors"
	"strings"

	"isling-be/internal/account/entity"
	"isling-be/internal/account/usecase"
	"isling-be/internal/account/usecase/request"
	common_entity "isling-be/internal/common/entity"
	"isling-be/pkg/postgres"

	pgx "github.com/jackc/pgx/v4"
)

type ProfileRepo struct {
	*postgres.Postgres
}

var _ usecase.ProfileRepository = (*ProfileRepo)(nil)

func NewProfileRepo(pg *postgres.Postgres) usecase.ProfileRepository {
	return &ProfileRepo{Postgres: pg}
}

func (repo *ProfileRepo) FindOneProfileByID(ctx context.Context, accountID common_entity.AccountID) (*entity.Profile, error) {
	sql := `
		SELECT
			a.id AS account_id,
			email,
			first_name,
			last_name,
			gender,
			date_of_birth,
			p.created_at AS created_at,
			p.updated_at AS updated_at
		FROM accounts AS a
			LEFT JOIN profiles AS p ON a.id = p.account_id
		WHERE a.id = $1
	`

	return rowToProfile(repo.Pool.QueryRow(ctx, sql, accountID))
}

func (repo *ProfileRepo) UpsertProfile(ctx context.Context, accountID common_entity.AccountID, createProfileReq *request.CreateProfileReq) (*entity.Profile, error) {
	sql := `
		INSERT INTO profiles (account_id, first_name, last_name, gender, date_of_birth)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (account_id)
		DO UPDATE SET first_name = $2, last_name = $3, gender = $4, date_of_birth = $5
	`

	_, err := repo.Pool.Exec(ctx, sql, accountID, createProfileReq.FirstName, createProfileReq.LastName, createProfileReq.Gender, createProfileReq.DateOfBirth)

	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return nil, common_entity.ErrAccountIDDuplicated
		}

		return nil, err
	}

	return repo.FindOneProfileByID(ctx, accountID)
}

func rowToProfile(row pgx.Row) (*entity.Profile, error) {
	profile := entity.Profile{}

	err := row.Scan(
		&profile.AccountID,
		&profile.Email,
		&profile.FirstName,
		&profile.LastName,
		&profile.Gender,
		&profile.DateOfBirth,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, common_entity.ErrAccountNotFound
	}

	if err != nil {
		return nil, err
	}

	return &profile, nil
}
