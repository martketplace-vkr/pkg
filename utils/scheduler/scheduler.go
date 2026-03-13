package scheduler

import (
	"github.com/robfig/cron"
)

type (
	Scheduler interface {
		Run() error
		Shutdown() error
		AddFunc(format string, cmd func()) error
	}

	// scheduler is a struct that represents a scheduler
	// wrapper around cron.Cron
	scheduler struct {
		cfg  Config
		cron *cron.Cron
	}
)

func New(cfg Config) *scheduler {
	worker := cron.New()
	if cfg.Location != nil {
		worker = cron.NewWithLocation(cfg.Location.Location)
	}

	return &scheduler{
		cfg:  cfg,
		cron: worker,
	}
}

func (s *scheduler) Run() error {
	s.cron.Start()

	return nil
}

func (s *scheduler) Shutdown() error {
	s.cron.Stop()

	return nil
}

func (s *scheduler) AddFunc(format string, cmd func()) error {
	return s.cron.AddFunc(format, cmd)
}
