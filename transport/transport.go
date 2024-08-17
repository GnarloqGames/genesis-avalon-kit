package transport

import (
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Transport interface {
	Request(subj string, data []byte, timeout time.Duration) (*nats.Msg, error)
}

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

func NewConn(config Configuration) (*Connection, error) {
	c, err := nats.Connect(config.URL, config.Options...)
	if err != nil {
		return nil, err
	}

	return &Connection{c}, nil
}

func Request[T protoreflect.ProtoMessage](bus Transport, subject string, request proto.Message, response T, timeout time.Duration) (*nats.Msg, error) {
	rawRequest, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	rawRes, err := bus.Request(subject, rawRequest, timeout)
	if err != nil {
		return nil, err
	}

	if err := proto.Unmarshal(rawRes.Data, response); err != nil {
		return rawRes, err
	}

	return rawRes, nil
}
