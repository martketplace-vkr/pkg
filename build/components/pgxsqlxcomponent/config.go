package pgxsqlxcomponent

import (
	"github.com/martketplace-vkr/pkg/build/components"
	"github.com/martketplace-vkr/pkg/pgxsqlxconnector"
	"github.com/martketplace-vkr/pkg/utils/migrate"
)

type Config struct {
	components.ComponentConfig `validate:"required"`
	pgxsqlxconnector.Config    `validate:"required"`
	MigrateConfig              *migrate.PgMigrateConfig
}
