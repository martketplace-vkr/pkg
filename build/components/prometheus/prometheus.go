package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	cmpName = "prometheus"
)

type Prometheus struct {
	cfg      Config
	mux      *http.ServeMux
	registry prometheus.Registerer
}

func New(
	cfg Config,
	registry prometheus.Registerer,
) *Prometheus {
	if registry == nil {
		registry = prometheus.DefaultRegisterer
	}

	return &Prometheus{
		cfg:      cfg,
		mux:      http.NewServeMux(),
		registry: registry,
	}
}

func (p Prometheus) Start(ctx context.Context) (err error) {
	go func() {
		p.mux.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}))

		err := http.ListenAndServe(p.cfg.Host, p.mux)
		if err != nil {
			log.Fatalf("failed to listen: %v", err.Error())
		}

		log.Infof("Metrics server is running on port: %s", p.cfg.Host)
		fmt.Println("Metrics server is running on port: ", p.cfg.Host)
	}()

	return err
}

func (p Prometheus) Stop(_ context.Context) (err error) {
	return err
}

func (p Prometheus) GetStartTimeout() time.Duration {
	return p.cfg.StartTimeout.Duration
}

func (p Prometheus) GetStopTimeout() time.Duration {
	return p.cfg.StopTimeout.Duration
}

func (p Prometheus) GetShutdownDelay() time.Duration {
	return p.cfg.ShutdownDelay.Duration
}

func (p Prometheus) GetName() string {
	return cmpName
}
