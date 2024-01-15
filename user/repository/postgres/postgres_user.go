package postgres

import (
	"context"
	"database/sql"

	"github.com/Kelvinmijaya/kelvin-rest-api/domain"
)

type postgresUserRepository struct {
	Conn *sql.DB
}

func NewPostgresUserRepository(conn *sql.DB) domain.UserRepository {
	return &postgresUserRepository{conn}
}

func (m *postgresUserRepository) Login(ctx context.Context, email string, password string) (err error) {
	query := `SELECT id, email, name FROM users WHERE email=$1 AND password=$2`

	rows, err := m.Conn.Query(query, email, password)

	if err != nil {
		return err
	}

	defer rows.Close()

	return
}
