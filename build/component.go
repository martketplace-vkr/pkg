package build

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/martketplace-vkr/pkg/logger/log"
)

type (
	Components []Component

	Component interface {
		// Start starts a component.
		Start(ctx context.Context) (err error)
		// Stop stops a component.
		Stop(ctx context.Context) (err error)
		// GetStartTimeout returns the start timeout.
		// each component must have start timeout. Otherwise, startup will be blocked.
		GetStartTimeout() time.Duration
		// GetStopTimeout returns the stop timeout.
		// each component must have stop timeout. Otherwise the graceful shutdown will be interrupted.
		GetStopTimeout() time.Duration
		// GetShutdownDelay returns the shutdown delay. It is used to wait for the component to stop.
		// In order to finish some work before stopping.
		// OPTIONAL. Default is 0.
		GetShutdownDelay() time.Duration
		// GetName returns the component name.
		GetName() string
	}
)

func (cmp *Components) StopAll(ctx context.Context) error {
	log.Info("stopping all components...")
	fmt.Println()

	for _, c := range *cmp {
		shutdownDelay := c.GetShutdownDelay()
		if shutdownDelay != 0 {
			log.Infof("waiting for component %s to stop with delay [%v]...", c.GetName(), shutdownDelay.String())
			time.Sleep(shutdownDelay)
		}

		stopCtx, cancel := context.WithTimeout(ctx, c.GetStopTimeout())
		errCh := make(chan error, 1)

		go func() {
			errCh <- c.Stop(ctx)
		}()

		select {
		case <-stopCtx.Done():
			cancel()
			return fmt.Errorf("timeout while stopping component %s", c.GetName())
		case err := <-errCh:
			cancel()
			if err != nil {
				return fmt.Errorf("error stopping component %s: %w", c.GetName(), err)
			}
			log.Infof("component %s stopped", c.GetName())
		}
	}

	return nil
}

func (cmp *Components) StartAll(ctx context.Context) error {
	fmt.Println()
	log.Info("starting all components...")
	fmt.Println()

	for _, c := range *cmp {
		log.Infof("component %s starting...", c.GetName())

		startCtx, cancel := context.WithTimeout(ctx, c.GetStartTimeout())
		errCh := make(chan error, 1)

		go func() {
			errCh <- c.Start(ctx)
		}()

		select {
		case <-startCtx.Done():
			cancel()
			return fmt.Errorf("timeout while starting component %s", c.GetName())
		case err := <-errCh:
			cancel()
			if err != nil {
				log.Infof("component %s failed to start ❌: %s", c.GetName(), err.Error())
				return fmt.Errorf("error starting component %s: %w", c.GetName(), err)
			}
			log.Infof("component %s started 🥇", c.GetName())
			fmt.Println()
		}
	}

	log.Infof("all components started, happy hunger games 💋")
	return nil
}

func (cmp *Components) GetComponentByName(cmpName string) Component {
	index := slices.IndexFunc(*cmp, func(c Component) bool {
		return c.GetName() == cmpName
	})
	if index == -1 {
		return nil
	}

	return (*cmp)[index]
}
