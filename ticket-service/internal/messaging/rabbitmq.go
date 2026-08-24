package messaging

import (
	"TMS/ticket-service/internal/models"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn     *amqp.Connection
	Channel  *amqp.Channel
	Confirms <-chan amqp.Confirmation
}

func NewRabbitMQ() (*RabbitMQ, error) { //database baglantisi ile ayni retry mantik
	rabbit := &RabbitMQ{}

	var err error

	for attempt := 1; attempt <= 5; attempt++ {
		err = rabbit.Connect()

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

func (r *RabbitMQ) QueueDeclare() error {
	queueName := os.Getenv("RABBITMQ_QUEUE")
	exchangeName := os.Getenv("RABBITMQ_EXCHANGE")
	dlxName := os.Getenv("RABBITMQ_DLX")
	dlqName := os.Getenv("RABBITMQ_DLQ")
	routingKey := os.Getenv("RABBITMQ_ROUTING_KEY")

	// Normal exchange
	err := r.Channel.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Dead Letter Exchange
	err = r.Channel.ExchangeDeclare(
		dlxName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Normal queue
	_, err = r.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange":    dlxName,
			"x-dead-letter-routing-key": routingKey,
		},
	)
	if err != nil {
		return err
	}

	// Dead Letter Queue
	_, err = r.Channel.QueueDeclare(
		dlqName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Normal exchange -> normal queue
	err = r.Channel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// DLX -> DLQ
	err = r.Channel.QueueBind(
		dlqName,
		routingKey,
		dlxName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = r.Channel.Confirm(false) //chaneli olusturdum confirm aliyorum ki mesajlarida confirmelyebileyim
	if err != nil {
		return err
	}

	r.Confirms = r.Channel.NotifyPublish(
		make(chan amqp.Confirmation, 1),
	)
	return nil
}
func (r *RabbitMQ) PublishNotification(notification models.NotificationRequest) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	err = r.publish(data)
	if err != nil {
		slog.Error("RabbitMQ publish failed, reconnecting", "error", err)

		if reconnectErr := r.Connect(); reconnectErr != nil {
			return reconnectErr
		}

		err = r.publish(data)
		if err != nil {
			return err
		}
	}

	slog.Info("RabbitMQ publish başarılı")
	return nil
}
func (r *RabbitMQ) publish(data []byte) error { //pakete ozel func
	err := r.Channel.Publish(
		os.Getenv("RABBITMQ_EXCHANGE"),
		os.Getenv("RABBITMQ_ROUTING_KEY"),
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return err
	}

	confirmation := <-r.Confirms //publish confirmation -> burada anlami publisher olarak mesaji gonderdim bir sorun yok, consumer tarafi bizi burada ilgilendirmiyor

	if !confirmation.Ack {
		return errors.New("RabbitMQ publisher confirm NACK") //noacknowledge ,confirm edilmedi
	}

	return nil
}
func (r *RabbitMQ) Connect() error {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}

	r.Conn = conn
	r.Channel = ch

	if err := r.QueueDeclare(); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	return nil
}
