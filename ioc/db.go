package ioc

import (
	"github.com/spf13/viper"
	"github.com/wureny/webook/webook/Internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*
func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	//	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:30002)/webook"))
	fmt.Println(config.Config.DB.DSN)
	fmt.Println("fuck")
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
*/

func InitDB() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var cfg = Config{
		//相当于设置了一个默认字段；
		DSN: "root:root@tcp(localhost:13316)/webook_default",
	}
	err := viper.UnmarshalKey("db", &cfg)
	db, err := gorm.Open(mysql.Open(cfg.DSN))
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
