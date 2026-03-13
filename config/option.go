package config

import "github.com/martketplace-vkr/pkg/config/vault"

const (
	defaultConfigPath = "./config/config.json"
)

type cfgOptions struct {
	cfgPath         string
	disableDefaults bool
	vaultOpts       []vault.Option
}

type Option interface {
	apply(*cfgOptions)
}

type cfgPathOption struct {
	cfgPath string
}

func (o *cfgPathOption) apply(opts *cfgOptions) {
	opts.cfgPath = o.cfgPath
}

func WithConfigPath(cfgPath string) Option {
	return &cfgPathOption{cfgPath: cfgPath}
}

type defaultValuesOpt struct {
	disableDefaults bool
}

func (o *defaultValuesOpt) apply(opts *cfgOptions) {
	opts.disableDefaults = o.disableDefaults
}

func DisableDefaults() Option {
	return &defaultValuesOpt{disableDefaults: true}
}

type vaultCfgOption struct {
	vaultOpts []vault.Option
}

func (o *vaultCfgOption) apply(opts *cfgOptions) {
	opts.vaultOpts = o.vaultOpts
}

func WithVaultOptions(vaultOpts ...vault.Option) Option {
	return &vaultCfgOption{vaultOpts: vaultOpts}
}

func newDefaultConfigOptions() *cfgOptions {
	return &cfgOptions{
		cfgPath:         defaultConfigPath,
		disableDefaults: false,
	}
}
