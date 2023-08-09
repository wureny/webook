package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/wureny/webook/webook/Internal/domain"
	"github.com/wureny/webook/webook/Internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	fmt.Println("service.test")
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码了
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// DEBUG
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) Edit(ctx context.Context, u domain.User) error {
	err := svc.repo.UpdateInfo(ctx, u)
	if err != nil {
		return err
	}
	return nil
}
func (svc *UserService) GetUser(ctx context.Context, id uint64) (domain.User, error) {
	u, err := svc.repo.GetUserInfo(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
