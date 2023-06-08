package entity

import (
	"html/template"
	"strings"
)

type EmailTemplateName string

const (
	EmailTemplate = `
To: {{.To}}\r\n
Subject: {{.Subject}}\r\n
\r\n
{{.Body}}\r\n
`

	SignUpTemplate EmailTemplateName = "SignUpTemplate"
)

type SMTPEmail interface {
	GetDestinationAddresses() []string
	GetMSG() ([]byte, error)
}

type SimpleEmail struct {
	To      string
	Subject string
	Body    string
}

var _ SMTPEmail = (*SimpleEmail)(nil)

func (se *SimpleEmail) GetDestinationAddresses() []string {
	return []string{se.To}
}

func (se *SimpleEmail) GetMSG() ([]byte, error) {
	t, err := template.New("simple email template").Parse(EmailTemplate)

	if err != nil {
		return []byte{}, err
	}

	strBuilder := new(strings.Builder)

	if err = t.Execute(strBuilder, se); err != nil {
		return []byte{}, err
	}

	return []byte(strBuilder.String()), nil
}
