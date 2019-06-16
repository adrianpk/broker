package rabbitmq

import (
	"context"

	"gitlab.com/mikrowezel/backend/log"
)

// NewRabbitMQ creates and return a new RabbitMQ broker.
func NewRabbitMQ(ctx context.Context, cfg *Config, log *log.Logger) (*RabbitMQ, error) {
	r, err := newRabbitMQ(ctx, cfg, log)
	return r, err
}
