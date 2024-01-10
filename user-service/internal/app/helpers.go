package app

import "time"

func (a *App) calculateMaxAuthCodeTTLForResend() time.Duration {
	return time.Duration(a.cfg.App.AuthCodeTTLMinutes-a.cfg.App.AuthCodeSendingMinRetryTimeMinutes) * time.Minute
}
