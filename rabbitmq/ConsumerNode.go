package rabbitmq

import (
	"context"
	"errors"

	"github.com/streadway/amqp"
	"github.com/syncfuture/go/u"
)

type ConsumerNode struct {
	Node *NodeConfig
}

// Declare declare exchanges, queues and bindings
func (x *ConsumerNode) Declare() error {
	// Build connection
	conn, err := amqp.Dial(x.Node.URL)
	if u.LogError(err) {
		return err
	}
	defer conn.Close()

	// Build channel
	ch, err := conn.Channel()
	if u.LogError(err) {
		return err
	}
	defer ch.Close()

	// Declare exchanges
	if !u.IsMissing(x.Node.Exchanges) {
		for _, exchange := range x.Node.Exchanges {
			err = ch.ExchangeDeclare(
				exchange.Name,
				exchange.Type,
				exchange.Durable,
				exchange.AutoDelete,
				exchange.Internal,
				exchange.NoWait,
				exchange.Args,
			)
			if u.LogError(err) {
				return err
			}
		}
	}

	// Declare queues
	if !u.IsMissing(x.Node.Queues) {
		for _, queue := range x.Node.Queues {
			_, err = ch.QueueDeclare(
				queue.Name,
				queue.Durable,
				queue.AutoDelete,
				queue.Exclusive,
				queue.NoWait,
				queue.Args,
			)
			if u.LogError(err) {
				return err
			}

			// QueueConfig binding
			if !u.IsMissing(queue.Bindings) {
				for _, binding := range queue.Bindings {
					ch.QueueBind(
						queue.Name,
						binding.RoutingKey,
						binding.ExchangeConfig,
						binding.NoWait,
						binding.Args,
					)
				}
			}
		}
	}

	return nil
}

func (x *ConsumerNode) Consume(ctx context.Context, receiver func(amqp.Delivery)) (err error) {
	if u.IsMissing(x.Node.Consumers) {
		err = errors.New("consumers is missing in configuration")
		u.LogError(err)
		return err
	}

	// Build connection
	conn, err := amqp.Dial(x.Node.URL)
	if u.LogError(err) {
		return err
	}
	defer conn.Close()

	// Declare consumers
	for _, consumerCfg := range x.Node.Consumers {
		go func(consumer *ConsumerConfig) {
			// Build channel
			ch, err := conn.Channel()
			if u.LogError(err) {
				return
			}
			defer ch.Close()

			msgs, err := ch.Consume(
				consumer.QueueConfig,
				consumer.Name,
				consumer.AutoAck,
				consumer.Exclusive,
				consumer.NoLocal,
				consumer.NoWait,
				consumer.Args,
			)
			if u.LogError(err) {
				return
			}

			for msg := range msgs {
				if isCanceled(ctx) {
					break
				}
				go receiver(msg) // DO NOT use pointer like &msg, https://www.jb51.net/article/138126.htm
			}
		}(consumerCfg)

	}

	return nil
}

func isCanceled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
