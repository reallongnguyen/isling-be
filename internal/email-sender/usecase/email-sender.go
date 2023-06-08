package usecase

import (
	"context"
	"isling-be/internal/email-sender/entity"
	"isling-be/pkg/logger"
)

type EmailSenderUC struct {
	log         logger.Interface
	emailClient EmailClient
}

var _ EmailSenderUsecase = (*EmailSenderUC)(nil)

func NewEmailSenderUC(log logger.Interface, emailClient EmailClient) *EmailSenderUC {
	return &EmailSenderUC{
		log:         log,
		emailClient: emailClient,
	}
}

func (uc *EmailSenderUC) SendTemplateMail(ctx context.Context, to, subject string, templateName entity.EmailTemplateName, data interface{}) error {
	uc.emailClient.SendMail(ctx, "no-reply@isling.me", &entity.SimpleEmail{
		To:      "reallongnguyen@gmail.com",
		Subject: "Create account confirmation",
		Body:    "Hi Long!",
	})

	return nil
}
