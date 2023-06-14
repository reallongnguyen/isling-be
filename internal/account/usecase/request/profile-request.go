package request

import (
	"isling-be/internal/account/entity"
	"strings"
)

type (
	CreateProfileReq struct {
		FirstName   string
		LastName    *string
		Gender      entity.GenderIdentity
		DateOfBirth string
	}
)

func (req *CreateProfileReq) Normalize() *CreateProfileReq {
	req.FirstName = strings.TrimSpace(req.FirstName)

	if req.LastName != nil {
		*req.LastName = strings.TrimSpace(*req.LastName)
	}

	return req
}
