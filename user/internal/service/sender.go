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
	authMailTemplate *template.Template
	mailSender       MailSender
}

func NewSenderService(mailSender MailSender, tmpl *template.Template) *SenderService {
	return &SenderService{
		authMailTemplate: tmpl,
		mailSender:       mailSender,
	}
}

func (f *SenderService) SendEmailWithAuthorizationCode(ctx context.Context, to []string, subject string, code int) error {
	buffer := &bytes.Buffer{}

	err := f.authMailTemplate.Execute(buffer, struct {
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
