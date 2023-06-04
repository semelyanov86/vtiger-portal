package service

import (
	"context"
	"github.com/pquerna/otp/totp"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"sync"
)

type AuthService struct {
	repo   repository.Users
	wg     *sync.WaitGroup
	cache  cache.Cache
	config config.Config
}

func NewAuthService(repo repository.Users, wg *sync.WaitGroup, cache cache.Cache, config config.Config) AuthService {
	return AuthService{repo: repo, wg: wg, cache: cache, config: config}
}

type OTPInput struct {
	UserId int64  `json:"id"`
	Token  string `json:"token"`
}

type OtpRegisterResult struct {
	Base32     string `json:"base32"`
	OtpauthUrl string `json:"otpauth_url"`
}

func (a AuthService) GenerateOtp(ctx context.Context, input OTPInput) (OtpRegisterResult, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      a.config.Otp.Issuer,
		AccountName: a.config.Otp.AccountName,
		SecretSize:  a.config.Otp.SecretSize,
	})

	if err != nil {
		return OtpRegisterResult{}, e.Wrap("can not generate OTP", err)
	}

	_, err = a.repo.GetById(ctx, input.UserId)

	if err != nil {
		return OtpRegisterResult{}, ErrUserNotFound
	}
	otpResult := OtpRegisterResult{
		Base32:     key.Secret(),
		OtpauthUrl: key.URL(),
	}
	err = a.repo.SaveOtp(ctx, key.Secret(), key.URL(), input.UserId)
	return otpResult, err
}
