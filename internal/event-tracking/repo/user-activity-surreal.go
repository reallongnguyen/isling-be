package repo

import (
	"isling-be/internal/event-tracking/entity"
	"isling-be/internal/event-tracking/usecase"
	"isling-be/pkg/surreal"
)

type UserActSurRepo struct {
	sur *surreal.Surreal
}

func NewUserActSurRepo(sur *surreal.Surreal) usecase.UserActRepository {
	return &UserActSurRepo{
		sur: sur,
	}
}

func (r *UserActSurRepo) InsertMany(items []entity.UserActivity[any]) error {
	if len(items) == 0 {
		return nil
	}

	ulid := "user_activities:ulid()"

	for idx := range items {
		items[idx].ID = ulid
	}

	sql := `
		INSERT INTO user_activities $items
	`

	_, err := r.sur.Query(sql, map[string]any{
		"items": items,
	})

	return err
}
