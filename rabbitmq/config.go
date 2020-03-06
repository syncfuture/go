package rabbitmq

import "github.com/streadway/amqp"

type (
	RabbitMQConfig struct {
		Nodes []*NodeConfig
	}
	NodeConfig struct {
		URL       string
		Exchanges []*ExchangeConfig
		Queues    []*QueueConfig
		Consumers []*ConsumerConfig
	}
	ExchangeConfig struct {
		Name       string
		Type       string
		Durable    bool
		AutoDelete bool
		Internal   bool
		NoWait     bool
		Args       amqp.Table
	}
	QueueConfig struct {
		Name       string
		Durable    bool
		AutoDelete bool
		Exclusive  bool
		NoWait     bool
		Args       amqp.Table
		Bindings   []*QueueBindingConfig
	}
	QueueBindingConfig struct {
		RoutingKey string
		Exchange   string
		NoWait     bool
		Args       amqp.Table
	}
	ConsumerConfig struct {
		Queue     string
		Name      string
		AutoAck   bool
		Exclusive bool
		NoLocal   bool
		NoWait    bool
		Args      amqp.Table
	}
)
