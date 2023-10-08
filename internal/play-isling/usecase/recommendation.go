package usecase

import (
	"context"
	"isling-be/config"
	account_entity "isling-be/internal/account/entity"
	"isling-be/internal/play-isling/entity"
	"isling-be/pkg/facade"
	"strconv"
	"strings"
	"time"

	"github.com/facebookgo/muster"
	"github.com/zhenghaoz/gorse/client"
	"golang.org/x/exp/slices"
)

type RecommendationUC struct {
	gorse         *client.GorseClient
	InsertFBBatch *InsertFeedbackBatch
}

var _ RecommendationUsecase = (*RecommendationUC)(nil)

var cfg, _ = config.NewConfig()

func NewRecommendationUC() *RecommendationUC {
	gorse := client.NewGorseClient(cfg.GORSE.URL, cfg.GORSE.APIKey)
	insertFBBatch := NewInsertFeedbackBatch(gorse)

	insertFBBatch.Start()

	return &RecommendationUC{
		gorse:         gorse,
		InsertFBBatch: insertFBBatch,
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
		Comment:    "insert room " + room.Name,
		Labels:     []string{},
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

func (uc *RecommendationUC) InsertFeedback(_ context.Context, feedback client.Feedback) error {
	uc.InsertFBBatch.Add(feedback)

	return nil
}

type InsertFeedbackBatch struct {
	MaxBatchSize        uint
	BatchTimeout        time.Duration
	PendingWorkCapacity uint
	muster              muster.Client
	gorse               *client.GorseClient
}

func NewInsertFeedbackBatch(gorse *client.GorseClient) *InsertFeedbackBatch {
	return &InsertFeedbackBatch{
		gorse:               gorse,
		MaxBatchSize:        1000,
		BatchTimeout:        60 * time.Second,
		PendingWorkCapacity: 8000,
	}
}

func (r *InsertFeedbackBatch) Start() error {
	r.muster.MaxBatchSize = r.MaxBatchSize
	r.muster.BatchTimeout = r.BatchTimeout
	r.muster.PendingWorkCapacity = r.PendingWorkCapacity
	r.muster.BatchMaker = func() muster.Batch { return &FeedbackBatch{client: r} }

	return r.muster.Start()
}

func (r *InsertFeedbackBatch) Stop() error {
	return r.muster.Stop()
}

func (r *InsertFeedbackBatch) Add(item client.Feedback) {
	r.muster.Work <- item
}

type FeedbackBatch struct {
	client *InsertFeedbackBatch
	Items  []client.Feedback
}

func (r *FeedbackBatch) Add(item any) {
	r.Items = append(r.Items, item.(client.Feedback))
}

func (r *FeedbackBatch) Fire(notifier muster.Notifier) {
	defer notifier.Done()

	slices.SortFunc(r.Items, func(f1, f2 client.Feedback) int {
		if f1.UserId != f2.UserId {
			return strings.Compare(f1.UserId, f2.UserId)
		}

		if f1.ItemId != f2.ItemId {
			return strings.Compare(f1.ItemId, f2.ItemId)
		}

		return strings.Compare(f1.FeedbackType, f2.FeedbackType)
	})

	uniqItems := slices.CompactFunc(r.Items, func(f1, f2 client.Feedback) bool {
		return f1.UserId == f2.UserId && f1.ItemId == f2.ItemId && f1.FeedbackType == f2.FeedbackType
	})

	_, err := r.client.gorse.PutFeedback(context.Background(), uniqItems)
	if err != nil {
		facade.Log().Error("feedbackBatch fire: %w", err)
	}
}
