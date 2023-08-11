package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	var rabbitConn *amqp.Connection
	var counts int64
	var backOff = 1 * time.Second

	for {
		connection, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("rabbitqm not ready...")
			counts++
		} else {
			fmt.Println()
			rabbitConn = connection
			break
		}
		if counts > 5 {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("Backing off for %d seconds...\n", int(math.Pow(float64(counts), 2)))
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
		continue
	}
	defer rabbitConn.Close()
	// app := Config{
	// 	Rabbit: rabbitConn,
	// }
	srv := &http.Server{
		Addr: ":80",
		//Handler: app.routes(),
	}
	srv.ListenAndServe()
}
