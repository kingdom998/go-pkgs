package cos

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type Config struct {
	Host      string `json:"host,omitempty"`
	Region    string `json:"region,omitempty"`
	Bucket    string `json:"bucket,omitempty"`
	SecretID  string `json:"secret_id,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
}

type storage struct {
	cos *cos.Client
}

func NewClient(config Config) *storage {
	u, _ := url.Parse(config.Host)
	b := &cos.BaseURL{BucketURL: u}
	return &storage{
		cos: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  config.SecretID,
				SecretKey: config.SecretKey,
			},
		}),
	}

}

func (s *storage) ListBuckets(ctx context.Context) (buckets []string, err error) {
	b, _, err := s.cos.Service.Get(ctx)
	if err != nil {
		return nil, err
	}
	for _, b := range b.Buckets {
		buckets = append(buckets, b.Name)
	}
	return
}

func (s *storage) UploadFromFile(ctx context.Context, cosFilePath, localFilePath string) (string, error) {
	result, err := s.cos.Object.PutFromFile(ctx, cosFilePath, localFilePath, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s/%s", result.Request.URL.Host, cosFilePath), err
}

func (s *storage) UploadFromBytes(ctx context.Context, cosFilePath string, body []byte) (string, error) {
	reader := bytes.NewReader(body)
	result, err := s.cos.Object.Put(ctx, cosFilePath, reader, nil)
	return fmt.Sprintf("https://%s/%s", result.Request.URL.Host, cosFilePath), err
}

func (s *storage) Download2File(ctx context.Context, filename, localpath string) (err error) {
	_, err = s.cos.Object.GetToFile(ctx, filename, filepath.Base(localpath), nil)
	return err
}
