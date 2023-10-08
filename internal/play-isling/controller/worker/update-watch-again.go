package worker

import (
	"context"
	"encoding/json"
	cm_entity "isling-be/internal/common/entity"
	et_entity "isling-be/internal/event-tracking/entity"
	"isling-be/internal/play-isling/usecase"
	"isling-be/pkg/facade"
	"strconv"
)

type WatchAgainUpdater struct {
	playUserUC usecase.PlayUserUsecase
}

func NewWatchAgainUpdater(playUserUC usecase.PlayUserUsecase) *WatchAgainUpdater {
	return &WatchAgainUpdater{
		playUserUC: playUserUC,
	}
}

// TODO: apply batch job when request amount is large
func (r *WatchAgainUpdater) Run() {
	insertFeedback := func(uuid string, payload []byte, metadata map[string]string) error {
		userAct := new(et_entity.UserActivity[et_entity.ActOnItemData])

		if err := json.Unmarshal(payload, userAct); err != nil {
			facade.Log().Debug("parse userActivity json error: %w", err)

			return nil
		}

		accountID, err := strconv.Atoi(userAct.UserID)
		if err != nil {
			facade.Log().Info("parse account id from userActivity: %w", err)

			return nil
		}

		roomID, err := strconv.Atoi(userAct.Data.ItemID)
		if err != nil {
			facade.Log().Info("parse room id from userActivity: %w", err)

			return nil
		}

		err = r.playUserUC.InsertRecentlyJoinedRoom(context.Background(), cm_entity.AccountID(accountID), int64(roomID))
		if err != nil {
			facade.Log().Error("update recently join room: %w", err)
		}

		return nil
	}

	if err := facade.MsgBus().Subscribe("room.watched", insertFeedback); err != nil {
		facade.Log().Error("subscribe 'room.watched': %w", err)
	}
}
