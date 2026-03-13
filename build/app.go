package build

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	unknown    = "unknown"
	versionEnv = "CI_COMMIT_TAG"
	comitEnv   = "CI_COMMIT_MESSAGE"
)

var (
	appBuild = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_build",
			Help: "app build info",
		}, []string{
			"version",
			"built_at",
			"commit",
			"domain",
			"team",
			"service_name",
		},
	)
)

type (
	App interface {
		Metadata() *Metadata
		Components() Components
	}
	app struct {
		metadata     *Metadata
		components   Components
		metadataPath string
		registry     prometheus.Registerer
	}
)

// Run starts an application as the graceful shutdown service.
func Run(ctx context.Context, a App) error {
	metadata := a.Metadata()
	if metadata == nil {
		return fmt.Errorf("project metadata were not specified")
	}

	cmps := a.Components()

	if err := cmps.StartAll(ctx); err != nil {
		return err
	}

	extraMetadata := loadExtraMetadata()

	appMeta := a.Metadata()
	appBuild.WithLabelValues(
		extraMetadata.Version,
		time.Now().String(),
		extraMetadata.CommitMessage,
		appMeta.Project.Domain,
		appMeta.Project.Team,
		appMeta.Project.Name,
	)

	// wait for OS signal for graceful shutdown
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	interruptSignal := <-quitCh

	log.Infof("Received signal: %s", interruptSignal)

	return cmps.StopAll(ctx)
}

func (app *app) Metadata() *Metadata {
	if app.metadata == nil {
		app.metadata = &Metadata{
			Project: Project{
				Name:   unknown,
				Team:   unknown,
				Domain: unknown,
			},
		}
	}

	return app.metadata
}

func (app *app) Components() Components {
	return app.components
}

func NewApp(cmp Components, opts ...Option) (*app, error) {
	const defaultMetadata = "./metadata.toml"

	if cmp == nil {
		cmp = Components{}
	}

	a := &app{
		components:   cmp,
		metadataPath: defaultMetadata,
		registry:     prometheus.DefaultRegisterer,
	}

	for _, opt := range opts {
		opt(a)
	}

	if err := a.registry.Register(appBuild); err != nil {
		return nil, fmt.Errorf("failed to register build info metric: %w", err)
	}

	if err := a.loadMetadata(); err != nil {
		log.Warnf("starting build with no metadata...")
	}

	return a, nil
}

func (app *app) loadMetadata() error {
	data, err := os.ReadFile(app.metadataPath)
	if err != nil {
		return fmt.Errorf("error reading metadata file %s: %w", app.metadataPath, err)
	}

	app.metadata = &Metadata{}

	if err := toml.Unmarshal(data, app.metadata); err != nil {
		return fmt.Errorf("could not load metadata: %w", err)
	}

	return nil
}
