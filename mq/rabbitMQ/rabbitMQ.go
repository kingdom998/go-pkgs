package rabbitMQ

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/kingdom998/go-pkgs/conf"
)

type rabbitMQ struct {
	log *log.Helper

	topic string
	ch    *amqp.Channel
}

func NewRabbitMQ(config *conf.RabbitMQ, logger log.Logger) *rabbitMQ {
	helper := log.NewHelper(log.With(logger, "module", "pkgs/mq/rabbitMQ"))
	url := fmt.Sprintf(config.Url, config.UserName, config.Password, config.Endpoint, config.Vhost)
	conn, err := amqp.Dial(url)
	if err != nil {
		helper.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	// 建立通道
	ch, err := conn.Channel()
	if err != nil {
		helper.Fatalf("Failed to open a channel: %s", err)
	}
	return &rabbitMQ{
		ch:    ch,
		topic: config.Route,
		log:   helper,
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
		r.log.Fatalf("Failed to declare a queue: %s", err)
		return err
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

func (r *rabbitMQ) Subscribe(ctx context.Context, topic string, callback func(context.Context, []byte) error) error {
	// 创建消费者并消费指定消息队列中的消息
	msgList, err := r.ch.Consume(
		topic, // message-queue
		"",    // consumer
		false, // 设置为非自动确认(可根据需求自己选择)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	for d := range msgList {
		err = callback(ctx, d.Body)
		if err != nil {
			r.log.Fatalf("Failed to handle %s with %s", d.Body, err)
		}
		// 手动回复ack
		d.Ack(false)
	}

	return err
}
