package domain

import (
	"context"
)

type File struct {
	Id   string `json:"id" prop:"id"`
	Name string `json:"name" prop:"name"`
}

type FileCreator interface {
	Create(ctx context.Context, file File) error
}

type FileGetter interface {
	GetById(ctx context.Context, id string) (File, error)
}

type FileUpdater interface {
	UpdateName(ctx context.Context, id string, name string) error
}

type FileDeleter interface {
	Delete(ctx context.Context, id string) error
}
