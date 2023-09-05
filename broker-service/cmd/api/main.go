 package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct {
	Rabbit *amqp.Connection
}

 func main() {
	// try to connect rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Printf("Error: \n", err)
		os.Exit(1)
	}

	defer rabbitConn.Close()
	
	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

 }
 
 func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Printf("RabbitMQ is not ready yet...\n")
			counts++
		} else {
			log.Printf("Connected to RabbitMQ\n")
			connection = c
			break	
		}
		if counts > 5 {
			log.Printf("Error: \n", err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Backing off for %v seconds", backOff)
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}