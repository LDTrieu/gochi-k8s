package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/ldtrieu/go-rabbit/lib/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	log.Println("Listening for and consuming RabbitMQ messages...")

	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

// connect tries to connect to RabbitMQ, and delays between attempts.
func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection
	var rabbitURL = os.Getenv("RABBIT_URL")

	for {
		c, err := amqp.Dial(rabbitURL)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {

			connection = c
			fmt.Println()
			break
		}

		if counts > 5 {
			// if we can't connect after five tries, something is wrong...
			fmt.Println(err)
			return nil, err
		}
		fmt.Printf("Backing off for %d seconds...\n", int(math.Pow(float64(counts), 2)))
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
