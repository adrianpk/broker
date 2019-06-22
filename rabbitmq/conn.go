package rabbitmq

import (
	"time"

	"github.com/cenkalti/backoff"
	"github.com/streadway/amqp"
)

// Connection tries to establish a connection to RabbitMQ.
// error if fails.
func (r *RabbitMQ) Connection() (*amqp.Connection, error) {
	url := r.cfg.RabbitMQURL()
	r.log.Info("Dialing to RabbitMQ broker", "host", url)
	return amqp.Dial(url)
}

// RetryConnection implements a backoff mechanism for establishing a connection
// to RabbitMQ; this is especially useful in containerize environments where
// components can be started out of order.
func (r *RabbitMQ) RetryConnection(cfg *Config) chan *amqp.Connection {
	result := make(chan *amqp.Connection)

	bo := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), uint64(cfg.BackoffMaxTries))

	go func() {
		defer close(result)

		url := cfg.RabbitMQURL()

		for i := 0; i <= cfg.BackoffMaxTries; i++ {

			r.log.Info("Dialing to RabbitMQ broker", "host", url)

			conn, err := amqp.Dial(url)
			if err == nil {
				r.log.Info("RabbitMQ connection established")
				result <- conn
				return
			}

			r.log.Error(err, "RabbitMQ connection error")

			// Backoff
			nb := bo.NextBackOff()
			if nb == backoff.Stop {
				result <- nil
				r.log.Info("Rabbit connection failed", "reason", "max number of tries reached")
				bo.Reset()
				return
			}

			r.log.Info("Rabbit connection failed", "retrying-in", nb.String(), "unit", "seconds")
			time.Sleep(nb)
		}
	}()

	return result
}