package models

import (
	"errors"

	"github.com/google/uuid"
)

var (
	errInvalidClientUUID = errors.New("invalid id for client")
)

type Client struct {
	id   uuid.UUID
	Name string
}

func CreateNewClient(name string) Client {
	return Client{
		Name: name,
	}
}

func (c *Client) GetID() uuid.UUID {
	return c.id
}

func (c *Client) SetID(id uuid.UUID) error {
	if id == uuid.Nil {
		return errInvalidClientUUID
	}
	c.id = id
	return nil
}
