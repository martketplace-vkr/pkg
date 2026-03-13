package kafka

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/IBM/sarama"
	"github.com/martketplace-vkr/pkg/inbox/dto"
	"github.com/martketplace-vkr/pkg/logger/log"
)

type Provider struct {
	cfg        MessageReceiverConfig
	cgFactory  func(groupID string) (sarama.ConsumerGroup, error)
	eventsChan chan dto.Event
	metrics    *kafkaMetrics
}

func New(
	cfg MessageReceiverConfig,
	cgFactory func(groupID string) (sarama.ConsumerGroup, error),
) *Provider {
	p := &Provider{
		cfg:        cfg,
		eventsChan: make(chan dto.Event, 1024),
		cgFactory:  cgFactory,
		// TODO сделать счтобы можно было не дефолт реджистри.
		metrics: newKafkaMetrics(cfg.MetricsNS, prometheus.DefaultRegisterer),
	}

	return p
}

func (p *Provider) GetEventsChan() chan dto.Event { return p.eventsChan }

func (p *Provider) Consume(ctx context.Context) {
	if p.cgFactory == nil {
		log.Errorf("fail to create kafka consumer factory, factory is not provided")
		return
	}

	var wg sync.WaitGroup
	errCh := make(chan error, p.cfg.ConsumersCount)

	startWorker := func(id int, cg sarama.ConsumerGroup) {
		defer wg.Done()
		defer func() {
			if cg != nil {
				_ = cg.Close()
			}
		}()

		h := &groupHandler{
			eventsChan:  p.eventsChan,
			metrics:     p.metrics,
			debugLog:    p.cfg.Debug,
			maxLogValue: p.cfg.MaxLogValueBytes,
			workerID:    id,
		}

		for {
			if ctx.Err() != nil {
				return
			}
			if err := cg.Consume(ctx, p.cfg.Topics, h); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) || errors.Is(err, context.Canceled) {
					log.Infof("kafka(worker=%d): consumer group closed; stopping", id)
					return
				}
				if p.metrics != nil {
					p.metrics.incKafkaConsumeErrors()
				}

				errCh <- fmt.Errorf("kafka(worker=%d): consume error: %w", id, err)
				log.Errorf("kafka(worker=%d): consume error: %v", id, err)
				time.Sleep(time.Duration(p.cfg.ConsumerJitterMs))
			}
		}
	}

	for i := 0; i < p.cfg.ConsumersCount; i++ {
		var cg sarama.ConsumerGroup
		var err error

		cg, err = p.cgFactory(p.cfg.GroupID)
		if err != nil {
			log.Errorf("kafka: failed to create consumer group for worker %d: %v", i, err)

			continue
		}

		wg.Add(1)
		go startWorker(i, cg)
	}

	go func() {
		wg.Wait()
		close(errCh)
		close(p.eventsChan)
	}()

	for err := range errCh {
		log.Errorf("error of kafka consume: %s", err.Error())
	}
}
