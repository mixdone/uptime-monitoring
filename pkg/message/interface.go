package message

import "context"

type MQ interface {
	Publish(queue string, body []byte) error
	Consume(ctx context.Context, queue string, handler func([]byte) error) error
	Close() error
}
