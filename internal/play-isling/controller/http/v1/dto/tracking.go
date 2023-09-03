package dto

import (
	common_entity "isling-be/internal/common/entity"
	"isling-be/internal/play-isling/usecase"
	"time"
)

type CreateAction struct {
	Type     string  `json:"type" validate:"required"`
	ObjectID *string `json:"objectId"`
}

func (action *CreateAction) ToCreateActionRequest(accountID common_entity.AccountID) *usecase.CreateActionRequest {
	return &usecase.CreateActionRequest{
		AccountID: accountID,
		Type:      action.Type,
		ObjectID:  action.ObjectID,
		Timestamp: time.Now(),
	}
}
