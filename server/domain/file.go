package domain

import (
	"context"

	"github.com/google/uuid"
)

type File struct {
	ID   uuid.UUID `json:"id" prop:"id"`
	Name string    `json:"name" prop:"name"`
}

type FileCreator interface {
	Create(ctx context.Context, file File) error
}

type FileGetter interface {
	GetByID(ctx context.Context, id uuid.UUID) (File, error)
}

type FileUpdater interface {
	UpdateName(ctx context.Context, id uuid.UUID, name string) error
}

type FileDeleter interface {
	Delete(ctx context.Context, id uuid.UUID) error
}
