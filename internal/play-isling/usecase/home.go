package usecase

import (
	"context"
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/entity"
	"isling-be/pkg/logger"
	"math"
	"strconv"

	"github.com/zhenghaoz/gorse/client"
	"golang.org/x/exp/slices"
)

type HomeUC struct {
	log          logger.Interface
	playUserRepo PlayUserRepository
	roomRepo     RoomRepository
	gorse        *client.GorseClient
}

var _ HomeUsecase = (*HomeUC)(nil)

func NewHomeUsecase(log logger.Interface, playUserRepo PlayUserRepository, roomRepo RoomRepository) HomeUsecase {
	return &HomeUC{
		log:          log,
		playUserRepo: playUserRepo,
		roomRepo:     roomRepo,
		gorse:        client.NewGorseClient("http://localhost:8087", ""),
	}
}

func (uc *HomeUC) Show(c context.Context, accountID common_entity.AccountID) (*HomePageResponse, error) {
	items, err := uc.gorse.GetItemRecommend(c, strconv.FormatInt(int64(accountID), 10), []string{}, "", "", 16, 0)
	if err != nil {
		return nil, err
	}

	recommendRoomID := make([]int64, 0, 16)
	scoreMap := make(map[int64]float64)

	for i, itemID := range items {
		id, convErr := strconv.Atoi(itemID)
		if convErr != nil {
			continue
		}

		recommendRoomID = append(recommendRoomID, int64(id))
		scoreMap[int64(id)] = 999 - float64(i)
	}

	scores, err := uc.gorse.GetItemPopular(c, "", 16, 0)
	if err != nil {
		return nil, err
	}

	for _, score := range scores {
		id, convErr := strconv.Atoi(score.Id)
		if convErr != nil {
			continue
		}

		recommendRoomID = append(recommendRoomID, int64(id))
		scoreMap[int64(id)] = math.Max(scoreMap[int64(id)], score.Score)
	}

	recommendRooms, err := uc.roomRepo.FindMany(c, &FindRoomFilter{IDIn: &recommendRoomID}, nil)
	if err != nil {
		return nil, err
	}

	recommendPublicRooms := make([]*entity.RoomPublic, 0, len(recommendRooms.Edges))

	for _, room := range recommendRooms.Edges {
		recommendPublicRooms = append(recommendPublicRooms, room.ToRoomPublic())
	}

	slices.SortFunc(recommendPublicRooms, func(a, b *entity.RoomPublic) int {
		return -int(math.Round(scoreMap[a.ID] - scoreMap[b.ID]))
	})

	collections := make([]*entity.RoomCollection, 0, 8)

	recentlyPublicRooms := uc.getRecentlyRooms(c, accountID)
	if len(recentlyPublicRooms) > 0 {
		collections = append(collections, &entity.RoomCollection{
			ID:    0,
			Name:  "Watch again",
			Rooms: recentlyPublicRooms,
		})
	}

	hpRes := HomePageResponse{
		ForYou: &entity.RoomCollection{
			ID:    0,
			Name:  "For You",
			Rooms: recommendPublicRooms,
		},
		Collections: collections,
	}

	return &hpRes, nil
}

func (uc *HomeUC) ShowGuest(c context.Context) (*HomePageResponse, error) {
	popularRoom, err := uc.getPopularRoom(c)
	if err != nil {
		return nil, err
	}

	hpRes := HomePageResponse{
		ForYou: &entity.RoomCollection{
			ID:    0,
			Name:  "For You",
			Rooms: *popularRoom,
		},
		Collections: []*entity.RoomCollection{},
	}

	return &hpRes, nil
}

func (uc *HomeUC) getPopularRoom(c context.Context) (*[]*entity.RoomPublic, error) {
	scores, err := uc.gorse.GetItemPopular(c, "", 16, 0)
	if err != nil {
		return nil, err
	}

	scoreMap := make(map[int64]float64)
	listRoomID := make([]int64, 0, 16)

	for _, score := range scores {
		id, convErr := strconv.Atoi(score.Id)
		if convErr != nil {
			continue
		}

		listRoomID = append(listRoomID, int64(id))
		scoreMap[int64(id)] = score.Score
	}

	popularRooms, err := uc.roomRepo.FindMany(c, &FindRoomFilter{IDIn: &listRoomID}, &Order{Field: "created_at", Direction: "desc"})
	if err != nil {
		return nil, err
	}

	newPublicRooms := make([]*entity.RoomPublic, 0, 16)

	for _, room := range popularRooms.Edges {
		newPublicRooms = append(newPublicRooms, room.ToRoomPublic())
	}

	slices.SortFunc(newPublicRooms, func(a, b *entity.RoomPublic) int {
		return -int(math.Round(scoreMap[a.ID] - scoreMap[b.ID]))
	})

	return &newPublicRooms, nil
}

func (uc *HomeUC) getRecentlyRooms(c context.Context, accountID common_entity.AccountID) []*entity.RoomPublic {
	playUser, err := uc.playUserRepo.GetOne(c, accountID)
	if err != nil {
		return nil
	}

	recentlyRooms, err := uc.roomRepo.FindMany(c, &FindRoomFilter{IDIn: &playUser.RecentlyJoinedRooms}, nil)
	if err != nil {
		return nil
	}

	recentlyPublicRooms := make([]*entity.RoomPublic, 0, len(recentlyRooms.Edges))

	for _, room := range recentlyRooms.Edges {
		recentlyPublicRooms = append(recentlyPublicRooms, room.ToRoomPublic())
	}

	return recentlyPublicRooms
}
