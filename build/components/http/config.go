package http

import (
	"github.com/martketplace-vkr/pkg/build/components"
	srv "github.com/martketplace-vkr/pkg/server/http"
)

type Config struct {
	srv.Config
	components.ComponentConfig
}
