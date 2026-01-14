package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type NotificationJob struct {
	TargetUserID string `json:"target_user_id"`
	Title        string `json:"title"`
	Message      string `json:"message"`
}

func main() {
	conn, _ := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()
	ch, _ := conn.Channel()
	defer ch.Close()

	// data to send
	job := NotificationJob{
		TargetUserID: "user123", // Must match the WS client ID
		Title:        "System Alert",
		Message:      "Your server usage is at 90%",
	}
	body, _ := json.Marshal(job)

	// Publish
	ch.PublishWithContext(context.Background(), "", "notifications_queue", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	log.Println("Message sent to Queue!")
}
