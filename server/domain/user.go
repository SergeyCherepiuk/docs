package domain

import "context"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserCreator interface {
	Create(ctx context.Context, user User) error
}

type UserGetter interface {
	Get(ctx context.Context, username string) (User, error)
}

type UserUpdater interface {
	UpdateUsername(ctx context.Context, username, newUsername string) error
	UpdatePassword(ctx context.Context, username, newPassword string) error
}

type UserDeleter interface {
	Delete(ctx context.Context, username string) error
}
