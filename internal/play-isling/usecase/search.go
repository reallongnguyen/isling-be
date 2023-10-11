package usecase

import (
	"context"
	common_en "isling-be/internal/common/entity"
	"isling-be/pkg/facade"
)

type SearchUC struct {
	searchRepo SearchRepository
}

var _ SearchUsecase = (*SearchUC)(nil)

func NewSearchUC(searchRepo SearchRepository) *SearchUC {
	return &SearchUC{
		searchRepo: searchRepo,
	}
}

func (r *SearchUC) Search(c context.Context, accountID common_en.AccountID, req *SearchRequest) (*SearchResponse, error) {
	errChan := make(chan error)
	totalChan := make(chan int)

	go func() {
		total, err := r.searchRepo.GetTotalRoomMatches(c, accountID, req)
		if err != nil {
			errChan <- err
		}

		totalChan <- total
	}()

	rooms, err := r.searchRepo.SearchRoom(c, accountID, req)
	if err != nil {
		facade.Log().Error("search room: %w", err)

		return nil, err
	}

	total := 0

	select {
	case err = <-errChan:
		facade.Log().Error("get total: %w", err)

		return nil, err
	case _total := <-totalChan:
		total = _total
	}

	roomCollection := common_en.NewCollection(rooms, req.Offset, req.Limit, total)

	searchRes := SearchResponse{
		Rooms: &roomCollection,
	}

	return &searchRes, nil
}
