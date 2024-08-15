package repository

import (
	"github.com/GnarloqGames/genesis-avalon-kit/database/cockroach/driver"
)

type Registry struct {
	conn *driver.Connection
}

func NewRegistry(c *driver.Connection) *Registry {
	return &Registry{
		conn: c,
	}
}
