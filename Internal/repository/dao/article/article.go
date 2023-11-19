package article

/*
import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (uint64, error)
	UpdateById(ctx context.Context, article Article) error
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) *GORMArticleDAO {
	return &GORMArticleDAO{db: db}
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (uint64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return uint64(art.Id), err
}
/*
func (dao *GORMArticleDAO) UpdateById(ctx context.Context, article Article) error {
	article.Utime = time.Now().UnixMilli()
	res := dao.db.WithContext(ctx).Model(&article).
		Where("id=? AND author_id=?", article.Id, article.AuthorId).
		Updates(map[string]interface{}{
			"title":   article.Title,
			"content": article.Content,
			"utime":   article.Utime,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		//dangerousDBOp.Count(1)
		// 补充一点日志
		return fmt.Errorf("更新失败，可能是创作者非法 id %d, author_id %d",
			article.Id, article.AuthorId)
	}
	return nil
}
*/
