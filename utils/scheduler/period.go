package scheduler

import (
	"fmt"
	"time"

	"github.com/martketplace-vkr/pkg/utils/duration"
)

func Period(t time.Duration) string {
	return fmt.Sprintf("@every %s", duration.Format(t))
}
