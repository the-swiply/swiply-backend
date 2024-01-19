package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/pkg/auf"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/user/internal/domain"
	"math/rand"
	"time"
)

const (
	authTemplatePath = "templates/auth_mail.html"
	authSubject      = "Swiply Authorization"
)

type AuthCodeCache interface {
	StoreAuthCode(ctx context.Context, email string, code int) error
	GetAuthCode(ctx context.Context, email string) (string, error)
	GetAuthCodeTTL(ctx context.Context, email string) (time.Duration, error)
	DeleteAuthCode(ctx context.Context, email string) error
}

type TokenStorage interface {
	StoreFingerprint(ctx context.Context, token string, fingerprint string) error
	GetFingerprint(ctx context.Context, token string) (string, error)
	DeleteFingerprint(ctx context.Context, token string) error
}

type TaskScheduler interface {
	ScheduleEmailSend(ctx context.Context, info domain.SendAuthCodeInfo) error
}

type UserService struct {
	cfg           UserConfig
	codeCache     AuthCodeCache
	tokenStorage  TokenStorage
	taskScheduler TaskScheduler

	invalidCodeAttempts *AttemptsRecorder
}

func NewUserService(cfg UserConfig, codeCache AuthCodeCache, tokenStorage TokenStorage, taskScheduler TaskScheduler) *UserService {
	auf.SetSecret([]byte(cfg.TokenSecret))

	return &UserService{
		cfg:                 cfg,
		codeCache:           codeCache,
		tokenStorage:        tokenStorage,
		taskScheduler:       taskScheduler,
		invalidCodeAttempts: newAttemptsRecorder(cfg.MaxInvalidCodeAttempts),
	}
}

func (u *UserService) SendAuthorizationCode(ctx context.Context, to string) error {
	ttl, err := u.codeCache.GetAuthCodeTTL(ctx, to)
	if err != nil {
		return fmt.Errorf("can't get auth code ttl: %w", err)
	}

	if ttl > 0 && ttl > u.cfg.MaxAuthCodeTTLForResend {
		return domain.ErrResendIsNotAllowed
	}

	code := int(rand.Int63n(9000) + 1000)

	err = u.codeCache.StoreAuthCode(ctx, to, code)
	if err != nil {
		return fmt.Errorf("can't store auth code: %w", err)
	}

	u.invalidCodeAttempts.clearAttempts(to)

	err = u.taskScheduler.ScheduleEmailSend(ctx, domain.SendAuthCodeInfo{
		To:      []string{to},
		Subject: authSubject,
		Code:    code,
	})
	if err != nil {
		if err := u.codeCache.DeleteAuthCode(ctx, to); err != nil {
			loggy.Errorln("can't delete code from cache:", err)
		}

		return fmt.Errorf("can't schedule code send: %w", err)
	}

	return nil
}

func (u *UserService) Login(ctx context.Context, email, code, fingerprint string) (domain.TokenPair, error) {
	storedCode, err := u.codeCache.GetAuthCode(ctx, email)
	if errors.Is(err, domain.ErrEntityIsNotExists) {
		return domain.TokenPair{}, domain.ErrCodeIsIncorrect
	}
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("can't get auth code: %w", err)
	}

	if storedCode != code {
		if ok := u.invalidCodeAttempts.addAttempt(email); !ok {
			err = u.codeCache.DeleteAuthCode(ctx, email)
			if err != nil {
				loggy.Errorln("can't delete auth code after too much attempts")
			}

			return domain.TokenPair{}, domain.ErrTooMuchAttempts
		}

		return domain.TokenPair{}, domain.ErrCodeIsIncorrect
	}

	u.invalidCodeAttempts.clearAttempts(email)

	err = u.codeCache.DeleteAuthCode(ctx, email)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("can't delete auth code: %w", err)
	}

	id := u.generateUUIDByEmail(email)

	tokenPair, err := u.generateTokenPair(id.String(), fingerprint)
	if err != nil {
		return domain.TokenPair{}, err
	}

	err = u.tokenStorage.StoreFingerprint(ctx, tokenPair.RefreshToken, fingerprint)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("can't store fingerprint: %w", err)
	}

	return tokenPair, nil
}

func (u *UserService) RefreshTokens(ctx context.Context, refreshToken string, fingerprint string) (domain.TokenPair, error) {
	claims, err := auf.ValidateJWTAndExtractClaims(refreshToken)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("%w: %w", domain.ErrValidateToken, err)
	}

	storedFingerprint, err := u.tokenStorage.GetFingerprint(ctx, refreshToken)
	if errors.Is(err, domain.ErrEntityIsNotExists) {
		return domain.TokenPair{}, domain.ErrValidateToken
	}
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("can't get fingerprint from storage: %w", err)
	}

	id, ok := claims["id"].(string)
	if !ok {
		return domain.TokenPair{}, fmt.Errorf("%w: no email presented", domain.ErrValidateToken)
	}

	if storedFingerprint != fingerprint {
		if err := u.tokenStorage.DeleteFingerprint(ctx, refreshToken); err != nil {
			loggy.Errorln("can't delete fingerprint from storage:", err)
		}

		return domain.TokenPair{}, domain.ErrValidateToken
	}

	tokenPair, err := u.generateTokenPair(id, fingerprint)
	if err != nil {
		return domain.TokenPair{}, err
	}

	err = u.tokenStorage.StoreFingerprint(ctx, tokenPair.RefreshToken, fingerprint)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("can't store fingerprint: %w", err)
	}

	err = u.tokenStorage.DeleteFingerprint(ctx, refreshToken)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("can't delete fingerprint from storage: %w", err)
	}

	return tokenPair, nil
}

func (u *UserService) generateTokenPair(id string, fingerprint string) (domain.TokenPair, error) {
	accessToken, err := auf.GenerateAccessJWT(auf.JWTAccessProperties{
		ID:          id,
		TTL:         u.cfg.AccessTokenTTL,
		Fingerprint: fingerprint,
	})
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("can't generate access token: %w", err)
	}

	refreshToken, err := auf.GenerateRefreshJWT(auf.JWTRefreshProperties{
		ID:  id,
		TTL: u.cfg.AccessTokenTTL,
	})
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("can't generate refresh token: %w", err)
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *UserService) generateUUIDByEmail(email string) uuid.UUID {
	return uuid.NewSHA1(u.cfg.UUIDNamespace, []byte(email))
}
