package worker

import (
	"fmt"
	"isling-be/pkg/facade"
	"isling-be/pkg/surreal"
	"strings"
	"time"

	"github.com/surrealdb/surrealdb.go"
)

type RoomAudCounter struct {
	sr *surreal.Surreal
}

func NewRoomAudCounter(sr *surreal.Surreal) *RoomAudCounter {
	return &RoomAudCounter{
		sr: sr,
	}
}

func (r *RoomAudCounter) Run() {
	go func() {
		for {
			facade.Log().Debug("RoomAudience: start delete zombie row in join table")

			_, err := r.sr.Query("DELETE join WHERE time::now() - time.pinged > 96s", map[string]string{})
			if err != nil {
				facade.Log().Error("RoomAudience: delete zombie row in join table: %w", err)
			}

			time.Sleep(60 * time.Second)
		}
	}()

	go func() {
		isFirstRun := true

		for {
			if isFirstRun {
				isFirstRun = false

				time.Sleep(10 * time.Second)
			} else {
				// if you change the duration here,
				// you must update the duration in internal/play-isling/repo/search_repo.go:29
				// and internal/play-isling/repo/room-surreal.go:92
				time.Sleep(20 * time.Second)
			}

			facade.Log().Debug("RoomAudience: start calculate audience")

			rawRes, err := r.sr.Query("SELECT count(in), out FROM join GROUP BY out", map[string]string{})
			if err != nil {
				facade.Log().Error("RoomAudience: calculate audience: %w", err)

				continue
			}

			res, err := surrealdb.SmartUnmarshal[[]struct {
				Count int    `json:"count"`
				Out   string `json:"out"`
			}](rawRes, err)
			if err != nil {
				facade.Log().Error("RoomAudience: calculate audience: %w", err)

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
				facade.Log().Error("RoomAudience: update audience count: %w", err)
			}
		}
	}()
}
