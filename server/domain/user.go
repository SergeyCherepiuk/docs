package domain

import "context"

type User struct {
	Username string `json:"username" prop:"username"`
	Password string `json:"password" prop:"password"`
}

type UserCreator interface {
	Create(ctx context.Context, user User) error
}

type UserGetter interface {
	GetByUsername(ctx context.Context, username string) (User, error)
}

type UserUpdater interface {
	UpdateUsername(ctx context.Context, user User, newUsername string) error
	UpdatePassword(ctx context.Context, user User, newPassword string) error
}

type UserDeleter interface {
	Delete(ctx context.Context, user User) error
}
