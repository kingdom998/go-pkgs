package bos

import (
	"context"
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/bos"
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
	client *bos.Client
}

func NewClient(config *Config) *storage {
	bosClient, err := bos.NewClient(config.Ak, config.Sk, config.Endpoint)
	if err != nil {
		panic(fmt.Sprintf("create bos client failed, err:%v", err))
	}
	return &storage{
		config: config,
		client: bosClient,
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

func (s *storage) UploadFromFile(ctx context.Context, cosFilePath, localFilePath string) (string, error) {
	result, err := s.client.ParallelUpload(s.config.Bucket, cosFilePath, localFilePath, "", nil)
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func (s *storage) Download2File(ctx context.Context, objectPath, localpath string) (err error) {
	return s.client.BasicGetObjectToFile(s.config.Bucket, objectPath, localpath)
}

func (s *storage) UploadFromBytes(ctx context.Context, objectPath string, body []byte) (string, error) {
	_, err := s.client.PutObjectFromBytes(s.config.Bucket, objectPath, body, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s/%s/%s", s.config.Endpoint, s.config.Bucket, objectPath), nil
}
