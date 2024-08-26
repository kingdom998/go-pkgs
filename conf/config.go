package conf

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HTTP     `json:"http"`
	Database Database `json:"database"`
	Redis    Redis    `json:"redis"`
	RocketMQ RocketMQ `json:"rocketMQ"`
	RabbitMQ RabbitMQ `json:"rabbitMQ"`
	COS      COS      `json:"cos"`
	WebUI    HTTP     `json:"webUI"`
	ComfyUI  HTTP     `json:"comfyUI"`
}

type HTTP struct {
	Addr string `json:"addr"`
}

type Database struct {
	Driver          string        `json:"driver,omitempty"`
	Source          string        `json:"source,omitempty"`
	MaxConnLifeTime time.Duration `json:"max_conn_life_time,omitempty"`
	MaxIdle         int32         `json:"max_idle,omitempty"`
	MaxOpen         int32         `json:"max_open,omitempty"`
}

type Redis struct {
	Addr         string
	Password     string
	Db           int64
	PoolSize     int32
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration
}

type RocketMQ struct {
	Endpoint   string `json:"endpoint,omitempty"`
	SecretKey  string `json:"secret_key,omitempty"`
	AccessKey  string `json:"access_key,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Topic      string `json:"topic,omitempty"`
	Group      string `json:"group,omitempty"`
	RetryCount int32  `json:"retry_count,omitempty"`
}

type RabbitMQ struct {
	Url        string
	Endpoint   string
	UserName   string
	Password   string
	Vhost      string
	Exchange   string
	Route      string
	Group      string
	RetryCount int32
}

type COS struct {
	Host      string `json:"host,omitempty"`
	Region    string `json:"region,omitempty"`
	Bucket    string `json:"bucket,omitempty"`
	SecretID  string `json:"secret_id,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
}

func New(confFile string) *Config {
	var config Config

	v := viper.New()
	v.SetConfigFile(confFile)
	//v.AddConfigPath("conf") // optionally look for conf in the working directory
	//v.SetConfigName("conf") // name of conf file (without extension)
	//v.SetConfigType("yaml")
	// set decode tag_name, default is mapstructure
	decoderConfigOption := func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	} // REQUIRED if the conf file does not have the extension in the name
	err := v.ReadInConfig() // Find and read the conf file
	if err != nil {
		log.Fatalf("parse %s failed", confFile)
	}
	if err = v.Unmarshal(&config, decoderConfigOption); err != nil {
		log.Fatalf("Unmarshal conf file error: %v", err)
	}
	return &config
}
