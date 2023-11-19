package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("手机或邮箱冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type DBProvider func() *gorm.DB

type GORMUserDAO struct {
	db *gorm.DB
	p  DBProvider
}

func NewUserDAOProvider(p DBProvider) UserDAO {
	return &GORMUserDAO{
		p: p,
	}
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	//创建时间：毫秒数
	Ctime int64
	//更新时间：毫秒数
	Utime    int64
	Birthday string
	UserName string
	Bio      string
	//phone不应该设置为unique，因为假如很多用户通过email来创建，未填phone，那都是空字符串，会冲突
	//email同理
	//所以没有就是null，要把空字符串转为null
	Phone sql.NullString `gorm:"unique"`
	//下面的做法问题：要解开引用，要判断是否为空
	//Phone *string
}

func (dao *GORMUserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	fmt.Println(u)
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突
			return ErrUserDuplicate
		}
	}
	return err
}
func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	//err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	if err == nil {
		return u, err
	}
	return u, nil
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	//err := dao.db.WithContext(ctx).First(&u, "email = ?", email).Error
	if err == nil {
		return u, err
	}
	return u, nil
}

func (dao *GORMUserDAO) Update(ctx context.Context, Bir string, Bio string, username string, id int64) error {
	var u User
	err := dao.db.WithContext(ctx).Where("Id=?", id).First(&u).Error
	if err != nil {
		return err
	}
	u.Birthday = Bir
	u.Bio = Bio
	u.UserName = username
	dao.db.Save(&u)
	return nil
}
func (dao *GORMUserDAO) GetUser(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("Id=?", id).First(&u).Error
	if err != nil {
		return User{}, err
	}
	return u, nil
}

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	Update(ctx context.Context, Bir string, Bio string, username string, id int64) error
	GetUser(ctx context.Context, id int64) (User, error)
}
