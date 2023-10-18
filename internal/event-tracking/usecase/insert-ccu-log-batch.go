package usecase

import (
	"isling-be/internal/event-tracking/entity"
	"isling-be/pkg/facade"
	"time"

	"github.com/facebookgo/muster"
)

type CCULogBatch struct {
	MaxBatchSize        uint
	BatchTimeout        time.Duration
	PendingWorkCapacity uint
	muster              muster.Client
	CCULogRepo          CCULogRepository
}

func (r *CCULogBatch) Start() error {
	r.muster.MaxBatchSize = r.MaxBatchSize
	r.muster.BatchTimeout = r.BatchTimeout
	r.muster.PendingWorkCapacity = r.PendingWorkCapacity
	r.muster.BatchMaker = func() muster.Batch { return &CCUBatch{Client: r} }

	return r.muster.Start()
}

func (r *CCULogBatch) Stop() error {
	return r.muster.Stop()
}

func (r *CCULogBatch) Add(item *entity.CCULog) error {
	r.muster.Work <- item

	return nil
}

type CCUBatch struct {
	Client *CCULogBatch
	Items  []*entity.CCULog
}

func (r *CCUBatch) Add(item any) {
	r.Items = append(r.Items, item.(*entity.CCULog))
}

func (r *CCUBatch) Fire(notifier muster.Notifier) {
	defer notifier.Done()

	err := r.Client.CCULogRepo.InsertMany(r.Items)
	if err != nil {
		facade.Log().Error("insert many CCU log: %w", err)
	}
}
