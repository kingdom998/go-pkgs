package mysql

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Driver          string        `json:"driver,omitempty"`
	Source          string        `json:"source,omitempty"`
	MaxConnLifeTime time.Duration `json:"max_conn_life_time,omitempty"`
	MaxIdle         int32         `json:"max_idle,omitempty"`
	MaxOpen         int32         `json:"max_open,omitempty"`
}

func NewClient(conf *Config, logger log.Logger) *gorm.DB {
	helper := log.NewHelper(log.With(logger, "module", "marvel-service/data/gorm"))
	db, err := gorm.Open(mysql.Open(conf.Source), &gorm.Config{})
	if err != nil {
		helper.Fatalf("failed opening connection to mysql: %v", err)
	}

	// db tracing init
	err = db.Use(otelgorm.NewPlugin())
	if err != nil {
		helper.Fatalf("marvel service orm tracing init error: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		helper.Fatalf("marvel service orm error: %v", err)
	}

	// 置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(int(conf.MaxIdle))
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(int(conf.MaxOpen))
	// 设置连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(conf.MaxConnLifeTime)
	err = sqlDB.Ping()
	if err != nil {
		helper.Fatalf("ping mysql error: %v", err)
	}

	return db
}
