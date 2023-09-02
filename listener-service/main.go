package main

import (
	"time"
	"fmt"
	"log"
	"os"
	"math"
	amqp "github.com/rabbitmq/amqp091-go"
	
	"listener/event"
)

func main() {
	// try to connect rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Printf("Error: \n", err)
		os.Exit(1)
	}

	defer rabbitConn.Close()

	// start listening for message
	log.Printf("Listening for message...\n")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume event
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Printf("Error: \n", err)
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