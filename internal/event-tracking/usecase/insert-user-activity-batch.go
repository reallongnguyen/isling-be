package usecase

import (
	"isling-be/internal/event-tracking/entity"
	"isling-be/pkg/facade"
	"time"

	"github.com/facebookgo/muster"
)

type UserActBatch struct {
	MaxBatchSize        uint
	BatchTimeout        time.Duration
	PendingWorkCapacity uint
	muster              muster.Client
	UserActRepo         UserActRepository
}

func (r *UserActBatch) Start() error {
	r.muster.MaxBatchSize = r.MaxBatchSize
	r.muster.BatchTimeout = r.BatchTimeout
	r.muster.PendingWorkCapacity = r.PendingWorkCapacity
	r.muster.BatchMaker = func() muster.Batch { return &Batch{Client: r} }

	return r.muster.Start()
}

func (r *UserActBatch) Stop() error {
	return r.muster.Stop()
}

func (r *UserActBatch) Add(item entity.UserActivity[any]) error {
	r.muster.Work <- item

	return nil
}

type Batch struct {
	Client *UserActBatch
	Items  []entity.UserActivity[any]
}

func (r *Batch) Add(item any) {
	r.Items = append(r.Items, item.(entity.UserActivity[any]))
}

func (r *Batch) Fire(notifier muster.Notifier) {
	defer notifier.Done()

	err := r.Client.UserActRepo.InsertMany(r.Items)
	if err != nil {
		facade.Log().Error("insert many user activities: %w", err)
	}
}
