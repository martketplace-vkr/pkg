package migrate

import (
	"embed"

	"github.com/cockroachdb/errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type PgMigrateConfig struct {
	FS      embed.FS
	MigPath string
}

func Migrate(
	fs embed.FS,
	migPath string,
	pgDSN string,
) (*migrate.Migrate, error) {
	d, err := iofs.New(fs, migPath)
	if err != nil {
		return nil, errors.Wrap(err, "embed postgres migrations")
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, pgDSN)
	if err != nil {
		return nil, errors.Wrap(err, "apply postgres migrations")
	}

	return m, nil
}
