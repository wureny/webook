package repository

import (
	"context"
	"github.com/wureny/webook/webook/Internal/domain"
	"github.com/wureny/webook/webook/Internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (re *UserRepository) Create(ctx context.Context, u domain.User) error {
	return re.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindById(int64) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	// 找到了回写 cache
}
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}
func (r *UserRepository) UpdateInfo(ctx context.Context, u domain.User) error {
	err := r.dao.Update(ctx, u.Birthday, u.Bio, u.UserName, u.Id)
	if err != nil {
		return err
	}
	return nil
}
func (r *UserRepository) GetUserInfo(ctx context.Context, id uint64) (domain.User, error) {
	s, err := r.dao.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	u.Email = s.Email
	u.Birthday = s.Birthday
	u.Bio = s.Bio
	u.UserName = s.UserName
	return u, nil
}
