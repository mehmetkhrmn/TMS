package messaging

import (
	"TMS/notification-service/internal/models"
	"TMS/notification-service/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/lib/pq"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn     *amqp091.Connection
	Channel  *amqp091.Channel
	Confirms chan amqp091.Confirmation
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
				if !IsRetryable(err) {
					slog.Error("Permanent error , sending to DLQ")
					_ = msg.Nack(false, false)
					continue
				}
				retryCount := getRetryCount(msg)

				if retryCount >= 3 {
					slog.Error("Notification failed after maximum retries")
					_ = msg.Nack(false, false)
					continue
					continue
				}

				retryCount++

				if err := r.publishRetry(msg, retryCount); err != nil {
					slog.Error("Retry publish failed", "error", err)
					_ = msg.Nack(false, true)
					continue
				}

				_ = msg.Ack(false)
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
	//Retry queues
	// Retry queue - 5 saniye
	_, err = r.Channel.QueueDeclare(
		"notifications.retry.5s",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             int32(5000), //time to live mesajin yasaycagi sure
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": os.Getenv("RABBITMQ_QUEUE"),
		},
	)
	if err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}

	// Retry queue - 10 saniye
	_, err = r.Channel.QueueDeclare(
		"notifications.retry.10s",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             int32(10000),
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": os.Getenv("RABBITMQ_QUEUE"),
		},
	)
	if err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}

	// Retry queue - 15 saniye
	_, err = r.Channel.QueueDeclare(
		"notifications.retry.15s",
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-message-ttl":             int32(15000),
			"x-dead-letter-exchange":    "",
			"x-dead-letter-routing-key": os.Getenv("RABBITMQ_QUEUE"),
		},
	)
	if err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return err
	}
	slog.Info("RabbitMQ connection established", "url", url)
	return nil
}
func (r *RabbitMQ) publishRetry(msg amqp091.Delivery, retryCount int) error {
	queue := ""

	switch retryCount {
	case 1:
		queue = "notifications.retry.5s"
	case 2:
		queue = "notifications.retry.10s"
	case 3:
		queue = "notifications.retry.15s"
	default:
		return errors.New("maximum retry count exceeded")
	}

	headers := amqp091.Table{
		"x-retry-count": retryCount,
	}

	err := r.Channel.Publish(
		"",
		queue,
		false,
		false,
		amqp091.Publishing{
			Headers:      headers,
			Body:         msg.Body,
			ContentType:  msg.ContentType,
			DeliveryMode: amqp091.Persistent,
		},
	)
	if err != nil {
		return err
	}
	confirmation := <-r.Confirms //brokera tarafindan publish confirm edildi
	if !confirmation.Ack {
		return errors.New("failed to acknowledge retry message")
	}
	return nil
}
func getRetryCount(msg amqp091.Delivery) int { //retry sayimizi header ile bulacagiz
	value, ok := msg.Headers["x-retry-count"]
	if !ok {
		return 0
	}

	count, ok := value.(int32)
	if !ok {
		return 0
	}

	return int(count)
}
func IsRetryable(err error) bool {
	// geçici DB hatasıysa true
	// kalıcı hataysa false
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code.Class() {
		case "08": // connection exception
			return true
		case "40": // transaction rollback
			return true
		default:
			return false
		}
	}

	// networkte olusabilcek hatalar de eklendi
	var netErr net.Error
	if errors.As(err, &netErr) {
		if	netErr.Timeout() || netErr.Temporary(){
			return true
		}
	}

	return false
}
