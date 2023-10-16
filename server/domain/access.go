package domain

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	RAcess   = iota
	RWAccess = iota
)

type Access struct {
	Granter  string `json:"granter"`
	Receiver string `json:"receiver"`
	Level    int
}

func (a Access) MarshalJSON() ([]byte, error) {
	var level string
	if a.Level == RWAccess {
		level = "RW"
	} else if a.Level == RAcess {
		level = "R"
	} else {
		return nil, fmt.Errorf("unknown access level value: %d", a.Level)
	}

	access := struct {
		Granter  string `json:"granter"`
		Receiver string `json:"receiver"`
		Level    string `json:"level"`
	}{a.Granter, a.Receiver, level}

	return json.Marshal(access)
}

func (a *Access) UnmarshalJSON(b []byte) error {
	var access struct {
		Granter  string `json:"granter"`
		Receiver string `json:"receiver"`
		Level    string `json:"level"`
	}
	if err := json.Unmarshal(b, &access); err != nil {
		return err
	}

	var level int
	if access.Level == "RW" {
		level = RWAccess
	} else if access.Level == "R" {
		level = RAcess
	} else {
		return fmt.Errorf("unknown access level value: %s", access.Level)
	}

	a.Granter = access.Granter
	a.Receiver = access.Receiver
	a.Level = level

	return nil
}

type AccessGranter interface {
	Grant(ctx context.Context, file File, access Access) error
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
