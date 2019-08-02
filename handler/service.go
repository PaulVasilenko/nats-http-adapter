package handler

import (
	"context"
	"encoding/json"

	"github.com/paulvasilenko/nats-http-adapter/nats"
	"github.com/pkg/errors"
)

type Service struct {
	NATS nats.Conn
}

type NATSMessage struct {
	Subject string          `json:"subject"`
	Type    string          `json:"type"`
	Data    json.RawMessage `json:"data"`
}

type NATSResult struct {
	Data interface{} `json:"data"`
}

func (s *Service) SendNATSMessage(ctx context.Context, m *NATSMessage) (*NATSResult, error) {
	if m.Subject == "" {
		return nil, BadRequest("subject is missing")
	}

	resp, err := s.NATS.SendNATSMessage(ctx, nats.SendNatsMessageParams{
		Type:    m.Type,
		Subject: m.Subject,
		Data:    m.Data,
	})
	if err != nil {
		if errors.Cause(err) == nats.ErrUnknownMessageType {
			return nil, BadRequest(err.Error())
		}

		return nil, err
	}

	if len(resp) < 0 {
		return nil, nil
	}

	if json.Valid(resp) {
		return &NATSResult{
			Data: json.RawMessage(resp),
		}, nil
	}

	return &NATSResult{
		Data: string(resp),
	}, nil
}
