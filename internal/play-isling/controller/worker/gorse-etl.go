package worker

import (
	"context"
	"encoding/json"
	et_entity "isling-be/internal/event-tracking/entity"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/facade"
	"time"

	"github.com/zhenghaoz/gorse/client"
)

type GorseETL struct {
	recommendationUC usecase.RecommendationUsecase
}

func NewGorseETL(recommendationUC usecase.RecommendationUsecase) *GorseETL {
	return &GorseETL{
		recommendationUC: recommendationUC,
	}
}

func (r *GorseETL) Run() {
	insertFeedback := func(uuid string, payload []byte, metadata map[string]string) error {
		userAct := new(et_entity.UserActivity[et_entity.ActOnItemData])

		if err := json.Unmarshal(payload, userAct); err != nil {
			facade.Log().Debug("parse userActivity json error: %w", err)

			return nil
		}

		feedback := client.Feedback{
			FeedbackType: userAct.EventName,
			UserId:       userAct.UserID,
			ItemId:       userAct.Data.ItemID,
			Timestamp:    userAct.Timestamp.Format(time.RFC3339),
		}

		return r.recommendationUC.InsertFeedback(context.Background(), feedback)
	}

	if err := facade.MsgBus().Subscribe("feedback-item", insertFeedback); err != nil {
		facade.Log().Error("subscribe 'feedback-item': %w", err)
	}
}
