package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"io/ioutil"
	"log"
)

func main() {
	//Get RabbitMQ server
	content, err := ioutil.ReadFile("./mq_server.txt")
	if err != nil {
		log.Fatal(err)
	}
	server := string(content)

	conn, err := amqp.Dial(server)
	handleError(err, "Failed to connect")
	defer conn.Close()

	chn, err := conn.Channel()
	handleError(err, "Failed to create a channel")
	defer chn.Close()

	q, err := chn.QueueDeclare(
		"new_queue", //queue name
		true,        //durable
		false,       //auto del
		false,       //exclusive
		false,       //no wait
		nil,         //arguments
	)
	handleError(err, "Failed to declare a queue")

	body := "This is a message"
	err = chn.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	handleError(err, "Failed to publish a message")
	log.Printf("Sent message: %s\n", body)
}

func handleError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
