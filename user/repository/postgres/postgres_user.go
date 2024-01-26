package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Kelvinmijaya/kelvin-rest-api/domain"
	"golang.org/x/crypto/bcrypt"
)

type postgresUserRepository struct {
	Conn *sql.DB
}

func NewPostgresUserRepository(conn *sql.DB) domain.UserRepository {
	return &postgresUserRepository{conn}
}

func (m *postgresUserRepository) Login(ctx context.Context, u *domain.User) (err error) {
	var rid int64
	var rname string
	var remail string
	var rPassword string

	sqlStatement := `SELECT id, email, name, password FROM users WHERE email=$1`

	err = m.Conn.QueryRowContext(ctx, sqlStatement, u.Email).Scan(&rid, &remail, &rname, &rPassword)
	if err != nil || rid == 0 {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(rPassword), []byte(u.Password))
	if err != nil {
		return errors.New("wrong password")
	}

	u.ID = rid
	u.Name = remail
	u.Password = ""
	return
}
