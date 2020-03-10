package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/syncfuture/go/u"
)

type Producer struct {
	Node *NodeConfig
	conn *amqp.Connection
}

func NewProducer(node *NodeConfig) (r *Producer, err error) {
	r = &Producer{
		Node: node,
	}

	// Build connection
	r.conn, err = amqp.Dial(r.Node.URL)
	if u.LogError(err) {
		return
	}
	err = declare(r.conn, r.Node)
	return
}

func (x *Producer) Publish(exchange, routingKey, msgType string, payload []byte, headers amqp.Table) error {
	return x.PublishRaw(exchange, routingKey, false, false, &amqp.Publishing{
		Headers: headers,
		Type:    msgType,
		Body:    payload,
	})
}

func (x *Producer) PublishRaw(exchange, routingKey string, mandatory, immediate bool, msg *amqp.Publishing) error {
	// Build channel
	ch, err := x.conn.Channel()
	if u.LogError(err) {
		return err
	}
	defer ch.Close()

	return ch.Publish(exchange, routingKey, false, false, *msg)
}

func (x *Producer) Close() error {
	return x.conn.Close()
}
