package domain

import (
	"context"
)

type File struct {
	Id   string `json:"id" prop:"id"`
	Name string `json:"name" prop:"name"`
}

type FileCreator interface {
	Create(ctx context.Context, file File, owner User) error
}

type FileGetter interface {
	GetById(ctx context.Context, id string) (File, error)
	GetOwner(ctx context.Context, file File) (User, error)
	GetAllForOwner(ctx context.Context, owner User) ([]File, error)
}

type FileUpdater interface {
	UpdateName(ctx context.Context, file File, name string) error
}

type FileDeleter interface {
	Delete(ctx context.Context, file File) error
	DeleteAllForOwner(ctx context.Context, owner User) error
}
