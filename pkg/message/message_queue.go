package message

import (
	"context"
	"errors"
	"sync"

	"github.com/mixdone/uptime-monitoring/pkg/logger"
)

type localMQ struct {
	queues map[string]chan []byte
	mutex  sync.RWMutex
	closed bool
	log    logger.Logger
}

func NewLocalMQ(log logger.Logger) MQ {
	return &localMQ{
		queues: make(map[string]chan []byte),
		log:    log.WithField("component", "message queue"),
	}
}

func (mq *localMQ) GetChan(name string) (chan []byte, error) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	if mq.closed {
		return nil, errors.New("message queues closed")
	}

	if ch, ok := mq.queues[name]; ok {
		return ch, nil
	}

	ch := make(chan []byte, 100)
	mq.queues[name] = ch
	return ch, nil
}

func (mq *localMQ) Publish(queue string, body []byte) error {
	ch, err := mq.GetChan(queue)
	if err != nil {
		return err
	}

	ch <- body
	return nil

}
func (mq *localMQ) Consume(ctx context.Context, queue string, handler func([]byte) error) error {
	ch, err := mq.GetChan(queue)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				err = handler(msg)
				mq.log.Error(err)
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (mq *localMQ) Close() error {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	if mq.closed {
		return errors.New("already closed")
	}

	for _, q := range mq.queues {
		close(q)
	}

	mq.closed = true
	return nil
}
