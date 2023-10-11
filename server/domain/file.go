package domain

import "github.com/google/uuid"

type File struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
}

// TODO: Write services for "File" type
