package rabbitmq

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"gitlab.com/mikrowezel/backend/broker/mapper"
	"gitlab.com/mikrowezel/backend/log"
)

// NewRabbitMQ creates and return a new RabbitMQ broker.
func NewRabbitMQ(ctx context.Context, cfg *Config, log *log.Logger) (*RabbitMQ, error) {
	r, err := newRabbitMQ(ctx, cfg, log)
	return r, err
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
	r.Channels = append(r.Channels, <-r.RetryChannel(10))

	return r, nil
}

// RabbitMQURL returns a RabbitMQ connection URL.
func (c *Config) RabbitMQURL() string {
	panic("TODO: Not implemented yet")
}

// SetConfig for RabbitMQ client.
func (r *RabbitMQ) SetConfig(cfg *Config) {
	r.cfg = cfg
}

// Connect to RabbitMQ.
func (r *RabbitMQ) Connect(retry bool) error {
	if r.cfg == nil {
		return errors.New("setup a configuration before connect")
	}
	if retry {
		r.conn = <-r.RetryConnection(r.cfg)
	}

	var err error
	r.conn, err = r.Connection()
	return err
}

func (r *RabbitMQ) Channel() (*amqp.Channel, error) {
	for _, ch := range r.Channels {
		if ch.IsOpen {
			return ch.Channel, nil
		}
	}
	return nil, errors.New("cannot get a channel")
}

// IsConnected returns true if broker
// connection is open.
func (r *RabbitMQ) IsConnected() bool {
	return !r.conn.IsClosed()
}

// AddExchange to the broker handler.
func (r *RabbitMQ) AddExchange(name, kind string, durable, autodelete, internal, nowait bool) error {
	if r.IsConnected() {
		return errors.New("no active connection")
	}

	ch, err := r.Channel()
	if err != nil {
		return err
	}

	ch.ExchangeDeclare(name, kind, durable, autodelete, internal, nowait, nil)

	r.Exchanges[name] = &Exchange{
		ID:         uuid.New(),
		Name:       name,
		Kind:       kind,
		Durable:    durable,
		AutoDelete: autodelete,
		Internal:   internal,
		NoWait:     nowait,
	}

	return nil
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

// AddEmitter to the broker.
// queue parameter is optional but if it is provided
// a queue and binding to the exchange will be created
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
