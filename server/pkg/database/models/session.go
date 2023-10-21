package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id        string    `json:"id" prop:"id"`
	Username  string    `json:"username" prop:"username"`
	CreatedAt time.Time `json:"createdAt" prop:"created_at"`
	ExpiresAt time.Time `json:"expiresAt" prop:"expires_at"`
}

func NewWeekSession(username string) Session {
	return Session{
		Id:        uuid.NewString(),
		Username:  username,
		CreatedAt: time.Now().In(time.UTC),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).In(time.UTC),
	}
}
