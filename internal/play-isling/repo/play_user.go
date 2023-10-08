package repo

import (
	"context"
	"errors"
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/postgres"

	"github.com/jackc/pgx/v4"
)

type PlayUserRepo struct {
	*postgres.Postgres
}

var _ usecase.PlayUserRepository = (*PlayUserRepo)(nil)

func NewPlayUserRepo(postgres *postgres.Postgres) usecase.PlayUserRepository {
	return &PlayUserRepo{
		Postgres: postgres,
	}
}

func (repo *PlayUserRepo) GetOne(c context.Context, accountID common_entity.AccountID) (*entity.PlayUser, error) {
	sql := `
		SELECT
			id,
			account_id,
			recently_joined_rooms
		FROM play_users pu
		WHERE account_id = $1
	`

	row := repo.Pool.QueryRow(c, sql, accountID)

	playUser, err := rowToPlayUser(row)

	if err == nil {
		return playUser, err
	}

	if errors.Is(err, entity.ErrPlayUserNotFound) {
		return repo.Create(c, accountID)
	}

	return nil, err
}

func (repo *PlayUserRepo) Create(c context.Context, accountID common_entity.AccountID) (*entity.PlayUser, error) {
	sql := `
		INSERT INTO play_users (account_id, recently_joined_rooms)
		VALUES ($1, $2::json)
		RETURNING id, account_id, recently_joined_rooms
	`

	row := repo.Pool.QueryRow(c, sql, accountID, []int64{})

	return rowToPlayUser(row)
}

func (repo *PlayUserRepo) Update(c context.Context, accountID common_entity.AccountID, playUser *entity.PlayUser) error {
	sql := `
		UPDATE play_users
		SET recently_joined_rooms = $2::json
		WHERE account_id = $1
	`

	_, err := repo.Pool.Exec(c, sql, accountID, playUser.RecentlyJoinedRooms)

	return err
}

func (repo *PlayUserRepo) InsertRecentlyJoinedRoom(c context.Context, accountID common_entity.AccountID, roomID int64) error {
	sql := `
		UPDATE play_users
		SET recently_joined_rooms = jsonb_insert(recently_joined_rooms, '{0}', $2)
		WHERE account_id = $1
	`

	_, err := repo.Pool.Exec(c, sql, accountID, roomID)

	return err
}

func rowToPlayUser(row pgx.Row) (*entity.PlayUser, error) {
	playUser := new(entity.PlayUser)

	err := row.Scan(
		&playUser.ID,
		&playUser.AccountID,
		&playUser.RecentlyJoinedRooms,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrPlayUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return playUser, nil
}
