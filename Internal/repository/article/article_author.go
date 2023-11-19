package article

import (
	"context"
	"github.com/wureny/webook/webook/Internal/domain"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}
