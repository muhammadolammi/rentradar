package notification

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

func publishMessage(body []byte, rabbitmq_url string) {
	conn, err := amqp.Dial(rabbitmq_url)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notifications", // queue name
		true,            // durable
		false,           // auto-delete
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key (queue name)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		log.Fatal(err)
	}
}

func worker(id int, queueName string, config *Config) {
	//    to consume message on the queue

	conn, err := amqp.Dial(config.RABBITMQUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue name
		"",        // consumer tag
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		log.Printf("Worker %d processing: %s", id, msg.Body)
		// Here you can call a function to send the actual notification
		// Unmarshal the body
		notification := Notification{}
		err = json.Unmarshal(msg.Body, &notification)
		if err != nil {
			log.Fatalf("error unmarshalling message body. err: %v", err)
		}

		SendNotification(notification, config.SMTPModel)
	}

}

func (config *Config) startWorkerPool(numWorkers int, queueName string) {
	for i := range numWorkers {
		go worker(i, queueName, config)
	}
}
