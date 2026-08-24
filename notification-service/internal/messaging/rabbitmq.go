package messaging

import (
	"TMS/notification-service/internal/models"
	"TMS/notification-service/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
}

func (r *RabbitMQ) Consume(ctx context.Context, repo *repository.Repository) error {
	err := r.Channel.Qos(
		10, // aynı anda maksimum 10 unacked mesaj
		0,  // prefetch size
		false,
	)
	if err != nil {
		return err
	}
	msgs, err := r.Channel.Consume(
		os.Getenv("RABBITMQ_QUEUE"),
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
	slog.Info("RabbitMQ consuming messages")
	for {
		select { //hangi caseden once veri gelirse o calisiyor
		case <-ctx.Done(): // context cancel ile iptal ediliyor ve done artik beklenmez oluyor o yuzden case si calisiyor
			return nil

		case msg, ok := <-msgs: //msgs kanalindan normal msg geliyor ve channel durumunu assignliyor gerisi eskisi ile ayni
			if !ok {
				return errors.New("RabbitMQ consumer stopped: deliveries channel closed")
			}

			var request models.Notification

			err := json.Unmarshal(msg.Body, &request)
			if err != nil {
				slog.Error("Invalid notification event", "error", err)
				_ = msg.Nack(false, false)
				continue
			}

			err = repo.CreateNotification(&request)
			if err != nil {
				slog.Error("Notification can't create", "error", err)
				_ = msg.Nack(false, true)
				continue
			}

			_ = msg.Ack(false)
		}
	}
}
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	rabbit := &RabbitMQ{}

	var err error

	for attempt := 1; attempt <= 5; attempt++ { //retry merkezi hale getirildi
		err = rabbit.Connect(url)

		if err == nil {
			return rabbit, nil
		}

		if attempt < 5 {
			delay := time.Duration(attempt*3) * time.Second
			time.Sleep(delay)
		}
	}

	return nil, err
}

func (r *RabbitMQ) Connect(url string) error {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil { //  channelde sorun varsa conn close
		_ = conn.Close()
		return err
	}

	r.Conn = conn //ampq deki conn ve cahaneli Rabbite bagladik
	r.Channel = channel

	slog.Info("RabbitMQ connection established", "url", url)
	return nil
}
