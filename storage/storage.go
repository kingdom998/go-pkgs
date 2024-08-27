package storage

import "context"

type Storage interface {
	ListBuckets(ctx context.Context) (buckets []string, err error)
	UploadFromFile(ctx context.Context, cosFilePath, localFilePath string) error
	UploadFromBytes(ctx context.Context, objectPath string, body []byte) error
	Download2File(ctx context.Context, objectPath, localpath string) (err error)
}
