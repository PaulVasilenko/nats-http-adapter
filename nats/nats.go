package nats

import (
	"context"
	"errors"
	"time"

	"github.com/nats-io/go-nats"
)

const (
	TypeRequest = "req"
	TypePublish = "pub"
)

var (
	ErrUnknownMessageType = errors.New("unknown type of message")
)

type Conn struct {
	Timeout time.Duration
	Conn    *nats.Conn
}

type SendNatsMessageParams struct {
	Type    string
	Subject string
	Data    []byte
}

func (c *Conn) SendNATSMessage(ctx context.Context, params SendNatsMessageParams) ([]byte, error) {
	switch params.Type {
	case TypeRequest:
		msg, err := c.Conn.Request(params.Subject, params.Data, c.Timeout)
		if err != nil {
			return nil, err
		}
		return msg.Data, nil
	case TypePublish:
		err := c.Conn.Publish(params.Subject, params.Data)
		if err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, ErrUnknownMessageType
	}
}
