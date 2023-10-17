package domain

import (
	"context"
)

const (
	RAcess   = "R"
	RWAccess = "RW"
)

type Access struct {
	Granter  string `json:"granter" prop:"granter"`
	Receiver string `json:"receiver" prop:"receiver"`
	Level    string `json:"level" prop:"level"`
}

type AccessGranter interface {
	Grant(ctx context.Context, file File, access Access) error
}

type AccessGetter interface {
	Get(ctx context.Context, file File, user User) (Access, error)
	GetAccesses(ctx context.Context, file File) ([]Access, error)
}

type AccessUpdater interface {
	UpdateLevel(ctx context.Context, file File, access Access, newLevel string) error
}

type AccessRevoker interface {
	Revoke(ctx context.Context, file File, user User) error
}
