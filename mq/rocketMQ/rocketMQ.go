package rocketMQ

import (
	"context"
	"os"
	"time"

	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	v2 "github.com/apache/rocketmq-clients/golang/v5/protocol/v2"
	"github.com/go-kratos/kratos/v2/log"
)

type Config struct {
	Driver     string `json:"driver,omitempty"`
	LogLevel   string `json:"log_level,omitempty"`
	Endpoint   string `json:"endpoint,omitempty"`
	AccessKey  string `json:"access_key,omitempty"`
	SecretKey  string `json:"secret_key,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Topic      string `json:"topic,omitempty"`
	Group      string `json:"group,omitempty"`
	RetryCount int32  `json:"retry_count,omitempty"`
}

type rocketMQ struct {
	log *log.Helper

	config   *Config
	producer rmq.Producer
	consumer rmq.SimpleConsumer
}

func New(conf *Config, logger log.Logger) *rocketMQ {
	helper := log.NewHelper(log.With(logger, "module", "pkgs/mq/rocketMQ"))
	os.Setenv("mq.consoleAppender.enabled", "true")
	os.Setenv("rocketmq.client.logLevel", "warn")
	if conf.LogLevel != "" {
		os.Setenv("rocketmq.client.logLevel", conf.LogLevel) // 设置日志等级
	}
	rmq.ResetLogger()

	r := &rocketMQ{
		config: conf,
		log:    helper,
	}
	r.producer = r.newProducer(conf)
	r.consumer = r.newConsumer(conf)

	return r
}

func (r *rocketMQ) Finalise() {
	if r.producer != nil {
		r.producer.GracefulStop()
	}
	if r.consumer != nil {
		r.consumer.GracefulStop()
	}
}

func (r *rocketMQ) newProducer(conf *Config) rmq.Producer {
	producer, err := rmq.NewProducer(&rmq.Config{
		Endpoint:      conf.Endpoint,
		NameSpace:     conf.Namespace,
		ConsumerGroup: conf.Group,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    conf.AccessKey,
			AccessSecret: conf.SecretKey,
		},
	},
		rmq.WithTopics(conf.Topic),
	)
	if err != nil {
		log.Fatal(err)
	}
	// start producer
	err = producer.Start()
	if err != nil {
		log.Fatal(err)
	}
	return producer
}

func (r *rocketMQ) newConsumer(conf *Config) rmq.SimpleConsumer {
	cosumerConf := &rmq.Config{
		Endpoint:      conf.Endpoint,
		ConsumerGroup: conf.Group,
		NameSpace:     conf.Namespace,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    conf.AccessKey,
			AccessSecret: conf.SecretKey,
		},
	}
	consumer, err := rmq.NewSimpleConsumer(cosumerConf,
		rmq.WithSubscriptionExpressions(map[string]*rmq.FilterExpression{
			conf.Topic: rmq.SUB_ALL,
		}),
		rmq.WithAwaitDuration(time.Second*5),
	)
	if err != nil {
		log.Fatal(err)
	}
	// start simpleConsumer
	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}
	return consumer
}

func (r *rocketMQ) Publish(ctx context.Context, topic string, body []byte) (err error) {
	// 构造消息内容
	msg := &rmq.Message{
		Topic: topic, // 设置topic名称
		Body:  body,
	}
	// 发送消息
	srs, err := r.producer.Send(ctx, msg)
	if err != nil {
		r.log.Warnw("msg", "send message failed", "err", err)
		return
	}
	for _, sr := range srs {
		r.log.Infow("msg", "send message success", "msg_id", sr.MessageID)
	}

	return

}

func (r *rocketMQ) Subscribe(ctx context.Context, callback func(context.Context, []byte) error) (err error) {
	for {
		mvs, err := r.consumer.Receive(ctx, 1, 10*time.Second)
		if err != nil {
			status, _ := rmq.AsErrRpcStatus(err)
			if status.GetCode() == int32(v2.Code_MESSAGE_NOT_FOUND) {
				continue
			}
			r.log.Warnw("msg", "receive msg error", "err", err)
			break
		}
		// ack message
		for _, mv := range mvs {
			callback(ctx, mv.GetBody())
			err = r.consumer.Ack(ctx, mv)
			if err != nil {
				r.log.Warnw("msg", "ack msg error", "err", err)
				break
			}
		}
		if err != nil {
			break
		}
	}

	return err
}
