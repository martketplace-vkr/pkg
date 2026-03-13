package tracer

import (
	"github.com/martketplace-vkr/pkg/build/components"
	"github.com/martketplace-vkr/pkg/tracer"
)

type Config struct {
	components.ComponentConfig `validate:"required"`
	tracer.Config              `validate:"required"`
}
