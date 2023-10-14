package domain

import (
	"context"
)

const (
	RAcess   = iota
	RWAccess = iota
)

type AccessGrander interface {
	Grand(ctx context.Context, file File, user User, access int) error
}

type AccessGetter interface {
	GetAccessors(ctx context.Context, file File) (map[User]int, error)
}

type AccessChecker interface {
	Check(ctx context.Context, file File, user User) (bool, error)
}

type AccessRevoker interface {
	Revoke(ctx context.Context, file File, user User) error
}
