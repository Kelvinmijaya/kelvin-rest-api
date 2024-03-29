package domain

import (
	"context"
	"time"
)

// Article is representing the Article data struct
type Article struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required"`
	URL       string    `json:"url" validate:"required"`
	Content   string    `json:"content" validate:"required"`
	Type      string    `json:"type" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ArticleUsecase represent the article's usecases
type ArticleUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Article, string, error)
	GetByID(ctx context.Context, id int64) (Article, error)
	Update(ctx context.Context, ar *Article) error
	Store(context.Context, *Article) error
	Delete(ctx context.Context, id int64) error
}

// ArticleRepository represent the article's repository contract
type ArticleRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Article, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (Article, error)
	Update(ctx context.Context, ar *Article) error
	Store(ctx context.Context, a *Article) error
	Delete(ctx context.Context, id int64) error
}
