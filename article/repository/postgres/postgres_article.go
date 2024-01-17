package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/Kelvinmijaya/kelvin-rest-api/article/repository"
	"github.com/Kelvinmijaya/kelvin-rest-api/domain"
	"github.com/sirupsen/logrus"
)

type postgresArticleRepository struct {
	Conn *sql.DB
}

// NewPostgresArticleRepository will create an object that represent the article.Repository interface
func NewPostgresArticleRepository(conn *sql.DB) domain.ArticleRepository {
	return &postgresArticleRepository{conn}
}

func (m *postgresArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Article, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Article, 0)
	for rows.Next() {
		t := domain.Article{}
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.URL,
			&t.Content,
			&t.Type,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *postgresArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Article, nextCursor string, err error) {
	query := `SELECT id, title, url, content, type, updated_at, created_at
  						FROM article WHERE created_at > $1 ORDER BY created_at LIMIT $2 `

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}
func (m *postgresArticleRepository) GetByID(ctx context.Context, id int64) (res domain.Article, err error) {
	query := `SELECT id, title ,url ,content, type, updated_at, created_at
  						FROM article WHERE id=$1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Article{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *postgresArticleRepository) Store(ctx context.Context, a *domain.Article) (err error) {
	sqlStatement := `INSERT INTO article (title, content, type, url, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, sqlStatement)
	if err != nil {
		return
	}

	var aID int64
	err = stmt.QueryRow(a.Title, a.Content, a.Type, a.URL, a.UpdatedAt, a.CreatedAt).Scan(&aID)
	if err != nil {
		log.Fatal(err)
	}

	a.ID = aID
	return
}

func (m *postgresArticleRepository) Update(ctx context.Context, ar *domain.Article) (err error) {
	query := `UPDATE article SET title=$1, url=$2, content=$3, type=$4, created_at=$5 WHERE id=$6`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Title, ar.URL, ar.Content, ar.Type, ar.UpdatedAt, ar.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (m *postgresArticleRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM article WHERE id=$1"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}
