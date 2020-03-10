package rabbitmq

import (
	"context"

	"github.com/syncfuture/go/slog"

	"github.com/streadway/amqp"
	"github.com/syncfuture/go/u"
)

type Consumer struct {
	Node     *NodeConfig
	Handlers map[string]func(amqp.Delivery)
	conn     *amqp.Connection
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewConsumer(node *NodeConfig) (r *Consumer, err error) {
	r = &Consumer{
		Node: node,
	}

	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.Handlers = make(map[string]func(amqp.Delivery), 0)

	// Build connection
	r.conn, err = amqp.Dial(r.Node.URL)
	if u.LogError(err) {
		return
	}
	err = declare(r.conn, r.Node)
	return
}

func (x *Consumer) RegisterHandler(name string, handler func(amqp.Delivery)) {
	x.Handlers[name] = handler
}

func (x *Consumer) Consume( /*handlers map[string]func(amqp.Delivery)*/ ) {
	if u.IsMissing(x.Node.Consumers) {
		panic("consumers is missing in configuration")
	}

	// Declare consumers
	for _, consumerCfg := range x.Node.Consumers {
		go func(consumer *ConsumerConfig) {
			if consumer.Name == "" {
				slog.Error("consumer name cannot be empty.")
				return
			}
			if consumer.Queue == "" {
				slog.Errorf("queue for consumer %s cannot be empty.", consumer.Name)
				return
			}
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
				go x.handle(msg) // DO NOT use pointer in for loop, like &msg, https://www.jb51.net/article/138126.htm
			}
		}(consumerCfg)
	}
}

func (x *Consumer) handle(msg amqp.Delivery) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error(err)
		}
	}()

	if handler, ok := x.Handlers[msg.Type]; ok {
		handler(msg)
	} else {
		slog.Warnf("cannnot find handler for message type '%s'", msg.Type)
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
