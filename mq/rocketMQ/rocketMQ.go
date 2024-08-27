package rocketMQ

import (
	"context"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/go-kratos/kratos/v2/log"
)

type Config struct {
	Endpoint   string `json:"endpoint,omitempty"`
	SecretKey  string `json:"secret_key,omitempty"`
	AccessKey  string `json:"access_key,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Topic      string `json:"topic,omitempty"`
	Group      string `json:"group,omitempty"`
	RetryCount int32  `json:"retry_count,omitempty"`
}

type rocketMQ struct {
	log *log.Helper

	config   *Config
	producer rocketmq.Producer
	consumer rocketmq.PullConsumer
}

func NewRocketMQ(config *Config, logger log.Logger) *rocketMQ {
	helper := log.NewHelper(log.With(logger, "module", "pkgs/mq/rocketMQ"))

	r := &rocketMQ{
		config: config,
		log:    helper,
	}
	r.producer = r.newProducer(config)
	r.consumer = r.newConsumer(config)

	return r
}

func (r *rocketMQ) newProducer(conf *Config) rocketmq.Producer {
	// 创建消息生产者
	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Endpoint})), // 设置服务地址
		// 设置acl权限
		producer.WithCredentials(primitive.Credentials{
			SecretKey: conf.SecretKey,
			AccessKey: conf.AccessKey,
		}),
		producer.WithNamespace(conf.Namespace),   // 设置命名空间名称
		producer.WithRetry(int(conf.RetryCount)), // 设置发送失败重试次数
	)
	if err != nil {
		r.log.Fatalf("init producer error: %v", err)
	}

	// 启动producer
	err = p.Start()
	if err != nil {
		r.log.Fatalf("start producer error: %v", err)
	}
	return p
}

func (r *rocketMQ) newConsumer(conf *Config) rocketmq.PullConsumer {
	var nameSrv, err = primitive.NewNamesrvAddr(conf.Endpoint)
	if err != nil {
		r.log.Fatalf("NewNamesrvAddr err: %v", err)
	}

	pullConsumer, err := rocketmq.NewPullConsumer(
		consumer.WithGroupName(conf.Group),
		consumer.WithNameServer(nameSrv),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: conf.AccessKey,
			SecretKey: conf.SecretKey,
		}),
		consumer.WithNamespace(conf.Namespace),
		consumer.WithMaxReconsumeTimes(conf.RetryCount),
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset), // 设置从起始位置开始消费
		consumer.WithConsumerModel(consumer.Clustering),                // 设置消费模式（默认集群模式）
	)

	c, err := rocketmq.NewPullConsumer(
		consumer.WithGroupName(conf.Group),                                                 // 设置消费者组
		consumer.WithNsResolver(primitive.NewPassthroughResolver([]string{conf.Endpoint})), // 设置服务地址
		// 设置acl权限
		consumer.WithCredentials(primitive.Credentials{
			SecretKey: conf.SecretKey,
			AccessKey: conf.AccessKey,
		}),
		consumer.WithNamespace(conf.Namespace),                         // 设置命名空间名称
		consumer.WithConsumeFromWhere(consumer.ConsumeFromFirstOffset), // 设置从起始位置开始消费
		consumer.WithConsumerModel(consumer.Clustering),                // 设置消费模式（默认集群模式）
	)
	if err != nil {
		r.log.Fatalf("init consumer error: %v", err)
	}
	err = pullConsumer.Subscribe(conf.Topic, consumer.MessageSelector{})
	if err != nil {
		r.log.Fatalf("fail to Subscribe: %v", err)
	}
	err = pullConsumer.Start()
	if err != nil {
		r.log.Fatalf("fail to Start: %v", err)
	}
	return c
}

func (r *rocketMQ) Publish(ctx context.Context, body []byte) (err error) {
	// 构造消息内容
	mq := &primitive.Message{
		Topic: r.config.Topic, // 设置topic名称
		Body:  body,
	}
	// 发送消息
	res, err := r.producer.SendSync(ctx, mq)
	if err != nil {
		r.log.Warnf("send message error: %v", err)
		return
	}
	r.log.Infof("send message success: result=%s\n", res.String())

	return

}

func (r *rocketMQ) Subscribe(ctx context.Context, callback func(context.Context, []byte) error) (err error) {
	r.consumer.Start()
	for {
		cr, err := r.consumer.Poll(ctx, time.Second*5)
		if consumer.IsNoNewMsgError(err) {
			continue
		}
		if err != nil {
			r.log.Warnw("msg", "pull consumer error", "err", err)
			time.Sleep(time.Second)
			continue
		}
		for _, msg := range cr.GetMsgList() {
			err = callback(ctx, msg.Body)
		}
		r.consumer.ACK(ctx, cr, consumer.ConsumeSuccess)
	}

	return

}
