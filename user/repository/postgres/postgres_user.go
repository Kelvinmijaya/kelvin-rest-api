package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Kelvinmijaya/kelvin-rest-api/domain"
)

type postgresUserRepository struct {
	Conn *sql.DB
}

func NewPostgresUserRepository(conn *sql.DB) domain.UserRepository {
	return &postgresUserRepository{conn}
}

func (m *postgresUserRepository) Login(ctx context.Context, email string, password string) (err error) {
	var rname string
	var rid int64
	var remail string

	sqlStatement := `SELECT id, email, name FROM users WHERE email=$1 AND password=$2`

	err = m.Conn.QueryRow(sqlStatement, email, password).Scan(&rid, &remail, &rname)
	if err != nil || rid == 0 {
		return errors.New("user not found")
	}

	return
}
