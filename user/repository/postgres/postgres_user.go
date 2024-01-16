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

func (m *postgresUserRepository) Login(ctx context.Context, email string, password string, u *domain.User) (err error) {
	var rid int64
	var rname string
	var remail string

	sqlStatement := `SELECT id, email, name FROM users WHERE email=$1 AND password=$2`

	err = m.Conn.QueryRowContext(ctx, sqlStatement, email, password).Scan(&rid, &remail, &rname)
	if err != nil || rid == 0 {
		return errors.New("user not found")
	}

	u.ID = rid
	u.Email = rname
	u.Name = remail
	return
}
