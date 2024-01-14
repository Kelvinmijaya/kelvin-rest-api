package postgres

import (
	"context"
	"database/sql"
	"log"

	"github.com/Kelvinmijaya/kelvin-rest-api/domain"
)

type postgresArticleRepository struct {
	Conn *sql.DB
}

// NewPostgresArticleRepository will create an object that represent the article.Repository interface
func NewPostgresArticleRepository(conn *sql.DB) domain.ArticleRepository {
	return &postgresArticleRepository{conn}
}

func (m *postgresArticleRepository) Store(ctx context.Context, a *domain.Article) (err error) {
	sqlStatement := `INSERT INTO article (title, content, type, url) VALUES ($1, $2, $3, $4) RETURNING id`
  stmt, err := m.Conn.PrepareContext(ctx, sqlStatement)
	if err != nil {
		return
	}

	var aID int64
	err = stmt.QueryRow(a.Title, a.Content, a.Type, a.URL).Scan(&aID)
	if err != nil {
		log.Fatal(err)
	}

	a.ID = aID
	return
}