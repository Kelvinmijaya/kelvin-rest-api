package domain

import "context"

// User representing the User data struct
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Passsword string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	LastActivity string `json:"last_activity"`
}

type UserUsecase interface {
	Fetch(ctx context.Context, id int64) (User, error)
	Login(ctx context.Context, email string, Password string) (User, error)
	Logout(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, ar *Article) error
	Delete(ctx context.Context, id int64) error
}

// AuthorRepository represent the author's repository contract
type UserRepository interface {
	GetByEmail(ctx context.Context, Email string) (User, error)
}
