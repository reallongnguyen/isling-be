package usecase

import (
	"context"
	account_entity "isling-be/internal/account/entity"
	"isling-be/internal/play-isling/entity"
	"strconv"
	"time"

	"github.com/zhenghaoz/gorse/client"
	"golang.org/x/exp/slices"
)

var feedbackTypes = []string{
	"read",
	"comment",
	"like",
	"share",
	"save",
	"reaction",
	"add-item",
	"watch-15min",
	"watch-1h",
}

type RecommendationUC struct {
	gorse *client.GorseClient
}

var _ RecommendationUsecase = (*RecommendationUC)(nil)

func NewRecommendationUC() RecommendationUsecase {
	return &RecommendationUC{
		gorse: client.NewGorseClient("http://localhost:8087", ""),
	}
}

func (uc *RecommendationUC) InsertUser(c context.Context, account *account_entity.Account) error {
	_, err := uc.gorse.InsertUser(c, client.User{
		Comment:   "insert user " + account.Email,
		Labels:    []string{},
		Subscribe: []string{},
		UserId:    strconv.FormatInt(int64(account.ID), 10),
	})

	return err
}

func (uc *RecommendationUC) InsertRoom(c context.Context, room *entity.Room) error {
	_, err := uc.gorse.InsertItem(c, client.Item{
		Comment: "insert room " + room.Name,
		Labels: []string{
			room.Name,
			room.Description,
		},
		Categories: []string{"room"},
		ItemId:     strconv.FormatInt(room.ID, 10),
		IsHidden:   room.Visibility != entity.VisibilityPublic,
		Timestamp:  time.Now().String(),
	})

	return err
}

func (uc *RecommendationUC) HideItem(c context.Context, itemID string) error {
	isHidden := true

	_, err := uc.gorse.UpdateItem(c, itemID, client.ItemPatch{
		IsHidden: &isHidden,
	})

	return err
}

func (uc *RecommendationUC) InsertFeedback(c context.Context, actions []CreateActionRequest) error {
	feedbacks := make([]client.Feedback, 0, 8)

	for _, action := range actions {
		if action.ObjectID == nil || !slices.Contains(feedbackTypes, action.Type) {
			continue
		}

		feedback := client.Feedback{
			FeedbackType: action.Type,
			UserId:       strconv.FormatInt(int64(action.AccountID), 10),
			ItemId:       *action.ObjectID,
			Timestamp:    action.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		}

		feedbacks = append(feedbacks, feedback)
	}

	_, err := uc.gorse.PutFeedback(c, feedbacks)

	return err
}
