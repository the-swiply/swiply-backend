package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
)

type MailSender interface {
	SendEmail(ctx context.Context, to []string, subject string, body []byte) error
}

type SenderService struct {
	mailSender MailSender
}

func NewSenderService(mailSender MailSender) *SenderService {
	return &SenderService{
		mailSender: mailSender,
	}
}

func (f *SenderService) SendEmailWithAuthorizationCode(ctx context.Context, to []string, subject string, code int) error {
	buffer := &bytes.Buffer{}
	tmpl, err := template.ParseFiles(authTemplatePath)
	if err != nil {
		return fmt.Errorf("can't parse template: %w", err)
	}

	err = tmpl.Execute(buffer, struct {
		Code int
	}{
		Code: code,
	})
	if err != nil {
		return fmt.Errorf("can't execute template: %w", err)
	}

	err = f.mailSender.SendEmail(ctx, to, subject, buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}
