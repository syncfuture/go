package rabbitmq

import (
	"github.com/streadway/amqp"
	"github.com/syncfuture/go/u"
)

func declare(conn *amqp.Connection, node *NodeConfig) error {
	// Build channel
	ch, err := conn.Channel()
	if u.LogError(err) {
		return err
	}
	defer ch.Close()

	// Declare exchanges
	if !u.IsMissing(node.Exchanges) {
		for _, exchange := range node.Exchanges {
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
	if !u.IsMissing(node.Queues) {
		for _, queue := range node.Queues {
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

			// Declare bindings
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
