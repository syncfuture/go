package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/syncfuture/go/u"
)

type Producer struct {
	Node *NodeConfig
}

func NewProducer(node *NodeConfig) *Producer {
	return &Producer{
		Node: node,
	}
}

func (x *Producer) Publish(exchange, routingKey string, payload []byte, headers amqp.Table) error {
	return x.PublishRaw(exchange, routingKey, false, false, &amqp.Publishing{
		Headers: headers,
		Body:    payload,
	})
}

func (x *Producer) PublishRaw(exchange, routingKey string, mandatory, immediate bool, msg *amqp.Publishing) error {
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

	return ch.Publish(exchange, routingKey, false, false, *msg)
}
