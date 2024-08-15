package transport

import (
	"github.com/nats-io/nats.go"
)

type Configuration struct {
	URL     string
	Options []nats.Option
}

var DefaultConfig = Configuration{
	URL:     "127.0.0.1:4222",
	Options: make([]nats.Option, 0),
}

type Connection struct {
	*nats.Conn
}

func NewEncodedConn(config Configuration) (*Connection, error) {
	c, err := nats.Connect(config.URL, config.Options...)
	if err != nil {
		return nil, err
	}

	return &Connection{c}, nil
}
