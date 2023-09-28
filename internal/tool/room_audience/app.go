package room_audience

import (
	"isling-be/internal/tool/model"
	"isling-be/pkg/logger"
	"isling-be/pkg/surreal"
	"time"
)

type RoomAudience struct {
	sr  *surreal.Surreal
	log logger.Interface
}

var _ model.Tool = (*RoomAudience)(nil)

func NewRoomAudience(log logger.Interface, sr *surreal.Surreal) *RoomAudience {
	return &RoomAudience{
		sr:  sr,
		log: log,
	}
}

func (r *RoomAudience) Start() error {
	for {
		r.log.Debug("RoomAudience: start delete zombie row in join table")

		_, err := r.sr.Query("DELETE join WHERE time::now() - time.pinged >= 120s", map[string]string{})
		if err != nil {
			r.log.Error("RoomAudience: delete zombie row in join table: %w", err)
		}

		time.Sleep(60 * time.Second)
	}
}
