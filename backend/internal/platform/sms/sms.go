package sms

import (
	"context"
	"log/slog"
)

// Sender delivers one-time passwords to a phone number.
type Sender interface {
	SendOTP(ctx context.Context, phone, code string) error
}

// DevSender logs the OTP instead of sending an SMS. Used in local development.
type DevSender struct {
	logger *slog.Logger
}

func NewDevSender(logger *slog.Logger) *DevSender {
	if logger == nil {
		logger = slog.Default()
	}
	return &DevSender{logger: logger}
}

func (d *DevSender) SendOTP(_ context.Context, phone, code string) error {
	d.logger.Info("dev sms stub: otp generated", "phone", phone, "otp", code)
	return nil
}
