package repo

import (
	"context"
	"errors"
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/postgres"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

const (
	initialSliceCap = 16
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

func (repo *RoomRepo) FindMany(c context.Context, filter *usecase.FindRoomFilter, order *usecase.Order) (*common_entity.Collection[*entity.Room], error) {
	binds := make([]interface{}, 0, 0)
	whereConditions := make([]string, 0, 0)

	if filter != nil {
		if filter.OwnerID != nil {
			binds = append(binds, filter.OwnerID)
			whereConditions = append(whereConditions, "owner_id = $"+strconv.Itoa(len(binds)))
		}

		if filter.IDIn != nil {
			binds = append(binds, filter.IDIn)
			whereConditions = append(whereConditions, "mr.id = ANY($"+strconv.Itoa(len(binds))+")")
		}
	}

	if len(whereConditions) == 0 {
		whereConditions = append(whereConditions, "TRUE")
	}

	whereClause := strings.Join(whereConditions, " AND ")

	orderClause := "mr.id DESC"

	if order != nil {
		orderClause = order.Field + " " + order.Direction + ", mr.id DESC"
	}

	sql := `
		SELECT
			mr.id,
			owner_id,
			p.account_id AS owner_id,
			p.first_name AS owner_first_name,
			p.last_name AS owner_last_name,
			p.avatar_url AS owner_avatar_url,
			visibility,
			invite_code,
			name,
			slug,
			description,
			cover,
			audience_count,
			audiences,
			mr.created_at,
			mr.updated_at
		FROM media_rooms mr
			LEFT JOIN profiles p ON (mr.owner_id = p.account_id)
		WHERE ` + whereClause + `
		ORDER BY ` + orderClause + `
	`

	rooms := make([]*entity.Room, 0, initialSliceCap)

	roomOwner := new(entity.RoomOwner)
	room := entity.Room{Owner: roomOwner}

	_, err := repo.Pool.QueryFunc(
		c,
		sql,
		binds,
		[]interface{}{
			&room.ID,
			&room.OwnerID,
			&room.Owner.ID,
			&room.Owner.FirstName,
			&room.Owner.LastName,
			&room.Owner.AvatarURL,
			&room.Visibility,
			&room.InviteCode,
			&room.Name,
			&room.Slug,
			&room.Description,
			&room.Cover,
			&room.AudienceCount,
			&room.Audiences,
			&room.CreatedAt,
			&room.UpdatedAt,
		},
		func(pgx.QueryFuncRow) error {
			newRoom := room
			newRoom.Audiences = append(newRoom.Audiences, room.Audiences...)
			newRoom.Owner = new(entity.RoomOwner)
			*newRoom.Owner = *room.Owner

			if newRoom.Owner.ID == 0 {
				newRoom.Owner = nil
			}

			rooms = append(rooms, &newRoom)

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	roomCollection := common_entity.NewCollection(rooms, 0, len(rooms), len(rooms))

	return &roomCollection, nil
}

func (repo *RoomRepo) FindOne(c context.Context, id int64) (*entity.Room, error) {
	sql := `
		SELECT
			mr.id,
			owner_id,
			p.account_id AS owner_id,
			p.first_name AS owner_first_name,
			p.last_name AS owner_last_name,
			p.avatar_url AS owner_avatar_url,
			visibility,
			invite_code,
			name,
			slug,
			description,
			cover,
			audience_count,
			audiences,
			mr.created_at,
			mr.updated_at
		FROM media_rooms mr
			LEFT JOIN profiles p ON (mr.owner_id = p.account_id)
		WHERE mr.id = $1
	`

	row := repo.Pool.QueryRow(c, sql, id)

	return rowToRoom(row)
}

func (repo *RoomRepo) FindOneBySlug(c context.Context, slug string) (*entity.Room, error) {
	sql := `
		SELECT
			mr.id,
			owner_id,
			p.account_id AS owner_id,
			p.first_name AS owner_first_name,
			p.last_name AS owner_last_name,
			p.avatar_url AS owner_avatar_url,
			visibility,
			invite_code,
			name,
			slug,
			description,
			cover,
			audience_count,
			audiences,
			mr.created_at,
			mr.updated_at
		FROM media_rooms mr
			LEFT JOIN profiles p ON (mr.owner_id = p.account_id)
		WHERE slug = $1
	`

	row := repo.Pool.QueryRow(c, sql, slug)

	return rowToRoom(row)
}

func (repo *RoomRepo) UpdateOne(c context.Context, room *entity.Room) (*entity.Room, error) {
	sql := `
		UPDATE media_rooms
		SET 
			owner_id = $2,
			visibility = $3,
			invite_code = $4,
			name = $5,
			slug = $6,
			description = $7,
			cover = $8,
			audience_count = $9,
			audiences = $10
		WHERE id = $1
		RETURNING
			updated_at
	`

	row := repo.Pool.QueryRow(
		c,
		sql,
		room.ID,
		room.OwnerID,
		room.Visibility,
		room.InviteCode,
		room.Name,
		room.Slug,
		room.Description,
		room.Cover,
		room.AudienceCount,
		room.Audiences,
	)

	newRoom := *room

	err := row.Scan(&newRoom.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &newRoom, nil
}

func (repo *RoomRepo) DeleteOne(c context.Context, id int64) error {
	sql := `
		DELETE FROM media_rooms
		WHERE id = $1
	`

	_, err := repo.Pool.Exec(c, sql, id)

	return err
}

func rowToRoom(row pgx.Row) (*entity.Room, error) {
	roomOwner := new(entity.RoomOwner)
	room := entity.Room{Owner: roomOwner}

	err := row.Scan(
		&room.ID,
		&room.OwnerID,
		&room.Owner.ID,
		&room.Owner.FirstName,
		&room.Owner.LastName,
		&room.Owner.AvatarURL,
		&room.Visibility,
		&room.InviteCode,
		&room.Name,
		&room.Slug,
		&room.Description,
		&room.Cover,
		&room.AudienceCount,
		&room.Audiences,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, entity.ErrRoomNotFound
	}

	if room.Owner.ID == 0 {
		room.Owner = nil
	}

	return &room, nil
}
