package mysql

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
		Driver:  os.Getenv("driver"),
		Source:  os.Getenv("source"),
		MaxIdle: 20,
		MaxOpen: 20,
	}
	log.Infof("config is %v", config)
}

func TestNewClient(t *testing.T) {
	logger := log.GetLogger()

	client := NewClient(&config, logger)
	_ = client
}
