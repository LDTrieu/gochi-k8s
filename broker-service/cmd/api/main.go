package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
	Etcd   *clientv3.Client
}

func main() {

	rabbitConn, err := connectToRabbit()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	log.Println("Starting broker service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToRabbit() (*amqp.Connection, error) {
	var rabbitConn *amqp.Connection
	var counts int64
	var rabbitURL = os.Getenv("RABBIT_URL")

	for {
		connection, err := amqp.Dial(rabbitURL)
		if err != nil {
			fmt.Println("rabbitmq not ready...")
			counts++
		} else {
			fmt.Println()
			rabbitConn = connection
			break
		}

		if counts > 15 {
			fmt.Println(err)
			return nil, errors.New("cannot connect to rabbit")
		}
		fmt.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
	fmt.Println("Connected to RabbitMQ!")
	return rabbitConn, nil
}
