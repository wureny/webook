package dao

import (
	"github.com/wureny/webook/webook/Internal/repository/dao/article"
	"gorm.io/gorm"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{},
		&Interactive{},
		&UserLikeBiz{},
		&Collection{},
		&UserCollectionBiz{},
		&article.Article{},
		&article.PublishedArticle{},
	)
}
