package mq

import "context"

type Queue interface {
	Publish(ctx context.Context, msg []byte) error
	Subscribe(ctx context.Context, callback func(context.Context, []byte) error) error
}
