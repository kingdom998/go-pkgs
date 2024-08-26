package cos

import (
	"context"
	"fmt"
	"os"
	"testing"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kingdom998/go-pkgs/conf"
)

var (
	config conf.COS
)

func init() {
	config = conf.COS{
		Host:      os.Getenv("host"),
		Region:    os.Getenv("region"),
		Bucket:    os.Getenv("bucket"),
		SecretID:  os.Getenv("secret_id"),
		SecretKey: os.Getenv("secret_key"),
	}
}

func TestListBuckets(t *testing.T) {
	ctx := context.Background()
	client := New(config)
	buckets, err := client.ListBuckets(ctx)
	if err != nil {
		t.Error(err)
	}
	for _, bucket := range buckets {
		fmt.Println(bucket)
	}

}
