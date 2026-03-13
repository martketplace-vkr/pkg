package pgxpoolcomponent

import (
	"github.com/martketplace-vkr/pkg/build/components"
	pgxpoolconnctor "github.com/martketplace-vkr/pkg/pgxpoolconnector"
	"github.com/martketplace-vkr/pkg/utils/migrate"
)

type Config struct {
	components.ComponentConfig `validate:"required"`
	pgxpoolconnctor.Config     `validate:"required"`
	MigrateConfig              *migrate.PgMigrateConfig
}
