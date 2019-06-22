package rabbitmq

import (
	"context"
	"errors"

	"github.com/streadway/amqp"
	"gitlab.com/mikrowezel/backend/broker"
	"gitlab.com/mikrowezel/backend/broker/mapper"
	"gitlab.com/mikrowezel/backend/log"
)

// Config for broker.
type Config struct {
	Name            string
	BackoffMaxTries int
}

// RabbitMQ is message broker handler.
type RabbitMQ struct {
	ctx       context.Context
	cfg       *Config
	log       *log.Logger
	name      string
	ready     bool
	alive     bool
	conn      *amqp.Connection
	Listeners map[string]*Listener
	Emitters  map[string]*Emitter
}

// Emitter is a RabbitMQ message emitter.
type Emitter struct {
	connection *amqp.Connection
	exchange   string
	events     chan *EmittedBaseMessage
	log        *log.Logger
}

// Listener is a RabbitMQ message listener.
type Listener struct {
	connection *amqp.Connection
	exchange   string
	queue      string
	mapper     mapper.BaseMessageMapper
	log        *log.Logger
}

// EmittedBaseMessage is an emitted base message.
type EmittedBaseMessage struct {
	event     broker.BaseMessage
	errorChan chan error
}
