package rediscomponent

import (
	"github.com/martketplace-vkr/pkg/build/components"
	"github.com/martketplace-vkr/pkg/redisconnector"
)

type Config struct {
	components.ComponentConfig `validate:"required"`
	redisconnector.Config      `validate:"required"`
}
