package request

import (
	"isling-be/internal/account/entity"
)

type (
	CreateProfileReq struct {
		FirstName   string
		LastName    string
		Gender      entity.GenderIdentity
		DateOfBirth string
	}
)
