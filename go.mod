module gitlab.com/mikrowezel/backend/broker

go 1.12

require (
	github.com/cenkalti/backoff v2.1.1+incompatible
	github.com/mitchellh/mapstructure v1.1.2
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	gitlab.com/mikrowezel/backend/log v0.0.0
)

replace gitlab.com/mikrowezel/backend/log => ../log
