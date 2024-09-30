package oss

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Config struct {
	Driver   string `json:"driver"`
	Endpoint string `json:"endpoint,omitempty"`
	Ak       string `json:"ak,omitempty"`
	Sk       string `json:"sk,omitempty"`
	Bucket   string `json:"bucket,omitempty"`
}

type storage struct {
	config *Config
	client *oss.Client
}

func NewClient(config *Config) *storage {
	client, err := oss.New(config.Endpoint, config.Ak, config.Sk)
	if err != nil {
		panic(fmt.Sprintf("create oss client failed, err:%v", err))
	}
	return &storage{
		config: config,
		client: client,
	}
}

func (s *storage) ListBuckets(ctx context.Context) (buckets []string, err error) {
	result, err := s.client.ListBuckets()
	if err != nil {
		return nil, err
	}
	for _, bucket := range result.Buckets {
		buckets = append(buckets, bucket.Name)
	}

	return
}

func (s *storage) UploadFromFile(ctx context.Context, objectPath, localObjectPath string) (string, error) {
	bucket, err := s.client.Bucket(s.config.Bucket)
	if err != nil {
		return "", err
	}

	err = bucket.PutObjectFromFile(objectPath, localObjectPath)
	if err != nil {
		return "", err
	}

	return objectPath, nil
}

func (s *storage) Download2File(ctx context.Context, objectPath, localpath string) (err error) {
	bucket, err := s.client.Bucket(s.config.Bucket)
	if err != nil {
		return
	}

	err = bucket.DownloadFile(objectPath, localpath, 100*1024, oss.Routines(3), oss.Checkpoint(true, ""))
	return
}

func (s *storage) UploadFromBytes(ctx context.Context, objectPath string, body []byte) (string, error) {
	bucket, err := s.client.Bucket(s.config.Bucket)
	if err != nil {
		return "", err
	}

	err = bucket.PutObject(objectPath, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	return objectPath, nil
}
