package provider

import (
	"fmt"

	"github.com/premchandkpc/large-microsrvices-system/services/notification-service/internal/config"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type SMTPProvider struct {
	cfg    *config.Config
	logger *zap.Logger
	dialer *gomail.Dialer
}

func NewSMTPProvider(cfg *config.Config, logger *zap.Logger) *SMTPProvider {
	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword)
	return &SMTPProvider{cfg: cfg, logger: logger, dialer: d}
}

func (p *SMTPProvider) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", p.cfg.SMTPFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := p.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("sending email: %w", err)
	}
	p.logger.Info("email sent", zap.String("to", to), zap.String("subject", subject))
	return nil
}
