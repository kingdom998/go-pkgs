package oss

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
		Endpoint: os.Getenv("endpoint"),
		Ak:       os.Getenv("ak"),
		Sk:       os.Getenv("sk"),
		Bucket:   os.Getenv("bucket"),
	}

	client = NewClient(&config)
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
	localFilePath := "oss_test.go"
	objectPath := "Stable-diffusion/" + localFilePath
	result, err := client.UploadFromFile(ctx, objectPath, localFilePath)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(result)
}

func TestDownloadFile(t *testing.T) {
	ctx := context.Background()
	localFilePath := "oss_test.go"
	objectPath := "Stable-diffusion/" + localFilePath
	err := client.Download2File(ctx, objectPath, "../"+localFilePath)
	if err != nil {
		t.Error(err)
	}
}
