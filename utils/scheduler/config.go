package scheduler

import "github.com/martketplace-vkr/pkg/utils/location"

type Config struct {
	Location *location.Location `validate:"required" default:"UTC"`
}
