package rabbitMQ

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
		Url:      os.Getenv("url"),
		Endpoint: os.Getenv("endpoint"),
		UserName: os.Getenv("username"),
		Password: os.Getenv("password"),
		Vhost:    os.Getenv("vhost"),
		Topic:    os.Getenv("route"),
	}
	log.Infof("load config %+v", config)
}

func TestSendMessage(t *testing.T) {
	ctx := context.Background()
	logger := log.NewStdLogger(os.Stdout)
	client := NewRabbitMQ(&config, logger)
	msg := "welcome at " + time.Now().Format("2006-01-02 15:04:05")
	for i := 0; i < 5; i++ {
		err := client.Publish(ctx, config.Topic, []byte(msg))
		if err != nil {
			log.Warnf("send message to rabbitmq error: %v", err)
		}
	}

}

func TestSubscribe(t *testing.T) {
	logger := log.GetLogger()
	client := NewRabbitMQ(&config, logger)
	ctx := context.Background()
	client.Subscribe(ctx, func(ctx context.Context, body []byte) error {
		log.Infof(string(body))
		return nil
	})
}
