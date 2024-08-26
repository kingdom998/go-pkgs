package data

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"net/url"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"

	"github.com/kingdom998/go-pkgs/conf"
)

type COS struct {
	Client *cos.Client
	Url    string
}

func NewDB(conf *conf.Database, logger log.Logger) *gorm.DB {
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

func NewCache(conf *conf.Redis, logger log.Logger) *redis.Client {
	helper := log.NewHelper(log.With(logger, "module", "ai-server/data/redis"))
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		DB:           int(conf.Db),
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		DialTimeout:  conf.DialTimeout,
		PoolSize:     int(conf.PoolSize),
	})
	cmd := rdb.Ping(context.Background())
	if cmd.Err() != nil {
		helper.Fatalf("ping redis error: %v", cmd.Err())
	}

	return rdb
}

func NewCOS(config *conf.COS) *COS {
	// 存储桶名称，由 bucketname-appid 组成，appid 必须填入，可以在 COS 控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
	// 替换为用户的 region，存储桶 region 可以在 COS 控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
	rawURL := fmt.Sprintf(config.Host, config.Bucket, config.Region)
	u, _ := url.Parse(rawURL)
	b := &cos.BaseURL{BucketURL: u}
	return &COS{
		Client: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  config.SecretID,
				SecretKey: config.SecretKey,
			},
		}),
		Url: rawURL,
	}
}
