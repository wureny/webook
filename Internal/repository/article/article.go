package article

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/wureny/webook/webook/Internal/domain"
	"github.com/wureny/webook/webook/Internal/repository"
	"github.com/wureny/webook/webook/Internal/repository/cache"
	"time"

	//"github.com/wureny/webook/webook/Internal/repository/dao"
	dao "github.com/wureny/webook/webook/Internal/repository/dao/article"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	// Sync 存储并同步数据
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
	//FindById(ctx context.Context, id int64) domain.Article
}

type CachedArticleRepository struct {
	dao      dao.ArticleDAO
	userRepo repository.UserRepository

	// v1 操作两个 DAO
	readerDAO dao.ReaderDAO
	authorDAO dao.AuthorDAO

	// 耦合了 DAO 操作的东西
	// 正常情况下，如果你要在 repository 层面上操作事务
	// 那么就只能利用 db 开始事务之后，创建基于事务的 DAO
	// 或者，直接去掉 DAO 这一层，在 repository 的实现中，直接操作 db
	cache cache.ArticleCache
	db    *gorm.DB
	//cache cache.ArticleCache
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	//TODO implement me
	// 1. 存储到数据库
	// 2. 存储到缓存
	// 3. 返回 ID
	panic("implement me")
}

func (c *CachedArticleRepository) SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error {
	return c.dao.SyncStatus(ctx, id, author, uint8(status))
}

func (c *CachedArticleRepository) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	if offset == 0 && limit <= 100 {
		arts, err := c.cache.GetFirstPage(ctx, uid)
		//TODO: 此处缺少precache实现
		if err != nil {
			return nil, err
		}
		return arts, nil
	}
	arts, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	data := slice.Map[dao.Article, domain.Article](arts, func(idx int, src dao.Article) domain.Article {
		return c.toDomain(src)
	})
	return data, nil
}

func (c *CachedArticleRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	data, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.toDomain(data), nil
}

func (c *CachedArticleRepository) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	user, err := c.userRepo.FindById(ctx, art.AuthorId)
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id:   user.Id,
			Name: user.NickName,
		},
		Status: 0,
		Ctime:  time.Time{},
		Utime:  time.Time{},
	}
	return res, nil
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: int64(art.Author.Id),
	})
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: int64(art.Author.Id),
	})
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (c *CachedArticleRepository) toEntity(domainarticle domain.Article) dao.Article {
	return dao.Article{
		Id:       domainarticle.Id,
		Title:    domainarticle.Title,
		Content:  domainarticle.Content,
		AuthorId: domainarticle.Author.Id,
	}
}

func (repo *CachedArticleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
}
