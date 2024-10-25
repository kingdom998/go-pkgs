package rocketMQ

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/joho/godotenv/autoload"
)

var (
	config Config
)

func init() {
	config = Config{
		Endpoint:   os.Getenv("endpoint"),
		AccessKey:  os.Getenv("access_key"),
		SecretKey:  os.Getenv("secret_key"),
		Namespace:  os.Getenv("namespace"),
		Topic:      os.Getenv("topic"),
		Group:      os.Getenv("group"),
		RetryCount: 3,
	}
	log.Infof("load config %+v", config)
}

func TestPublish(t *testing.T) {
	ctx := context.Background()
	logger := log.NewStdLogger(os.Stdout)
	client := New(&config, logger)
	msg := "this is rockemq " + time.Now().Format("2006-01-02 15:04:05")
	for i := 0; i < 5; i++ {
		err := client.Publish(ctx, config.Topic, []byte(msg))
		if err != nil {
			log.Warnf("send message to rabbitmq error: %v", err)
		}
	}

}

func TestSubscribe(t *testing.T) {
	logger := log.GetLogger()
	client := New(&config, logger)
	ctx := context.Background()
	client.Subscribe(ctx, func(ctx context.Context, body []byte) error {
		log.Infof(string(body))
		return nil
	})
	defer client.Finalise()
}
