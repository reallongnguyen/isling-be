package entity

import cme "isling-be/internal/common/entity"

type PlayUser struct {
	ID                  int64
	AccountID           cme.AccountID
	RecentlyJoinedRooms []int64
}
