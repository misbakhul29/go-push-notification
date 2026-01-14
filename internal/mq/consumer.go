package mq

import (
	"context"
	"encoding/json"
	"go-push-service/internal/config"
	"go-push-service/internal/models"
	"go-push-service/pkg/logger"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func StartConsumer(cfg *config.Config, rdb *redis.Client) {
	conn, err := amqp091.Dial(cfg.RabbitMQURL)
	if err != nil {
		logger.Fatal(err, "Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal(err, "Failed to open a channel")
	}

	q, err := ch.QueueDeclare(
		cfg.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatal(err, "Failed to declare a queue")
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatal(err, "Failed to register a consumer")
	}

	go func() {
		ctx := context.Background()
		for d := range msgs {
			var job models.NotificationJob
			if err := json.Unmarshal(d.Body, &job); err != nil {
				logger.Error(err, "Error unmarshalling message")
				d.Nack(false, false)
				continue
			}

			logger.Infof("job_id", job.ID, "Processing notification job")

			channelName := "user:" + job.TargetUserID + ":notify"
			err := rdb.Publish(ctx, channelName, d.Body).Err()

			if err != nil {
				logger.Error(err, "Failed to publish to Redis")
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	logger.Info("RabbitMQ Consumer started and listening...")
}
