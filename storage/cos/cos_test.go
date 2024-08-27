package cos

import (
	"context"
	"fmt"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
)

var (
	config Config
	client *storage
)

func init() {
	config = Config{
		Host:      os.Getenv("host"),
		Region:    os.Getenv("region"),
		Bucket:    os.Getenv("bucket"),
		SecretID:  os.Getenv("secret_id"),
		SecretKey: os.Getenv("secret_key"),
	}
	client = NewClient(config)
}

func TestListBuckets(t *testing.T) {
	ctx := context.Background()
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		t.Error(err)
	}
	for _, bucket := range buckets {
		fmt.Println(bucket)
	}
}

func TestUploadFromFile(t *testing.T) {
	ctx := context.Background()
	localFilePath := "cos_test.go"
	objectPath := "Stable-diffusion/" + localFilePath
	result, err := client.UploadFromFile(ctx, objectPath, localFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestUploadFromBytes(t *testing.T) {
	ctx := context.Background()
	localFilePath := "cos_test.go"
	objectPath := "Stable-diffusion/" + localFilePath
	body, err := os.ReadFile(localFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	result, err := client.UploadFromBytes(ctx, objectPath, body)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestDownloadFile(t *testing.T) {
	ctx := context.Background()
	client := NewClient(config)
	filepath := "webui/011c0b9e-dd65-47d1-b8d7-df0e708c1401.png"
	err := client.Download2File(ctx, filepath, "ai.png")
	if err != nil {
		t.Error(err)
	}
}
