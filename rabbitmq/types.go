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

type ConnStatus string

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

// newRabbitMQ create a new RabbitMQ broker handler.
func newRabbitMQ(ctx context.Context, cfg *Config, log *log.Logger) (*RabbitMQ, error) {
	r := &RabbitMQ{
		// Generic
		ctx:       ctx,
		cfg:       cfg,
		name:      cfg.Name,
		ready:     false,
		alive:     false,
		log:       log,
		Listeners: make(map[string]*Listener),
		Emitters:  make(map[string]*Emitter),
	}

	r.conn = <-r.RetryConnection(cfg)

	return r, nil
}

// RabbitMQURL returns a RabbitMQ connection URL.
func (c *Config) RabbitMQURL() string {
	panic("TODO: Not implemented yet")
}

// ConStatus returns true if broker
// connection is open.
func (r *RabbitMQ) ConnStatus() bool {
	return !r.conn.IsClosed
}

// ConStratus returns true if broker
// connection is open.
func (r *RabbitMQ) ConnStatus() bool {
	return !r.conn.IsClosed
}

// AddListener to the broker
func (r *RabbitMQ) AddListener(name, exchange, queue string) error {
	l, err := r.NewListener(exchange, queue)
	if err != nil {
		return err
	}
	r.Listeners[name] = l
	return nil
}

// AddEmitter to the brokergo i
// queue parameter is optional but if it is provided
// A queue and binding to the exchange will be created
// for each provided name.
func (r *RabbitMQ) AddEmitter(name, exchange string, queue ...string) error {
	if len(queue) < 1 {
		queue = []string{""} //TODO: Implement (multiple) queue binding
	}

	e, err := r.NewEmitter(exchange, queue[0])
	if err != nil {
		return err
	}

	r.Emitters[name] = e
	return nil
}

// NewListener returns a new RabbitMQ broker listener.
func (r *RabbitMQ) NewListener(exchange, queue string) (*Listener, error) {
	if r.conn == nil {
		return nil, errors.New("broker has no connection")
	}

	return &Listener{
		connection: r.conn,
		exchange:   exchange,
		queue:      queue,
		mapper:     mapper.NewMessageMapper(),
		log:        r.log,
	}, nil
}

// NewEmitter returns a new RabbitMQ broker emitter.
func (r *RabbitMQ) NewEmitter(exchange, queue string) (*Emitter, error) {
	if r.conn == nil {
		return nil, errors.New("broker has no connection")
	}

	return &Emitter{
		connection: r.conn,
		exchange:   exchange,
		log:        r.log,
	}, nil
}
