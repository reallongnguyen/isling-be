package usecase

import (
	"isling-be/internal/event-tracking/entity"
	"time"
)

type (
	UserActRepository interface {
		InsertMany(items []*entity.UserActivity[any]) error
	}

	CCULogRepository interface {
		InsertMany(items []*entity.CCULog) error
		CountCCU(timestamp time.Time, windowSize uint) (int64, error)
	}
)
