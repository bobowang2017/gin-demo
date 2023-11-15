package mq

import (
	"context"
	"github.com/wagslane/go-rabbitmq"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// https://github.com/wagslane/go-rabbitmq

type RabbitMqClient struct {
	conn      *rabbitmq.Conn
	publisher *rabbitmq.Publisher
}

func (r *RabbitMqClient) Publish(exchange, routingKey string, msg []byte) error {
	return r.publisher.PublishWithContext(
		context.Background(),
		msg,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange(exchange),
	)
}

var (
	rabbitMqClient     *RabbitMqClient
	rabbitMqClientOnce sync.Once
)

func NewRabbitMqClient() *RabbitMqClient {
	rabbitMqClientOnce.Do(func() {
		conn, err := rabbitmq.NewConn(
			"amqp://guest:guest@localhost",
			rabbitmq.WithConnectionOptionsLogging,
		)
		if err != nil {
			panic(err)
		}
		publisher, err := rabbitmq.NewPublisher(
			conn,
			rabbitmq.WithPublisherOptionsLogging,
			rabbitmq.WithPublisherOptionsExchangeDeclare,
		)
		if err != nil {
			panic(err)
		}
		rabbitMqClient = &RabbitMqClient{
			conn:      conn,
			publisher: publisher,
		}
	})
	return rabbitMqClient
}

func StartConsumer(handler rabbitmq.Handler) {
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@localhost",
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	consumer, err := rabbitmq.NewConsumer(
		conn,
		handler,
		//func(d rabbitmq.Delivery) rabbitmq.Action {
		//	log.Logger.Infof("consumed: %v", string(d.Body))
		//	// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
		//	return rabbitmq.Ack
		//},
		"my_queue",
		rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
