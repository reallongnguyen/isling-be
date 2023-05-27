package entity

import (
	common_entity "isling-be/internal/common/entity"
	"time"
)

type RefreshTokens struct {
	ID             int
	AccountID      common_entity.AccountID
	EncryptedToken string
	Revoked        bool
	CreatedAt      time.Time
}
