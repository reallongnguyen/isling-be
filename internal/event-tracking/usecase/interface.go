package usecase

import (
	"isling-be/internal/event-tracking/entity"
)

type (
	UserActRepository interface {
		InsertMany(items []entity.UserActivity[any]) error
	}
)
