package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/syncfuture/go/u"
)

type ProducerNode struct {
	Node *NodeConfig
}

func (x *ProducerNode) Publish(exchange, routingKey string) error {
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

	return ch.Publish("", "", false, false, amqp.Publishing{})
}
