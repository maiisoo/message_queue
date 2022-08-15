package main

import (
	"io/ioutil"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//Get RabbitMQ server from file
	content, err := ioutil.ReadFile("./mq_server.txt")
	if err != nil {
		log.Fatal(err)
	}
	server := string(content)

	conn, err := amqp.Dial(server)
	handleError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	chn, err := conn.Channel()
	handleError(err, "Failed to open a channel")
	defer chn.Close()

	q, err := chn.QueueDeclare(
		"new_queue", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	handleError(err, "Failed to declare a queue")

	msgs, err := chn.Consume(
		q.Name,
		"",    // consumer
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // arguments
	)
	handleError(err, "Failed to consume messages")

	go func() {
		for m := range msgs {
			log.Printf("Received a message: %s", m.Body)
		}
	}()

}

func handleError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
