package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler_Run(t *testing.T) {
	cfg := Config{}
	s := New(cfg)

	// Call the Run method and assert that it starts the cron service
	err := s.Run()
	assert.NoError(t, err, "expected no error when running the scheduler")
	assert.NotNil(t, s.cron, "cron object should not be nil")
	s.Shutdown()
}

func TestScheduler_Shutdown(t *testing.T) {
	cfg := Config{}
	s := New(cfg)

	// Start the cron and then shut it down
	s.Run()
	err := s.Shutdown()
	assert.NoError(t, err, "expected no error when shutting down the scheduler")
}

func TestScheduler_AddFunc(t *testing.T) {
	cfg := Config{}
	s := New(cfg)

	triggeredCount := 0
	cmd := func() { triggeredCount++ }

	// Add a function to be scheduled immediately
	err := s.AddFunc(Period(time.Second), cmd)
	assert.NoError(t, err, "expected no error when adding a function to scheduler")

	// Run the scheduler and wait to see if the command is triggered
	s.Run()
	time.Sleep(2 * time.Second)

	// Assert the command was triggered
	assert.Equal(t, triggeredCount, 2)

	s.Shutdown()
}
