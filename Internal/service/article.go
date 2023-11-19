package service

import (
	"context"
	"github.com/wureny/webook/webook/Internal/domain"
	"github.com/wureny/webook/webook/Internal/repository/article"
	"time"
)

type ArticleService interface {
	Publish(ctx context.Context, article domain.Article) (int64, error)
	PublishV1(ctx context.Context, article domain.Article) (int64, error)
	Save(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
}

type articleService struct {
	repo article.ArticleRepository

	// V1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
	//	l      logger.LoggerV1
}

func (svc *articleService) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	// 另一个选项，在这里组装 Author，调用 UserService
	return svc.repo.GetPublishedById(ctx, id)
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *articleService) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, uid, offset, limit)
}

func (a *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	// art.Status = domain.ArticleStatusPrivate 然后直接把整个 art 往下传
	return a.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	// 制作库
	//id, err := a.repo.Create(ctx, art)
	//// 线上库呢？
	//a.repo.SyncToLiveDB(ctx, art)
	return a.repo.Sync(ctx, art)
}

func (a articleService) PublishV1(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	if id > 0 {
		err = a.author.Update(ctx, article)
	} else {
		id, err = a.author.Create(ctx, article)
	}
	if err != nil {
		return 0, err
	}
	article.Id = id
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * time.Duration(i))
		id, err = a.reader.Save(ctx, article)
		if err == nil {
			return id, nil
		}
		//TODO 日志部分
	}
	if err != nil {
		//TODO 日志部分
		// 接入你的告警系统，手工处理一下
		// 走异步，我直接保存到本地文件
		// 走 Canal
		// 打 MQ
		// 以上是面试技巧
	}
	return id, err
}

func (a articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

func NewArticleService(repo article.ArticleRepository, auth article.ArticleAuthorRepository, reader article.ArticleReaderRepository) ArticleService {
	return &articleService{
		repo:   repo,
		author: auth,
		reader: reader,
	}
}

func (a *articleService) update(ctx context.Context, art domain.Article) error {
	// 只要你不更新 author_id
	// 但是性能比较差
	//artInDB := a.repo.FindById(ctx, art.Id)
	//if art.Author.Id != artInDB.Author.Id {
	//	return errors.New("更新别人的数据")
	//}
	return a.repo.Update(ctx, art)
}
