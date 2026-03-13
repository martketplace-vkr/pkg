package redisconnector

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg Config) (client *redis.Client, err error) {
	var (
		opts   *redis.Options
		tlsCfg *tls.Config
	)

	if cfg.UseCertificates {
		tlsCfg, err = setupTLS(cfg)
		if err != nil {
			return client, errors.Wrap(err, "setupTLS")
		}
	}

	opts = &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		Password:     cfg.Password,
		DB:           cfg.DB,
		TLSConfig:    tlsCfg,
	}

	client = redis.NewClient(opts)
	if err = client.Ping(context.Background()).Err(); err != nil {
		return nil, errors.Wrapf(err, "ping")
	}

	if cfg.WithMetricsHook {
		client.AddHook(NewHook())
	}

	if cfg.EnableTracing {
		if err := redisotel.InstrumentTracing(client); err != nil {
			return client, errors.Wrap(err, "instrument tracing")
		}
	}

	return client, nil
}

func setupTLS(cfg Config) (*tls.Config, error) {
	certs := make([]tls.Certificate, 0, 0)
	if cfg.CertificatesPaths.Cert != "" && cfg.CertificatesPaths.Key != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertificatesPaths.Cert, cfg.CertificatesPaths.Key)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"certPath: %v, keyPath: %v",
				cfg.CertificatesPaths.Cert,
				cfg.CertificatesPaths.Key,
			)
		}

		certs = append(certs, cert)
	}

	caCert, err := os.ReadFile(cfg.CertificatesPaths.Ca)
	if err != nil {
		return nil, errors.Wrapf(err, "ca load path: %v", cfg.CertificatesPaths.Ca)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		InsecureSkipVerify: cfg.InsecureSkipVerify,
		Certificates:       certs,
		RootCAs:            caCertPool,
	}, nil
}
