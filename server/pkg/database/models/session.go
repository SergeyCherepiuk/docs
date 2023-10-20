package models

import "time"

// TODO: Check if there is a better way to represent a session
//  1. CreatedAt and ExpiresAt is int64 in Neo4j
//  2. Querying substructure (User) might be challenging
type Session struct {
	Id        string    `json:"id" prop:"id"`
	User      User      `json:"user" prop:"user"`
	CreatedAt time.Time `json:"createdAt" prop:"created_at"`
	ExpiresAt time.Time `json:"expiresAt" prop:"expires_at"`
}
