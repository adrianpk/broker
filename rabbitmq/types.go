package rabbitmq

import (
	"context"
	"sync"

	"github.com/google/uuid"
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
	Channels  []*Channel
	Exchanges map[string]*Exchange
	Queues    map[string]*Queue
	Bindings  map[string]*Binding
	Listeners map[string]*Listener
	Emitters  map[string]*Emitter
}

// Channel lets the broker client
// handle RabbitMQ channels.
type Channel struct {
	mutex   *sync.Mutex
	ID      uuid.UUID
	Name    string
	Channel *amqp.Channel
	IsOpen  bool
}

// Exchange lets the broker client
// handle RabbitMQ exchanges.
type Exchange struct {
	ID         uuid.UUID
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	ArgsTable  map[string]interface{}
	log        *log.Logger
}

// Queue lets the broker client
// handle RabbitMQ queues.
type Queue struct {
	ID        string
	Name      string
	Messages  int
	Consumers int
	log       *log.Logger
}

// Binding lets the broker client
// handle bindings between exchanges and queues.
type Binding struct {
	ID       string
	Name     string
	Exchange *Exchange
	Queue    *Queue
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
