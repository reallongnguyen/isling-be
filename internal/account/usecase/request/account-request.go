package request

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
