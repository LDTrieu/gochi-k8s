package event

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Emitter for publishing AMQP events
type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		panic(e)
	}
	defer channel.Close()
	//return declareExchange(channel)
	return nil
}

func (e Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	//log.Println("Pushing to ", getExchangename())
	err = channel.PublishWithContext(
		context.Background(),
		getExchangeName(),
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Sending message: %s -> %s", event, getExchangeName())
	return nil
}
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
