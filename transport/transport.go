package transport

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
)

type Encoder string

const (
	EncoderJSON     Encoder = "json"
	EncoderProtobuf Encoder = "protobuf"
)

type Configuration struct {
	URL     string
	Encoder Encoder
	Options []nats.Option
}

var DefaultConfig = Configuration{
	URL:     "127.0.0.1:4222",
	Encoder: EncoderJSON,
	Options: make([]nats.Option, 0),
}

type Connection struct {
	*nats.EncodedConn
}

func NewEncodedConn(config Configuration) (*Connection, error) {
	c, err := nats.Connect(config.URL, config.Options...)
	if err != nil {
		return nil, err
	}

	var encoder string
	switch config.Encoder {
	case EncoderJSON:
		encoder = nats.JSON_ENCODER
	case EncoderProtobuf:
		encoder = protobuf.PROTOBUF_ENCODER
	default:
		return nil, fmt.Errorf("invalid encoder: %s", config.Encoder)
	}

	nc, err := nats.NewEncodedConn(c, encoder)
	if err != nil {
		return nil, err
	}

	return &Connection{nc}, nil
}

func ParseEncoder(str string) Encoder {
	encoder := EncoderJSON
	switch str {
	case "json":
	case "protobuf":
		encoder = EncoderProtobuf
	}

	return encoder
}
