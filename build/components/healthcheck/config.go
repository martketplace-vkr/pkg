package healthcheck

import "github.com/martketplace-vkr/pkg/build/components"

type Config struct {
	Address string `validate:"required" default:":8000"`
	components.ComponentConfig
}
