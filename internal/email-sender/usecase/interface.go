package usecase

import (
	"context"
	"isling-be/internal/email-sender/entity"
)

type (
	EmailSenderUsecase interface {
		SendTemplateMail(context.Context, string, string, entity.EmailTemplateName, interface{}) error
	}

	EmailClient interface {
		SendMail(context.Context, string, *entity.SimpleEmail) error
	}
)
