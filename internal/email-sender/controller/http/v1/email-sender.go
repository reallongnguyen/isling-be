package v1

import (
	"isling-be/internal/email-sender/entity"
	"isling-be/internal/email-sender/usecase"
	"isling-be/pkg/logger"

	"github.com/labstack/echo/v4"
)

type EmailSenderRouter struct {
	log           logger.Interface
	emailSenderUC usecase.EmailSenderUsecase
}

func NewEmailSenderRouter(g *echo.Group, log logger.Interface, emailSenderUC usecase.EmailSenderUsecase) *EmailSenderRouter {
	router := EmailSenderRouter{
		emailSenderUC: emailSenderUC,
		log:           log,
	}

	group := g.Group("/email-senders")
	group.GET("", router.sendEmail)

	return &router
}

func (router *EmailSenderRouter) sendEmail(c echo.Context) error {
	router.emailSenderUC.SendTemplateMail(c.Request().Context(), "reallongnguyen@gmail.com", "Test", entity.SignUpTemplate, nil)

	return nil
}
