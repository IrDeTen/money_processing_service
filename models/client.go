package models

import "github.com/google/uuid"

type Client struct {
	id   uuid.UUID
	Name string
}

func CreateNewClient(name string) Client {
	return Client{
		id:   uuid.New(),
		Name: name,
	}
}

func (c *Client) GetID() uuid.UUID {
	return c.id
}
