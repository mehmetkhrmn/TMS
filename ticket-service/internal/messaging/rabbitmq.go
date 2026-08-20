package messaging

import (
	"TMS/ticket-service/internal/notification"
	"encoding/json"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ() (*RabbitMQ, error) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{Conn: conn, Channel: ch}, nil

}
func (r *RabbitMQ) QueueDeclare() error {
	_, err := r.Channel.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE"),
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		return err
	}
	return nil
}
func (r *RabbitMQ) PublishNotification(notification notification.NotificationRequest) error {
	data, err := json.Marshal(notification) //yine isteği struct yapısına göre düzenliyoruz json olarak
	if err != nil {
		return err
	}
	err = r.Channel.Publish(os.Getenv("RABBITMQ_EXCHANGE"),
		os.Getenv("RABBITMQ_ROUTING_KEY"),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
func (r *RabbitMQ) Publish(event notification.NotificationRequest) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return r.Channel.Publish("",
		"notifcations",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
