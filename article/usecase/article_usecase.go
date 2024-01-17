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

func (a *articleUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	return
}

func (a *articleUsecase) GetByID(c context.Context, id int64) (res domain.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (a *articleUsecase) Update(c context.Context, ar *domain.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *articleUsecase) Store(c context.Context, ar *domain.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	ar.CreatedAt = time.Now()
	err = a.articleRepo.Store(ctx, ar)
	return
}

func (a *articleUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (domain.Article{}) {
		return domain.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}
