package usecase

import (
	"context"
	"errors"
	"net/mail"
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

func (u *userUsecase) Login(c context.Context, m *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// check valid email and password
	var ok bool
	if ok, err = loginValidator(m); !ok {
		return err
	}

	err = u.userRepo.Login(ctx, m)
	if err != nil {
		return err
	}

	return
}

func loginValidator(m *domain.User) (bool, error) {
	_, err := mail.ParseAddress(m.Email)
	if err != nil {
		return false, errors.New("email is not valid")
	}

	if len([]rune(m.Password)) < 4 {
		return false, errors.New("password minimum 4 character")
	}
	return true, nil
}
