package redis

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/joho/godotenv/autoload"
)

var (
	config Config
)

func init() {
	config = Config{
		Addr:     "172.16.10.163:6379",
		UserName: "",
		Password: "",
	}
}

func TestNewClient(t *testing.T) {
	logger := log.GetLogger()

	client := NewClient(&config, logger)
	_ = client
}
