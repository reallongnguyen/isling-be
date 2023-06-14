package request

import "strings"

type (
	CreateAccountReq struct {
		Email    string
		Password string
	}

	ChangePasswordReq struct {
		OldPassword string
		NewPassword string
	}
)

func (req *CreateAccountReq) Normalize() *CreateAccountReq {
	req.Email = strings.ToLower(req.Email)

	return req
}
