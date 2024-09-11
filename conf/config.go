package conf

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kingdom998/go-pkgs/cache/redis"
	"github.com/kingdom998/go-pkgs/http"
	"github.com/kingdom998/go-pkgs/mq/rabbitMQ"
	"github.com/kingdom998/go-pkgs/storage/bos"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP    http.Config     `json:"http"`
	Redis   redis.Config    `json:"redis"`
	MQ      rabbitMQ.Config `json:"mq"`
	Storage bos.Config      `json:"storage"`
	WebUI   http.Config     `json:"webUI"`
	ComfyUI http.Config     `json:"comfyUI"`
	Worker  http.Config     `json:"worker"`
}

func New(confFile string) *Config {
	var config Config

	v := viper.New()
	v.SetConfigFile(confFile)
	// set decode tag_name, default is mapstructure
	decoderConfigOption := func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	} // REQUIRED if the conf file does not have the extension in the name
	err := v.ReadInConfig() // Find and read the conf file
	if err != nil {
		log.Fatalf("parse %s failed, with err: %+v", confFile, err)
	}
	if err = v.Unmarshal(&config, decoderConfigOption); err != nil {
		log.Fatalf("Unmarshal conf file error: %v", err)
	}
	return &config
}
