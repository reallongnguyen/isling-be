package repo

import (
	"context"
	"fmt"
	"isling-be/internal/email-sender/entity"
	"isling-be/internal/email-sender/usecase"
	"strings"
)

type EmailClientFaker struct {
}

var _ usecase.EmailClient = (*EmailClientFaker)(nil)

func NewEmailClientFaker() *EmailClientFaker {
	return &EmailClientFaker{}
}

func (client *EmailClientFaker) SendMail(ctx context.Context, from string, email *entity.SimpleEmail) error {
	toAdds := email.GetDestinationAddresses()
	body, err := email.GetMSG()

	if err != nil {
		return err
	}

	bodyString := string(body)

	fmt.Println("--------------------------------")
	fmt.Println("Send an email")
	fmt.Println("From: " + from)
	fmt.Println("To: " + strings.Join(toAdds, ","))
	fmt.Println(bodyString)
	fmt.Println("--------------------------------")

	return nil
}
