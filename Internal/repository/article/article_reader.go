package article

import (
	"context"
	"github.com/wureny/webook/webook/Internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}
