package worker

import (
	"context"
	"encoding/json"
	et_entity "isling-be/internal/event-tracking/entity"
	"isling-be/internal/play-isling/usecase"
	"time"

	"github.com/zhenghaoz/gorse/client"
)

type GorseETL struct {
	userActChan      <-chan string
	recommendationUC usecase.RecommendationUsecase
}

func NewGorseETL(
	userActChan <-chan string,
	recommendationUC usecase.RecommendationUsecase,
) *GorseETL {
	return &GorseETL{
		userActChan:      userActChan,
		recommendationUC: recommendationUC,
	}
}

func (r *GorseETL) Run() {
	go func() {
		for act := range r.userActChan {
			userAct := new(et_entity.UserActivity[et_entity.ActOnItemData])

			if err := json.Unmarshal([]byte(act), userAct); err != nil {
				continue
			}

			feedback := client.Feedback{
				FeedbackType: userAct.EventName,
				UserId:       userAct.UserID,
				ItemId:       userAct.Data.ItemID,
				Timestamp:    userAct.Timestamp.Format(time.RFC3339),
			}

			r.recommendationUC.InsertFeedback(context.Background(), feedback)
		}
	}()
}
