package messaging

import (
	"TMS/notification-service/internal/models"
	"TMS/notification-service/internal/repository"
	"encoding/json"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
}

func (r *RabbitMQ) Consume(repo *repository.Repository) error {
	msgs, err := r.Channel.Consume(
		"notifications",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		var request models.Notification

		err := json.Unmarshal(msg.Body, &request)
		if err != nil {
			slog.Error("Invalid notification event", "error", err)
			_ = msg.Nack(false, false) //json bozuk bir daha deneme
			continue
		}

		err = repo.CreateNotification(&request)
		if err != nil {
			slog.Error("Notification can't create", "error", err)
			_ = msg.Nack(false, true) //mesaj başarısız queue bir daha ekle
			continue
		}

		_ = msg.Ack(false) //mesaj başarılı
	}

	return nil
}
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = channel.QueueDeclare(
		"notifications",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: channel,
	}, nil
}
