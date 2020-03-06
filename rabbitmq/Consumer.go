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
	err = declare(r.conn, r.Node)
	return
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
