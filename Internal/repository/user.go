package repository

import (
	"context"
	"github.com/wureny/webook/webook/Internal/domain"
	"github.com/wureny/webook/webook/Internal/repository/cache"
	"github.com/wureny/webook/webook/Internal/repository/dao"
	"sync"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}

func (re *UserRepository) Create(ctx context.Context, u domain.User) error {
	return re.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindById(ctx context.Context, id uint64) (domain.User, error) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	// 找到了回写 cache
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		// 必然是有数据
		return u, nil
	}
	ue, err := r.dao.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 我这里怎么办？
			// 打日志，做监控
			//return domain.User{}, err
		}
	}()
	wg.Wait()
	return u, err
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
	e := r.cache.Set(ctx, u)
	if e != nil {
		return e
	}
	return nil
}
func (r *UserRepository) GetUserInfo(ctx context.Context, id uint64) (domain.User, error) {
	us, err := r.cache.Get(ctx, id)
	if err == nil {
		// 必然是有数据
		return us, nil
	}
	s, err := r.dao.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	u.Email = s.Email
	u.Birthday = s.Birthday
	u.Bio = s.Bio
	u.UserName = s.UserName
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 我这里怎么办？
			// 打日志，做监控
			//return domain.User{}, err
		}
	}()
	wg.Wait()
	return u, nil
}
