package models

import "github.com/google/uuid"

type Client struct {
	ID   uuid.UUID
	Name string
}
