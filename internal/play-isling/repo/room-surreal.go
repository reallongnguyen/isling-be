package repo

import (
	"context"
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/postgres"
	"isling-be/pkg/surreal"
	"strings"

	"github.com/surrealdb/surrealdb.go"
)

type RoomRepoSurreal struct {
	*postgres.Postgres
	*surreal.Surreal
}

var _ usecase.RoomRepository = (*RoomRepoSurreal)(nil)

func NewRoomRepoSurreal(sur *surreal.Surreal, pg *postgres.Postgres) *RoomRepoSurreal {
	return &RoomRepoSurreal{
		Surreal:  sur,
		Postgres: pg,
	}
}

func (repo *RoomRepoSurreal) Create(c context.Context, room *entity.Room) (*entity.Room, error) {
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

func (repo *RoomRepoSurreal) FindMany(c context.Context, filter *usecase.FindRoomFilter, order *usecase.Order) (*common_entity.Collection[entity.Room], error) {
	vars := make(map[string]interface{})
	whereConditions := make([]string, 0, 0)

	if filter != nil {
		if filter.OwnerID != nil {
			vars["ownerID"] = filter.OwnerID
			whereConditions = append(whereConditions, "ownerID = $ownerID")
		}

		if filter.IDIn != nil {
			vars["IDIn"] = filter.IDIn
			whereConditions = append(whereConditions, "originalID IN $IDIn")
		}
	}

	if len(whereConditions) == 0 {
		whereConditions = append(whereConditions, "TRUE")
	}

	whereClause := strings.Join(whereConditions, " AND ")

	orderClause := "originalID DESC"

	if order != nil {
		orderClause = order.Field + " " + order.Direction + ", originalID DESC"
	}

	sql := `
		SELECT
			*,
			ownerID as owner,
			(
				IF audCountUpdatedAt AND time::now() - audCountUpdatedAt < 39s
				THEN audienceCount
				ELSE 0
				END
			) AS audienceCount
		FROM media_rooms
		WHERE ` + whereClause + `
		ORDER ` + orderClause + `
		FETCH owner
	`

	raw, err := repo.Surreal.Query(sql, vars)

	roomSurreals, err := surrealdb.SmartUnmarshal[[]entity.RoomSurreal](raw, err)
	if err != nil {
		return nil, err
	}

	rooms := make([]entity.Room, len(roomSurreals))

	for idx := range roomSurreals {
		rooms[idx] = *roomSurreals[idx].ToRoom()
	}

	roomCollection := common_entity.NewCollection(rooms, 0, len(rooms), len(rooms))

	return &roomCollection, nil
}

func (repo *RoomRepoSurreal) FindOne(c context.Context, id int64) (*entity.Room, error) {
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

func (repo *RoomRepoSurreal) FindOneBySlug(c context.Context, slug string) (*entity.Room, error) {
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

func (repo *RoomRepoSurreal) UpdateOne(c context.Context, room *entity.Room) (*entity.Room, error) {
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

func (repo *RoomRepoSurreal) DeleteOne(c context.Context, id int64) error {
	sql := `
		DELETE FROM media_rooms
		WHERE id = $1
	`

	_, err := repo.Pool.Exec(c, sql, id)

	return err
}
