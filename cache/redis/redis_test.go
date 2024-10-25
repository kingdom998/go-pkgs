package redis

import (
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/joho/godotenv/autoload"
)

var (
	config Config
)

func init() {
	config = Config{
		Addr:     os.Getenv("addr"),
		Username:  os.Getenv("user"),
		Password: os.Getenv("password"),
		Db: 1,
	}
	log.Infof("config is %v\n", config)
}

func TestNewClient(t *testing.T) {
	logger := log.GetLogger()

	client := NewClient(&config, logger)
	_ = client
}
