package domain

import (
	"time"
)

type User struct {
	Id         int64
	Email      string
	Password   string
	Ctime      time.Time
	Birthday   string
	UserName   string
	Bio        string
	Phone      string
	WechatInfo WechatInfo
	NickName   string
}

/*type Address struct {
}
*/
