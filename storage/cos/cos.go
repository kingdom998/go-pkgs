package cos

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kingdom998/go-pkgs/conf"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"path/filepath"
)

type cosClient struct {
	cos *cos.Client
}

func New(config conf.COS) *cosClient {
	rawURL := fmt.Sprintf(config.Host, config.Bucket, config.Region)
	fmt.Println(rawURL)
	u, _ := url.Parse(rawURL)
	b := &cos.BaseURL{BucketURL: u}
	return &cosClient{
		cos: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  config.SecretID,
				SecretKey: config.SecretKey,
			},
		}),
	}

}

func (c *cosClient) ListBuckets(ctx context.Context) (buckets []string, err error) {
	s, _, err := c.cos.Service.Get(ctx)
	if err != nil {
		return nil, err
	}
	for _, b := range s.Buckets {
		buckets = append(buckets, b.Name)
	}
	return
}

func (c *cosClient) UploadFromFile(ctx context.Context, cosFilePath, localFilePath string) error {
	_, err := c.cos.Object.PutFromFile(ctx, cosFilePath, localFilePath, nil)
	return err
}

func (c *cosClient) UploadFromBytes(ctx context.Context, cosFilePath string, body []byte) error {
	reader := bytes.NewReader(body)
	_, err := c.cos.Object.Put(ctx, cosFilePath, reader, nil)
	return err
}

func (c *cosClient) Download2File(ctx context.Context, filename, localpath string) (err error) {
	_, err = c.cos.Object.GetToFile(ctx, filename, filepath.Base(localpath), nil)
	return err
}
