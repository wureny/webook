package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{db: db}
}

func (G GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	t := time.Now().UnixMilli()
	art.Ctime = t
	art.Utime = t
	err := G.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

// UpdateById 只更新标题、内容和状态
func (G GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	t := time.Now().UnixMilli()
	res := G.db.Model(&Article{}).WithContext(ctx).Where("id=? AND author_id=?", art.Id, art.AuthorId).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
		"utime":   t,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		fmt.Println("状态未更新！")
	}
	return nil
}

func (G GORMArticleDAO) GetByAuthor(ctx context.Context, author int64, offset int, limit int) ([]Article, error) {
	var arts []Article
	err := G.db.Model(&Article{}).WithContext(ctx).Where("author_id=?", author).Offset(offset).Limit(limit).Order("utime DESC").Find(&arts).Error
	if err != nil {
		return nil, err
	}
	return arts, nil
}

func (G GORMArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := G.db.Model(&Article{}).WithContext(ctx).Where("id=?", id).Find(&art).Error
	if err != nil {
		return Article{}, err
	}
	return art, nil
}

func (G GORMArticleDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	var pub PublishedArticle
	err := G.db.WithContext(ctx).
		Where("id = ?", id).
		First(&pub).Error
	return pub, err
}

func (G GORMArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (G GORMArticleDAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}
