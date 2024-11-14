package oss

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
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
	log.Infof("config is : %#v", config)

	client = NewClient(&config)
	log.Infof("config is : %#v", config)
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

func TestUploadModels(t *testing.T) {
	localDir := os.Getenv("localModelsDir")
	ossDir := os.Getenv("ossModelsDir")
	err := filepath.Walk(localDir, func(localFilePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relativePath, inErr := filepath.Rel(localDir, localFilePath)
			if inErr != nil {
				return inErr
			}
			cosFilePath := filepath.Join(ossDir, relativePath)
			cosFilePath = filepath.ToSlash(cosFilePath) // 转换为 UNIX 风格路径
			fmt.Println("uploading ", localFilePath)
			_, err = client.UploadFromFile(context.Background(), cosFilePath, localFilePath)
			if err != nil {
				return fmt.Errorf("failed to upload %s to %s: %w", localFilePath, cosFilePath, err)
			}
			log.Infof("Uploaded %s to %s successfully.", localFilePath, cosFilePath)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("upload objects failed with %v", err)
	}
}
