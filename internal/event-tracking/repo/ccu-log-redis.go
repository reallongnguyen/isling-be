package repo

import (
	"context"
	"isling-be/internal/event-tracking/entity"
	"isling-be/internal/event-tracking/usecase"
	"isling-be/pkg/redis"
	"strconv"
	"time"

	"golang.org/x/exp/slices"
)

type CCULogRepo struct {
	red *redis.Redis
}

func NewCCULogRepo(red *redis.Redis) usecase.CCULogRepository {
	return &CCULogRepo{
		red: red,
	}
}

func (r *CCULogRepo) InsertMany(items []*entity.CCULog) error {
	if len(items) == 0 {
		return nil
	}

	ctx := context.Background()

	cluster := make(map[string][]string)

	for idx := range items {
		bucket := "ccu:" + strconv.FormatInt(items[idx].GetUnixMin(), 10)

		if _, found := cluster[bucket]; !found {
			cluster[bucket] = make([]string, 0)
		}

		cluster[bucket] = append(cluster[bucket], items[idx].GetKey())
	}

	for k, v := range cluster {
		slices.Sort(v)
		set := slices.Compact(v)

		r.red.PFAdd(ctx, k, set)
	}

	return nil
}

func (r *CCULogRepo) CountCCU(timestamp time.Time, windowSize uint) (int64, error) {
	min := timestamp.Unix() / 60

	ctx := context.Background()

	buckets := make([]string, windowSize)

	for i := 0; i < int(windowSize); i++ {
		buckets[i] = "ccu:" + strconv.FormatInt(min-int64(i), 10)
	}

	return r.red.PFCount(ctx, buckets...).Result()
}
