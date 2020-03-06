package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"
	"github.com/syncfuture/go/u"
)

type Consumer struct {
	Node   *NodeConfig
	conn   *amqp.Connection
	ctx    context.Context
	cancel context.CancelFunc
}

func NewConsumer(node *NodeConfig) (r *Consumer, err error) {
	r = &Consumer{
		Node: node,
	}

	r.ctx, r.cancel = context.WithCancel(context.Background())

	// Build connection
	r.conn, err = amqp.Dial(r.Node.URL)
	if u.LogError(err) {
		return
	}

	r.declare()
	return
}

// Declare declare exchanges, queues and bindings
func (x *Consumer) declare() error {
	// Build channel
	ch, err := x.conn.Channel()
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
						binding.Exchange,
						binding.NoWait,
						binding.Args,
					)
				}
			}
		}
	}

	return nil
}

func (x *Consumer) Consume(receiver func(amqp.Delivery)) {
	if u.IsMissing(x.Node.Consumers) {
		panic("consumers is missing in configuration")
	}

	// Declare consumers
	for _, consumerCfg := range x.Node.Consumers {
		go func(consumer *ConsumerConfig) {
			// Build channel
			ch, err := x.conn.Channel()
			if u.LogError(err) {
				return
			}
			defer ch.Close()

			msgs, err := ch.Consume(
				consumer.Queue,
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
				if x.isCanceled() {
					break
				}
				go receiver(msg) // DO NOT use pointer like &msg, https://www.jb51.net/article/138126.htm
			}
		}(consumerCfg)
	}
}

func (x *Consumer) isCanceled() bool {
	select {
	case <-x.ctx.Done():
		return true
	default:
		return false
	}
}

func (x *Consumer) Close() error {
	x.cancel()
	return x.conn.Close()
}
