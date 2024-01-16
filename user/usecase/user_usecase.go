package usecase

import (
	"context"
	"time"

	"github.com/Kelvinmijaya/kelvin-rest-api/domain"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(a domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       a,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) Login(c context.Context, email string, password string, m *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err = u.userRepo.Login(ctx, email, password, m)
	if err != nil {
		return
	}
	return
}
