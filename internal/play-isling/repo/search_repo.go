package repo

import (
	"context"
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/surreal"

	"github.com/surrealdb/surrealdb.go"
)

type SearchRepo struct {
	*surreal.Surreal
}

var _ usecase.SearchRepository = (*SearchRepo)(nil)

func NewSearchRepo(sur *surreal.Surreal) *SearchRepo {
	return &SearchRepo{Surreal: sur}
}

func (r *SearchRepo) SearchRoom(_ context.Context, _ common_entity.AccountID, req *usecase.SearchRequest) ([]entity.RoomSearchResult, error) {
	sql := `
		SELECT
			*,
			ownerID AS owner,
			(
				IF audCountUpdatedAt AND time::now() - audCountUpdatedAt < 39s
			  THEN audienceCount
				ELSE 0
				END
			) AS audienceCount,
			math::mean([search::score(1),search::score(2)]) AS score
		FROM media_rooms
		WHERE name @1@ $query OR ownerFullName @2@ $query
		ORDER score DESC
		LIMIT $limit
		START $start
		FETCH owner
	`

	raw, err := r.Surreal.Query(sql, map[string]interface{}{
		"query": req.Query,
		"limit": req.Limit,
		"start": req.Offset,
	})

	rooms, err := surrealdb.SmartUnmarshal[[]entity.RoomSearchResult](raw, err)
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *SearchRepo) GetTotalRoomMatches(_ context.Context, _ common_entity.AccountID, req *usecase.SearchRequest) (int, error) {
	sql := `
		SELECT COUNT() AS total
		FROM media_rooms
		WHERE name @@ $query OR ownerFullName @@ $query
	`

	raw, err := r.Surreal.Query(sql, map[string]interface{}{
		"query": req.Query,
	})

	count, err := surrealdb.SmartUnmarshal[[]struct {
		Total int `json:"total"`
	}](raw, err)
	if err != nil {
		return 0, err
	}

	if len(count) == 0 {
		return 0, nil
	}

	return count[0].Total, nil
}
