package repo

import (
	"context"
	"isling-be/internal/play-isling/entity"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/postgres"
)

type RoomRepo struct {
	*postgres.Postgres
}

var _ usecase.RoomRepository = (*RoomRepo)(nil)

func NewRoomRepo(pg *postgres.Postgres) *RoomRepo {
	return &RoomRepo{Postgres: pg}
}

func (repo *RoomRepo) Create(c context.Context, room *entity.Room) (*entity.Room, error) {
	sql := `
		INSERT INTO media_rooms (owner_id, visibility, invite_code, name, slug, description, cover)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, audiences, created_at, updated_at
	`

	err := repo.Pool.QueryRow(
		c,
		sql,
		room.OwnerID,
		room.Visibility,
		room.InviteCode,
		room.Name,
		room.Slug,
		room.Description,
		room.Cover,
	).Scan(
		&room.ID,
		&room.Audiences,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return room, nil
}
