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
	GetAccesses(ctx context.Context, file File) ([]Access, error)
}

type AccessChecker interface {
	Check(ctx context.Context, file File, user User) (bool, error)
}

type AccessRevoker interface {
	Revoke(ctx context.Context, file File, user User) error
}
