package domain

import "context"

// User representing the User data struct
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" query:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type UserUsecase interface {
	// Fetch(ctx context.Context, id int64) (User, error)
	Login(ctx context.Context, u *User) error
	// Delete(ctx context.Context, id int64) error
}

// AuthorRepository represent the author's repository contract
type UserRepository interface {
	Login(ctx context.Context, u *User) error
}
