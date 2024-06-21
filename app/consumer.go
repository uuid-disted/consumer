package app

import (
	"github.com/uuid-disted/consumer/internal/services/broker"
)

func (app *Application) ConsumeMessage(broker *broker.RabbitMQBroker, queueName string) error {
	return broker.Consume(queueName, app.uuids)
}

func (app *Application) GetUUIDs(n int) []string {
	uuids := make([]string, n)
	for i := 0; i < n; i++ {
		uuids[i] = <-app.uuids
	}
	return uuids
}
