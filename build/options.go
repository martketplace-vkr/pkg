package build

import "github.com/prometheus/client_golang/prometheus"

type Option func(*app)

func WithMetadataPath(path string) Option {
	return func(a *app) {
		a.metadataPath = path
	}
}

func WithPromRegistry(reg prometheus.Registerer) Option {
	return func(a *app) {
		a.registry = reg
	}
}
