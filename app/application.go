package app

import (
	"fmt"

	"github.com/uuid-disted/consumer/internal/services/broker"
)

type Application struct {
	brokers []*broker.RabbitMQBroker
	uuids   chan string
}

func NewApplication(brokersHost []string) *Application {
	app := &Application{
		brokers: make([]*broker.RabbitMQBroker, len(brokersHost)),
		uuids:   make(chan string, 1000000),
	}

	for i, host := range brokersHost {
		newBroker := broker.NewRabbitMQBroker(host)
		err := newBroker.Connect()
		if err != nil {
			fmt.Printf("Error connecting to broker %s: %v\n", host, err)
			continue
		}
		app.brokers[i] = newBroker
	}

	return app
}

func (app *Application) Run(queueName string) {
	worker := func(id int, broker *broker.RabbitMQBroker) {
		err := app.ConsumeMessage(broker, queueName)
		if err != nil {
			fmt.Printf("Error occurred on worker (%d): %v\n", id, err)
		}
	}

	for i := 0; i < len(app.brokers); i++ {
		go worker(i, app.brokers[i])
	}
}
