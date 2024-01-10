package mailer

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
)

const (
	defaultHeaders = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
)

type SMTPClient struct {
	cfg SMTPConfig
}

func NewSMTPClient(cfg SMTPConfig) (*SMTPClient, error) {
	_, _, err := net.SplitHostPort(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("invalid addr: %w", err)
	}

	return &SMTPClient{
		cfg: cfg,
	}, nil
}

func (s *SMTPClient) SendEmail(ctx context.Context, to []string, subject string, body []byte) error {
	smtpMsg := fmt.Sprintf(
		`Subject: %s
%s

%s
`, subject, defaultHeaders, body)

	host, _, _ := net.SplitHostPort(s.cfg.Addr)

	auth := smtp.PlainAuth("", s.cfg.SenderEmail, s.cfg.SenderPassword, host)

	errCh := make(chan error)
	go func() {
		errCh <- smtp.SendMail(s.cfg.Addr, auth, s.cfg.SenderEmail, to, []byte(smtpMsg))
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		if err != nil {
			return err
		}
	}

	return nil
}
