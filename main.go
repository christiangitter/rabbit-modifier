package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		fmt.Println("RABBITMQ_URL is not set")
		return
	}

	// Use rabbitMQURL
	fmt.Println(rabbitMQURL)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the key to be modified: ")
	//The _ is used to ignore the error returned
	key, _ := reader.ReadString('\n')

	fmt.Print("Enter the new value: ")
	value, _ := reader.ReadString('\n')

	// Remove newline characters
	key = strings.Replace(key, "\n", "", -1)
	key = strings.Replace(key, "\r", "", -1)
	value = strings.Replace(value, "\n", "", -1)
	value = strings.Replace(value, "\r", "", -1)

	// Connect to RabbitMQ server
	// This is a short variable declaration in Go. Conn set up the connection and the err hold any error that occurred during the Consume call.
	conn, err := amqp.Dial(rabbitMQURL)
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
			payLoad[key] = value
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
