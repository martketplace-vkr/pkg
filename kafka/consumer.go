package kafka

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"
	"sync"

	"github.com/IBM/sarama"
	"github.com/martketplace-vkr/pkg/logger/log"
)

type (
	ConsumerGroup interface {
		ConsumeFunc(ctx context.Context, topics []string, handlerFunc HandlerFunc, options ...ConsumerHandlerOption) error
		Close() error
	}

	HandlerFunc func(context.Context, *sarama.ConsumerMessage) error

	ConsumerHandlerOption func(c *ConsumerHandler)

	closerError struct {
		err []error
	}

	ConsumerHandler struct {
		Handle    HandlerFunc
		OnSetup   []func(sarama.ConsumerGroupSession) error
		OnCleanup []func(sarama.ConsumerGroupSession) error
	}

	consumerGroup struct {
		closers      []io.Closer
		saramaConfig *sarama.Config
		addresses    []string
		groupID      string
	}
)

func NewConsumerGroup(
	addresses []string,
	config *sarama.Config,
	groupID string,
) ConsumerGroup {
	config.Consumer.Return.Errors = true

	return &consumerGroup{
		saramaConfig: config,
		addresses:    addresses,
		groupID:      groupID,
	}
}

func (c *consumerGroup) ConsumeFunc(
	ctx context.Context,
	topics []string,
	handlerFunc HandlerFunc,
	options ...ConsumerHandlerOption,
) error {
	saramaConsumer, err := sarama.NewConsumerGroup(c.addresses, c.groupID, c.saramaConfig)
	if err != nil {
		return fmt.Errorf("failed to create kafka consumer group from client: %w", err)
	}

	c.appendCloser(saramaConsumer)

	go func() {
		for err := range saramaConsumer.Errors() {
			log.Errorf("consumer group error: %v", err)
		}
	}()

	handler := &ConsumerHandler{Handle: handlerFunc}

	for _, option := range options {
		option(handler)
	}

	go func() {
		for {
			if err := saramaConsumer.Consume(ctx, topics, handler); err != nil {
				log.Errorf("failed to consume by consumer group %s: %v", c.groupID, err)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

func (c *consumerGroup) Close() error {
	errs := make([]error, 0, len(c.closers))

	for _, closer := range c.closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return &closerError{err: errs}
	}

	return nil
}

func (c *consumerGroup) appendCloser(closer io.Closer) {
	c.closers = append(c.closers, closer)
}

func (h *ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	for _, onSetup := range h.OnSetup {
		if err := onSetup(session); err != nil {
			return err
		}
	}

	return nil
}

func (h *ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	session.Commit()

	for _, onCleanup := range h.OnCleanup {
		if err := onCleanup(session); err != nil {
			return err
		}
	}

	return nil
}

func (h *ConsumerHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			switch v := r.(type) {
			case string:
				err = fmt.Errorf("panic: %s\n%s", v, stack)
			case error:
				err = fmt.Errorf("panic: %w\n%s", v, stack)
			default:
				err = fmt.Errorf("panic: %v\n%s", v, stack)
			}
		}
	}()

	for message := range claim.Messages() {
		session.MarkMessage(message, "")
		if err := h.Handle(session.Context(), message); err != nil {
			log.Errorf("failed to handle message: %v", err)

			continue
		}

	}

	return err
}

func (c closerError) Error() string {
	return fmt.Sprintf("closers error: %v", c.err)
}

func WithResetToOldestOffset(partition int32, offset int64, topic, meta string) ConsumerHandlerOption {
	once := sync.Once{}

	return func(h *ConsumerHandler) {
		h.OnSetup = append(h.OnSetup, func(session sarama.ConsumerGroupSession) error {
			once.Do(func() {
				session.ResetOffset(topic, partition, offset, meta)
				session.Commit()
			})

			return nil
		})
	}
}

func WithResetToNewestOffset(partition int32, offset int64, topic, meta string) ConsumerHandlerOption {
	once := sync.Once{}

	return func(h *ConsumerHandler) {
		h.OnSetup = append(h.OnSetup, func(session sarama.ConsumerGroupSession) error {
			once.Do(func() {
				session.MarkOffset(topic, partition, offset, meta)
				session.Commit()
			})

			return nil
		})
	}
}
