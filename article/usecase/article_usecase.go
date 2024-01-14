package usecase

import (
	"context"
	"time"

	"github.com/Kelvinmijaya/kelvin-rest-api/domain"
)

type articleUsecase struct {
	articleRepo    domain.ArticleRepository
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of domain.ArticleUsecase interface
func NewArticleUsecase(a domain.ArticleRepository, timeout time.Duration) domain.ArticleUsecase {
	return &articleUsecase{
		articleRepo:    a,
		contextTimeout: timeout,
	}
}

func (a *articleUsecase) Store(c context.Context, m *domain.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	// validation article title exist or not
	// existedArticle, _ := a.GetByTitle(ctx, m.Title)
	// if existedArticle != (domain.Article{}) {
	// 	return domain.ErrConflict
	// }

	err = a.articleRepo.Store(ctx, m)
	return
}