package room_audience

import (
	"fmt"
	"isling-be/internal/tool/model"
	"isling-be/pkg/logger"
	"isling-be/pkg/surreal"
	"strings"
	"time"

	"github.com/surrealdb/surrealdb.go"
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
	go func() {
		for {
			r.log.Debug("RoomAudience: start delete zombie row in join table")

			_, err := r.sr.Query("DELETE join WHERE time::now() - time.pinged >= 120s", map[string]string{})
			if err != nil {
				r.log.Error("RoomAudience: delete zombie row in join table: %w", err)
			}

			time.Sleep(60 * time.Second)
		}
	}()

	go func() {
		isFirstRun := true

		for {
			if isFirstRun {
				isFirstRun = false

				time.Sleep(30 * time.Second)
			} else {
				time.Sleep(60 * time.Second)
			}

			r.log.Debug("RoomAudience: start calculate audience")

			rawRes, err := r.sr.Query("SELECT count(in), out FROM join GROUP BY out", map[string]string{})
			if err != nil {
				r.log.Error("RoomAudience: calculate audience: %w", err)

				continue
			}

			res, err := surrealdb.SmartUnmarshal[[]struct {
				Count int    `json:"count"`
				Out   string `json:"out"`
			}](rawRes, err)
			if err != nil {
				r.log.Error("RoomAudience: calculate audience: %w", err)

				continue
			}

			if len(res) == 0 {
				continue
			}

			roomIDList := make([]string, len(res))
			audCountClauses := make([]string, len(res))

			for idx, count := range res {
				roomIDList[idx] = count.Out
				audCountClauses[idx] = fmt.Sprintf("IF id = %s THEN %d", count.Out, count.Count)
			}

			audCountClause := strings.Join(audCountClauses, " ELSE ")

			surql := `
			UPDATE media_rooms MERGE {
				audienceCount: (` + audCountClause + ` END),
				audCountUpdatedAt: time::now(),
			}
			WHERE id IN $idList
		`

			_, err = r.sr.Query(surql, map[string]interface{}{
				"idList": roomIDList,
			})
			if err != nil {
				r.log.Error("RoomAudience: update audience count: %w", err)
			}
		}
	}()

	return nil
}
