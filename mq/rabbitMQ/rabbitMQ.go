package rabbitMQ

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Endpoint     string
	Port         string
	UserName     string
	Password     string
	Vhost        string
	Exchange     string
	ExchangeType string
	Topic        string
	RetryCount   int32
}

type rabbitMQ struct {
	log *log.Helper

	config *Config
	ch     *amqp.Channel
}

func NewRabbitMQ(config *Config, logger log.Logger) *rabbitMQ {
	helper := log.NewHelper(log.With(logger, "module", "pkgs/mq/rabbitMQ"))
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", config.UserName, config.Password, config.Endpoint, config.Port, config.Vhost)
	conn, err := amqp.Dial(url)
	if err != nil {
		helper.Fatalf("Failed to connect to rabbitMQ: %s", err)
	}

	// 建立通道
	ch, err := conn.Channel()
	if err != nil {
		helper.Fatalf("Failed to open a channel: %s", err)
	}
	return &rabbitMQ{
		ch:     ch,
		config: config,
		log:    helper,
	}
}


func (r *rabbitMQ) Finalise() {
	if r.ch != nil {
		r.ch.Close()
	}
}

func (r *rabbitMQ) Publish(ctx context.Context, topic string, msg []byte) error {
	// 声明消息队列
	_, err := r.ch.QueueDeclare(
		topic,
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	// 发布消息到指定的消息队列
	err = r.ch.Publish(
		"",    // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	return err
}

func (r *rabbitMQ) PublishWithEx(ctx context.Context, topic string, msg []byte) error {
	err := r.ch.ExchangeDeclare(
		r.config.Exchange,
		r.config.ExchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %s", err)
	}

	// 发布消息到指定的消息队列
	err = r.ch.Publish(
		r.config.Exchange,
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	return err
}

func (r *rabbitMQ) Subscribe(ctx context.Context, callback func(context.Context, []byte) error) error {
	// 创建消费者并消费指定消息队列中的消息
	msgList, err := r.ch.Consume(
		r.config.Topic, // message-queue
		"consumer",     // todo 后续进行传参
		false,          // 设置为非自动确认(可根据需求自己选择)
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		r.log.Fatalf("new consume client failed with err %v", err)
	}

	for d := range msgList {
		err = callback(ctx, d.Body)
		if err != nil {
			r.log.Fatalf("Failed to handle %s with %s", d.Body, err)
		}
		// 手动回复ack
		d.Ack(true)
	}

	return err
}
