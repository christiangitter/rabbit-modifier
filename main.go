package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// simple error handling helper function
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// simple JSON validation helper function
func isValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func main() {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Inspect the queue to get the current message count
	queue_name := "test_queue"
	queue, err := ch.QueueInspect(queue_name)
	failOnError(err, "Failed to inspect the queue")

	messageCount := queue.Messages
	log.Printf("Queue %s has %d messages.\n", queue_name, messageCount)

	// register a consumer
	msgs, err := ch.Consume(
		queue_name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	receivedCount := 0

	for msg := range msgs {
		if receivedCount >= messageCount {
			log.Println("Processed the expected number of messages, stopping.")
			break
		}
		log.Printf("Received a message: %s", msg.Body)

		// Unmarshal the JSON into a map
		var messageMap map[string]interface{}
		err := json.Unmarshal(msg.Body, &messageMap)
		failOnError(err, "Failed to unmarshal JSON")

		// Navigate to the "fullyQualifiedClassName" key
		if payLoad, ok := messageMap["payLoad"].(map[string]interface{}); ok {
			payLoad["fullyQualifiedClassName"] = "de.test.test.TS.endpoint.TESTEndpoint.TEST.TestInformation"
		}

		// Marshal the map back into JSON
		modifiedMessage, err := json.Marshal(messageMap)
		failOnError(err, "Failed to marshal JSON")

		// json validation
		isValid := isValidJSON(string(modifiedMessage))
		fmt.Println("The provided JSON is valid:", isValid)

		// Re-queue the modified message to the same or a different queue
		err = ch.Publish(
			"",         // exchange
			queue_name, // routing key (queue)
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        modifiedMessage,
			})
		err = msg.Ack(false)
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", modifiedMessage)

		receivedCount++
	}
}
