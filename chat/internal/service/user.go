package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/the-swiply/swiply-backend/pkg/auf"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/user/internal/entity"
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
	ScheduleEmailSend(ctx context.Context, info entity.SendAuthCodeInfo) error
}

type UserService struct {
	cfg           UserConfig
	codeCache     AuthCodeCache
	tokenStorage  TokenStorage
	taskScheduler TaskScheduler
}

func NewUserService(cfg UserConfig, codeCache AuthCodeCache, tokenStorage TokenStorage, taskScheduler TaskScheduler) *UserService {
	return &UserService{
		cfg:           cfg,
		codeCache:     codeCache,
		tokenStorage:  tokenStorage,
		taskScheduler: taskScheduler,
	}
}

func (u *UserService) SendAuthorizationCode(ctx context.Context, to string) error {
	code := int(rand.Int63n(9000) + 1000)

	ttl, err := u.codeCache.GetAuthCodeTTL(ctx, to)
	if err != nil {
		return fmt.Errorf("can't get auth code ttl: %w", err)
	}

	if ttl > 0 && ttl > u.cfg.MaxAuthCodeTTLForResend {
		return ErrResendIsNotAllowed
	}

	err = u.codeCache.StoreAuthCode(ctx, to, code)
	if err != nil {
		return fmt.Errorf("can't store auth code: %w", err)
	}

	err = u.taskScheduler.ScheduleEmailSend(ctx, entity.SendAuthCodeInfo{
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

func (u *UserService) Login(ctx context.Context, email, code, fingerprint string) (entity.TokenPair, error) {
	storedCode, err := u.codeCache.GetAuthCode(ctx, email)
	if errors.Is(err, ErrEntityIsNotExists) {
		return entity.TokenPair{}, ErrCodeIsIncorrect
	}
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("can't get auth code: %w", err)
	}

	if storedCode != code {
		return entity.TokenPair{}, ErrCodeIsIncorrect
	}

	err = u.codeCache.DeleteAuthCode(ctx, email)
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("can't delete auth code: %w", err)
	}

	tokenPair, err := u.generateTokenPair(email, fingerprint)
	if err != nil {
		return entity.TokenPair{}, err
	}

	err = u.tokenStorage.StoreFingerprint(ctx, tokenPair.RefreshToken, fingerprint+email)
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("can't store fingerprint: %w", err)
	}

	return tokenPair, nil
}

func (u *UserService) RefreshTokens(ctx context.Context, refreshToken string, fingerprint string) (entity.TokenPair, error) {
	claims, err := auf.ValidateJWTAndExtractClaims(refreshToken, []byte(u.cfg.TokenSecret))
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("%w: %w", ErrValidateToken, err)
	}

	storedFingerprint, err := u.tokenStorage.GetFingerprint(ctx, refreshToken)
	if errors.Is(err, ErrEntityIsNotExists) {
		return entity.TokenPair{}, ErrValidateToken
	}
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("can't get fingerprint from storage: %w", err)
	}

	email, ok := claims["user"].(string)
	if !ok {
		return entity.TokenPair{}, fmt.Errorf("%w: no email presented", ErrValidateToken)
	}

	if storedFingerprint != fingerprint+email {
		if err := u.tokenStorage.DeleteFingerprint(ctx, refreshToken); err != nil {
			loggy.Errorln("can't delete fingerprint from storage:", err)
		}

		return entity.TokenPair{}, ErrValidateToken
	}

	tokenPair, err := u.generateTokenPair(email, fingerprint)
	if err != nil {
		return entity.TokenPair{}, err
	}

	err = u.tokenStorage.StoreFingerprint(ctx, tokenPair.RefreshToken, fingerprint+email)
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("can't store fingerprint: %w", err)
	}

	err = u.tokenStorage.DeleteFingerprint(ctx, refreshToken)
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("can't delete fingerprint from storage: %w", err)
	}

	return tokenPair, nil
}

func (u *UserService) generateTokenPair(email string, fingerprint string) (entity.TokenPair, error) {
	accessToken, err := auf.GenerateAccessJWT(auf.JWTAccessProperties{
		User:        email,
		TTL:         u.cfg.AccessTokenTTL,
		Secret:      []byte(u.cfg.TokenSecret),
		Fingerprint: fingerprint,
	})
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("can't generate access token: %w", err)
	}

	refreshToken, err := auf.GenerateRefreshJWT(auf.JWTRefreshProperties{
		User:   email,
		TTL:    u.cfg.AccessTokenTTL,
		Secret: []byte(u.cfg.TokenSecret),
	})
	if err != nil {
		return entity.TokenPair{}, fmt.Errorf("can't generate refresh token: %w", err)
	}

	return entity.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
